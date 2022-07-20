// Copyright 2019-2022 Charles Korn.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// and the Commons Clause License Condition v1.0 (the "Condition");
// you may not use this file except in compliance with both the License and Condition.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// You may obtain a copy of the Condition at
//
//     https://commonsclause.com/
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License and the Condition is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See both the License and the Condition for the specific language governing permissions and
// limitations under the License and the Condition.

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	cloudstorage "cloud.google.com/go/storage"
	"github.com/batect/services-common/graceful"
	"github.com/batect/services-common/middleware"
	"github.com/batect/services-common/startup"
	"github.com/batect/services-common/tracing"
	"github.com/batect/updates.batect.dev/server/api"
	"github.com/batect/updates.batect.dev/server/events"
	"github.com/batect/updates.batect.dev/server/storage"
	"github.com/sirupsen/logrus"
	"github.com/unrolled/secure"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"google.golang.org/api/option"
	htransport "google.golang.org/api/transport/http"
)

func main() {
	config, err := getConfig()

	if err != nil {
		logrus.WithError(err).Error("Could not load application configuration.")
		os.Exit(1)
	}

	flush, err := startup.InitialiseObservability(config.ServiceName, config.ServiceVersion, config.ProjectID, config.HoneycombAPIKey)

	if err != nil {
		logrus.WithError(err).Error("Could not initialise observability tooling.")
		os.Exit(1)
	}

	defer flush()

	runServer(config)
}

func runServer(config *serviceConfig) {
	srv, err := createServer(config)

	if err != nil {
		logrus.WithError(err).Error("Could not create server.")
		os.Exit(1)
	}

	if err := graceful.RunServerWithGracefulShutdown(srv); err != nil {
		logrus.WithError(err).Error("Could not run server.")
		os.Exit(1)
	}
}

func createServer(config *serviceConfig) (*http.Server, error) {
	cloudStorageClient, err := createCloudStorageClient()

	if err != nil {
		return nil, fmt.Errorf("could not create Cloud Storage client: %w", err)
	}

	eventSink := createEventSink(cloudStorageClient, config)

	mux := http.NewServeMux()
	mux.Handle("/", otelhttp.WithRouteTag("/", http.HandlerFunc(api.Home)))
	mux.Handle("/ping", otelhttp.WithRouteTag("/ping", http.HandlerFunc(api.Ping)))
	mux.Handle("/v1/latest", otelhttp.WithRouteTag("/v1/latest", createLatestHandler(cloudStorageClient, eventSink, config)))
	mux.Handle("/v1/files/", otelhttp.WithRouteTag("/v1/files", api.NewFilesHandler(eventSink)))

	securityHeaders := secure.New(secure.Options{
		FrameDeny:             true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'none'; frame-ancestors 'none'",
		ReferrerPolicy:        "no-referrer",
	})

	wrappedMux := middleware.TraceIDExtractionMiddleware(
		middleware.LoggerMiddleware(
			logrus.StandardLogger(),
			config.ProjectID,
			securityHeaders.Handler(mux),
		),
	)

	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", config.Port),
		Handler: otelhttp.NewHandler(
			wrappedMux,
			"Updates API",
			otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents),
			otelhttp.WithSpanNameFormatter(tracing.NameHTTPRequestSpan),
		),
		ReadHeaderTimeout: 10 * time.Second,
	}

	return srv, nil
}

func createEventSink(cloudStorageClient *cloudstorage.Client, config *serviceConfig) events.EventSink {
	bucketName := fmt.Sprintf("%v-events", config.ProjectID)

	return events.NewCloudStorageEventSink(bucketName, cloudStorageClient)
}

func createLatestHandler(cloudStorageClient *cloudstorage.Client, eventSink events.EventSink, config *serviceConfig) http.Handler {
	bucketName := fmt.Sprintf("%v-public", config.ProjectID)
	store := storage.NewCloudStorageLatestVersionStore(bucketName, cloudStorageClient)

	return api.NewLatestHandler(store, eventSink)
}

func createCloudStorageClient() (*cloudstorage.Client, error) {
	scopesOption := option.WithScopes(cloudstorage.ScopeReadWrite)
	credsOption := option.WithCredentialsFile(getCredentialsFilePath())
	tracingClientOption, err := withTracingClient(scopesOption, credsOption)

	if err != nil {
		return nil, fmt.Errorf("could not create tracing client: %w", err)
	}

	cloudStorageClient, err := cloudstorage.NewClient(context.Background(), tracingClientOption)

	if err != nil {
		return nil, fmt.Errorf("could not create Cloud Storage client: %w", err)
	}

	return cloudStorageClient, nil
}

func withTracingClient(opts ...option.ClientOption) (option.ClientOption, error) {
	// We have to do this because setting http.DefaultTransport to a non-default implementation causes something deep in the bowels of the
	// Google Cloud SDK to ignore it and create a fresh transport with many of the settings copied across from DefaultTransport.
	// Being explicit about the client forces the SDK to use the transport.
	trans, err := htransport.NewTransport(context.Background(), http.DefaultTransport, opts...)

	if err != nil {
		return nil, fmt.Errorf("could not create transport: %w", err)
	}

	httpClient := http.Client{
		Transport: trans,
	}

	return option.WithHTTPClient(&httpClient), nil
}
