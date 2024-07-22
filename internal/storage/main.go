package storage

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/url"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
)

func (s *Storage) GetImageBase64(object *url.URL) (*string, error) {
	spacesURL, err := ParseDOSpacesURL(object)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url [%s]: %w", object.String(), err)
	}

	output, err := s.client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(spacesURL.Bucket),
		Key:    aws.String(spacesURL.Key),
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

func (s *Storage) ValidateImage(object *url.URL) error {
	spacesURL, err := ParseDOSpacesURL(object)
	if err != nil {
		return fmt.Errorf("failed to parse url [%s]: %w", object.String(), err)
	}

	if func() error {
		for _, bucket := range s.allowedBuckets {
			if spacesURL.Bucket == bucket {
				return nil
			}
		}
		return ErrBucketNotAllowed
	}() != nil {
		return fmt.Errorf("bucket=%s: %w", spacesURL.Bucket, err)
	}

	// output can't be nil
	output, err := s.client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(spacesURL.Bucket),
		Key:    aws.String(spacesURL.Key),
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

func ParseDOSpacesURL(object *url.URL) (*SpacesURL, error) {
	spacesURL := &SpacesURL{
		URL: object,
	}

	components := DOSpacesURLRegexp.FindStringSubmatch(object.String())
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
	if errors.Is(err, ErrImageTooLarge) &&
		errors.Is(err, ErrIncorrectImageType) &&
		errors.Is(err, ErrURLRegexp) &&
		errors.Is(err, ErrBucketNotAllowed) {
		return true
	}
	return false
}
