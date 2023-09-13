//nolint:lll
//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./../../../mock/auth/$GOPACKAGE/$GOFILE
package database

import (
	"context"
	"errors"

	"github.com/and-period/furumane/internal/auth/entity"
)

var (
	ErrInvalidArgument    = errors.New("database: invalid argument")
	ErrNotFound           = errors.New("database: not found")
	ErrAlreadyExists      = errors.New("database: already exists")
	ErrFailedPrecondition = errors.New("database: failed precondition")
	ErrCanceled           = errors.New("database: canceled")
	ErrDeadlineExceeded   = errors.New("database: deadline exceeded")
	ErrInternal           = errors.New("database: internal error")
	ErrUnknown            = errors.New("database: unknown")
)

type Database struct {
	Admin Admin
}

type Admin interface {
	Get(ctx context.Context, adminID string, fields ...string) (*entity.Admin, error)
	GetByCognitoID(ctx context.Context, cognitoID string, fields ...string) (*entity.Admin, error)
	GetByEmail(ctx context.Context, email string, fields ...string) (*entity.Admin, error)
	Create(ctx context.Context, admin *entity.Admin, auth func(context.Context) error) error
	UpdateEmail(ctx context.Context, adminID, email string) error
	UpdateVerifiedAt(ctx context.Context, adminID string) error
	Delete(ctx context.Context, adminID string, auth func(context.Context) error) error
}
