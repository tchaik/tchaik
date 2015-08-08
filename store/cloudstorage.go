// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package store

import (
	"fmt"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"

	"google.golang.org/cloud"
	"google.golang.org/cloud/storage"
)

// CloudStorageClient implements Client and handles fetching Files from Google
// Cloud Storage buckets.
type CloudStorageClient struct {
	projID string
	name   string
}

// NewCloudStorageClient creates a new Client implementation which will proxy filesystem calls to
// Google Cloud Storage bucket.
func NewCloudStorageClient(projID, bucket string) *CloudStorageClient {
	return &CloudStorageClient{
		projID: projID,
		name:   bucket,
	}
}

func (c *CloudStorageClient) Get(ctx context.Context, path string) (*File, error) {
	client, err := google.DefaultClient(ctx, storage.ScopeReadOnly)
	if err != nil {
		return nil, fmt.Errorf("unable to get default client: %v", err)
	}

	// TODO(dhowden): This will panic if c.projID is empty, though a wrapper just for
	// this seems over the top...
	ctx = cloud.WithContext(ctx, c.projID, client)

	obj, err := storage.StatObject(ctx, c.name, path)
	if err != nil {
		return nil, fmt.Errorf("unable to stat object: %v", err)
	}

	r, err := storage.NewReader(ctx, c.name, path)
	if err != nil {
		return nil, fmt.Errorf("error fetching '%v' from '%v': %v", path, c.name, err)
	}

	return &File{
		ReadCloser: r,
		Name:       obj.Name,
		ModTime:    obj.Updated,
		Size:       int64(obj.Size),
	}, nil
}
