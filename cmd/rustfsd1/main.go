package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	log "github.com/sirupsen/logrus"
)

func main() {
	region := "jpa"
	accessKeyID := "MBtRkFcP7v3rdNq2hs1X"
	secretAccessKey := "isPwIxO3vKBkcDQmHR8LfErN5C14qoXVt7yAuGYn"
	endpointURL := "http://127.0.0.1:9000"

	if accessKeyID == "" || secretAccessKey == "" || region == "" || endpointURL == "" {
		log.Fatal("missing the env: RUSTFS_ACCESS_KEY_ID / RUSTFS_SECRET_ACCESS_KEY / RUSTFS_REGION / RUSTFS_ENDPOINT_URL")
	}

	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region), // 指定东京区域
		config.WithCredentialsProvider(
			aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
				accessKeyID,
				secretAccessKey,
				"", // session token 可选
			)),
		),
	)

	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Create an S3 client
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = &endpointURL
		o.UsePathStyle = true
	})

	resp, err := client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		log.Fatalf("list buckets failed: %v", err)
	}

	fmt.Println("Buckets:")
	for _, b := range resp.Buckets {
		log.Println(" -", *b.Name)
	}

	f, err := os.Open("/opt/data/apps/docker-compose.yaml")
	if err != nil {
		log.Errorf("cat not open file %v", err.Error())
		return
	}
	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String("main1"),
		Key:    aws.String("test.txt"),
		Body:   f,
	})
	if err != nil {
		log.Fatalf("upload object failed: %v", err)
	}

	url, err := GetS3PresignedURL(ctx, client, "main1", "test.txt", time.Hour)
	if err != nil {
		log.Errorf("cat not get presigned url for preview %w", err)
	}
	log.Infof("%v", url)
}

func GetS3PresignedURL(ctx context.Context, client *s3.Client, bucket, key string, expires time.Duration) (string, error) {

	presignClient := s3.NewPresignClient(client)

	out, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(expires))
	if err != nil {
		return "", err
	}
	return out.URL, nil
}
