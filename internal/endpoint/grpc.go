package endpoint

import (
	"context"

	v1 "github.com/dzakaammar/event-scheduling-example/gen/go/proto/v1"
	"github.com/dzakaammar/event-scheduling-example/internal"
)

type GRPCEndpoint struct {
	svc internal.EventService
}

func NewGRPCEndpoint(svc internal.EventService) *GRPCEndpoint {
	return &GRPCEndpoint{
		svc: svc,
	}
}

func (g *GRPCEndpoint) CreateEvent(ctx context.Context, req *v1.CreateEventRequest) (*v1.CreateEventResponse, error) {
	return nil, nil
}

func (g *GRPCEndpoint) DeleteEventByID(ctx context.Context, req *v1.DeleteEventByIDRequest) (*v1.DeleteEventByIDResponse, error) {
	return nil, nil
}

func (g *GRPCEndpoint) UpdateEvent(ctx context.Context, req *v1.UpdateEventRequest) (*v1.UpdateEventResponse, error) {
	return nil, nil
}

func (g *GRPCEndpoint) FindEventByID(ctx context.Context, req *v1.FindEventByIDRequest) (*v1.FindEventByIDResponse, error) {
	return nil, nil
}
