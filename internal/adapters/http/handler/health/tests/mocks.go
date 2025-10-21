package tests

import (
	"context"
	"database/sql"

	"github.com/redis/go-redis/v9"
)

type MockDB struct {
	PingFunc func(ctx context.Context) error
}

func (m *MockDB) PingContext(ctx context.Context) error {
	if m.PingFunc != nil {
		return m.PingFunc(ctx)
	}
	return nil
}

type MockRedisClient struct {
	PingResult *redis.StatusCmd
	PingErr    error
}

func (m *MockRedisClient) Ping(ctx context.Context) *redis.StatusCmd {
	if m.PingResult != nil {
		return m.PingResult
	}
	
	cmd := redis.NewStatusCmd(ctx)
	if m.PingErr != nil {
		cmd.SetErr(m.PingErr)
	} else {
		cmd.SetVal("PONG")
	}
	return cmd
}

type DBPinger interface {
	PingContext(ctx context.Context) error
}

type RedisPinger interface {
	Ping(ctx context.Context) *redis.StatusCmd
}

func NewMockDB(pingFunc func(ctx context.Context) error) *MockDB {
	return &MockDB{PingFunc: pingFunc}
}

func NewMockRedisClient(err error) *MockRedisClient {
	return &MockRedisClient{PingErr: err}
}

func DBFromMock(mock *MockDB) *sql.DB {
	return nil
}

type TestHealthHandler struct {
	db          DBPinger
	redisClient RedisPinger
	logger      interface{}
	version     string
}

