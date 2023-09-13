package mysql

import (
	"context"
	"time"

	"github.com/and-period/furumane/internal/auth/database"
	"github.com/and-period/furumane/internal/auth/entity"
	"github.com/and-period/furumane/pkg/jst"
	"github.com/and-period/furumane/pkg/mysql"
	"gorm.io/gorm"
)

const adminTable = "admins"

type admin struct {
	db  *mysql.Client
	now func() time.Time
}

func newAdmin(db *mysql.Client) database.Admin {
	return &admin{
		db:  db,
		now: jst.Now,
	}
}

func (a *admin) Get(ctx context.Context, adminID string, fields ...string) (*entity.Admin, error) {
	var admin *entity.Admin

	stmt := a.db.
		Statement(ctx, a.db.DB, adminTable, fields...).
		Where("id = ?", adminID)

	if err := stmt.First(&admin).Error; err != nil {
		return nil, dbError(err)
	}
	return admin, nil
}

func (a *admin) GetByCognitoID(ctx context.Context, cognitoID string, fields ...string) (*entity.Admin, error) {
	var admin *entity.Admin

	stmt := a.db.
		Statement(ctx, a.db.DB, adminTable, fields...).
		Where("cognito_id = ?", cognitoID)

	if err := stmt.First(&admin).Error; err != nil {
		return nil, dbError(err)
	}
	return admin, nil
}

func (a *admin) GetByEmail(ctx context.Context, email string, fields ...string) (*entity.Admin, error) {
	var admin *entity.Admin

	stmt := a.db.
		Statement(ctx, a.db.DB, adminTable, fields...).
		Where("email = ?", email)

	if err := stmt.First(&admin).Error; err != nil {
		return nil, dbError(err)
	}
	return admin, nil
}

func (a *admin) Create(ctx context.Context, admin *entity.Admin, auth func(context.Context) error) error {
	err := a.db.Transaction(ctx, func(tx *gorm.DB) error {
		now := a.now()
		admin.CreatedAt, admin.UpdatedAt = now, now

		if err := tx.WithContext(ctx).Create(&admin).Error; err != nil {
			return err
		}
		return auth(ctx)
	})
	return dbError(err)
}

func (a *admin) UpdateEmail(ctx context.Context, adminID, email string) error {
	updates := map[string]interface{}{
		"email":      email,
		"updated_at": a.now(),
	}
	stmt := a.db.DB.WithContext(ctx).
		Table(adminTable).
		Where("id = ?", adminID)

	err := stmt.Updates(updates).Error
	return dbError(err)
}

func (a *admin) UpdateVerifiedAt(ctx context.Context, adminID string) error {
	now := a.now()
	updates := map[string]interface{}{
		"verified_at": now,
		"updated_at":  now,
	}
	stmt := a.db.DB.WithContext(ctx).
		Table(adminTable).
		Where("id = ?", adminID)

	err := stmt.Updates(updates).Error
	return dbError(err)
}

func (a *admin) Delete(ctx context.Context, adminID string, auth func(context.Context) error) error {
	err := a.db.Transaction(ctx, func(tx *gorm.DB) error {
		now := a.now()
		updates := map[string]interface{}{
			"exists":     nil,
			"updated_at": now,
			"deleted_at": now,
		}
		stmt := tx.WithContext(ctx).
			Table(adminTable).
			Where("id = ?", adminID)

		if err := stmt.Updates(updates).Error; err != nil {
			return err
		}
		return auth(ctx)
	})
	return dbError(err)
}
