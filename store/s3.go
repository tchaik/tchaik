// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package store

import (
	"net/http"
	"time"

	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
)

// S3Client implements Client and handles fetching Files from S3 buckets.
type S3Client struct {
	name   string
	auth   aws.Auth
	region aws.Region
}

func NewS3Client(bucket string, auth aws.Auth, region aws.Region) *S3Client {
	return &S3Client{
		name:   bucket,
		auth:   auth,
		region: region,
	}
}

// Implements Client.
func (c *S3Client) Get(path string) (*File, error) {
	s3 := s3.New(c.auth, c.region)
	b := s3.Bucket(c.name)

	k, err := b.GetKey(path)
	if err != nil {
		return nil, err
	}

	rc, err := b.GetReader(path)
	if err != nil {
		return nil, err
	}

	modTime, _ := time.Parse(http.TimeFormat, k.LastModified)
	return &File{
		ReadCloser: rc,
		Name:       k.Key,
		ModTime:    modTime,
		Size:       k.Size,
	}, nil
}
