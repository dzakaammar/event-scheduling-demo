package endpoint_test

import (
	"context"
	"reflect"
	"testing"

	v1 "github.com/dzakaammar/event-scheduling-example/gen/go/proto/v1"
	"github.com/dzakaammar/event-scheduling-example/internal"
	"github.com/dzakaammar/event-scheduling-example/internal/endpoint"
)

func TestNewGRPCEndpoint(t *testing.T) {
	type args struct {
		svc internal.EventService
	}
	tests := []struct {
		name string
		args args
		want *endpoint.GRPCEndpoint
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := endpoint.NewGRPCEndpoint(tt.args.svc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGRPCEndpoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGRPCEndpoint_CreateEvent(t *testing.T) {
	type fields struct {
		svc internal.EventService
	}
	type args struct {
		ctx context.Context
		req *v1.CreateEventRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *v1.CreateEventResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := endpoint.NewGRPCEndpoint(tt.fields.svc)
			got, err := g.CreateEvent(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GRPCEndpoint.CreateEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCEndpoint.CreateEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGRPCEndpoint_DeleteEventByID(t *testing.T) {
	type fields struct {
		svc internal.EventService
	}
	type args struct {
		ctx context.Context
		req *v1.DeleteEventByIDRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *v1.DeleteEventByIDResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := endpoint.NewGRPCEndpoint(tt.fields.svc)
			got, err := g.DeleteEventByID(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GRPCEndpoint.DeleteEventByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCEndpoint.DeleteEventByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGRPCEndpoint_UpdateEvent(t *testing.T) {
	type fields struct {
		svc internal.EventService
	}
	type args struct {
		ctx context.Context
		req *v1.UpdateEventRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *v1.UpdateEventResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := endpoint.NewGRPCEndpoint(tt.fields.svc)
			got, err := g.UpdateEvent(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GRPCEndpoint.UpdateEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCEndpoint.UpdateEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGRPCEndpoint_FindEventByID(t *testing.T) {
	type fields struct {
		svc internal.EventService
	}
	type args struct {
		ctx context.Context
		req *v1.FindEventByIDRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *v1.FindEventByIDResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := endpoint.NewGRPCEndpoint(tt.fields.svc)
			got, err := g.FindEventByID(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GRPCEndpoint.FindEventByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCEndpoint.FindEventByID() = %v, want %v", got, tt.want)
			}
		})
	}
}
