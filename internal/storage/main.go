package storage

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/url"

	"github.com/aws/amazon-ssm-agent/agent/s3util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (s *Storage) GetImageBase64(object *url.URL) (*string, error) {
	bucket, key, err := s.bucketAndKey(object)
	if err != nil {
		return nil, fmt.Errorf("failed to get bucket and key: %w", err)
	}

	output, err := s.client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object meta: %w", err)
	}
	defer output.Body.Close()

	imageBytes, err := io.ReadAll(output.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read image: %w", err)
	}

	imageBase64 := base64.StdEncoding.EncodeToString(imageBytes)

	return &imageBase64, nil
}

func (s *Storage) ValidateImage(object *url.URL, id string) error {
	bucket, key, err := s.bucketAndKey(object)
	if err != nil {
		return fmt.Errorf("failed to get bucket and key: %w", err)
	}

	if bucket != s.bucket {
		return ErrInvalidBucket
	}

	if key != id {
		return ErrInvalidKey
	}

	// output can't be nil
	output, err := s.client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to get image meta: %w", err)
	}

	if *output.ContentType != "image/jpeg" && *output.ContentType != "image/png" {
		return ErrIncorrectImageType
	}

	if *output.ContentLength > maxImageSize {
		return ErrImageTooLarge
	}

	return nil
}

func (s *Storage) bucketAndKey(link *url.URL) (bucket, key string, err error) {
	switch s.backend {
	case digitalOceanBackend:
		spacesURL, err := parseDOSpacesURL(link)
		if err != nil {
			return "", "", fmt.Errorf("failed to parse url [%s]: %w", link, err)
		}
		return spacesURL.Bucket, spacesURL.Key, nil
	case awsBackend:
		s3URL := s3util.ParseAmazonS3URL(nil, link)
		if s3URL.Region != s.region {
			return "", "", ErrRegionMismatched
		}
		return s3URL.Bucket, s3URL.Key, nil
		// should be never happened
	default:
		return "", "", errors.New("invalid backend")
	}
}

func (s *Storage) GeneratePutURL(fileName, contentType string, contentLength int64) (signedURL, key string, err error) {
	key = uuid.New().String()
	if fileName != "" {
		key = fileName
	}
	req, _ := s.client.PutObjectRequest(&s3.PutObjectInput{
		Bucket:        &s.bucket,
		Key:           &key,
		ContentType:   &contentType,
		ContentLength: &contentLength,
	})

	signedURL, err = req.Presign(s.presignedURLExpiration)
	if err != nil {
		return "", "", fmt.Errorf("failed to sign request: %w", err)
	}

	return signedURL, key, nil
}

func parseDOSpacesURL(object *url.URL) (*SpacesURL, error) {
	spacesURL := &SpacesURL{
		URL: object,
	}

	components := doSpacesURLRegexp.FindStringSubmatch(object.String())
	if components == nil {
		return nil, ErrURLRegexp
	}

	// never panic because of regexp validation
	spacesURL.Bucket = components[1]
	spacesURL.Region = components[2]
	spacesURL.Key = components[3]

	return spacesURL, nil
}

func IsBadRequestError(err error) bool {
	if errors.Is(err, ErrImageTooLarge) ||
		errors.Is(err, ErrIncorrectImageType) ||
		errors.Is(err, ErrURLRegexp) ||
		errors.Is(err, ErrInvalidBucket) ||
		errors.Is(err, ErrInvalidKey) {
		return true
	}
	return false
}