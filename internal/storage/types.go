package storage

import (
	"fmt"
	"net/url"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
)

var DOSpacesURLRegexp = regexp.MustCompile(`^https:\/\/(.+?)\.(.+?)(?:\.cdn)?\.digitaloceanspaces\.com\/(.+)$`)

const maxImageSize = 1 << 22 // 4mb

var (
	ErrImageTooLarge      = fmt.Errorf("too large image, must be not greater than %d bytes", maxImageSize)
	ErrIncorrectImageType = fmt.Errorf("incorrect object type, must be image/png or image/jpeg")
	ErrURLRegexp          = fmt.Errorf("url don't match regexp")
	ErrInvalidBucket      = fmt.Errorf("invalid bucket")
	ErrInvalidKey         = fmt.Errorf("invalid key")

	defaultPresignedURLExpiration = 5 * time.Minute
)

type Storage struct {
	client                 *s3.S3
	bucket                 string
	presignedURLExpiration time.Duration
}

type SpacesURL struct {
	URL    *url.URL
	Bucket string
	Key    string
	Region string
}
