package auth

import "context"

type contextKey string

const (
	userIDKey   contextKey = "user_id"
	nicknameKey contextKey = "nickname"
)

func SetUserID(ctx context.Context, id uint) context.Context {
	return context.WithValue(ctx, userIDKey, id)
}

func GetUserID(ctx context.Context) (uint, bool) {
	id, ok := ctx.Value(userIDKey).(uint)
	return id, ok
}

func SetNickname(ctx context.Context, nickname string) context.Context {
	return context.WithValue(ctx, nicknameKey, nickname)
}

func GetNickname(ctx context.Context) (string, bool) {
	nickname, ok := ctx.Value(nicknameKey).(string)
	return nickname, ok
}
