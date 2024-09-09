package service

import (
	handle "github.com/xtls/xray-core/app/proxyman/command"
	states "github.com/xtls/xray-core/app/stats/command"
	"google.golang.org/grpc"
)

func NewHandleService(cc *grpc.ClientConn) HandleService {
	return HandleService{
		cc: handle.NewHandlerServiceClient(cc),
	}
}

func NewStatesService(cc *grpc.ClientConn) StateService {
	return StateService{
		cc: states.NewStatsServiceClient(cc),
	}
}
