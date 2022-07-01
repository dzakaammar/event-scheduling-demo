//go:build integration

package endpoint_test

import (
	"context"

	v1 "github.com/dzakaammar/event-scheduling-example/gen/go/proto/v1"
	"github.com/dzakaammar/event-scheduling-example/internal/core"
	"github.com/dzakaammar/event-scheduling-example/internal/endpoint"
	"github.com/dzakaammar/event-scheduling-example/internal/postgresql"
	"github.com/dzakaammar/event-scheduling-example/internal/scheduling"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/metadata"
)

var _ = Describe("Creating an Event", func() {
	eventRepo := postgresql.NewEventRepository(db)
	schedulingSvc := scheduling.NewService(eventRepo)
	endpoint := endpoint.NewGRPCEndpoint(schedulingSvc)

	var basedReq *v1.CreateEventRequest
	BeforeEach(func() {
		basedReq = &v1.CreateEventRequest{
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
		}
	})

	When("user is unauthorized", func() {
		It("returns error", func() {
			res, err := endpoint.CreateEvent(context.Background(), basedReq)
			Expect(err).ShouldNot(BeNil())
			Expect(res).Should(BeNil())
		})
	})

	When("user is authorized", func() {
		var (
			ctx     context.Context
			actorID = "1"
		)
		BeforeEach(func() {
			ctx = metadata.NewIncomingContext(context.Background(), metadata.MD{
				"Authorization": []string{actorID},
			})
		})

		When("the start time format is invalid", func() {
			It("returns error", func() {
				basedReq.Event.Schedule[0].StartTime = "invalid time"
				res, err := endpoint.CreateEvent(ctx, basedReq)
				Expect(err).ShouldNot(BeNil())
				Expect(res).Should(BeNil())
			})
		})

		When("the end time format is invalid", func() {
			It("returns error", func() {
				basedReq.Event.Schedule[0].EndTime = "invalid time"
				res, err := endpoint.CreateEvent(ctx, basedReq)
				Expect(err).ShouldNot(BeNil())
				Expect(res).Should(BeNil())
			})
		})

		When("the timezone format is invalid", func() {
			It("returns error", func() {
				basedReq.Event.Timezone = "invalid"
				res, err := endpoint.CreateEvent(ctx, basedReq)
				Expect(err).ShouldNot(BeNil())
				Expect(res).Should(BeNil())
			})
		})

		When("there's no schedule", func() {
			It("returns error", func() {
				basedReq.Event.Schedule = nil
				res, err := endpoint.CreateEvent(ctx, basedReq)
				Expect(err).ShouldNot(BeNil())
				Expect(res).Should(BeNil())
			})
		})

		When("the event data is valid", func() {
			It("creates the event", func() {
				res, err := endpoint.CreateEvent(ctx, basedReq)
				Expect(err).Should(BeNil())
				Expect(res).ShouldNot(BeNil())

				e, err := eventRepo.FindByID(context.Background(), res.GetId())
				Expect(err).Should(BeNil())
				Expect(e).ShouldNot(BeNil())

				Expect(e.ID).To(Equal(res.GetId()))
				Expect(e.Title).To(Equal(basedReq.Event.Title))
				Expect(e.Description).To(Equal(basedReq.Event.Description))
				Expect(e.Timezone).To(Equal(basedReq.Event.Timezone))
				Expect(len(e.Schedules)).To(Equal(len(basedReq.Event.Schedule)))
				Expect(len(e.Invitations)).To(Equal(len(basedReq.Event.Attendees)))
				Expect(e.CreatedBy).To(Equal(actorID))
				Expect(e.CreatedAt).NotTo(BeNil())
			})
		})

		When("the event has no attendees", func() {
			It("creates the event", func() {
				basedReq.Event.Attendees = nil
				res, err := endpoint.CreateEvent(ctx, basedReq)
				Expect(err).Should(BeNil())
				Expect(res).ShouldNot(BeNil())

				e, err := eventRepo.FindByID(context.Background(), res.GetId())
				Expect(err).Should(BeNil())
				Expect(e).ShouldNot(BeNil())

				Expect(e.ID).To(Equal(res.GetId()))
			})
		})
	})
})

