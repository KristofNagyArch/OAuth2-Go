package cache

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"oauth2-go/domain"
	"strings"
	"time"
)

const (
	keySeparator          = "::"
	refresh               = "refresh"
	access                = "access"
	accessCacheRetention  = 60     // 1 hour in minute
	refreshCacheRetention = 144000 // 100 day in minute
)

type TokenRedisRepository struct {
	client      *redis.Client
	appEnv      string
	serviceName string
}

var _ domain.RedisRepository = &redisConnection{}

type redisConnection struct {
	client      *redis.Client
	appEnv      string
	serviceName string
}

func InitTokenRedisRepository(client *redis.Client, appEnv, serviceName string) {
	domain.RedisRepositoryImpl = &redisConnection{client: client, appEnv: appEnv, serviceName: serviceName}
}

func (r *redisConnection) FindAccessTokenByPartnerID(ctx context.Context, partnerID string) string {
	result, err := r.client.Get(ctx, r.buildAccessKey(partnerID)).Result()
	if err == redis.Nil {
		return ""
	}
	if err != nil {
		return ""
	}
	var token string
	err = json.Unmarshal([]byte(result), &token)
	if err != nil {
		return ""
	}

	return token
}

func (r *redisConnection) StoreAccessTokenToPartnerID(ctx context.Context, partnerID, token string) {
	b, err := json.Marshal(token)
	if err != nil {
		return
	}
	err = r.CacheAccessToken(ctx, partnerID, b, accessCacheRetention)
	if err != nil {
		return
	}
}

func (r *redisConnection) InvalidateAccessToken(ctx context.Context, partnerID string) {
	_, err := r.client.Del(ctx, r.buildAccessKey(partnerID)).Result()
	if err != nil {
		return
	}
}

func (r *redisConnection) FindRefreshTokenByPartnerID(ctx context.Context, partnerID string) string {
	result, err := r.client.Get(ctx, r.buildRefreshKey(partnerID)).Result()
	if err == redis.Nil {
		return ""
	}
	if err != nil {
		return ""
	}
	var token string
	err = json.Unmarshal([]byte(result), &token)
	if err != nil {
		return ""
	}

	return token
}

func (r *redisConnection) StoreRefreshTokenToPartnerID(ctx context.Context, partnerID, token string) {
	b, err := json.Marshal(token)
	if err != nil {
		return
	}
	err = r.CacheRefreshToken(ctx, partnerID, b, refreshCacheRetention)
	if err != nil {
		return
	}
}

func (r *redisConnection) InvalidateRefreshToken(ctx context.Context, partnerID string) {
	_, err := r.client.Del(ctx, r.buildRefreshKey(partnerID)).Result()
	if err != nil {
		return
	}
}

func (r *redisConnection) CacheAccessToken(ctx context.Context, partnerID string, b []byte, cacheRetention time.Duration) error {
	return r.cacheWithNoExpiration(ctx, r.buildAccessKey(partnerID), b, cacheRetention)
}

func (r *redisConnection) CacheRefreshToken(ctx context.Context, partnerID string, b []byte, cacheRetention time.Duration) error {
	return r.cacheWithNoExpiration(ctx, r.buildRefreshKey(partnerID), b, cacheRetention)
}

func (r *redisConnection) cacheWithNoExpiration(ctx context.Context, key string, value []byte, cacheRetention time.Duration) error {
	_, err := r.client.Set(ctx, key, string(value), cacheRetention*time.Minute).Result()
	if err != nil {
		return err
	}
	return nil
}

func (r *redisConnection) buildAccessKey(partnerID string) string {
	// quickbooks-auth::local::accessToken::partnerID -> accessToken
	return strings.Join([]string{r.serviceName, r.appEnv, access, partnerID}, keySeparator)
}

func (r *redisConnection) buildRefreshKey(partnerID string) string {
	// quickbooks-auth::local::refreshToken::partnerID -> refreshToken
	return strings.Join([]string{r.serviceName, r.appEnv, refresh, partnerID}, keySeparator)
}
