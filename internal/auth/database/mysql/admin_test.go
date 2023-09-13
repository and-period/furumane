package mysql

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumane/internal/auth/database"
	"github.com/and-period/furumane/internal/auth/entity"
	"github.com/and-period/furumane/pkg/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdmin(t *testing.T) {
	t.Parallel()
	assert.NotNil(t, newAdmin(nil))
}

func TestAdmin_Get(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db := dbClient
	now := func() time.Time {
		return current
	}

	err := deleteAll(ctx)
	require.NoError(t, err)

	a := fakeAdmin("admin-id", "cognito-id", "test@example.com", now())
	err = db.DB.WithContext(ctx).Create(&a).Error
	require.NoError(t, err)

	type args struct {
		adminID string
	}
	type want struct {
		admin *entity.Admin
		err   error
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				adminID: "admin-id",
			},
			want: want{
				admin: a,
				err:   nil,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				adminID: "",
			},
			want: want{
				admin: nil,
				err:   database.ErrNotFound,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			tt.setup(ctx, t, db)

			db := &admin{db: db, now: now}
			actual, err := db.Get(ctx, tt.args.adminID)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.admin, actual)
		})
	}
}

func fakeAdmin(adminID, cognitoID, email string, now time.Time) *entity.Admin {
	return &entity.Admin{
		ID:           adminID,
		CognitoID:    cognitoID,
		ProviderType: entity.ProviderTypeEmail,
		Email:        email,
		CreatedAt:    now,
		UpdatedAt:    now,
		VerifiedAt:   now,
	}
}