var _ = Describe("Deleting an Event", func() {
	eventRepo := postgresql.NewEventRepository(db)
	schedulingSvc := scheduling.NewService(eventRepo)
	endpoint := endpoint.NewGRPCEndpoint(schedulingSvc)

	var event *core.Event
	BeforeEach(func() {
		event = core.NewEvent("test_actor")
		err := eventRepo.Store(context.Background(), event)
		Expect(err).Should(BeNil())
	})

	When("user is unauthorized", func() {
		It("returns an error", func() {
			empty, err := endpoint.DeleteEventByID(context.Background(), &v1.DeleteEventByIDRequest{
				Id: event.ID,
			})
			Expect(err).ShouldNot(BeNil())
			Expect(empty).Should(BeNil())

			e, err := eventRepo.FindByID(context.Background(), event.ID)
			Expect(err).Should(BeNil())
			Expect(e).ShouldNot(BeNil())
		})
	})

	When("user is authorized", func() {
		var ctx context.Context
		BeforeEach(func() {
			ctx = metadata.NewIncomingContext(context.Background(), metadata.MD{
				"Authorization": []string{"test_actor"},
			})
		})

		When("the event is exists", func() {
			It("deletes the data", func() {
				empty, err := endpoint.DeleteEventByID(ctx, &v1.DeleteEventByIDRequest{
					Id: event.ID,
				})
				Expect(err).Should(BeNil())
				Expect(empty).ShouldNot(BeNil())

				_, err = eventRepo.FindByID(context.Background(), event.ID)
				Expect(err).ShouldNot(BeNil())
			})
		})

		When("the event is not exists", func() {
			It("returns an error", func() {
				empty, err := endpoint.DeleteEventByID(ctx, &v1.DeleteEventByIDRequest{
					Id: "invalid id",
				})
				Expect(err).ShouldNot(BeNil())
				Expect(empty).Should(BeNil())
			})
		})
	})
})

// func TestGRPCEndpoint_UpdateEvent(t *testing.T) {
// 	type fields struct {
// 		svcMock func(ctrl *gomock.Controller) core.SchedulingService
// 	}
// 	type args struct {
// 		ctx context.Context
// 		req *v1.UpdateEventRequest
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		want    *emptypb.Empty
// 		wantErr bool
// 	}{
// 		{
// 			name: "OK",
// 			fields: fields{
// 				svcMock: func(ctrl *gomock.Controller) core.SchedulingService {
// 					svcMock := mock.NewMockSchedulingService(ctrl)
// 					svcMock.EXPECT().UpdateEvent(gomock.Any(), gomock.Any()).Times(1).
// 						Return(nil)

// 					return svcMock
// 				},
// 			},
// 			args: args{
// 				ctx: metadata.NewIncomingContext(context.Background(), metadata.MD{
// 					"Authorization": []string{"1"},
// 				}),
// 				req: &v1.UpdateEventRequest{
// 					Id: "test123",
// 					Event: &v1.Event{
// 						Title:       "test",
// 						Description: "test",
// 						Timezone:    "Asia/Jakarta",
// 						Attendees:   []string{"2", "3"},
// 						Schedule: []*v1.Schedule{
// 							{
// 								StartTime:     "2022-01-01T00:00:00+07:00",
// 								EndTime:       "2022-01-01T01:00:00+07:00",
// 								IsFullDay:     false,
// 								RecurringType: v1.RecurringType_NONE,
// 							},
// 							{
// 								StartTime:     "2022-01-01T01:00:00+07:00",
// 								EndTime:       "2022-01-01T02:00:00+07:00",
// 								IsFullDay:     false,
// 								RecurringType: v1.RecurringType_NONE,
// 							},
// 						},
// 					},
// 				},
// 			},
// 			want: &emptypb.Empty{},
// 		},
// 		{
// 			name: "Not OK - error from service",
// 			fields: fields{
// 				svcMock: func(ctrl *gomock.Controller) core.SchedulingService {
// 					svcMock := mock.NewMockSchedulingService(ctrl)
// 					svcMock.EXPECT().UpdateEvent(gomock.Any(), gomock.Any()).Times(1).
// 						Return(internal.ErrInvalidRequest)

// 					return svcMock
// 				},
// 			},
// 			args: args{
// 				ctx: metadata.NewIncomingContext(context.Background(), metadata.MD{
// 					"Authorization": []string{"1"},
// 				}),
// 				req: &v1.UpdateEventRequest{
// 					Id: "test123",
// 					Event: &v1.Event{
// 						Title:       "test",
// 						Description: "test",
// 						Timezone:    "Asia/Jakarta",
// 						Attendees:   []string{"2", "3"},
// 						Schedule: []*v1.Schedule{
// 							{
// 								StartTime:     "2022-01-01T00:00:00+07:00",
// 								EndTime:       "2022-01-01T01:00:00+07:00",
// 								IsFullDay:     false,
// 								RecurringType: v1.RecurringType_NONE,
// 							},
// 							{
// 								StartTime:     "2022-01-01T01:00:00+07:00",
// 								EndTime:       "2022-01-01T02:00:00+07:00",
// 								IsFullDay:     false,
// 								RecurringType: v1.RecurringType_NONE,
// 							},
// 						},
// 					},
// 				},
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			g := endpoint.NewGRPCEndpoint(tt.fields.svcMock(ctrl))
// 			got, err := g.UpdateEvent(tt.args.ctx, tt.args.req)
// 			if tt.wantErr {
// 				assert.Error(t, err)
// 				return
// 			}
// 			assert.Equal(t, tt.want, got)
// 		})
// 	}
// }

