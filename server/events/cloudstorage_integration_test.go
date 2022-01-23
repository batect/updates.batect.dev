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

package events_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"

	cloudstorage "cloud.google.com/go/storage"
	"github.com/batect/service-observability/middleware/testutils"
	"github.com/batect/updates.batect.dev/server/events"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gomegatypes "github.com/onsi/gomega/types"
	"github.com/sirupsen/logrus/hooks/test"
	"google.golang.org/api/option"
)

var _ = Describe("Posting events to Cloud Storage", func() {
	var bucket *cloudstorage.BucketHandle
	var sink events.EventSink

	BeforeEach(func() {
		project := "my-project"
		bucketName := "test-events-store-" + uuid.New().String()

		// Note that we also have to set the STORAGE_EMULATOR_HOST environment variable so that object downloads
		// are done from the correct host and over HTTP (rather than HTTPS).
		opts := []option.ClientOption{
			option.WithEndpoint("http://cloud-storage/storage/v1/"),
		}

		client, err := cloudstorage.NewClient(context.Background(), opts...)
		Expect(err).ToNot(HaveOccurred())

		bucket = client.Bucket(bucketName)
		err = bucket.Create(context.Background(), project, nil)
		Expect(err).ToNot(HaveOccurred())

		// I don't understand why, but if we reuse the same client for the upload and inspection as part of the tests,
		// the '/storage/v1' part of the endpoint configured above is dropped. So we have to recreate it.
		client, err = cloudstorage.NewClient(context.Background(), opts...)
		Expect(err).ToNot(HaveOccurred())

		timeSource := func() time.Time { return time.Date(2021, 3, 1, 9, 54, 40, 123456789, time.UTC) }
		uuidSource := func() uuid.UUID { return uuid.MustParse("11112222-3333-4444-5555-666677778888") }
		sink = events.NewCloudStorageEventSinkWithSpecificDependencies(bucketName, client, timeSource, uuidSource)
	})

	Context("posting latest version check events", func() {
		var hook *test.Hook

		BeforeEach(func() {
			ctx := context.Background()
			ctx, hook = testutils.ContextWithTestLogger(ctx)

			sink.PostLatestVersionCheck(ctx, "MyCoolThing/1.2.3")
		})

		It("logs no messages", func() {
			Expect(hook.Entries).To(BeEmpty())
		})

		It("stores the event in the bucket at the expected path", func() {
			Expect(bucket.Object("v1/latest/2021/03/01/11112222-3333-4444-5555-666677778888.json")).To(HaveContent(MatchJSON(`
				{
					"timestamp": "2021-03-01T09:54:40.123456789Z",
					"eventId": "11112222-3333-4444-5555-666677778888",
					"userAgent": "MyCoolThing/1.2.3"
				}
			`)))
		})

		It("stores the event in the bucket with the JSON media type", func() {
			Expect(bucket.Object("v1/latest/2021/03/01/11112222-3333-4444-5555-666677778888.json")).To(HaveContentType("application/json"))
		})

		It("stores the event in the bucket compressed", func() {
			Expect(bucket.Object("v1/latest/2021/03/01/11112222-3333-4444-5555-666677778888.json")).To(HaveContentEncoding("gzip"))
		})
	})

	Context("posting file download events", func() {
		var hook *test.Hook

		BeforeEach(func() {
			ctx := context.Background()
			ctx, hook = testutils.ContextWithTestLogger(ctx)

			sink.PostFileDownload(ctx, "MyCoolThing/1.2.3", "4.5.6", "batect-7.8.9.jar")
		})

		It("logs no messages", func() {
			Expect(hook.Entries).To(BeEmpty())
		})

		It("stores the event in the bucket at the expected path", func() {
			Expect(bucket.Object("v1/files/2021/03/01/11112222-3333-4444-5555-666677778888.json")).To(HaveContent(MatchJSON(`
				{
					"timestamp": "2021-03-01T09:54:40.123456789Z",
					"eventId": "11112222-3333-4444-5555-666677778888",
					"userAgent": "MyCoolThing/1.2.3",
					"version": "4.5.6",
					"fileName": "batect-7.8.9.jar"
				}
			`)))
		})

		It("stores the event in the bucket with the JSON media type", func() {
			Expect(bucket.Object("v1/files/2021/03/01/11112222-3333-4444-5555-666677778888.json")).To(HaveContentType("application/json"))
		})

		It("stores the event in the bucket compressed", func() {
			Expect(bucket.Object("v1/files/2021/03/01/11112222-3333-4444-5555-666677778888.json")).To(HaveContentEncoding("gzip"))
		})
	})
})

