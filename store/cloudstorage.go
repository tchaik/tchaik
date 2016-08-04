// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package store

import (
	"fmt"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"

	"google.golang.org/api/option"

	"cloud.google.com/go/storage"
)

// CloudStorageClient implements Client and handles fetching Files from Google
// Cloud Storage buckets.
type CloudStorageClient struct {
	bucket string
}

// NewCloudStorageClient creates a new Client implementation which will proxy filesystem calls to
// Google Cloud Storage bucket.
func NewCloudStorageClient(bucket string) *CloudStorageClient {
	return &CloudStorageClient{
		bucket: bucket,
	}
}

func (c *CloudStorageClient) Get(ctx context.Context, path string) (*File, error) {
	ts, err := google.DefaultTokenSource(ctx, storage.ScopeReadOnly)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve default token source: %v", err)
	}

	client, err := storage.NewClient(ctx, option.WithTokenSource(ts))
	if err != nil {
		return nil, fmt.Errorf("unable to get default client: %v", err)
	}

	bh := client.Bucket(c.bucket)
	obj := bh.Object(path)

	attrs, err := obj.Attrs(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch object attributes: %v", err)
	}

	r, err := obj.NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetching '%v' from '%v': %v", path, c.bucket, err)
	}

	return &File{
		ReadCloser: r,
		Name:       attrs.Name,
		ModTime:    attrs.Updated,
		Size:       attrs.Size,
	}, nil
}
