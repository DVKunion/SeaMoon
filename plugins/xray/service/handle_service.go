package service

import (
	"context"

	"github.com/xtls/xray-core/app/proxyman/command"

	"github.com/DVKunion/SeaMoon/plugins/xray/config"
)

type HandleService struct {
	cc command.HandlerServiceClient
}

func (h *HandleService) AddInbound(ctx context.Context, opts ...config.Options) error {
	cf, err := config.Render(opts...)
	if err != nil {
		return err
	}
	for _, ibc := range cf.Inbound {
		if _, err := h.cc.AddInbound(ctx, &command.AddInboundRequest{
			Inbound: ibc,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (h *HandleService) AddOutbound(ctx context.Context, opts ...config.Options) error {
	cf, err := config.Render(opts...)
	if err != nil {
		return err
	}
	for _, obc := range cf.Outbound {
		if _, err := h.cc.AddOutbound(ctx, &command.AddOutboundRequest{
			Outbound: obc,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (h *HandleService) AddInboundUser() {
	// todo
}

func (h *HandleService) AddOutboundUser() {
	// todo
}

func (h *HandleService) RemoveInbound(ctx context.Context, opts ...config.Options) error {
	cf, err := config.Render(opts...)
	if err != nil {
		return err
	}
	for _, ibc := range cf.Inbound {
		if _, err := h.cc.RemoveInbound(ctx, &command.RemoveInboundRequest{
			Tag: ibc.Tag,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (h *HandleService) RemoveOutbound(ctx context.Context, opts ...config.Options) error {
	cf, err := config.Render(opts...)
	if err != nil {
		return err
	}
	for _, obc := range cf.Inbound {
		if _, err := h.cc.RemoveOutbound(ctx, &command.RemoveOutboundRequest{
			Tag: obc.Tag,
		}); err != nil {
			return err
		}
	}
	return nil
}
