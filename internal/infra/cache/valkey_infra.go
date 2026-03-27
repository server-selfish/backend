package cache_infra

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/valkey-io/valkey-go"
	"github.com/valkey-io/valkey-go/valkeycompat"
)

type (
	ValkeyInfra interface {
		Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
		Get(ctx context.Context, key string) (string, error)
		Delete(ctx context.Context, key string) error
	}
	valkeInfra struct {
		valkeyClient valkeycompat.Cmdable
		logger       zerolog.Logger
	}
)

func NewValkeyCache(valkeyClient valkey.Client, logger zerolog.Logger) ValkeyInfra {
	compat := valkeycompat.NewAdapter(valkeyClient)
	if err := compat.Ping(context.Background()).Err(); err != nil {
		logger.Fatal().Err(err).Msg("failed to ping valkey client")
	}
	return &valkeInfra{
		valkeyClient: compat,
		logger:       logger,
	}
}

func (v *valkeInfra) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	err := v.valkeyClient.Set(ctx, key, value, expiration).Err()
	if err != nil {
		v.logger.Error().Err(err).Str("key", key).Msg("failed to set value in redis")
		return err
	}
	return nil
}

func (v *valkeInfra) Get(ctx context.Context, key string) (string, error) {
	val, err := v.valkeyClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			v.logger.Warn().Str("key", key).Msg("key not found in redis")
			return "", err
		}
		v.logger.Error().Err(err).Str("key", key).Msg("failed to get value from redis")
		return "", err
	}
	return val, nil
}

func (v *valkeInfra) Delete(ctx context.Context, key string) error {
	err := v.valkeyClient.Del(ctx, key).Err()
	if err != nil {
		v.logger.Error().Err(err).Str("key", key).Msg("failed to delete key from redis")
		return err
	}
	return nil
}
