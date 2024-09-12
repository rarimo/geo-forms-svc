package storage

import (
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"gitlab.com/distributed_lab/dig"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
)

type Storager interface {
	Storage() *Storage
}

func NewStorager(getter kv.Getter) Storager {
	return &storager{
		getter: getter,
	}
}

type storager struct {
	once   comfig.Once
	getter kv.Getter
}

// Storage support DigitalOceanSpaces and AWS S3
// Other providers are not supported.
func (c *storager) Storage() *Storage {
	return c.once.Do(func() interface{} {
		var envCfg struct {
			S3Key    string `dig:"S3_KEY,clear"`
			S3Secret string `dig:"S3_SECRET,clear"`
		}

		err := dig.Out(&envCfg).Now()
		if err != nil {
			panic(fmt.Errorf("failed to dig out spaces key and secret: %w", err))
		}

		var cfg struct {
			Backend                string         `fig:"backend,required"`
			Endpoint               string         `fig:"endpoint,required"`
			Bucket                 string         `fig:"bucket,required"`
			PresignedURLExpiration *time.Duration `fig:"presigned_url_expiration"`
			Region                 *string        `fig:"region"`
			APIKey                 string         `fig:"api_key,required"`
		}

		err = figure.Out(&cfg).
			From(kv.MustGetStringMap(c.getter, "storage")).
			Please()
		if err != nil {
			panic(fmt.Errorf("failed to figure out s3 storage config: %w", err))
		}

		switch cfg.Backend {
		case digitalOceanBackend, awsBackend:
		default:
			panic(errors.New("invalid backend provided"))
		}

		if cfg.PresignedURLExpiration == nil {
			cfg.PresignedURLExpiration = &defaultPresignedURLExpiration
		}

		if cfg.Region == nil {
			cfg.Region = aws.String(defaultRegion)
		}

		s3Config := &aws.Config{
			Credentials:      credentials.NewStaticCredentials(envCfg.S3Key, envCfg.S3Secret, ""),
			Endpoint:         aws.String(cfg.Endpoint),
			Region:           cfg.Region,
			S3ForcePathStyle: aws.Bool(false),
		}

		newSession, err := session.NewSession(s3Config)
		if err != nil {
			panic(fmt.Errorf("failed to create session: %w", err))
		}

		s3Client := s3.New(newSession)

		return &Storage{
			client:                 s3Client,
			bucket:                 cfg.Bucket,
			presignedURLExpiration: *cfg.PresignedURLExpiration,
			backend:                cfg.Backend,
			region:                 *cfg.Region,
			APIKey:                 cfg.APIKey,
		}
	}).(*Storage)
}
