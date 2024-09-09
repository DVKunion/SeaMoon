package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xtls/xray-core/app/stats/command"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestSysState(t *testing.T) {
	cc, err := grpc.NewClient(fmt.Sprintf("%s:%d", "127.0.0.1", 10086),
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		))
	assert.NoError(t, err)
	sc := StateService{
		cc: command.NewStatsServiceClient(cc),
	}
	a, err := sc.GetSysStates(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, true, a.NumGoroutine > 0)
	assert.Equal(t, true, a.NumGC > 0)
	assert.Equal(t, true, a.Alloc > 0)
	assert.Equal(t, true, a.TotalAlloc > 0)
}
