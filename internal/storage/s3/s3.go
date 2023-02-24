package s3

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/polyxia-org/morty-registry/internal/storage"
	"github.com/polyxia-org/morty-registry/pkg/config"
)

const (
	urlExpiresIn = 5 * time.Minute
)

type Storage struct {
	bucketName string

	client *s3.Client
	ctx    context.Context
}

var _ storage.Storage = &Storage{}

func NewStorage(c config.S3Config) (*Storage, error) {

	log.Println("bootstrapping 's3' backend engine")
	ctx := context.Background()

	cfg, err := awsConfig.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	if _, err := cfg.Credentials.Retrieve(ctx); err != nil {
		return nil, err
	}

	// If endpoint isn't empty, it means that the user has configured the backend with
	// an S3 compliant endpoint like MinIO. To work with the AWS SDK, we need to
	// define a custom endpoint resolver.
	if c.Endpoint != "" {
		log.Printf("custom endpoint detected for configuration: %s\n", c.Endpoint)
		cfg.EndpointResolverWithOptions = aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:               c.Endpoint,
				SigningRegion:     c.Region,
				Source:            aws.EndpointSourceCustom,
				HostnameImmutable: true,
			}, nil
		})
	}

	return &Storage{
		client:     s3.NewFromConfig(cfg),
		bucketName: c.Bucket,
		ctx:        ctx,
	}, nil
}

func (s *Storage) GetUploadLink(key string) (string, string, error) {

	presignClient := s3.NewPresignClient(s.client)

	res, err := presignClient.PresignPutObject(s.ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	}, func(po *s3.PresignOptions) {
		po.Expires = urlExpiresIn
	})

	if err != nil {
		log.Println(err)
		return "", "", fmt.Errorf("failed to pre-sign the get object request for key %s in bucketName %s. please check the logs for more details", key, s.bucketName)
	}

	log.Printf("successfully generated pre-signed URL for key %s with a validity period of %f minutes\n", key, urlExpiresIn.Minutes())

	return res.Method, res.URL, nil
}

func (s *Storage) Healthcheck() error {
	ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
	defer cancel()

	input := &s3.HeadBucketInput{
		Bucket: &s.bucketName,
	}

	if _, err := s.client.HeadBucket(ctx, input); err != nil {
		return fmt.Errorf("failed to connect to S3 bucket %s", s.bucketName)
	}

	return nil
}
