// Copyright 2019-2023 Charles Korn.
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

package events

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"time"

	cloudstorage "cloud.google.com/go/storage"
	"github.com/batect/services-common/middleware"
	"github.com/google/uuid"
)

type cloudStorageEventSink struct {
	client     *cloudstorage.Client
	bucket     *cloudstorage.BucketHandle
	timeSource func() time.Time
	uuidSource func() uuid.UUID
}

func NewCloudStorageEventSink(bucketName string, client *cloudstorage.Client) EventSink {
	timeSource := func() time.Time { return time.Now().UTC() }

	return NewCloudStorageEventSinkWithSpecificDependencies(bucketName, client, timeSource, uuid.New)
}

func NewCloudStorageEventSinkWithSpecificDependencies(bucketName string, client *cloudstorage.Client, timeSource func() time.Time, uuidSource func() uuid.UUID) EventSink {
	return &cloudStorageEventSink{
		client:     client,
		bucket:     client.Bucket(bucketName),
		timeSource: timeSource,
		uuidSource: uuidSource,
	}
}

func (c *cloudStorageEventSink) PostLatestVersionCheck(ctx context.Context, userAgent string) {
	timestamp := c.timeSource()
	eventID := c.uuidSource()

	event := map[string]interface{}{
		"eventId":   eventID,
		"timestamp": timestamp,
		"userAgent": userAgent,
	}

	if err := c.postEvent(ctx, "v1/latest", timestamp, eventID, event); err != nil {
		log := middleware.LoggerFromContext(ctx)
		log.WithError(err).Error("Failed to post latest version check event.")
	}
}

func (c *cloudStorageEventSink) PostFileDownload(ctx context.Context, userAgent string, version string, fileName string) {
	timestamp := c.timeSource()
	eventID := c.uuidSource()

	event := map[string]interface{}{
		"eventId":   eventID,
		"timestamp": timestamp,
		"userAgent": userAgent,
		"version":   version,
		"fileName":  fileName,
	}

	if err := c.postEvent(ctx, "v1/files", timestamp, eventID, event); err != nil {
		log := middleware.LoggerFromContext(ctx)
		log.WithError(err).Error("Failed to post file download event.")
	}
}

func (c *cloudStorageEventSink) postEvent(ctx context.Context, prefix string, timestamp time.Time, eventID uuid.UUID, event map[string]interface{}) error {
	w := c.bucket.
		Object(fmt.Sprintf("%v/%v/%02d/%02d/%v.json", prefix, timestamp.Year(), timestamp.Month(), timestamp.Day(), eventID)).
		If(cloudstorage.Conditions{DoesNotExist: true}).
		NewWriter(ctx)

	w.ContentType = "application/json"
	w.ContentEncoding = "gzip"
	gzipper := gzip.NewWriter(w)

	bytes, err := json.Marshal(event)

	if err != nil {
		return fmt.Errorf("converting event to JSON failed: %w", err)
	}

	if _, err := gzipper.Write(bytes); err != nil {
		return fmt.Errorf("writing to Cloud Storage failed: %w", err)
	}

	if err := gzipper.Close(); err != nil {
		return fmt.Errorf("closing gzip stream failed: %w", err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("storing event in Cloud Storage failed: %w", err)
	}

	return nil
}
