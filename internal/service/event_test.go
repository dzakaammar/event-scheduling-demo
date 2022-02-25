package service_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/dzakaammar/event-scheduling-example/internal"
	"github.com/dzakaammar/event-scheduling-example/internal/mock"
	"github.com/dzakaammar/event-scheduling-example/internal/service"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewEventService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		eventRepo internal.EventRepository
	}
	tests := []struct {
		name string
		args args
		want *service.EventService
	}{
		{
			name: "OK",
			args: args{
				eventRepo: mock.NewMockEventRepository(ctrl),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := service.NewEventService(tt.args.eventRepo)
			assert.NotNil(t, got)
		})
	}
}

func TestEventService_CreateEvent(t *testing.T) {
	type fields struct {
		eventRepo internal.EventRepository
	}
	type args struct {
		ctx context.Context
		req *internal.CreateEventRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := service.NewEventService(tt.fields.eventRepo)
			if err := e.CreateEvent(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("EventService.CreateEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEventService_DeleteEventByID(t *testing.T) {
	type fields struct {
		eventRepo internal.EventRepository
	}
	type args struct {
		ctx context.Context
		req *internal.DeleteEventByIDRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := service.NewEventService(tt.fields.eventRepo)
			if err := e.DeleteEventByID(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("EventService.DeleteEventByID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEventService_UpdateEvent(t *testing.T) {
	type fields struct {
		eventRepo internal.EventRepository
	}
	type args struct {
		ctx context.Context
		req *internal.UpdateEventRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := service.NewEventService(tt.fields.eventRepo)
			if err := e.UpdateEvent(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("EventService.UpdateEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEventService_FindEventByID(t *testing.T) {
	type fields struct {
		eventRepo internal.EventRepository
	}
	type args struct {
		ctx context.Context
		req *internal.FindEventByIDRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *internal.Event
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := service.NewEventService(tt.fields.eventRepo)
			got, err := e.FindEventByID(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("EventService.FindEventByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EventService.FindEventByID() = %v, want %v", got, tt.want)
			}
		})
	}
}
