package aws

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Client struct {
	*s3.Client
	Endpoint string
	Bucket   string
}

func NewClient(key, secret, region, endpoint, bucket string) (*Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
			key,
			secret,
			"",
		))),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, err
	}
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})

	return &Client{client, endpoint, bucket}, nil
}

func (c *Client) Buckets(ctx context.Context) (*s3.ListBucketsOutput, error) {
	return c.Client.ListBuckets(ctx, &s3.ListBucketsInput{})
}

// put file
func (c *Client) Put(ctx context.Context, key string, data []byte) (*s3.PutObjectOutput, error) {
	return c.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(c.Bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	})
}

// get file
func (c *Client) Get(ctx context.Context, key string) (*s3.GetObjectOutput, error) {
	return c.Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.Bucket),
		Key:    aws.String(key),
	})
}

// delete file
func (c *Client) Delete(ctx context.Context, key string) (*s3.DeleteObjectOutput, error) {
	return c.Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(c.Bucket),
		Key:    aws.String(key),
	})
}

// Get Url
func (c *Client) Url(ctx context.Context, key string) string {
	return fmt.Sprintf("%s/%s/%s", c.Endpoint, c.Bucket, key)
}

// list files
func (c *Client) List(ctx context.Context) (*s3.ListObjectsV2Output, error) {
	return c.Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(c.Bucket),
	})
}
func ToString(v *string) string {
	return aws.ToString(v)
}
