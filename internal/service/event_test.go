package service_test

import (
	"context"
	"testing"
	"time"

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
		eventRepoMock func(ctrl *gomock.Controller) internal.EventRepository
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
		{
			name: "OK",
			fields: fields{
				eventRepoMock: func(ctrl *gomock.Controller) internal.EventRepository {
					repo := mock.NewMockEventRepository(ctrl)
					repo.EXPECT().Store(gomock.Any(), gomock.Any()).Return(nil)
					return repo
				},
			},
			args: args{
				ctx: context.Background(),
				req: &internal.CreateEventRequest{
					ActorID: "123",
					Event: &internal.Event{
						ID:          "123",
						Title:       "test",
						Description: "test123",
						Timezone:    "Asia/Jakarta",
						Schedules: []internal.Schedule{
							{
								ID:            "test123",
								EventID:       "123",
								StartTime:     time.Now().Unix(),
								Duration:      120,
								IsFullDay:     false,
								RecurringType: internal.RecurringType_None,
							},
							{
								ID:            "test1234",
								EventID:       "123",
								StartTime:     time.Now().Unix(),
								Duration:      120,
								IsFullDay:     false,
								RecurringType: internal.RecurringType_None,
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Not OK - error from repo",
			fields: fields{
				eventRepoMock: func(ctrl *gomock.Controller) internal.EventRepository {
					repo := mock.NewMockEventRepository(ctrl)
					repo.EXPECT().Store(gomock.Any(), gomock.Any()).Return(internal.ErrValidationFailed)
					return repo
				},
			},
			args: args{
				ctx: context.Background(),
				req: &internal.CreateEventRequest{
					ActorID: "123",
					Event: &internal.Event{
						ID:          "123",
						Title:       "test",
						Description: "test123",
						Timezone:    "Asia/Jakarta",
						Schedules: []internal.Schedule{
							{
								ID:            "test123",
								EventID:       "123",
								StartTime:     time.Now().Unix(),
								Duration:      120,
								IsFullDay:     false,
								RecurringType: internal.RecurringType_None,
							},
							{
								ID:            "test1234",
								EventID:       "123",
								StartTime:     time.Now().Unix(),
								Duration:      120,
								IsFullDay:     false,
								RecurringType: internal.RecurringType_None,
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Not OK - validation error",
			fields: fields{
				eventRepoMock: func(ctrl *gomock.Controller) internal.EventRepository {
					repo := mock.NewMockEventRepository(ctrl)
					return repo
				},
			},
			args: args{
				ctx: context.Background(),
				req: &internal.CreateEventRequest{
					ActorID: "123",
					Event: &internal.Event{
						ID:          "",
						Title:       "",
						Description: "",
						Timezone:    "",
						Schedules: []internal.Schedule{
							{
								ID:            "test123",
								EventID:       "123",
								StartTime:     time.Now().Unix(),
								Duration:      120,
								IsFullDay:     false,
								RecurringType: internal.RecurringType_None,
							},
							{
								ID:            "test1234",
								EventID:       "123",
								StartTime:     time.Now().Unix(),
								Duration:      120,
								IsFullDay:     false,
								RecurringType: internal.RecurringType_None,
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

			e := service.NewEventService(tt.fields.eventRepoMock(ctrl))
			err := e.CreateEvent(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestEventService_DeleteEventByID(t *testing.T) {
	type fields struct {
		eventRepoMock func(ctrl *gomock.Controller) internal.EventRepository
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
		{
			name: "OK",
			fields: fields{
				eventRepoMock: func(ctrl *gomock.Controller) internal.EventRepository {
					repo := mock.NewMockEventRepository(ctrl)
					repo.EXPECT().DeleteByID(gomock.Any(), gomock.Any()).Times(1).
						Return(nil)
					return repo
				},
			},
			args: args{
				ctx: context.Background(),
				req: &internal.DeleteEventByIDRequest{
					ActorID: "test123",
					EventID: "123",
				},
			},
			wantErr: false,
		},
		{
			name: "Not OK - error from repo",
			fields: fields{
				eventRepoMock: func(ctrl *gomock.Controller) internal.EventRepository {
					repo := mock.NewMockEventRepository(ctrl)
					repo.EXPECT().DeleteByID(gomock.Any(), gomock.Any()).Times(1).
						Return(internal.ErrInvalidRequest)
					return repo
				},
			},
			args: args{
				ctx: context.Background(),
				req: &internal.DeleteEventByIDRequest{
					ActorID: "test123",
					EventID: "123",
				},
			},
			wantErr: true,
		},
		{
			name: "Not OK - invalid request",
			fields: fields{
				eventRepoMock: func(ctrl *gomock.Controller) internal.EventRepository {
					repo := mock.NewMockEventRepository(ctrl)
					return repo
				},
			},
			args: args{
				ctx: context.Background(),
				req: &internal.DeleteEventByIDRequest{
					ActorID: "",
					EventID: "",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			e := service.NewEventService(tt.fields.eventRepoMock(ctrl))
			err := e.DeleteEventByID(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestEventService_UpdateEvent(t *testing.T) {
	type fields struct {
		eventRepoMock func(ctrl *gomock.Controller) internal.EventRepository
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
		{
			name: "OK",
			fields: fields{
				eventRepoMock: func(ctrl *gomock.Controller) internal.EventRepository {
					repo := mock.NewMockEventRepository(ctrl)
					repo.EXPECT().Update(gomock.Any(), gomock.Any()).Times(1).
						Return(nil)
					return repo
				},
			},
			args: args{
				ctx: context.Background(),
				req: &internal.UpdateEventRequest{
					ID:      "test123",
					ActorID: "test123",
					Event: &internal.Event{
						ID:          "test123",
						Title:       "updated",
						Description: "description",
						Timezone:    "Asia/Jakarta",
						Schedules: []internal.Schedule{
							{
								ID:            "sch1",
								EventID:       "test123",
								StartTime:     time.Now().Unix(),
								Duration:      120,
								IsFullDay:     false,
								RecurringType: internal.RecurringType_None,
							},
							{
								ID:            "sch1",
								EventID:       "test123",
								StartTime:     time.Now().Unix(),
								Duration:      120,
								IsFullDay:     false,
								RecurringType: internal.RecurringType_None,
							},
						},
						Invitations: []internal.Invitation{
							{
								ID:      "inv1",
								EventID: "test123",
								UserID:  "2",
								Status:  internal.InvitationStatus_Confirmed,
								Token:   "123",
							},
						},
					},
				},
			},
		},
		{
			name: "Not OK - error from repo",
			fields: fields{
				eventRepoMock: func(ctrl *gomock.Controller) internal.EventRepository {
					repo := mock.NewMockEventRepository(ctrl)
					repo.EXPECT().Update(gomock.Any(), gomock.Any()).Times(1).
						Return(internal.ErrInvalidRequest)
					return repo
				},
			},
			args: args{
				ctx: context.Background(),
				req: &internal.UpdateEventRequest{
					ID:      "test123",
					ActorID: "test123",
					Event: &internal.Event{
						ID:          "test123",
						Title:       "updated",
						Description: "description",
						Timezone:    "Asia/Jakarta",
						Schedules: []internal.Schedule{
							{
								ID:            "sch1",
								EventID:       "test123",
								StartTime:     time.Now().Unix(),
								Duration:      120,
								IsFullDay:     false,
								RecurringType: internal.RecurringType_None,
							},
							{
								ID:            "sch1",
								EventID:       "test123",
								StartTime:     time.Now().Unix(),
								Duration:      120,
								IsFullDay:     false,
								RecurringType: internal.RecurringType_None,
							},
						},
						Invitations: []internal.Invitation{
							{
								ID:      "inv1",
								EventID: "test123",
								UserID:  "2",
								Status:  internal.InvitationStatus_Confirmed,
								Token:   "123",
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Not OK - validation error",
			fields: fields{
				eventRepoMock: func(ctrl *gomock.Controller) internal.EventRepository {
					repo := mock.NewMockEventRepository(ctrl)
					return repo
				},
			},
			args: args{
				ctx: context.Background(),
				req: &internal.UpdateEventRequest{
					ID:      "",
					ActorID: "",
					Event:   nil,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			e := service.NewEventService(tt.fields.eventRepoMock(ctrl))
			err := e.UpdateEvent(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestEventService_FindEventByID(t *testing.T) {
	type fields struct {
		eventRepoMock func(ctrl *gomock.Controller) internal.EventRepository
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
		{
			name: "OK",
			fields: fields{
				eventRepoMock: func(ctrl *gomock.Controller) internal.EventRepository {
					repo := mock.NewMockEventRepository(ctrl)
					repo.EXPECT().FindByID(gomock.Any(), gomock.Any()).Times(1).
						Return(&internal.Event{
							ID:          "123",
							Title:       "test123",
							Description: "test123",
							Timezone:    "Asia/Jakarta",
							Schedules: []internal.Schedule{
								{
									ID:            "sch1",
									EventID:       "123",
									StartTime:     time.Now().Unix(),
									Duration:      120,
									IsFullDay:     false,
									RecurringType: internal.RecurringType_None,
								},
							},
							Invitations: []internal.Invitation{
								{
									ID:      "invitation1",
									EventID: "123",
									UserID:  "2",
									Status:  internal.InvitationStatus_Unknown,
									Token:   "123",
								},
								{
									ID:      "invitation2",
									EventID: "123",
									UserID:  "3",
									Status:  internal.InvitationStatus_Unknown,
									Token:   "123",
								},
							},
						}, nil)
					return repo
				},
			},
			args: args{
				ctx: context.Background(),
				req: &internal.FindEventByIDRequest{
					EventID: "123",
				},
			},
			want: &internal.Event{
				ID:          "123",
				Title:       "test123",
				Description: "test123",
				Timezone:    "Asia/Jakarta",
				Schedules: []internal.Schedule{
					{
						ID:            "sch1",
						EventID:       "123",
						StartTime:     time.Now().Unix(),
						Duration:      120,
						IsFullDay:     false,
						RecurringType: internal.RecurringType_None,
					},
				},
				Invitations: []internal.Invitation{
					{
						ID:      "invitation1",
						EventID: "123",
						UserID:  "2",
						Status:  internal.InvitationStatus_Unknown,
						Token:   "123",
					},
					{
						ID:      "invitation2",
						EventID: "123",
						UserID:  "3",
						Status:  internal.InvitationStatus_Unknown,
						Token:   "123",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Not OK - error from repo",
			fields: fields{
				eventRepoMock: func(ctrl *gomock.Controller) internal.EventRepository {
					repo := mock.NewMockEventRepository(ctrl)
					repo.EXPECT().FindByID(gomock.Any(), gomock.Any()).Times(1).
						Return(nil, internal.ErrInvalidRequest)
					return repo
				},
			},
			args: args{
				ctx: context.Background(),
				req: &internal.FindEventByIDRequest{
					EventID: "123",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Not OK - invalid request",
			fields: fields{
				eventRepoMock: func(ctrl *gomock.Controller) internal.EventRepository {
					repo := mock.NewMockEventRepository(ctrl)
					return repo
				},
			},
			args: args{
				ctx: context.Background(),
				req: &internal.FindEventByIDRequest{
					EventID: "",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			e := service.NewEventService(tt.fields.eventRepoMock(ctrl))
			got, err := e.FindEventByID(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.EqualValues(t, tt.want, got)
		})
	}
}
