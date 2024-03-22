package floodControl

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"task/floodControl/redisDB"
	"time"
)

type FloodController struct {
	repo   *redisDB.FloodControlRepo
	dur    time.Duration
	numReq int64
}

func (c *FloodController) Check(ctx context.Context, userID int64) (bool, error) {
	limit, err := c.repo.GetRequestLimit(ctx, userID)
	if errors.Is(err, redis.Nil) {
		err = c.repo.PutUserRestriction(ctx, userID, c.dur, c.numReq-1)
		return true, nil
	} else if err != nil {
		return false, err
	}
	if limit <= 0 {
		return false, nil
	}
	err = c.repo.DecreaseUserRestriction(ctx, userID)
	if err != nil {
		return false, err
	}
	return true, nil
}
