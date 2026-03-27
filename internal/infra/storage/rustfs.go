package storage_infra

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
)

type (
	RustfsInfra interface {
		Set(ctx context.Context, bucketName, objectPath string, reader io.Reader, objectSize int64) error
		Get(ctx context.Context, bucketName, objectPath string) (io.ReadCloser, error)
		Delete(ctx context.Context, bucketName, fileName string) error
		CreateBucketIfNotExist(ctx context.Context, bucketName string) error
		SetPolicy(ctx context.Context, bucketName, policy string) error
		GetPolicy(ctx context.Context, bucketName string) (string, error)
	}
	rustfsInfra struct {
		client *minio.Client
	}
)

func NewRustfsInfra(client *minio.Client) RustfsInfra {
	return &rustfsInfra{
		client: client,
	}
}

func (r *rustfsInfra) Set(ctx context.Context, bucketName, objectPath string, reader io.Reader, objectSize int64) error {
	_, err := r.client.PutObject(ctx, bucketName, objectPath, reader, objectSize, minio.PutObjectOptions{})
	return err
}

func (r *rustfsInfra) Get(ctx context.Context, bucketName, objectPath string) (io.ReadCloser, error) {
	object, err := r.client.GetObject(ctx, bucketName, objectPath, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return object, nil
}

func (r *rustfsInfra) Delete(ctx context.Context, bucketName, fileName string) error {
	return r.client.RemoveObject(ctx, bucketName, fileName, minio.RemoveObjectOptions{})
}

func (r *rustfsInfra) CreateBucketIfNotExist(ctx context.Context, bucketName string) error {
	exists, err := r.client.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}
	if !exists {
		return r.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	}
	return nil
}
func (r *rustfsInfra) SetPolicy(ctx context.Context, bucketName, policy string) error {
	return r.client.SetBucketPolicy(ctx, bucketName, policy)
}

func (r *rustfsInfra) GetPolicy(ctx context.Context, bucketName string) (string, error) {
	return r.client.GetBucketPolicy(ctx, bucketName)
}
