package postgresql_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/dzakaammar/event-scheduling-example/internal"
	"github.com/dzakaammar/event-scheduling-example/internal/postgresql"
	"github.com/jmoiron/sqlx"
)

func TestNewEventRepository(t *testing.T) {
	type args struct {
		db *sqlx.DB
	}
	tests := []struct {
		name string
		args args
		want *postgresql.EventRepository
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := postgresql.NewEventRepository(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEventRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventRepository_Store(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx   context.Context
		event *internal.Event
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
			e := postgresql.NewEventRepository(tt.fields.db)
			if err := e.Store(tt.args.ctx, tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("EventRepository.Store() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEventRepository_DeleteByID(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx context.Context
		id  string
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
			e := postgresql.NewEventRepository(tt.fields.db)
			if err := e.DeleteByID(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("EventRepository.DeleteByID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEventRepository_Update(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx   context.Context
		event *internal.Event
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
			e := postgresql.NewEventRepository(tt.fields.db)
			if err := e.Update(tt.args.ctx, tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("EventRepository.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEventRepository_FindByID(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx context.Context
		id  string
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
			e := postgresql.NewEventRepository(tt.fields.db)
			got, err := e.FindByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("EventRepository.FindByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EventRepository.FindByID() = %v, want %v", got, tt.want)
			}
		})
	}
}