// func TestGRPCEndpoint_FindEventByID(t *testing.T) {
// 	type fields struct {
// 		svcMock func(ctrl *gomock.Controller) core.SchedulingService
// 	}
// 	type args struct {
// 		ctx context.Context
// 		req *v1.FindEventByIDRequest
// 	}

// 	updatedAt := time.Date(2022, 0o1, 0o1, 0o0, 0o0, 0o0, 0o0, time.UTC)
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		want    *v1.FindEventByIDResponse
// 		wantErr bool
// 	}{
// 		{
// 			name: "OK",
// 			fields: fields{
// 				svcMock: func(ctrl *gomock.Controller) core.SchedulingService {
// 					svcMock := mock.NewMockSchedulingService(ctrl)
// 					svcMock.EXPECT().FindEventByID(gomock.Any(), &core.FindEventByIDRequest{
// 						EventID: "test123",
// 					}).Times(1).Return(&core.Event{
// 						ID:          "test123",
// 						Title:       "test123",
// 						Description: "test123",
// 						Timezone:    "Asia/Jakarta",
// 						Schedules: []core.Schedule{
// 							{
// 								ID:                "sch1",
// 								EventID:           "test123",
// 								StartTime:         time.Date(2022, 0o1, 0o1, 0o0, 0o0, 0o0, 0o0, time.UTC).Unix(),
// 								Duration:          120,
// 								IsFullDay:         false,
// 								RecurringType:     core.RecurringType_None,
// 								RecurringInterval: 0,
// 							},
// 						},
// 						CreatedBy: "1",
// 						Invitations: []core.Invitation{
// 							{
// 								ID:      "123",
// 								EventID: "test123",
// 								UserID:  "2",
// 								Status:  core.InvitationStatus_Unknown,
// 								Token:   "test",
// 							},
// 						},
// 						CreatedAt: time.Date(2022, 0o1, 0o1, 0o0, 0o0, 0o0, 0o0, time.UTC),
// 						UpdatedAt: &updatedAt,
// 					}, nil)

// 					return svcMock
// 				},
// 			},
// 			args: args{
// 				ctx: context.Background(),
// 				req: &v1.FindEventByIDRequest{
// 					Id: "test123",
// 				},
// 			},
// 			want: &v1.FindEventByIDResponse{
// 				Event: &v1.Event{
// 					Id:          "test123",
// 					Title:       "test123",
// 					Description: "test123",
// 					Attendees:   []string{"2"},
// 					Timezone:    "Asia/Jakarta",
// 					Schedule: []*v1.Schedule{
// 						{
// 							Id:            "sch1",
// 							StartTime:     "2022-01-01T07:00:00+07:00",
// 							EndTime:       "2022-01-01T09:00:00+07:00",
// 							IsFullDay:     false,
// 							RecurringType: v1.RecurringType_NONE,
// 						},
// 					},
// 					CreatedBy:     "1",
// 					CreatedAt:     "2022-01-01T00:00:00Z",
// 					LastUpdatedAt: "2022-01-01T00:00:00Z",
// 				},
// 			},
// 		},
// 		{
// 			name: "Not OK- error from service",
// 			fields: fields{
// 				svcMock: func(ctrl *gomock.Controller) core.SchedulingService {
// 					svcMock := mock.NewMockSchedulingService(ctrl)
// 					svcMock.EXPECT().FindEventByID(gomock.Any(), &core.FindEventByIDRequest{
// 						EventID: "test123",
// 					}).Times(1).Return(nil, internal.ErrInvalidRequest)

// 					return svcMock
// 				},
// 			},
// 			args: args{
// 				ctx: context.Background(),
// 				req: &v1.FindEventByIDRequest{
// 					Id: "test123",
// 				},
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			g := endpoint.NewGRPCEndpoint(tt.fields.svcMock(ctrl))
// 			got, err := g.FindEventByID(tt.args.ctx, tt.args.req)
// 			if tt.wantErr {
// 				assert.Error(t, err)
// 				return
// 			}
// 			assert.Equal(t, tt.want, got)
// 		})
// 	}
// }
