package endpoint

import (
	"context"
	"time"

	v1 "github.com/dzakaammar/event-scheduling-example/gen/go/proto/v1"
	"github.com/dzakaammar/event-scheduling-example/internal/core"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GRPCEndpoint struct {
	v1.UnimplementedAPIServer
	svc core.SchedulingService
}

func NewGRPCEndpoint(svc core.SchedulingService) *GRPCEndpoint {
	return &GRPCEndpoint{
		svc: svc,
	}
}

func (g *GRPCEndpoint) CreateEvent(ctx context.Context, req *v1.CreateEventRequest) (*v1.CreateEventResponse, error) {
	createReq, err := parseCreateEventRequest(ctx, req)
	if err != nil {
		log.Error(err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = g.svc.CreateEvent(ctx, createReq)
	if err != nil {
		log.Error(err)
		return nil, mapErrToStatusCode(err)
	}
	return &v1.CreateEventResponse{
		Id: createReq.Event.ID,
	}, nil
}

func (g *GRPCEndpoint) DeleteEventByID(ctx context.Context, req *v1.DeleteEventByIDRequest) (*emptypb.Empty, error) {
	delReq, err := parseDeleteEventByIDtRequest(ctx, req)
	if err != nil {
		log.Error(err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = g.svc.DeleteEventByID(ctx, delReq)
	if err != nil {
		log.Error(err)
		return nil, mapErrToStatusCode(err)
	}
	return &emptypb.Empty{}, nil
}

func (g *GRPCEndpoint) UpdateEvent(ctx context.Context, req *v1.UpdateEventRequest) (*emptypb.Empty, error) {
	updateReq, err := parseUpdateEventByIDRequest(ctx, req)
	if err != nil {
		log.Error(err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = g.svc.UpdateEvent(ctx, updateReq)
	if err != nil {
		log.Error(err)
		return nil, mapErrToStatusCode(err)
	}
	return &emptypb.Empty{}, nil
}

func (g *GRPCEndpoint) FindEventByID(ctx context.Context, req *v1.FindEventByIDRequest) (*v1.FindEventByIDResponse, error) {
	event, err := g.svc.FindEventByID(ctx, &core.FindEventByIDRequest{EventID: req.GetId()})
	if err != nil {
		log.Error(err)
		return nil, mapErrToStatusCode(err)
	}

	res, err := parseEventToPB(event)
	if err != nil {
		return nil, mapErrToStatusCode(err)
	}

	return &v1.FindEventByIDResponse{
		Event: res,
	}, nil
}

func (g *GRPCEndpoint) Check(ctx context.Context, _ *v1.HealthCheckRequest) (*v1.HealthCheckResponse, error) {
	return &v1.HealthCheckResponse{Status: v1.HealthCheckResponse_SERVING}, nil
}

func (g *GRPCEndpoint) Watch(_ *v1.HealthCheckRequest, _ v1.API_WatchServer) error {
	return status.Error(codes.Unimplemented, "unimplemented")
}

func extractAuthorization(ctx context.Context) string {
	m, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	a := m.Get("Authorization")
	if len(a) == 0 {
		return ""
	}
	return a[0]
}

func parseCreateEventRequest(ctx context.Context, req *v1.CreateEventRequest) (*core.CreateEventRequest, error) {
	if req == nil || req.GetEvent() == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	actorID := extractAuthorization(ctx)

	event := core.NewEvent()
	event.Title = req.GetEvent().GetTitle()
	event.Description = req.GetEvent().GetDescription()
	event.Timezone = req.GetEvent().GetTimezone()
	event.CreatedBy = actorID
	event.CreatedAt = time.Now()

	sch, err := parseSchedules(req.GetEvent().GetSchedule(), event.ID)
	if err != nil {
		return nil, err
	}
	event.Schedules = sch

	event.Invitations = parseInvitations(req.GetEvent().GetAttendees(), event.ID)

	return &core.CreateEventRequest{
		ActorID: actorID,
		Event:   event,
	}, nil
}

func parseDeleteEventByIDtRequest(ctx context.Context, req *v1.DeleteEventByIDRequest) (*core.DeleteEventByIDRequest, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	return &core.DeleteEventByIDRequest{
		ActorID: extractAuthorization(ctx),
		EventID: req.GetId(),
	}, nil
}

func parseUpdateEventByIDRequest(ctx context.Context, req *v1.UpdateEventRequest) (*core.UpdateEventRequest, error) {
	if req == nil || req.GetEvent() == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	now := time.Now()
	event := core.Event{
		ID:          req.GetId(),
		Title:       req.GetEvent().GetTitle(),
		Description: req.GetEvent().GetDescription(),
		Timezone:    req.GetEvent().GetTimezone(),
		UpdatedAt:   &now,
	}

	sch, err := parseSchedules(req.GetEvent().GetSchedule(), event.ID)
	if err != nil {
		return nil, err
	}
	event.Schedules = sch
	event.Invitations = parseInvitations(req.GetEvent().GetAttendees(), event.ID)

	return &core.UpdateEventRequest{
		ID:      req.GetId(),
		ActorID: extractAuthorization(ctx),
		Event:   &event,
	}, nil
}

func parseSchedules(sch []*v1.Schedule, eventID string) ([]core.Schedule, error) {
	schedules := make([]core.Schedule, len(sch))
	for index, sch := range sch {
		s, err := core.NewSchedule(
			eventID,
			sch.GetStartTime(),
			sch.GetEndTime(),
			sch.GetIsFullDay(),
			mapRecurringType(sch.GetRecurringType()),
		)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		schedules[index] = s
	}
	return schedules, nil
}

func parseInvitations(attendees []string, eventID string) []core.Invitation {
	invitations := make([]core.Invitation, len(attendees))
	for index, userID := range attendees {
		i := core.NewInvitation(eventID, userID)
		invitations[index] = i
	}
	return invitations
}

func parseEventToPB(event *core.Event) (*v1.Event, error) {
	e := &v1.Event{
		Id:            event.ID,
		Title:         event.Title,
		Description:   event.Description,
		Timezone:      event.Timezone,
		CreatedAt:     event.CreatedAt.Format(time.RFC3339),
		CreatedBy:     event.CreatedBy,
		LastUpdatedAt: event.GetUpdatedAt(),
	}

	schedules := make([]*v1.Schedule, len(event.Schedules))
	for index, sch := range event.Schedules {
		st, err := sch.StartTimeIn(event.Timezone)
		if err != nil {
			return nil, err
		}

		s := &v1.Schedule{
			Id:            sch.ID,
			StartTime:     st.Format(time.RFC3339),
			EndTime:       sch.EndTimeFrom(st).Format(time.RFC3339),
			IsFullDay:     sch.IsFullDay,
			RecurringType: mapRecurringTypeToPB(sch.RecurringType),
		}
		schedules[index] = s
	}
	e.Schedule = schedules

	attendees := make([]string, len(event.Invitations))
	for index, inv := range event.Invitations {
		attendees[index] = inv.UserID
	}
	e.Attendees = attendees

	return e, nil
}

func mapRecurringType(rt v1.RecurringType) core.RecurringType {
	switch rt {
	case v1.RecurringType_DAILY:
		return core.RecurringType_Daily
	case v1.RecurringType_EVERY_WEEK:
		return core.RecurringType_Every_Week
	default:
		return core.RecurringType_None
	}
}

func mapRecurringTypeToPB(rt core.RecurringType) v1.RecurringType {
	switch rt {
	case core.RecurringType_Daily:
		return v1.RecurringType_DAILY
	case core.RecurringType_Every_Week:
		return v1.RecurringType_EVERY_WEEK
	default:
		return v1.RecurringType_NONE
	}
}
