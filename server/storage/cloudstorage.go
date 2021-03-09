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

package storage

import (
	"context"
	"fmt"
	"io/ioutil"

	cloudstorage "cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type cloudStorageLatestVersionStore struct {
	client *cloudstorage.Client
	bucket *cloudstorage.BucketHandle
}

func NewCloudStorageLatestVersionStore(bucketName string, opts ...option.ClientOption) (LatestVersionStore, error) {
	client, err := cloudstorage.NewClient(context.Background(), opts...)

	if err != nil {
		return nil, fmt.Errorf("could not create Cloud Storage client: %w", err)
	}

	store := cloudStorageLatestVersionStore{
		client: client,
		bucket: client.Bucket(bucketName),
	}

	return &store, nil
}

func (c *cloudStorageLatestVersionStore) GetLatestVersionDescriptor(ctx context.Context) (VersionDescriptor, error) {
	reader, err := c.bucket.Object("v1/latest.json").NewReader(ctx)

	if err != nil {
		return VersionDescriptor{}, fmt.Errorf("could not get latest version descriptor: %w", err)
	}

	defer reader.Close()

	content, err := ioutil.ReadAll(reader)

	if err != nil {
		return VersionDescriptor{}, fmt.Errorf("could not read file content: %w", err)
	}

	descriptor := VersionDescriptor{
		Content:     content,
		ContentType: reader.Attrs.ContentType,
	}

	return descriptor, nil
}
