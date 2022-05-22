package endpoint_test

import (
	"context"
	"testing"
	"time"

	v1 "github.com/dzakaammar/event-scheduling-example/gen/go/proto/v1"
	"github.com/dzakaammar/event-scheduling-example/internal"
	"github.com/dzakaammar/event-scheduling-example/internal/endpoint"
	"github.com/dzakaammar/event-scheduling-example/internal/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestNewGRPCEndpoint(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		svc internal.EventService
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "OK",
			args: args{
				svc: mock.NewMockEventService(ctrl),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := endpoint.NewGRPCEndpoint(tt.args.svc)
			assert.NotNil(t, svc)
		})
	}
}

func TestGRPCEndpoint_CreateEvent(t *testing.T) {
	type fields struct {
		svcMock func(ctrl *gomock.Controller) internal.EventService
	}
	type args struct {
		ctx context.Context
		req *v1.CreateEventRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "OK",
			fields: fields{
				svcMock: func(ctrl *gomock.Controller) internal.EventService {
					svc := mock.NewMockEventService(ctrl)
					svc.EXPECT().CreateEvent(gomock.Any(), gomock.Any()).Times(1).Return(nil)
					return svc
				},
			},
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(), metadata.MD{
					"Authorization": []string{"1"},
				}),
				req: &v1.CreateEventRequest{
					Event: &v1.Event{
						Title:       "test",
						Description: "test description",
						Timezone:    "Asia/Jakarta",
						Attendees:   []string{"2", "3"},
						Schedule: []*v1.Schedule{
							{
								StartTime:     "2022-01-01T00:00:00+07:00",
								EndTime:       "2022-01-01T01:00:00+07:00",
								IsFullDay:     false,
								RecurringType: v1.RecurringType_NONE,
							},
							{
								StartTime:     "2022-01-01T01:00:00+07:00",
								EndTime:       "2022-01-01T02:00:00+07:00",
								IsFullDay:     false,
								RecurringType: v1.RecurringType_NONE,
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Not OK - service error",
			fields: fields{
				svcMock: func(ctrl *gomock.Controller) internal.EventService {
					svc := mock.NewMockEventService(ctrl)
					svc.EXPECT().CreateEvent(gomock.Any(), gomock.Any()).Times(1).Return(internal.ErrInvalidRequest)
					return svc
				},
			},
			args: args{
				ctx: context.Background(),
				req: &v1.CreateEventRequest{
					Event: &v1.Event{
						Title:       "test",
						Description: "test description",
						Timezone:    "Asia/Jakarta",
						Attendees:   []string{"2", "3"},
						Schedule: []*v1.Schedule{
							{
								StartTime:     "2022-01-01T00:00:00+07:00",
								EndTime:       "2022-01-01T01:00:00+07:00",
								IsFullDay:     false,
								RecurringType: v1.RecurringType_NONE,
							},
							{
								StartTime:     "2022-01-01T01:00:00+07:00",
								EndTime:       "2022-01-01T02:00:00+07:00",
								IsFullDay:     false,
								RecurringType: v1.RecurringType_NONE,
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Not OK - failed to parse data",
			fields: fields{
				svcMock: func(ctrl *gomock.Controller) internal.EventService {
					svc := mock.NewMockEventService(ctrl)
					return svc
				},
			},
			args: args{
				ctx: context.Background(),
				req: &v1.CreateEventRequest{
					Event: &v1.Event{
						Title:       "test",
						Description: "test description",
						Timezone:    "Asia/Jakarta",
						Attendees:   []string{"2", "3"},
						Schedule: []*v1.Schedule{
							{
								StartTime:     "invalid time",
								EndTime:       "2022-01-01T01:00:00+07:00",
								IsFullDay:     false,
								RecurringType: v1.RecurringType_NONE,
							},
							{
								StartTime:     "2022-01-01T01:00:00+07:00",
								EndTime:       "2022-01-01T02:00:00+07:00",
								IsFullDay:     false,
								RecurringType: v1.RecurringType_NONE,
							},
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			g := endpoint.NewGRPCEndpoint(tt.fields.svcMock(ctrl))
			got, err := g.CreateEvent(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NotNil(t, got)
			assert.NotEmpty(t, got.GetId())
		})
	}
}

func TestGRPCEndpoint_DeleteEventByID(t *testing.T) {
	type fields struct {
		svcMock func(ctrl *gomock.Controller) internal.EventService
	}
	type args struct {
		ctx context.Context
		req *v1.DeleteEventByIDRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *emptypb.Empty
		wantErr bool
	}{
		{
			name: "OK",
			fields: fields{
				svcMock: func(ctrl *gomock.Controller) internal.EventService {
					svcMock := mock.NewMockEventService(ctrl)
					svcMock.EXPECT().DeleteEventByID(gomock.Any(), &internal.DeleteEventByIDRequest{
						ActorID: "1",
						EventID: "test123",
					}).Times(1).Return(nil)
					return svcMock
				},
			},
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(), metadata.MD{
					"Authorization": []string{"1"},
				}),
				req: &v1.DeleteEventByIDRequest{
					Id: "test123",
				},
			},
			want: &emptypb.Empty{},
		},
		{
			name: "Not OK - failed from service",
			fields: fields{
				svcMock: func(ctrl *gomock.Controller) internal.EventService {
					svcMock := mock.NewMockEventService(ctrl)
					svcMock.EXPECT().DeleteEventByID(gomock.Any(), &internal.DeleteEventByIDRequest{
						ActorID: "1",
						EventID: "test123",
					}).Times(1).Return(internal.ErrInvalidRequest)
					return svcMock
				},
			},
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(), metadata.MD{
					"Authorization": []string{"1"},
				}),
				req: &v1.DeleteEventByIDRequest{
					Id: "test123",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			g := endpoint.NewGRPCEndpoint(tt.fields.svcMock(ctrl))
			got, err := g.DeleteEventByID(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGRPCEndpoint_UpdateEvent(t *testing.T) {
	type fields struct {
		svcMock func(ctrl *gomock.Controller) internal.EventService
	}
	type args struct {
		ctx context.Context
		req *v1.UpdateEventRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *emptypb.Empty
		wantErr bool
	}{
		{
			name: "OK",
			fields: fields{
				svcMock: func(ctrl *gomock.Controller) internal.EventService {
					svcMock := mock.NewMockEventService(ctrl)
					svcMock.EXPECT().UpdateEvent(gomock.Any(), gomock.Any()).Times(1).
						Return(nil)

					return svcMock
				},
			},
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(), metadata.MD{
					"Authorization": []string{"1"},
				}),
				req: &v1.UpdateEventRequest{
					Id: "test123",
					Event: &v1.Event{
						Title:       "test",
						Description: "test",
						Timezone:    "Asia/Jakarta",
						Attendees:   []string{"2", "3"},
						Schedule: []*v1.Schedule{
							{
								StartTime:     "2022-01-01T00:00:00+07:00",
								EndTime:       "2022-01-01T01:00:00+07:00",
								IsFullDay:     false,
								RecurringType: v1.RecurringType_NONE,
							},
							{
								StartTime:     "2022-01-01T01:00:00+07:00",
								EndTime:       "2022-01-01T02:00:00+07:00",
								IsFullDay:     false,
								RecurringType: v1.RecurringType_NONE,
							},
						},
					},
				},
			},
			want: &emptypb.Empty{},
		},
		{
			name: "Not OK - error from service",
			fields: fields{
				svcMock: func(ctrl *gomock.Controller) internal.EventService {
					svcMock := mock.NewMockEventService(ctrl)
					svcMock.EXPECT().UpdateEvent(gomock.Any(), gomock.Any()).Times(1).
						Return(internal.ErrInvalidRequest)

					return svcMock
				},
			},
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(), metadata.MD{
					"Authorization": []string{"1"},
				}),
				req: &v1.UpdateEventRequest{
					Id: "test123",
					Event: &v1.Event{
						Title:       "test",
						Description: "test",
						Timezone:    "Asia/Jakarta",
						Attendees:   []string{"2", "3"},
						Schedule: []*v1.Schedule{
							{
								StartTime:     "2022-01-01T00:00:00+07:00",
								EndTime:       "2022-01-01T01:00:00+07:00",
								IsFullDay:     false,
								RecurringType: v1.RecurringType_NONE,
							},
							{
								StartTime:     "2022-01-01T01:00:00+07:00",
								EndTime:       "2022-01-01T02:00:00+07:00",
								IsFullDay:     false,
								RecurringType: v1.RecurringType_NONE,
							},
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			g := endpoint.NewGRPCEndpoint(tt.fields.svcMock(ctrl))
			got, err := g.UpdateEvent(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGRPCEndpoint_FindEventByID(t *testing.T) {
	type fields struct {
		svcMock func(ctrl *gomock.Controller) internal.EventService
	}
	type args struct {
		ctx context.Context
		req *v1.FindEventByIDRequest
	}

	updatedAt := time.Date(2022, 0o1, 0o1, 0o0, 0o0, 0o0, 0o0, time.UTC)
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *v1.FindEventByIDResponse
		wantErr bool
	}{
		{
			name: "OK",
			fields: fields{
				svcMock: func(ctrl *gomock.Controller) internal.EventService {
					svcMock := mock.NewMockEventService(ctrl)
					svcMock.EXPECT().FindEventByID(gomock.Any(), &internal.FindEventByIDRequest{
						EventID: "test123",
					}).Times(1).Return(&internal.Event{
						ID:          "test123",
						Title:       "test123",
						Description: "test123",
						Timezone:    "Asia/Jakarta",
						Schedules: []internal.Schedule{
							{
								ID:                "sch1",
								EventID:           "test123",
								StartTime:         time.Date(2022, 0o1, 0o1, 0o0, 0o0, 0o0, 0o0, time.UTC).Unix(),
								Duration:          120,
								IsFullDay:         false,
								RecurringType:     internal.RecurringType_None,
								RecurringInterval: 0,
							},
						},
						CreatedBy: "1",
						Invitations: []internal.Invitation{
							{
								ID:      "123",
								EventID: "test123",
								UserID:  "2",
								Status:  internal.InvitationStatus_Unknown,
								Token:   "test",
							},
						},
						CreatedAt: time.Date(2022, 0o1, 0o1, 0o0, 0o0, 0o0, 0o0, time.UTC),
						UpdatedAt: &updatedAt,
					}, nil)

					return svcMock
				},
			},
			args: args{
				ctx: context.Background(),
				req: &v1.FindEventByIDRequest{
					Id: "test123",
				},
			},
			want: &v1.FindEventByIDResponse{
				Event: &v1.Event{
					Id:          "test123",
					Title:       "test123",
					Description: "test123",
					Attendees:   []string{"2"},
					Timezone:    "Asia/Jakarta",
					Schedule: []*v1.Schedule{
						{
							Id:            "sch1",
							StartTime:     "2022-01-01T07:00:00+07:00",
							EndTime:       "2022-01-01T09:00:00+07:00",
							IsFullDay:     false,
							RecurringType: v1.RecurringType_NONE,
						},
					},
					CreatedBy:     "1",
					CreatedAt:     "2022-01-01T00:00:00Z",
					LastUpdatedAt: "2022-01-01T00:00:00Z",
				},
			},
		},
		{
			name: "Not OK- error from service",
			fields: fields{
				svcMock: func(ctrl *gomock.Controller) internal.EventService {
					svcMock := mock.NewMockEventService(ctrl)
					svcMock.EXPECT().FindEventByID(gomock.Any(), &internal.FindEventByIDRequest{
						EventID: "test123",
					}).Times(1).Return(nil, internal.ErrInvalidRequest)

					return svcMock
				},
			},
			args: args{
				ctx: context.Background(),
				req: &v1.FindEventByIDRequest{
					Id: "test123",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			g := endpoint.NewGRPCEndpoint(tt.fields.svcMock(ctrl))
			got, err := g.FindEventByID(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
