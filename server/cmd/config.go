// Copyright 2019-2021 Charles Korn.
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
	"os"

	"github.com/sirupsen/logrus"
)

func getServiceName() string {
	return getEnvOrDefault("K_SERVICE", "updates")
}

func getVersion() string {
	return getEnvOrDefault("K_REVISION", "local")
}

func getEnvOrDefault(name string, fallback string) string {
	if value, ok := os.LookupEnv(name); ok {
		return value
	}

	return fallback
}

func getPort() string {
	return getEnvOrExit("PORT")
}

func getProjectID() string {
	return getEnvOrExit("GOOGLE_PROJECT")
}

func getCredentialsFilePath() string {
	variableName := "GOOGLE_APPLICATION_CREDENTIALS"
	value := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")

	if value == "" {
		logrus.WithField("variable", variableName).Info("Credentials file environment variable is not set, will fallback to default credential sources for GCP connections.")
	}

	return value
}

func getEnvOrExit(name string) string {
	value := os.Getenv(name)

	if value == "" {
		logrus.WithField("variable", name).Fatal("Environment variable is not set.")
	}

	return value
}
