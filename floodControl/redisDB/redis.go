package redisDB

import (
	"context"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

type FloodControlRepo struct {
	db redis.Client
}

func (r *FloodControlRepo) GetRequestLimit(ctx context.Context, userID int64) (int64, error) {
	result, err := r.db.Get(ctx, strconv.FormatInt(userID, 10)).Int64()
	return result, err
}
func (r *FloodControlRepo) PutUserRestriction(ctx context.Context, userID int64, dur time.Duration, numReq int64) error {
	return r.db.SetEx(ctx, strconv.FormatInt(userID, 10), numReq, dur).Err()
}
func (r *FloodControlRepo) DecreaseUserRestriction(ctx context.Context, userID int64) error {
	return r.db.Decr(ctx, strconv.FormatInt(userID, 10)).Err()
}
