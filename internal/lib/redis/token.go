package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (r *RedisProvider) CheckFromBlackList(
	ctx context.Context,
	userID int,
	expiredAt string,
	token string,
) error {

	ttlBefore, err := time.Parse(time.RFC3339Nano, expiredAt)
	if err != nil {
		r.Log.Error("Failed to parse ttl", "error", err)
		return status.Errorf(codes.Internal, "server error")
	}

	pipe := r.Client.Pipeline()
	resPipe := make(map[string]*redis.SliceCmd)
	fields := map[string][]string{
		"user_mass": {
			"ttl",
		},
		"user_once": {
			"ttl",
			"token",
		},
	}

	for key, val := range fields {
		resPipe[key] = pipe.HMGet(ctx, fmt.Sprintf("%s:%d", key, userID), val...)
	}

	if _, err := pipe.Exec(ctx); err != nil {
		r.Log.Error("Failed to get user from blacklist", "error", err)
		return status.Errorf(codes.Internal, "server error")
	}

	for key, val := range resPipe {

		if list, err := val.Result(); err == nil && list != nil {
			if list[0] != nil {
				ttl, err := time.Parse(time.RFC3339Nano, list[0].(string))
				if err != nil {
					r.Log.Error("Failed to parse ttl", "error", err)
					return status.Errorf(codes.Internal, "server error")
				}
				if !ttl.Before(ttlBefore) {
					return status.Errorf(codes.Unauthenticated, "unauthorized token")
				}
			}
			if key == "user_once" && list[1] != nil {
				if list[1] == token {
					return status.Errorf(codes.Unauthenticated, "unauthorized token")
				}
			}
		} else if err != nil {
			r.Log.Error("Failed to get user from blacklist", "error", err)
			return status.Errorf(codes.Internal, "server error")
		}
	}

	return nil
}

func (r *RedisProvider) AddToBlackList(
	ctx context.Context,
	userID int,
	token string,
	multitype string,
) error {
	userKey := fmt.Sprintf("user_%s:%d", multitype, userID)
	ttl := time.Now()
	switch multitype {
	case "once":
		if err := r.Client.HSet(ctx, userKey, "token", token, "ttl", ttl).Err(); err != nil {
			r.Log.Error("Failed to add user to blacklist", "error", err)
			return err
		}
	case "mass":
		if err := r.Client.HSet(ctx, userKey, "ttl", ttl).Err(); err != nil {
			r.Log.Error("Failed to add user to blacklist", "error", err)
			return err
		}
	default:
		r.Log.Error("Failed to add user to blacklist", "error", fmt.Sprintf("invalid multitype: %s", multitype))
		return status.Errorf(codes.Internal, "server error")
	}

	if err := r.Client.Expire(ctx, userKey, time.Minute*30).Err(); err != nil {
		r.Log.Error("Failed to set expiration time for user", "error", err)
		return err
	}
	return nil
}
