package storage

import (
	"fmt"

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

func (c *storager) Storage() *Storage {
	return c.once.Do(func() interface{} {
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
