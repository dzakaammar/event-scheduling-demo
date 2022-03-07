package endpoint

import (
	"context"
	"time"

	v1 "github.com/dzakaammar/event-scheduling-example/gen/go/proto/v1"
	"github.com/dzakaammar/event-scheduling-example/internal"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
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

func (g *GRPCEndpoint) DeleteEventByID(ctx context.Context, req *v1.DeleteEventByIDRequest) (*v1.DeleteEventByIDResponse, error) {
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
	return &v1.DeleteEventByIDResponse{}, nil
}

func (g *GRPCEndpoint) UpdateEvent(ctx context.Context, req *v1.UpdateEventRequest) (*v1.UpdateEventResponse, error) {
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
	return &v1.UpdateEventResponse{}, nil
}

func (g *GRPCEndpoint) FindEventByID(ctx context.Context, req *v1.FindEventByIDRequest) (*v1.FindEventByIDResponse, error) {
	event, err := g.svc.FindEventByID(ctx, &internal.FindEventByIDRequest{EventID: req.GetId()})
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

func parseCreateEventRequest(ctx context.Context, req *v1.CreateEventRequest) (*internal.CreateEventRequest, error) {
	if req == nil || req.GetEvent() == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	actorID := extractAuthorization(ctx)

	event := internal.NewEvent()
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

	return &internal.CreateEventRequest{
		ActorID: actorID,
		Event:   event,
	}, nil
}

func parseDeleteEventByIDtRequest(ctx context.Context, req *v1.DeleteEventByIDRequest) (*internal.DeleteEventByIDRequest, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	return &internal.DeleteEventByIDRequest{
		ActorID: extractAuthorization(ctx),
		EventID: req.GetId(),
	}, nil
}

func parseUpdateEventByIDRequest(ctx context.Context, req *v1.UpdateEventRequest) (*internal.UpdateEventRequest, error) {
	if req == nil || req.GetEvent() == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	now := time.Now()
	event := internal.Event{
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

	return &internal.UpdateEventRequest{
		ID:      req.GetId(),
		ActorID: extractAuthorization(ctx),
		Event:   &event,
	}, nil
}

func parseSchedules(sch []*v1.Schedule, eventID string) ([]internal.Schedule, error) {
	var schedules []internal.Schedule
	for _, sch := range sch {
		s, err := internal.NewSchedule(
			eventID,
			sch.GetStartTime(),
			sch.GetEndTime(),
			sch.GetIsFullDay(),
			mapRecurringType(sch.GetRecurringType()),
		)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		schedules = append(schedules, s)
	}
	return schedules, nil
}

func parseInvitations(attendees []string, eventID string) []internal.Invitation {
	var invitations []internal.Invitation
	for _, userID := range attendees {
		i := internal.NewInvitation(eventID, userID)
		invitations = append(invitations, i)
	}
	return invitations
}

func parseEventToPB(event *internal.Event) (*v1.Event, error) {
	e := &v1.Event{
		Id:            event.ID,
		Title:         event.Title,
		Description:   event.Description,
		Timezone:      event.Timezone,
		CreatedAt:     event.CreatedAt.Format(time.RFC3339),
		CreatedBy:     event.CreatedBy,
		LastUpdatedAt: event.GetUpdatedAt(),
	}

	var schedules []*v1.Schedule
	for _, sch := range event.Schedules {
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
		schedules = append(schedules, s)
	}
	e.Schedule = schedules

	var attendees []string
	for _, inv := range event.Invitations {
		attendees = append(attendees, inv.UserID)
	}
	e.Attendees = attendees

	return e, nil
}

func mapRecurringType(rt v1.RecurringType) internal.RecurringType {
	switch rt {
	case v1.RecurringType_DAILY:
		return internal.RecurringType_Daily
	case v1.RecurringType_EVERY_WEEK:
		return internal.RecurringType_Every_Week
	default:
		return internal.RecurringType_None
	}
}

func mapRecurringTypeToPB(rt internal.RecurringType) v1.RecurringType {
	switch rt {
	case internal.RecurringType_Daily:
		return v1.RecurringType_DAILY
	case internal.RecurringType_Every_Week:
		return v1.RecurringType_EVERY_WEEK
	default:
		return v1.RecurringType_NONE
	}
}
