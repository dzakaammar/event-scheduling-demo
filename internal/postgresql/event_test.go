package postgresql_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dzakaammar/event-scheduling-example/internal"
	"github.com/dzakaammar/event-scheduling-example/internal/postgresql"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestNewEventRepository(t *testing.T) {
	db, _, _ := sqlmock.New()
	conn, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}))
	type args struct {
		db *gorm.DB
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "OK",
			args: args{
				db: conn,
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
		dbMock func(t *testing.T) *gorm.DB
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
		{
			name: "OK",
			fields: fields{
				dbMock: func(t *testing.T) *gorm.DB {
					db, mock, _ := sqlmock.New()
					g, _ := gorm.Open(postgres.New(postgres.Config{
						Conn: db,
					}))

					mock.ExpectBegin()
					mock.ExpectExec(`INSERT INTO "event"`).WillReturnResult(sqlmock.NewResult(1, 1))
					mock.ExpectCommit()
					mock.MatchExpectationsInOrder(true)
					return g
				},
			},
			args: args{
				ctx: context.Background(),
				event: &internal.Event{
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
				dbMock: func(t *testing.T) *gorm.DB {
					db, mock, _ := sqlmock.New()
					g, _ := gorm.Open(postgres.New(postgres.Config{
						Conn: db,
					}))

					mock.ExpectBegin()
					mock.ExpectExec(`INSERT INTO "event"`).WillReturnError(errors.New("error"))
					mock.ExpectRollback()
					mock.MatchExpectationsInOrder(true)
					return g
				},
			},
			args: args{
				ctx: context.Background(),
				event: &internal.Event{
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
		dbMock func(t *testing.T) *gorm.DB
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
				dbMock: func(t *testing.T) *gorm.DB {
					db, mock, _ := sqlmock.New()
					g, _ := gorm.Open(postgres.New(postgres.Config{
						Conn: db,
					}))

					mock.ExpectBegin()
					mock.ExpectExec(`DELETE FROM "event"`).WillReturnResult(sqlmock.NewResult(1, 1))
					mock.ExpectCommit()
					mock.MatchExpectationsInOrder(true)

					return g
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
				dbMock: func(t *testing.T) *gorm.DB {
					db, mock, _ := sqlmock.New()
					g, _ := gorm.Open(postgres.New(postgres.Config{
						Conn: db,
					}))

					mock.ExpectBegin()
					mock.ExpectExec(`DELETE FROM "event"`).WillReturnError(errors.New("error"))
					mock.ExpectRollback()
					mock.MatchExpectationsInOrder(true)

					return g
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
		dbMock func(t *testing.T) *gorm.DB
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
		{
			name: "OK",
			fields: fields{
				dbMock: func(t *testing.T) *gorm.DB {
					db, mock, _ := sqlmock.New()
					g, _ := gorm.Open(postgres.New(postgres.Config{
						Conn: db,
					}))

					mock.ExpectBegin()
					mock.ExpectExec(`UPDATE "event" SET`).WillReturnResult(sqlmock.NewResult(1, 1))
					mock.ExpectExec(`INSERT INTO "schedule"`).WillReturnResult(sqlmock.NewResult(1, 1))
					mock.ExpectExec(`INSERT INTO "invitation"`).WillReturnResult(sqlmock.NewResult(1, 1))
					mock.ExpectCommit()
					mock.MatchExpectationsInOrder(true)

					return g
				},
			},
			args: args{
				ctx: context.Background(),
				event: &internal.Event{
					ID: "123",
					Schedules: []internal.Schedule{
						{
							ID:                "123",
							EventID:           "123",
							StartTime:         time.Now().Unix(),
							Duration:          120,
							IsFullDay:         false,
							RecurringType:     internal.RecurringType_None,
							RecurringInterval: 0,
						},
					},
					Invitations: []internal.Invitation{
						{
							ID:      "123",
							EventID: "123",
							UserID:  "123",
							Status:  internal.InvitationStatus_Unknown,
							Token:   "123",
						},
					},
				},
			},
		},
		{
			name: "Not OK - error",
			fields: fields{
				dbMock: func(t *testing.T) *gorm.DB {
					db, mock, _ := sqlmock.New()
					g, _ := gorm.Open(postgres.New(postgres.Config{
						Conn: db,
					}))

					mock.ExpectBegin()
					mock.ExpectExec(`UPDATE "event" SET`).WillReturnResult(sqlmock.NewResult(1, 1))
					mock.ExpectExec(`INSERT INTO "schedule"`).WillReturnResult(sqlmock.NewResult(1, 1))
					mock.ExpectExec(`INSERT INTO "invitation"`).WillReturnError(errors.New("error"))
					mock.ExpectRollback()
					mock.MatchExpectationsInOrder(true)

					return g
				},
			},
			args: args{
				ctx: context.Background(),
				event: &internal.Event{
					ID: "123",
					Schedules: []internal.Schedule{
						{
							ID:                "123",
							EventID:           "123",
							StartTime:         time.Now().Unix(),
							Duration:          120,
							IsFullDay:         false,
							RecurringType:     internal.RecurringType_None,
							RecurringInterval: 0,
						},
					},
					Invitations: []internal.Invitation{
						{
							ID:      "123",
							EventID: "123",
							UserID:  "123",
							Status:  internal.InvitationStatus_Unknown,
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
		dbMock func(t *testing.T) *gorm.DB
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
		{
			name: "OK",
			fields: fields{
				dbMock: func(t *testing.T) *gorm.DB {
					db, mock, _ := sqlmock.New()
					g, _ := gorm.Open(postgres.New(postgres.Config{
						Conn: db,
					}))

					mock.ExpectQuery(`^SELECT .+ FROM "event"`).WithArgs("123").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("123"))
					mock.ExpectQuery(`^SELECT .+ FROM "invitation"`).WithArgs("123").WillReturnRows(sqlmock.NewRows([]string{"id"}))
					mock.ExpectQuery(`^SELECT .+ FROM "schedule"`).WithArgs("123").WillReturnRows(sqlmock.NewRows([]string{"id"}))
					mock.MatchExpectationsInOrder(true)

					return g
				},
			},
			args: args{
				ctx: context.Background(),
				id:  "123",
			},
			want: &internal.Event{
				ID:          "123",
				Invitations: []internal.Invitation{},
				Schedules:   []internal.Schedule{},
			},
		},
		{
			name: "Not OK - error",
			fields: fields{
				dbMock: func(t *testing.T) *gorm.DB {
					db, mock, _ := sqlmock.New()
					g, _ := gorm.Open(postgres.New(postgres.Config{
						Conn: db,
					}))

					mock.ExpectQuery(`^SELECT .+ FROM "event"`).WithArgs("123").WillReturnError(errors.New("error"))
					mock.MatchExpectationsInOrder(true)

					return g
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
