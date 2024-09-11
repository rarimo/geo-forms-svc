package storage

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/aws/amazon-ssm-agent/agent/s3util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (s *Storage) ValidateImage(object *url.URL, id string) error {
	bucket, key, err := s.BucketAndKey(object)
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

func (s *Storage) BucketAndKey(link *url.URL) (bucket, key string, err error) {
	switch s.backend {
	case digitalOceanBackend:
		spacesURL, err := ParseDOSpacesURL(link)
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
		return "", "", fmt.Errorf("unknown backend: %s", s.backend)
	}
}

func (s *Storage) UploadB64Image(imageB64 *string) (key string, err error) {
	imageBytes, err := base64.StdEncoding.DecodeString(*imageB64)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	contentType := http.DetectContentType(imageBytes)
	if contentType != "image/jpeg" && contentType != "image/png" {
		return "", fmt.Errorf("incorrect file type")
	}

	key = uuid.New().String()
	_, err = s.client.PutObject(&s3.PutObjectInput{
		Bucket:             &s.bucket,
		Key:                &key,
		Body:               bytes.NewReader(imageBytes),
		ContentDisposition: aws.String("attachment"),
		ContentType:        aws.String(http.DetectContentType(imageBytes)),
		ContentLength:      aws.Int64(int64(len(imageBytes))),
	})
	if err != nil {
		return "", fmt.Errorf("failed to put image in s3: %w", err)
	}

	return key, nil
}

func (s *Storage) GetURL(key string) string {
	req, _ := s.client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: &s.bucket,
		Key:    &key,
	})

	// will never error
	signedReq, _ := req.Presign(time.Minute)
	components := regexp.MustCompile(`(.+?)\?`).FindStringSubmatch(signedReq)
	if components == nil {
		return ""
	}
	return components[1]
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

func (s *Storage) GenerateGetURL(link *url.URL) (signedURL string, err error) {
	bucket, key, err := s.BucketAndKey(link)
	if err != nil {
		return "", fmt.Errorf("failed to get bucket and key: %w", err)
	}

	req, _ := s.client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})

	signedURL, err = req.Presign(time.Hour * 164)
	if err != nil {
		return "", fmt.Errorf("failed to sign request: %w", err)
	}

	return signedURL, nil
}

func (s *Storage) RawSignedGetURL(key string) (signedURL string, err error) {
	req, _ := s.client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: &s.bucket,
		Key:    &key,
	})

	signedURL, err = req.Presign(time.Hour * 164)
	if err != nil {
		return "", fmt.Errorf("failed to sign request: %w", err)
	}

	return signedURL, nil
}

func ParseDOSpacesURL(object *url.URL) (*SpacesURL, error) {
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
