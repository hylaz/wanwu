package oauth2_util

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	_redis         *redis.Client
	expirationCode = time.Minute * 5
)

type CodePayload struct {
	ClientID string `json:"client_id"`
	UserID   string `json:"user_id"`
}

func SaveCode(ctx context.Context, code string, payload CodePayload) error {
	b, _ := json.Marshal(payload)
	if err := _redis.Set(ctx, getRedisCodeKey(code), string(b), expirationCode).Err(); err != nil {
		return fmt.Errorf("save code %v client_id %v err: %v", code, payload.ClientID, err)
	}
	return nil
}

func ValidateCode(ctx context.Context, code, clientID string) (CodePayload, error) {
	ret := _redis.Get(ctx, getRedisCodeKey(code))
	if err := ret.Err(); err != nil {
		return CodePayload{}, fmt.Errorf("validate code %v err: %v", code, err)
	}
	var payload CodePayload
	if err := json.Unmarshal([]byte(ret.Val()), &payload); err != nil {
		return CodePayload{}, fmt.Errorf("validate code %v unmarshal err: %v", code, err)
	}
	if payload.ClientID != clientID {
		return CodePayload{}, fmt.Errorf("validate code %v client_id %v err: invalid client_id %v", code, payload.ClientID, clientID)
	}
	if err := _redis.Del(ctx, getRedisCodeKey(code)).Err(); err != nil {
		return CodePayload{}, fmt.Errorf("validate code %v client_id %v delete err %v", code, payload.ClientID, err)
	}
	return payload, nil
}

func SaveRefreshToken(ctx context.Context, refreshToken string, expiration time.Duration) error {
	if err := _redis.Set(ctx, getRedisRefreshTokenKey(refreshToken), 1, expiration).Err(); err != nil {
		return fmt.Errorf("save refresh token %v err: %v", refreshToken, err)
	}
	return nil
}

func ValidateRefreshToken(ctx context.Context, refreshToken string) error {
	if err := _redis.Get(ctx, getRedisRefreshTokenKey(refreshToken)).Err(); err != nil {
		return fmt.Errorf("validate refresh token %v err: %v", refreshToken, err)
	}
	if err := _redis.Del(ctx, getRedisRefreshTokenKey(refreshToken)).Err(); err != nil {
		return fmt.Errorf("validate refresh token %v delete err: %v", refreshToken, err)
	}
	return nil
}

// --- internal ---

func getRedisCodeKey(code string) string {
	return fmt.Sprintf("oauth2-code:%v", code)
}

func getRedisRefreshTokenKey(refreshToken string) string {
	return fmt.Sprintf("oauth2-refresh-token:%v", refreshToken)
}
