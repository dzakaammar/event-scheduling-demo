package scheduling_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/dzakaammar/event-scheduling-example/internal"
	"github.com/dzakaammar/event-scheduling-example/internal/core"
	"github.com/dzakaammar/event-scheduling-example/internal/mock"
	"github.com/dzakaammar/event-scheduling-example/internal/scheduling"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewEventService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		eventRepo core.EventRepository
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
			got := scheduling.NewService(tt.args.eventRepo)
			assert.NotNil(t, got)
		})
	}
}

func TestEventService_CreateEvent(t *testing.T) {
	type fields struct {
		eventRepoMock func(ctrl *gomock.Controller) core.EventRepository
	}
	type args struct {
		ctx context.Context
		req *core.CreateEventRequest
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
				eventRepoMock: func(ctrl *gomock.Controller) core.EventRepository {
					repo := mock.NewMockEventRepository(ctrl)
					repo.EXPECT().Store(gomock.Any(), gomock.Any()).Return(nil)
					return repo
				},
			},
			args: args{
				ctx: context.Background(),
				req: &core.CreateEventRequest{
					ActorID: "123",
					Event: &core.Event{
						ID:          "123",
						Title:       "test",
						Description: "test123",
						Timezone:    "Asia/Jakarta",
						Schedules: []core.Schedule{
							{
								ID:                "test123",
								EventID:           "123",
								StartTime:         time.Now().Unix(),
								DurationInMinutes: 120,
								IsFullDay:         false,
								RecurringType:     core.RecurringType_None,
							},
							{
								ID:                "test1234",
								EventID:           "123",
								StartTime:         time.Now().Unix(),
								DurationInMinutes: 120,
								IsFullDay:         false,
								RecurringType:     core.RecurringType_None,
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
				eventRepoMock: func(ctrl *gomock.Controller) core.EventRepository {
					repo := mock.NewMockEventRepository(ctrl)
					repo.EXPECT().Store(gomock.Any(), gomock.Any()).Return(internal.ErrValidationFailed)
					return repo
				},
			},
			args: args{
				ctx: context.Background(),
				req: &core.CreateEventRequest{
					ActorID: "123",
					Event: &core.Event{
						ID:          "123",
						Title:       "test",
						Description: "test123",
						Timezone:    "Asia/Jakarta",
						Schedules: []core.Schedule{
							{
								ID:                "test123",
								EventID:           "123",
								StartTime:         time.Now().Unix(),
								DurationInMinutes: 120,
								IsFullDay:         false,
								RecurringType:     core.RecurringType_None,
							},
							{
								ID:                "test1234",
								EventID:           "123",
								StartTime:         time.Now().Unix(),
								DurationInMinutes: 120,
								IsFullDay:         false,
								RecurringType:     core.RecurringType_None,
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
				eventRepoMock: func(ctrl *gomock.Controller) core.EventRepository {
					repo := mock.NewMockEventRepository(ctrl)
					return repo
				},
			},
			args: args{
				ctx: context.Background(),
				req: &core.CreateEventRequest{
					ActorID: "123",
					Event: &core.Event{
						ID:          "",
						Title:       "",
						Description: "",
						Timezone:    "",
						Schedules: []core.Schedule{
							{
								ID:                "test123",
								EventID:           "123",
								StartTime:         time.Now().Unix(),
								DurationInMinutes: 120,
								IsFullDay:         false,
								RecurringType:     core.RecurringType_None,
							},
							{
								ID:                "test1234",
								EventID:           "123",
								StartTime:         time.Now().Unix(),
								DurationInMinutes: 120,
								IsFullDay:         false,
								RecurringType:     core.RecurringType_None,
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

			e := scheduling.NewService(tt.fields.eventRepoMock(ctrl))
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
		eventRepoMock func(ctrl *gomock.Controller) core.EventRepository
	}
	type args struct {
		ctx context.Context
		req *core.DeleteEventByIDRequest
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
				eventRepoMock: func(ctrl *gomock.Controller) core.EventRepository {
					repo := mock.NewMockEventRepository(ctrl)
					repo.EXPECT().FindByID(gomock.Any(), gomock.Any()).Times(1).
						Return(&core.Event{}, nil)
					repo.EXPECT().DeleteByID(gomock.Any(), gomock.Any()).Times(1).
						Return(nil)
					return repo
				},
			},
			args: args{
				ctx: context.Background(),
				req: &core.DeleteEventByIDRequest{
					ActorID: "test123",
					EventID: "123",
				},
			},
			wantErr: false,
		},
		{
			name: "Not OK - data not found",
			fields: fields{
				eventRepoMock: func(ctrl *gomock.Controller) core.EventRepository {
					repo := mock.NewMockEventRepository(ctrl)
					repo.EXPECT().FindByID(gomock.Any(), gomock.Any()).Times(1).
						Return(nil, errors.New("data not found")) //nolint:goerr113
					return repo
				},
			},
			args: args{
				ctx: context.Background(),
				req: &core.DeleteEventByIDRequest{
					ActorID: "test123",
					EventID: "123",
				},
			},
			wantErr: true,
		},
		{
			name: "Not OK - error from repo",
			fields: fields{
				eventRepoMock: func(ctrl *gomock.Controller) core.EventRepository {
					repo := mock.NewMockEventRepository(ctrl)
					repo.EXPECT().FindByID(gomock.Any(), gomock.Any()).Times(1).
						Return(&core.Event{}, nil)
					repo.EXPECT().DeleteByID(gomock.Any(), gomock.Any()).Times(1).
						Return(internal.ErrInvalidRequest)
					return repo
				},
			},
			args: args{
				ctx: context.Background(),
				req: &core.DeleteEventByIDRequest{
					ActorID: "test123",
					EventID: "123",
				},
			},
			wantErr: true,
		},
		{
			name: "Not OK - invalid request",
			fields: fields{
				eventRepoMock: func(ctrl *gomock.Controller) core.EventRepository {
					repo := mock.NewMockEventRepository(ctrl)
					return repo
				},
			},
			args: args{
				ctx: context.Background(),
				req: &core.DeleteEventByIDRequest{
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

			e := scheduling.NewService(tt.fields.eventRepoMock(ctrl))
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
		eventRepoMock func(ctrl *gomock.Controller) core.EventRepository
	}
	type args struct {
		ctx context.Context
		req *core.UpdateEventRequest
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
				eventRepoMock: func(ctrl *gomock.Controller) core.EventRepository {
					repo := mock.NewMockEventRepository(ctrl)
					repo.EXPECT().Update(gomock.Any(), gomock.Any()).Times(1).
						Return(nil)
					return repo
				},
			},
			args: args{
				ctx: context.Background(),
				req: &core.UpdateEventRequest{
					ID:      "test123",
					ActorID: "test123",
					Event: &core.Event{
						ID:          "test123",
						Title:       "updated",
						Description: "description",
						Timezone:    "Asia/Jakarta",
						Schedules: []core.Schedule{
							{
								ID:                "sch1",
								EventID:           "test123",
								StartTime:         time.Now().Unix(),
								DurationInMinutes: 120,
								IsFullDay:         false,
								RecurringType:     core.RecurringType_None,
							},
							{
								ID:                "sch1",
								EventID:           "test123",
								StartTime:         time.Now().Unix(),
								DurationInMinutes: 120,
								IsFullDay:         false,
								RecurringType:     core.RecurringType_None,
							},
						},
						Invitations: []core.Invitation{
							{
								ID:      "inv1",
								EventID: "test123",
								UserID:  "2",
								Status:  core.InvitationStatus_Confirmed,
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
				eventRepoMock: func(ctrl *gomock.Controller) core.EventRepository {
					repo := mock.NewMockEventRepository(ctrl)
					repo.EXPECT().Update(gomock.Any(), gomock.Any()).Times(1).
						Return(internal.ErrInvalidRequest)
					return repo
				},
			},
			args: args{
				ctx: context.Background(),
				req: &core.UpdateEventRequest{
					ID:      "test123",
					ActorID: "test123",
					Event: &core.Event{
						ID:          "test123",
						Title:       "updated",
						Description: "description",
						Timezone:    "Asia/Jakarta",
						Schedules: []core.Schedule{
							{
								ID:                "sch1",
								EventID:           "test123",
								StartTime:         time.Now().Unix(),
								DurationInMinutes: 120,
								IsFullDay:         false,
								RecurringType:     core.RecurringType_None,
							},
							{
								ID:                "sch1",
								EventID:           "test123",
								StartTime:         time.Now().Unix(),
								DurationInMinutes: 120,
								IsFullDay:         false,
								RecurringType:     core.RecurringType_None,
							},
						},
						Invitations: []core.Invitation{
							{
								ID:      "inv1",
								EventID: "test123",
								UserID:  "2",
								Status:  core.InvitationStatus_Confirmed,
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
				eventRepoMock: func(ctrl *gomock.Controller) core.EventRepository {
					repo := mock.NewMockEventRepository(ctrl)
					return repo
				},
			},
			args: args{
				ctx: context.Background(),
				req: &core.UpdateEventRequest{
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

			e := scheduling.NewService(tt.fields.eventRepoMock(ctrl))
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
		eventRepoMock func(ctrl *gomock.Controller) core.EventRepository
	}
	type args struct {
		ctx context.Context
		req *core.FindEventByIDRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *core.Event
		wantErr bool
	}{
		{
			name: "OK",
			fields: fields{
				eventRepoMock: func(ctrl *gomock.Controller) core.EventRepository {
					repo := mock.NewMockEventRepository(ctrl)
					repo.EXPECT().FindByID(gomock.Any(), gomock.Any()).Times(1).
						Return(&core.Event{
							ID:          "123",
							Title:       "test123",
							Description: "test123",
							Timezone:    "Asia/Jakarta",
							Schedules: []core.Schedule{
								{
									ID:                "sch1",
									EventID:           "123",
									StartTime:         time.Now().Unix(),
									DurationInMinutes: 120,
									IsFullDay:         false,
									RecurringType:     core.RecurringType_None,
								},
							},
							Invitations: []core.Invitation{
								{
									ID:      "invitation1",
									EventID: "123",
									UserID:  "2",
									Status:  core.InvitationStatus_Unknown,
									Token:   "123",
								},
								{
									ID:      "invitation2",
									EventID: "123",
									UserID:  "3",
									Status:  core.InvitationStatus_Unknown,
									Token:   "123",
								},
							},
						}, nil)
					return repo
				},
			},
			args: args{
				ctx: context.Background(),
				req: &core.FindEventByIDRequest{
					EventID: "123",
				},
			},
			want: &core.Event{
				ID:          "123",
				Title:       "test123",
				Description: "test123",
				Timezone:    "Asia/Jakarta",
				Schedules: []core.Schedule{
					{
						ID:                "sch1",
						EventID:           "123",
						StartTime:         time.Now().Unix(),
						DurationInMinutes: 120,
						IsFullDay:         false,
						RecurringType:     core.RecurringType_None,
					},
				},
				Invitations: []core.Invitation{
					{
						ID:      "invitation1",
						EventID: "123",
						UserID:  "2",
						Status:  core.InvitationStatus_Unknown,
						Token:   "123",
					},
					{
						ID:      "invitation2",
						EventID: "123",
						UserID:  "3",
						Status:  core.InvitationStatus_Unknown,
						Token:   "123",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Not OK - error from repo",
			fields: fields{
				eventRepoMock: func(ctrl *gomock.Controller) core.EventRepository {
					repo := mock.NewMockEventRepository(ctrl)
					repo.EXPECT().FindByID(gomock.Any(), gomock.Any()).Times(1).
						Return(nil, internal.ErrInvalidRequest)
					return repo
				},
			},
			args: args{
				ctx: context.Background(),
				req: &core.FindEventByIDRequest{
					EventID: "123",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Not OK - invalid request",
			fields: fields{
				eventRepoMock: func(ctrl *gomock.Controller) core.EventRepository {
					repo := mock.NewMockEventRepository(ctrl)
					return repo
				},
			},
			args: args{
				ctx: context.Background(),
				req: &core.FindEventByIDRequest{
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

			e := scheduling.NewService(tt.fields.eventRepoMock(ctrl))
			got, err := e.FindEventByID(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.EqualValues(t, tt.want, got)
		})
	}
}
