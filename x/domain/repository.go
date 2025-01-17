package domain

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/totegamma/concurrent/core"
)

// Repository is the interface for host repository
type Repository interface {
	GetByFQDN(ctx context.Context, key string) (core.Domain, error)
	GetByCCID(ctx context.Context, ccid string) (core.Domain, error)
	GetByCSID(ctx context.Context, ccid string) (core.Domain, error)
	Upsert(ctx context.Context, host core.Domain) (core.Domain, error)
	GetList(ctx context.Context) ([]core.Domain, error)
	Delete(ctx context.Context, id string) error
	UpdateScrapeTime(ctx context.Context, id string, scrapeTime time.Time) error
	Update(ctx context.Context, host core.Domain) error
}

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new host repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// GetByFQDN returns a host by FQDN
func (r *repository) GetByFQDN(ctx context.Context, key string) (core.Domain, error) {
	ctx, span := tracer.Start(ctx, "Domain.Repository.GetByFQDN")
	defer span.End()

	var host core.Domain
	err := r.db.WithContext(ctx).First(&host, "id = ?", key).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return core.Domain{}, core.NewErrorNotFound()
		}
		span.RecordError(err)
		return host, err
	}

	return host, nil
}

// GetByCCID returns a host by CCID
func (r *repository) GetByCCID(ctx context.Context, ccid string) (core.Domain, error) {
	ctx, span := tracer.Start(ctx, "Domain.Repository.GetByCCID")
	defer span.End()

	var host core.Domain
	err := r.db.WithContext(ctx).First(&host, "cc_id = ?", ccid).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return core.Domain{}, core.NewErrorNotFound()
		}
		span.RecordError(err)
		return host, err
	}

	return host, nil
}

func (r *repository) GetByCSID(ctx context.Context, csid string) (core.Domain, error) {
	ctx, span := tracer.Start(ctx, "Domain.Repository.GetByCSID")
	defer span.End()

	var host core.Domain
	err := r.db.WithContext(ctx).First(&host, "cs_id = ?", csid).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return core.Domain{}, core.NewErrorNotFound()
		}
		span.RecordError(err)
		return host, err
	}

	return host, nil
}

// Upsert creates new host
func (r *repository) Upsert(ctx context.Context, host core.Domain) (core.Domain, error) {
	ctx, span := tracer.Start(ctx, "Domain.Repository.Upsert")
	defer span.End()

	err := r.db.WithContext(ctx).Save(&host).Error

	return host, err
}

// GetList returns list of schemas by schema
func (r *repository) GetList(ctx context.Context) ([]core.Domain, error) {
	ctx, span := tracer.Start(ctx, "Domain.Repository.GetList")
	defer span.End()

	var hosts []core.Domain
	err := r.db.WithContext(ctx).Find(&hosts).Error
	return hosts, err
}

// Delete deletes a host
func (r *repository) Delete(ctx context.Context, id string) error {
	ctx, span := tracer.Start(ctx, "Domain.Repository.Delete")
	defer span.End()

	return r.db.WithContext(ctx).Delete(&core.Domain{}, "id = ?", id).Error
}

// UpdateScrapeTime updates scrape time
func (r *repository) UpdateScrapeTime(ctx context.Context, id string, scrapeTime time.Time) error {
	ctx, span := tracer.Start(ctx, "Domain.Repository.UpdateScrapeTime")
	defer span.End()

	return r.db.WithContext(ctx).Model(&core.Domain{}).Where("id = ?", id).Update("last_scraped", scrapeTime).Error
}

// Update updates a host
func (r *repository) Update(ctx context.Context, host core.Domain) error {
	ctx, span := tracer.Start(ctx, "Domain.Repository.Update")
	defer span.End()

	return r.db.WithContext(ctx).Model(&core.Domain{}).Where("id = ?", host.ID).Updates(&host).Error
}
