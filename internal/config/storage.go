package config

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/url"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"gitlab.com/distributed_lab/dig"
	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/kv"
)

var DOSpacesURLRegexp = regexp.MustCompile(`^https:\/\/(.+?)\.(.+?)(?:\.cdn)?\.digitaloceanspaces\.com\/(.+)$`)

const maxImageSize = 1 << 22 // 4mb

var (
	ErrImageTooLarge      = fmt.Errorf("too large image, must be not greater than %d bytes", maxImageSize)
	ErrIncorrectImageType = fmt.Errorf("incorrect object type, must be image/png or image/jpeg")
	ErrURLRegexp          = fmt.Errorf("url don't match regexp")
	ErrBucketNotAllowed   = fmt.Errorf("bucket not allowed")
)

type Storage struct {
	client         *s3.S3
	allowedBuckets []string
}

func (c *config) Storage() *Storage {
	return c.storage.Do(func() interface{} {
		var envCfg struct {
			SpacesKey    string `dig:"SPACES_KEY,clear"`
			SpacesSecret string `dig:"SPACES_SECRET,clear"`
		}

		err := dig.Out(&envCfg).Now()
		if err != nil {
			panic(fmt.Errorf("failed to dig out spaces key and secret: %w", err))
		}

		var cfg struct {
			Endpoint       string   `fig:"endpoint,required"`
			AllowedBuckets []string `fig:"allowed_buckets,required"`
		}

		err = figure.Out(&cfg).
			From(kv.MustGetStringMap(c.getter, "storage")).
			Please()
		if err != nil {
			panic(fmt.Errorf("failed to figure out s3 storage config: %w", err))
		}

		s3Config := &aws.Config{
			Credentials:      credentials.NewStaticCredentials(envCfg.SpacesKey, envCfg.SpacesSecret, ""),
			Endpoint:         aws.String(cfg.Endpoint),
			Region:           aws.String("us-east-1"),
			S3ForcePathStyle: aws.Bool(false),
		}

		newSession, err := session.NewSession(s3Config)
		if err != nil {
			panic(fmt.Errorf("failed to create session: %w", err))
		}

		s3Client := s3.New(newSession)

		return &Storage{
			client:         s3Client,
			allowedBuckets: cfg.AllowedBuckets,
		}
	}).(*Storage)
}

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

type SpacesURL struct {
	URL    *url.URL
	Bucket string
	Key    string
	Region string
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
