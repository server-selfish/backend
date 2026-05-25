package storage_infra

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/rs/zerolog"
)

type (
	RustfsInfra interface {
		Set(ctx context.Context, bucketName, objectPath string, reader io.Reader, objectSize int64) error
		Get(ctx context.Context, bucketName, objectPath string) (io.ReadCloser, error)
		Delete(ctx context.Context, bucketName, objectPath string) error
		CreateBucketIfNotExist(ctx context.Context, bucketName string) error
		SetPolicy(ctx context.Context, bucketName, policy string) error
		GetPolicy(ctx context.Context, bucketName string) (string, error)
	}
	rustfsInfra struct {
		client *minio.Client
		logger *zerolog.Logger
	}
)

func NewRustfsInfra(client *minio.Client, logger zerolog.Logger) RustfsInfra {
	return &rustfsInfra{
		client: client,
		logger: &logger,
	}
}

func (r *rustfsInfra) Set(ctx context.Context, bucketName, objectPath string, reader io.Reader, objectSize int64) error {
	_, err := r.client.PutObject(ctx, bucketName, objectPath, reader, objectSize, minio.PutObjectOptions{})
	if err != nil {
		r.logger.Error().Err(err).Str("bucket-name", bucketName).Str("object-path", objectPath).Msg("put object failed")
	}
	return err
}

func (r *rustfsInfra) Get(ctx context.Context, bucketName, objectPath string) (io.ReadCloser, error) {
	object, err := r.client.GetObject(ctx, bucketName, objectPath, minio.GetObjectOptions{})
	if err != nil {
		r.logger.Error().Err(err).Str("bucket-name", bucketName).Str("object-path", objectPath).Msg("get object failed")
		return nil, err
	}
	return object, nil
}

func (r *rustfsInfra) Delete(ctx context.Context, bucketName, objectPath string) error {
	if err := r.client.RemoveObject(ctx, bucketName, objectPath, minio.RemoveObjectOptions{}); err != nil {
		r.logger.Error().Err(err).Str("bucket-name", bucketName).Str("object-path", objectPath).Msg("delete object failed")
		return err
	}
	return nil
}

func (r *rustfsInfra) CreateBucketIfNotExist(ctx context.Context, bucketName string) error {
	exists, err := r.client.BucketExists(ctx, bucketName)
	if err != nil {
		r.logger.Error().Err(err).Str("bucket-name", bucketName).Msg("check bucket name existance failed")
		return err
	}
	if !exists {
		if err := r.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
			r.logger.Error().Err(err).Str("bucket-name", bucketName).Msg("make bucket failed")
			return err
		}
		return nil
	}
	return nil
}
func (r *rustfsInfra) SetPolicy(ctx context.Context, bucketName, policy string) error {
	if err := r.client.SetBucketPolicy(ctx, bucketName, policy); err != nil {
		r.logger.Error().Err(err).Str("bucket-name", bucketName).Str("policy", policy).Msg("set bucket policy failed")
		return err
	}
	return nil
}

func (r *rustfsInfra) GetPolicy(ctx context.Context, bucketName string) (string, error) {
	pol, err := r.client.GetBucketPolicy(ctx, bucketName)
	if err != nil {
		r.logger.Error().Err(err).Str("bucket-name", bucketName).Msg("get bucket policy failed")
		return "", err
	}
	return pol, nil
}
