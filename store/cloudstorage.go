// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package store

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/storage/v1"
)

// CloudStorageClient implements Client and handles fetching Files from Google
// Cloud Storage buckets.
type CloudStorageClient struct {
	name string
}

// NewCloudStorageClient creates a new Client implementation which will proxy filesystem calls to
// Google Cloud Storage bucket.
func NewCloudStorageClient(bucket string) *CloudStorageClient {
	return &CloudStorageClient{
		name: bucket,
	}
}

func (c *CloudStorageClient) Get(ctx context.Context, path string) (*File, error) {
	// Authentication is provided by the gcloud tool when running locally, and
	// by the associated service account when running on Compute Engine.
	client, err := google.DefaultClient(ctx, storage.DevstorageReadOnlyScope)
	if err != nil {
		return nil, fmt.Errorf("unable to get default client: %v", err)
	}
	service, err := storage.New(client)
	if err != nil {
		return nil, fmt.Errorf("unable to create storage service: %v", err)
	}

	res, err := service.Objects.Get(c.name, path).Do()
	if err != nil {
		return nil, err
	}

	resp, err := client.Get(res.MediaLink)
	if err != nil {
		return nil, err
	}

	modTime, _ := time.Parse(http.TimeFormat, res.Updated)
	return &File{
		ReadCloser: resp.Body,
		Name:       res.Name,
		ModTime:    modTime,
		Size:       int64(res.Size),
	}, nil
}