type haveContentMatcher struct {
	expectedContentMatcher gomegatypes.GomegaMatcher
	actualContent          string
}

func HaveContent(expectedContentMatcher gomegatypes.GomegaMatcher) gomegatypes.GomegaMatcher {
	return &haveContentMatcher{expectedContentMatcher, ""}
}

func (c *haveContentMatcher) Match(actual interface{}) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	reader, err := actual.(*cloudstorage.ObjectHandle).NewReader(ctx)

	if err != nil {
		return false, fmt.Errorf("could not get content of object: %w", err)
	}

	defer reader.Close()

	actualBytes, err := ioutil.ReadAll(reader)

	if err != nil {
		return false, fmt.Errorf("could not read content of object: %w", err)
	}

	c.actualContent = string(actualBytes)

	return c.expectedContentMatcher.Match(c.actualContent)
}

func (c *haveContentMatcher) FailureMessage(actual interface{}) string {
	return c.expectedContentMatcher.FailureMessage(c.actualContent)
}

func (c *haveContentMatcher) NegatedFailureMessage(actual interface{}) string {
	return c.expectedContentMatcher.NegatedFailureMessage(c.actualContent)
}

type haveContentTypeMatcher struct {
	expectedContentType string
	actualContentType   string
}

func HaveContentType(expectedContentType string) gomegatypes.GomegaMatcher {
	return &haveContentTypeMatcher{expectedContentType, ""}
}

func (c *haveContentTypeMatcher) Match(actual interface{}) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	attrs, err := actual.(*cloudstorage.ObjectHandle).Attrs(ctx)

	if err != nil {
		return false, fmt.Errorf("could not get attributes of object: %w", err)
	}

	c.actualContentType = attrs.ContentType

	return c.expectedContentType == c.actualContentType, nil
}

func (c *haveContentTypeMatcher) FailureMessage(actual interface{}) string {
	return fmt.Sprintf("Expected object '%v' to have content type '%v', but it was '%v'", actual.(*cloudstorage.ObjectHandle).ObjectName(), c.expectedContentType, c.actualContentType)
}

func (c *haveContentTypeMatcher) NegatedFailureMessage(actual interface{}) string {
	return fmt.Sprintf(
		"Expected object '%v' to not have content type '%v', but it was '%v'",
		actual.(*cloudstorage.ObjectHandle).ObjectName(),
		c.expectedContentType,
		c.actualContentType,
	)
}

type haveContentEncodingMatcher struct {
	expectedContentEncoding string
	actualContentEncoding   string
}

func HaveContentEncoding(expectedContentEncoding string) gomegatypes.GomegaMatcher {
	return &haveContentEncodingMatcher{expectedContentEncoding, ""}
}

func (c *haveContentEncodingMatcher) Match(actual interface{}) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	attrs, err := actual.(*cloudstorage.ObjectHandle).Attrs(ctx)

	if err != nil {
		return false, fmt.Errorf("could not get attributes of object: %w", err)
	}

	c.actualContentEncoding = attrs.ContentEncoding

	return c.expectedContentEncoding == c.actualContentEncoding, nil
}

func (c *haveContentEncodingMatcher) FailureMessage(actual interface{}) string {
	return fmt.Sprintf(
		"Expected object '%v' to have content encoding '%v', but it was '%v'",
		actual.(*cloudstorage.ObjectHandle).ObjectName(),
		c.expectedContentEncoding,
		c.actualContentEncoding,
	)
}

func (c *haveContentEncodingMatcher) NegatedFailureMessage(actual interface{}) string {
	return fmt.Sprintf(
		"Expected object '%v' to not have content encoding '%v', but it was '%v'",
		actual.(*cloudstorage.ObjectHandle).ObjectName(),
		c.expectedContentEncoding,
		c.actualContentEncoding,
	)
}
