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

package storage_test

import (
	"context"

	cloudstorage "cloud.google.com/go/storage"
	"github.com/batect/updates.batect.dev/server/storage"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/api/option"
)

var _ = Describe("Getting version information from Cloud Storage", func() {
	var bucket *cloudstorage.BucketHandle
	var store storage.LatestVersionStore

	BeforeEach(func() {
		project := "my-project"
		bucketName := "test-version-store-" + uuid.New().String()

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

		store = storage.NewCloudStorageLatestVersionStore(bucketName, client)
	})

	Describe("given the version information file does not exist in the bucket", func() {
		var descriptor storage.VersionDescriptor
		var err error

		BeforeEach(func() {
			descriptor, err = store.GetLatestVersionDescriptor(context.Background())
		})

		It("returns an empty version descriptor", func() {
			Expect(descriptor).To(Equal(storage.VersionDescriptor{}))
		})

		It("returns an appropriate error", func() {
			Expect(err).To(MatchError("could not get latest version descriptor: storage: object doesn't exist"))
		})
	})

	Describe("given the version information file exists in the bucket", func() {
		var descriptor storage.VersionDescriptor
		var err error

		BeforeEach(func() {
			w := bucket.Object("v1/latest.json").NewWriter(context.Background())
			w.ContentType = "application/json+descriptor"
			_, writeError := w.Write([]byte(`{"some":"descriptor"}`))
			Expect(writeError).ToNot(HaveOccurred())
			writeError = w.Close()
			Expect(writeError).ToNot(HaveOccurred())

			descriptor, err = store.GetLatestVersionDescriptor(context.Background())
		})

		It("returns a version descriptor with the details from the bucket", func() {
			Expect(descriptor).To(Equal(storage.VersionDescriptor{
				Content:     []byte(`{"some":"descriptor"}`),
				ContentType: "application/json+descriptor",
			}))
		})

		It("does not return an error", func() {
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
