package postgresql_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dzakaammar/event-scheduling-example/internal/core"
	"github.com/dzakaammar/event-scheduling-example/internal/postgresql"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestNewEventRepository(t *testing.T) {
	db, _, _ := sqlmock.New()
	type args struct {
		db *sqlx.DB
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "OK",
			args: args{
				db: sqlx.NewDb(db, "postgres"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := postgresql.NewEventRepository(tt.args.db)
			assert.NotNil(t, got)
		})
	}
}

func TestEventRepository_Store(t *testing.T) {
	type fields struct {
		dbMock func(t *testing.T) *sqlx.DB
	}
	type args struct {
		ctx   context.Context
		event *core.Event
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
				dbMock: func(t *testing.T) *sqlx.DB {
					db, mock, _ := sqlmock.New()

					mock.ExpectBegin()
					mock.ExpectExec(`INSERT INTO event`).WillReturnResult(sqlmock.NewResult(1, 1))
					mock.ExpectCommit()
					mock.MatchExpectationsInOrder(true)
					return sqlx.NewDb(db, "postgres")
				},
			},
			args: args{
				ctx: context.Background(),
				event: &core.Event{
					ID:          "test",
					Title:       "test123",
					Description: "test123",
					Timezone:    "Asia/Jakarta",
				},
			},
		},
		{
			name: "Not OK - error",
			fields: fields{
				dbMock: func(t *testing.T) *sqlx.DB {
					db, mock, _ := sqlmock.New()
					mock.ExpectBegin()
					mock.ExpectExec(`INSERT INTO event`).WillReturnError(errors.New("error"))
					mock.ExpectRollback()
					mock.MatchExpectationsInOrder(true)
					return sqlx.NewDb(db, "postgres")
				},
			},
			args: args{
				ctx: context.Background(),
				event: &core.Event{
					ID:          "test",
					Title:       "test123",
					Description: "test123",
					Timezone:    "Asia/Jakarta",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := postgresql.NewEventRepository(tt.fields.dbMock(t))
			err := e.Store(tt.args.ctx, tt.args.event)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestEventRepository_DeleteByID(t *testing.T) {
	type fields struct {
		dbMock func(t *testing.T) *sqlx.DB
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
		{
			name: "OK",
			fields: fields{
				dbMock: func(t *testing.T) *sqlx.DB {
					db, mock, _ := sqlmock.New()
					mock.ExpectExec(`DELETE FROM event`).WithArgs("test123").WillReturnResult(sqlmock.NewResult(1, 1))
					mock.MatchExpectationsInOrder(true)

					return sqlx.NewDb(db, "postgres")
				},
			},
			args: args{
				ctx: context.Background(),
				id:  "test123",
			},
		},
		{
			name: "Not OK - error",
			fields: fields{
				dbMock: func(t *testing.T) *sqlx.DB {
					db, mock, _ := sqlmock.New()

					mock.ExpectExec(`DELETE FROM event`).WithArgs("test123").WillReturnError(errors.New("error"))
					mock.MatchExpectationsInOrder(true)

					return sqlx.NewDb(db, "postgres")
				},
			},
			args: args{
				ctx: context.Background(),
				id:  "test123",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := postgresql.NewEventRepository(tt.fields.dbMock(t))
			err := e.DeleteByID(tt.args.ctx, tt.args.id)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestEventRepository_Update(t *testing.T) {
	type fields struct {
		dbMock func(t *testing.T) *sqlx.DB
	}
	type args struct {
		ctx   context.Context
		event *core.Event
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
				dbMock: func(t *testing.T) *sqlx.DB {
					db, mock, _ := sqlmock.New()

					mock.ExpectBegin()
					mock.ExpectExec(`UPDATE event SET`).WillReturnResult(sqlmock.NewResult(1, 1))
					mock.ExpectExec(`INSERT INTO schedule`).WillReturnResult(sqlmock.NewResult(1, 1))
					mock.ExpectExec(`INSERT INTO invitation`).WillReturnResult(sqlmock.NewResult(1, 1))
					mock.ExpectCommit()
					mock.MatchExpectationsInOrder(true)

					return sqlx.NewDb(db, "postgres")
				},
			},
			args: args{
				ctx: context.Background(),
				event: &core.Event{
					ID: "123",
					Schedules: []core.Schedule{
						{
							ID:                "123",
							EventID:           "123",
							StartTime:         time.Now().Unix(),
							Duration:          120,
							IsFullDay:         false,
							RecurringType:     core.RecurringType_None,
							RecurringInterval: 0,
						},
					},
					Invitations: []core.Invitation{
						{
							ID:      "123",
							EventID: "123",
							UserID:  "123",
							Status:  core.InvitationStatus_Unknown,
							Token:   "123",
						},
					},
				},
			},
		},
		{
			name: "Not OK - error",
			fields: fields{
				dbMock: func(t *testing.T) *sqlx.DB {
					db, mock, _ := sqlmock.New()

					mock.ExpectBegin()
					mock.ExpectExec(`UPDATE event SET`).WillReturnResult(sqlmock.NewResult(1, 1))
					mock.ExpectExec(`INSERT INTO schedule`).WillReturnResult(sqlmock.NewResult(1, 1))
					mock.ExpectExec(`INSERT INTO invitation`).WillReturnError(errors.New("error"))
					mock.ExpectRollback()
					mock.MatchExpectationsInOrder(true)

					return sqlx.NewDb(db, "postgres")
				},
			},
			args: args{
				ctx: context.Background(),
				event: &core.Event{
					ID: "123",
					Schedules: []core.Schedule{
						{
							ID:                "123",
							EventID:           "123",
							StartTime:         time.Now().Unix(),
							Duration:          120,
							IsFullDay:         false,
							RecurringType:     core.RecurringType_None,
							RecurringInterval: 0,
						},
					},
					Invitations: []core.Invitation{
						{
							ID:      "123",
							EventID: "123",
							UserID:  "123",
							Status:  core.InvitationStatus_Unknown,
							Token:   "123",
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := postgresql.NewEventRepository(tt.fields.dbMock(t))
			err := e.Update(tt.args.ctx, tt.args.event)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestEventRepository_FindByID(t *testing.T) {
	type fields struct {
		dbMock func(t *testing.T) *sqlx.DB
	}
	type args struct {
		ctx context.Context
		id  string
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
				dbMock: func(t *testing.T) *sqlx.DB {
					db, mock, _ := sqlmock.New()

					mock.ExpectQuery(`^SELECT .+ FROM event`).WithArgs("123").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("123"))
					mock.ExpectQuery(`^SELECT .+ FROM schedule`).WithArgs("123").WillReturnRows(sqlmock.NewRows([]string{"id"}))
					mock.ExpectQuery(`^SELECT .+ FROM invitation`).WithArgs("123").WillReturnRows(sqlmock.NewRows([]string{"id"}))
					mock.MatchExpectationsInOrder(true)

					return sqlx.NewDb(db, "postgres")
				},
			},
			args: args{
				ctx: context.Background(),
				id:  "123",
			},
			want: &core.Event{
				ID:          "123",
				Invitations: []core.Invitation(nil),
				Schedules:   []core.Schedule(nil),
			},
		},
		{
			name: "Not OK - error",
			fields: fields{
				dbMock: func(t *testing.T) *sqlx.DB {
					db, mock, _ := sqlmock.New()

					mock.ExpectQuery(`^SELECT .+ FROM event`).WithArgs("123").WillReturnError(errors.New("error"))
					mock.MatchExpectationsInOrder(true)

					return sqlx.NewDb(db, "postgres")
				},
			},
			args: args{
				ctx: context.Background(),
				id:  "123",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := postgresql.NewEventRepository(tt.fields.dbMock(t))
			got, err := e.FindByID(tt.args.ctx, tt.args.id)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.EqualValues(t, tt.want, got)
		})
	}
}
