package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xtls/xray-core/app/proxyman/command"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestAddInbound(t *testing.T) {
	cc, err := grpc.NewClient(fmt.Sprintf("%s:%d", "127.0.0.1", 10086),
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		))
	assert.NoError(t, err)
	sc := HandleService{
		cc: command.NewHandlerServiceClient(cc),
	}
	err = sc.AddInbound(context.Background())
	assert.NoError(t, err)
}

func TestAddOutbound(t *testing.T) {
	cc, err := grpc.NewClient(fmt.Sprintf("%s:%d", "127.0.0.1", 10086),
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		))
	assert.NoError(t, err)
	sc := HandleService{
		cc: command.NewHandlerServiceClient(cc),
	}
	err = sc.AddInbound(context.Background())
	assert.NoError(t, err)
}
