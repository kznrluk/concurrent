package entity

import (
	"context"
	"github.com/totegamma/concurrent/x/core"
	"gorm.io/gorm"
	"time"
)

// Repository is the interface for host repository
type Repository interface {
	GetEntity(ctx context.Context, key string) (core.Entity, error)
	CreateEntity(ctx context.Context, entity *core.Entity, meta *core.EntityMeta) error
	UpdateEntity(ctx context.Context, entity *core.Entity) error
	GetList(ctx context.Context) ([]core.Entity, error)
	ListModified(ctx context.Context, modified time.Time) ([]core.Entity, error)
	Delete(ctx context.Context, key string) error
	Ack(ctx context.Context, ack *core.Ack) error
	Unack(ctx context.Context, ack *core.Ack) error
	Total(ctx context.Context) (int64, error)
	GetAcker(ctx context.Context, key string) ([]core.Ack, error)
	GetAcking(ctx context.Context, key string) ([]core.Ack, error)
	GetAddress(ctx context.Context, ccid string) (core.Address, error)
	UpdateAddress(ctx context.Context, ccid string, domain string, signedAt time.Time) error
}

type repository struct {
	db *gorm.DB
}

// Total returns the total number of entities
func (r *repository) Total(ctx context.Context) (int64, error) {
	ctx, span := tracer.Start(ctx, "RepositoryTotal")
	defer span.End()

	var count int64
	err := r.db.WithContext(ctx).Model(&core.Entity{}).Where("domain IS NULL or domain = ''").Count(&count).Error
	return count, err
}

// NewRepository creates a new host repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// GetAddress returns the address of a entity
func (r *repository) GetAddress(ctx context.Context, ccid string) (core.Address, error) {
	ctx, span := tracer.Start(ctx, "RepositoryGetAddress")
	defer span.End()

	var addr core.Address
	err := r.db.WithContext(ctx).First(&addr, "id = ?", ccid).Error
	return addr, err
}

// SetAddress sets the address of a entity
func (r *repository) SetAddress(ctx context.Context, ccid string, address string) error {
	ctx, span := tracer.Start(ctx, "RepositorySetAddress")
	defer span.End()

	return r.db.WithContext(ctx).Model(&core.Entity{}).Where("id = ?", ccid).Update("address", address).Error
}

// UpdateAddress updates the address of a entity
func (r *repository) UpdateAddress(ctx context.Context, ccid string, domain string, signedAt time.Time) error {
	ctx, span := tracer.Start(ctx, "RepositoryUpdateAddress")
	defer span.End()

	// create if not exists
	var addr core.Address
	err := r.db.WithContext(ctx).First(&addr, "id = ?", ccid).Error
	if err != nil {
		return r.db.WithContext(ctx).Create(&core.Address{
			ID:       ccid,
			Domain:   domain,
			SignedAt: signedAt,
		}).Error
	}

	return r.db.WithContext(ctx).Model(&core.Address{}).Where("id = ?", ccid).Update("domain", domain).Error
}

// Get returns a entity by key
func (r *repository) GetEntity(ctx context.Context, key string) (core.Entity, error) {
	ctx, span := tracer.Start(ctx, "RepositoryGet")
	defer span.End()

	var entity core.Entity
	err := r.db.WithContext(ctx).First(&entity, "id = ?", key).Error
	return entity, err
}

// Create creates new entity
func (r *repository) CreateEntity(ctx context.Context, entity *core.Entity, meta *core.EntityMeta) error {
	ctx, span := tracer.Start(ctx, "RepositoryCreate")
	defer span.End()

	err := r.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.WithContext(ctx).Create(&entity).Error; err != nil {
			return err
		}

		if err := tx.WithContext(ctx).Create(&meta).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}

// GetList returns all entities
func (r *repository) GetList(ctx context.Context) ([]core.Entity, error) {
	ctx, span := tracer.Start(ctx, "RepositoryGetList")
	defer span.End()

	var entities []core.Entity
	err := r.db.WithContext(ctx).Model(&core.Entity{}).Find(&entities).Error
	return entities, err
}

// ListModified returns all entities which modified after given time
func (r *repository) ListModified(ctx context.Context, time time.Time) ([]core.Entity, error) {
	ctx, span := tracer.Start(ctx, "RepositoryListModified")
	defer span.End()

	var entities []core.Entity
	err := r.db.WithContext(ctx).Model(&core.Entity{}).Where("m_date > ?", time).Find(&entities).Error
	return entities, err
}

// Delete deletes a entity
func (r *repository) Delete(ctx context.Context, id string) error {
	ctx, span := tracer.Start(ctx, "RepositoryDelete")
	defer span.End()

	return r.db.WithContext(ctx).Delete(&core.Entity{}, "id = ?", id).Error
}

// Update updates a entity
func (r *repository) UpdateEntity(ctx context.Context, entity *core.Entity) error {
	ctx, span := tracer.Start(ctx, "RepositoryUpdate")
	defer span.End()

	return r.db.WithContext(ctx).Where("id = ?", entity.ID).Updates(&entity).Error
}

// Ack creates a new ack
func (r *repository) Ack(ctx context.Context, ack *core.Ack) error {
	ctx, span := tracer.Start(ctx, "RepositoryAck")
	defer span.End()

	ack.Valid = true

	return r.db.WithContext(ctx).Save(&ack).Error
}

// Unack deletes a ack
func (r *repository) Unack(ctx context.Context, ack *core.Ack) error {
	ctx, span := tracer.Start(ctx, "RepositoryUnack")
	defer span.End()

	ack.Valid = false

	return r.db.WithContext(ctx).Save(&ack).Error
}

// GetAcker returns all acks for a entity
func (r *repository) GetAcker(ctx context.Context, key string) ([]core.Ack, error) {
	ctx, span := tracer.Start(ctx, "RepositoryGetAcker")
	defer span.End()

	var acks []core.Ack
	err := r.db.WithContext(ctx).Where("valid = true and \"to\" = ?", key).Find(&acks).Error
	return acks, err
}

// GetAcking returns all acks for a entity
func (r *repository) GetAcking(ctx context.Context, key string) ([]core.Ack, error) {
	ctx, span := tracer.Start(ctx, "RepositoryGetAcking")
	defer span.End()

	var acks []core.Ack
	err := r.db.WithContext(ctx).Where("valid = true \"from\" = ?", key).Find(&acks).Error
	return acks, err
}
