package pkg

import (
	"context"
	"errors"
	"github.com/linuxsuren/api-testing/pkg/testing/remote"
	"github.com/redis/go-redis/v9"
	"time"
)

func (s *remoteserver) getClient(ctx context.Context) (cli *redis.Client, err error) {
	store := remote.GetStoreFromContext(ctx)
	if store == nil {
		err = errors.New("no connect to redis server")
	} else {
		cli = redis.NewClient(&redis.Options{
			Addr:        store.URL,
			DialTimeout: 5 * time.Second,
			Username:    store.Username,
			Password:    store.Password,
		})
	}
	return
}
