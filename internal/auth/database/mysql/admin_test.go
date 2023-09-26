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

func TestAdmin_GetByCognitoID(t *testing.T) {
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
		cognitoID string
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
				cognitoID: "cognito-id",
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
				cognitoID: "",
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
			actual, err := db.GetByCognitoID(ctx, tt.args.cognitoID)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.admin, actual)
		})
	}
}

func TestAdmin_GetByEmail(t *testing.T) {
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
		email string
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
				email: "test@example.com",
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
				email: "",
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
			actual, err := db.GetByEmail(ctx, tt.args.email)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.admin, actual)
		})
	}
}

func TestAdmin_Create(t *testing.T) {
	db := dbClient
	now := func() time.Time {
		return current
	}

	a := fakeAdmin("admin-id", "cognito-id", "test@example.com", now())

	type args struct {
		admin *entity.Admin
		fn    func(ctx context.Context) error
	}
	type want struct {
		err error
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
				admin: a,
				fn: func(ctx context.Context) error {
					return nil
				},
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "already exists",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				a := fakeAdmin("admin-id", "cognito-id", "test@example.com", now())
				err := db.DB.WithContext(ctx).Create(&a).Error
				require.NoError(t, err)
			},
			args: args{
				admin: a,
				fn:    nil,
			},
			want: want{
				err: database.ErrAlreadyExists,
			},
		},
		{
			name:  "failed to callback",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				admin: a,
				fn: func(ctx context.Context) error {
					return assert.AnError
				},
			},
			want: want{
				err: database.ErrUnknown,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			err := deleteAll(ctx)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &admin{db: db, now: now}
			err = db.Create(ctx, tt.args.admin, tt.args.fn)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestAdmin_UpdateEmail(t *testing.T) {
	db := dbClient
	now := func() time.Time {
		return current
	}

	type args struct {
		adminID string
		email   string
	}
	type want struct {
		err error
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name: "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				admin := fakeAdmin("admin-id", "cognito-id", "test@example.com", now())
				err := db.DB.WithContext(ctx).Create(&admin).Error
				require.NoError(t, err)
			},
			args: args{
				adminID: "admin-id",
				email:   "test@example.com",
			},
			want: want{
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			err := deleteAll(ctx)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &admin{db: db, now: now}
			err = db.UpdateEmail(ctx, tt.args.adminID, tt.args.email)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestAdmin_UpdateVerifiedAt(t *testing.T) {
	db := dbClient
	now := func() time.Time {
		return current
	}

	type args struct {
		adminID string
	}
	type want struct {
		err error
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name: "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				admin := fakeAdmin("admin-id", "cognito-id", "test@example.com", now())
				err := db.DB.WithContext(ctx).Create(&admin).Error
				require.NoError(t, err)
			},
			args: args{
				adminID: "admin-id",
			},
			want: want{
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			err := deleteAll(ctx)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &admin{db: db, now: now}
			err = db.UpdateVerifiedAt(ctx, tt.args.adminID)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestAdmin_Delete(t *testing.T) {
	db := dbClient
	now := func() time.Time {
		return current
	}

	type args struct {
		adminID string
		fn      func(ctx context.Context) error
	}
	type want struct {
		err error
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name: "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				admin := fakeAdmin("admin-id", "cognito-id", "test@example.com", now())
				err := db.DB.WithContext(ctx).Create(&admin).Error
				require.NoError(t, err)
			},
			args: args{
				adminID: "admin-id",
				fn: func(ctx context.Context) error {
					return nil
				},
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "failed to callback",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				admin := fakeAdmin("admin-id", "cognito-id", "test@example.com", now())
				err := db.DB.WithContext(ctx).Create(&admin).Error
				require.NoError(t, err)
			},
			args: args{
				adminID: "admin-id",
				fn: func(ctx context.Context) error {
					return assert.AnError
				},
			},
			want: want{
				err: database.ErrUnknown,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			err := deleteAll(ctx)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &admin{db: db, now: now}
			err = db.Delete(ctx, tt.args.adminID, tt.args.fn)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func fakeAdmin(adminID, cognitoID, email string, now time.Time) *entity.Admin {
	return &entity.Admin{
		ID:           adminID,
		CognitoID:    cognitoID,
		ProviderType: entity.ProviderTypeEmail,
		Email:        email,
		PhoneNumber:  "09012341234",
		CreatedAt:    now,
		UpdatedAt:    now,
		VerifiedAt:   now,
	}
}
