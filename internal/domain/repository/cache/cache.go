package cache_repository

import (
	"context"
	"encoding/json"
	"time"

	cache_infra "github.com/server-selfish/backend/internal/infra/cache"
	"github.com/valkey-io/valkey-go"
)

type (
	CacheRepository interface {
		Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
		SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error
		Get(ctx context.Context, key string) (string, error)
		GetJSON(ctx context.Context, key string, dest interface{}) error
		Delete(ctx context.Context, key string) error
	}
	cacheRepository struct {
		ci cache_infra.ValkeyInfra
	}
)

func NewCacheRepository(ci cache_infra.ValkeyInfra) CacheRepository {
	return cacheRepository{
		ci: ci,
	}
}

// Delete implements [CacheRepository].
func (c cacheRepository) Delete(ctx context.Context, key string) error {
	err := c.ci.Del(ctx, key).Err()
	if err != nil {
		// c.logger.Error().Err(err).Str("key", key).Msg("failed to delete key from valkey")
		return err
	}
	return nil
}

// Get implements [CacheRepository].
func (c cacheRepository) Get(ctx context.Context, key string) (string, error) {
	val, err := c.ci.Get(ctx, key).Result()
	if err != nil {
		if err == valkey.Nil {
			// v.logger.Warn().Str("key", key).Msg("key not found in valkey")
			return "", err
		}
		// v.logger.Error().Err(err).Str("key", key).Msg("failed to get value from redis")
		return "", err
	}
	return val, nil
}

// GetJSON implements [CacheRepository].
func (c cacheRepository) GetJSON(ctx context.Context, key string, dest interface{}) error {
	result, err := c.ci.JSONGet(ctx, key, "$").Result()
	if err != nil {
		if err == valkey.Nil {
			// v.logger.Warn().Str("key", key).Msg("key not found in valkey")
			return err
		}
		// v.logger.Error().Err(err).Str("key", key).Msg("failed to get JSON from valkey")
		return err
	}

	var raw []json.RawMessage
	if err := json.Unmarshal([]byte(result), &raw); err != nil {
		return err
	}
	if len(raw) == 0 {
		return valkey.Nil
	}

	return json.Unmarshal(raw[0], dest)
}

// Set implements [CacheRepository].
func (c cacheRepository) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	err := c.ci.Set(ctx, key, value, expiration).Err()
	if err != nil {
		// v.logger.Error().Err(err).Str("key", key).Msg("failed to set value in valkey")
		return err
	}
	return nil
}

// SetJSON implements [CacheRepository].
func (c cacheRepository) SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		// v.logger.Error().Err(err).Str("key", key).Msg("failed to marshal JSON")
		return err
	}

	err = c.ci.JSONSet(ctx, key, "$", string(data)).Err()
	if err != nil {
		// v.logger.Error().Err(err).Str("key", key).Msg("failed to set JSON in valkey")
		return err
	}

	if expiration > 0 {
		err = c.ci.Expire(ctx, key, expiration).Err()
		if err != nil {
			// v.logger.Error().Err(err).Str("key", key).Msg("failed to set expiration")
			return err
		}
	}
	return nil
}
