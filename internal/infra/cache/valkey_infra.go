package cache_infra

import (
	"context"
	"encoding/json"
	"time"

	"github.com/rs/zerolog"
	"github.com/valkey-io/valkey-go"
	"github.com/valkey-io/valkey-go/valkeycompat"
)

type (
	ValkeyInfra interface {
		Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
		SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error
		Get(ctx context.Context, key string) (string, error)
		GetJSON(ctx context.Context, key string, dest interface{}) error
		Delete(ctx context.Context, key string) error
	}
	valkeyInfra struct {
		valkeyClient valkeycompat.Cmdable
		logger       zerolog.Logger
	}
)

func NewValkeyCache(valkeyClient valkey.Client, logger zerolog.Logger) ValkeyInfra {
	compat := valkeycompat.NewAdapter(valkeyClient)
	if err := compat.Ping(context.Background()).Err(); err != nil {
		logger.Fatal().Err(err).Msg("failed to ping valkey client")
	}
	return &valkeyInfra{
		valkeyClient: compat,
		logger:       logger,
	}
}

func (v *valkeyInfra) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	err := v.valkeyClient.Set(ctx, key, value, expiration).Err()
	if err != nil {
		v.logger.Error().Err(err).Str("key", key).Msg("failed to set value in valkey")
		return err
	}
	return nil
}

func (v *valkeyInfra) SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		v.logger.Error().Err(err).Str("key", key).Msg("failed to marshal JSON")
		return err
	}

	err = v.valkeyClient.JSONSet(ctx, key, "$", string(data)).Err()
	if err != nil {
		v.logger.Error().Err(err).Str("key", key).Msg("failed to set JSON in valkey")
		return err
	}

	if expiration > 0 {
		err = v.valkeyClient.Expire(ctx, key, expiration).Err()
		if err != nil {
			v.logger.Error().Err(err).Str("key", key).Msg("failed to set expiration")
			return err
		}
	}

	return nil
}

func (v *valkeyInfra) Get(ctx context.Context, key string) (string, error) {
	val, err := v.valkeyClient.Get(ctx, key).Result()
	if err != nil {
		if err == valkey.Nil {
			v.logger.Warn().Str("key", key).Msg("key not found in valkey")
			return "", err
		}
		v.logger.Error().Err(err).Str("key", key).Msg("failed to get value from redis")
		return "", err
	}
	return val, nil
}

func (v *valkeyInfra) GetJSON(ctx context.Context, key string, dest interface{}) error {
	result, err := v.valkeyClient.JSONGet(ctx, key, "$").Result()
	if err != nil {
		if err == valkey.Nil {
			v.logger.Warn().Str("key", key).Msg("key not found in valkey")
			return err
		}
		v.logger.Error().Err(err).Str("key", key).Msg("failed to get JSON from valkey")
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

func (v *valkeyInfra) Delete(ctx context.Context, key string) error {
	err := v.valkeyClient.Del(ctx, key).Err()
	if err != nil {
		v.logger.Error().Err(err).Str("key", key).Msg("failed to delete key from valkey")
		return err
	}
	return nil
}
