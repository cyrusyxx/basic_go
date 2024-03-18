package dao

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	realmysql "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestGORMUserDAO_Insert(t *testing.T) {
	testCases := []struct {
		name string
		mock func(t *testing.T) *sql.DB

		ctx  context.Context
		user User

		wantErr error
	}{
		{
			name: "success",
			mock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectExec(".*").
					WillReturnResult(sqlmock.NewResult(1, 1))
				return db
			},
			ctx:     context.Background(),
			user:    User{},
			wantErr: nil,
		},
		{
			name: "email duplicate",
			mock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectExec(".*").
					WillReturnError(&realmysql.MySQLError{Number: 1062})
				return db
			},
			ctx:     context.Background(),
			user:    User{},
			wantErr: ErrDuplicateUser,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sqldb := tc.mock(t)
			db, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      sqldb,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				SkipDefaultTransaction: true,
				DisableAutomaticPing:   true,
			})
			assert.NoError(t, err)
			dao := NewGORMUserDAO(db)
			err = dao.Insert(tc.ctx, tc.user)

			assert.Equal(t, tc.wantErr, err)
		})
	}
}
