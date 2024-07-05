package domain

import (
	"context"
)

var (
	// RedisRepositoryImpl implement RedisRepository
	RedisRepositoryImpl RedisRepository
)

// RedisRepository is an interface for RedisRepositoryImpl
type RedisRepository interface {
	FindAccessTokenByPartnerID(ctx context.Context, partnerID string) string
	StoreAccessTokenToPartnerID(ctx context.Context, partnerID, token string)
	InvalidateAccessToken(ctx context.Context, partnerID string)

	FindRefreshTokenByPartnerID(ctx context.Context, partnerID string) string
	StoreRefreshTokenToPartnerID(ctx context.Context, partnerID, token string)
	InvalidateRefreshToken(ctx context.Context, partnerID string)
}
