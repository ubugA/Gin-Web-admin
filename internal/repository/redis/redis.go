package redis

import (
	"context"
	"errors"
	"log"
	"runtime"
	"strconv"
	"sync"
	"time"

	"gin-api-admin/configs"
	"gin-api-admin/internal/pkg/core"
	"gin-api-admin/internal/pkg/timeutil"
	"gin-api-admin/internal/pkg/trace"

	redisV8 "github.com/go-redis/redis/v8"
)

var (
	client *redisV8.Client
	once   sync.Once
)

type loggingHook struct {
	ts time.Time
}

func (h *loggingHook) BeforeProcess(ctx context.Context, cmd redisV8.Cmder) (context.Context, error) {
	h.ts = time.Now()
	return ctx, nil
}

func (h *loggingHook) AfterProcess(ctx context.Context, cmd redisV8.Cmder) error {
	monoCtx, ok := ctx.(core.StdContext)
	if !ok {
		log.Println("redis hook AfterProcess ctx illegal")
		return errors.New("ctx illegal")
	}

	if monoCtx.Trace != nil {
		redisInfo := new(trace.Redis)
		redisInfo.Time = timeutil.CSTLayoutString()
		redisInfo.Stack = fileWithLineNum()
		redisInfo.Cmd = cmd.String()
		redisInfo.CostSeconds = time.Since(h.ts).Seconds()

		monoCtx.Trace.AppendRedis(redisInfo)
	}

	return nil
}

func (h *loggingHook) BeforeProcessPipeline(ctx context.Context, cmds []redisV8.Cmder) (context.Context, error) {
	return ctx, nil
}

func (h *loggingHook) AfterProcessPipeline(ctx context.Context, cmds []redisV8.Cmder) error {
	return nil
}

func GetRedisClient() *redisV8.Client {
	once.Do(func() {
		cfg := configs.Get().Redis
		client = redisV8.NewClient(&redisV8.Options{
			Addr:         cfg.Addr,
			Password:     cfg.Pass,
			DB:           cfg.Db,
			MaxRetries:   3,
			DialTimeout:  time.Second * 5,
			ReadTimeout:  time.Second * 20,
			WriteTimeout: time.Second * 20,
			PoolSize:     50,
			MinIdleConns: 2,
			PoolTimeout:  time.Minute,
		})

		if err := client.Ping(context.Background()).Err(); err != nil {
			panic(err)
		}

		client.AddHook(&loggingHook{})

	})

	return client
}

func fileWithLineNum() string {
	_, file, line, ok := runtime.Caller(5)
	if ok {
		return file + ":" + strconv.FormatInt(int64(line), 10)
	}

	return ""
}
