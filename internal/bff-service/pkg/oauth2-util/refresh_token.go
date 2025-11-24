package oauth2_util

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const RefreshTokenExpiration = time.Hour * 24 * 7 // 7day

func GenerateRefreshToken(ctx context.Context, userID, clientID string, expiration time.Duration) (string, error) {
	if err := checkInit(); err != nil {
		return "", err
	}
	refreshToken := uuid.NewString()
	//save refresh token to redis
	err := saveRefreshToken(ctx, refreshToken, expiration, RefreshTokenPayload{
		UserID:   userID,
		ClientID: clientID,
	})
	if err != nil {
		return "", err
	}
	return refreshToken, nil
}

func ValidateRefreshToken(ctx context.Context, refreshToken, clientID string) (RefreshTokenPayload, error) {
	if err := checkInit(); err != nil {
		return RefreshTokenPayload{}, err
	}
	ret := _redis.Get(ctx, getRedisRefreshTokenKey(refreshToken))
	if err := ret.Err(); err != nil {
		return RefreshTokenPayload{}, fmt.Errorf("validate refresh token %v err: %v", refreshToken, err)
	}
	var payload RefreshTokenPayload
	if err := json.Unmarshal([]byte(ret.Val()), &payload); err != nil {
		return RefreshTokenPayload{}, fmt.Errorf("validate refresh token %v unmarshal err: %v", refreshToken, err)
	}
	if payload.ClientID != clientID {
		return RefreshTokenPayload{}, fmt.Errorf("validate refresh token %v client_id %v err: invalid client_id %v", refreshToken, payload.ClientID, clientID)
	}
	if err := _redis.Del(ctx, getRedisRefreshTokenKey(refreshToken)).Err(); err != nil {
		return RefreshTokenPayload{}, fmt.Errorf("validate refresh token %v delete err: %v", refreshToken, err)
	}
	return payload, nil
}

// --- internal ---

func getRedisRefreshTokenKey(refreshToken string) string {
	return fmt.Sprintf("oauth2-refresh-token:%v", refreshToken)
}

func saveRefreshToken(ctx context.Context, refreshToken string, expiration time.Duration, payload RefreshTokenPayload) error {
	b, _ := json.Marshal(payload)
	if err := _redis.Set(ctx, getRedisRefreshTokenKey(refreshToken), b, expiration).Err(); err != nil {
		return fmt.Errorf("save refresh token %v err: %v", refreshToken, err)
	}
	return nil
}
