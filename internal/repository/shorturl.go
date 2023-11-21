package repository

import (
	"context"
	"gorm.io/gorm"
	"nunu_ginblog/internal/model"
)

type ShorturlRepository interface {
	Create(ctx context.Context, shorturl *model.ShortUrl) error
	Update(ctx context.Context, shorturl *model.ShortUrl) error
	FirstByDesturl(ctx context.Context, url string) (*model.ShortUrl, error)
	DeleteByurl(ctx context.Context, url string) error
	FindShortUrlList(ctx context.Context, pageNum int, pageSize int) ([]model.ShortUrl, error)
}

func NewShorturlRepository(repository *Repository) ShorturlRepository {
	return &shorturlRepository{
		Repository: repository,
	}
}

type shorturlRepository struct {
	*Repository
}

// FindShortUrlList 查询所有的短链接
func (r *shorturlRepository) FindShortUrlList(ctx context.Context, pageNum int, pageSize int) ([]model.ShortUrl, error) {
	var (
		shortUrlList []model.ShortUrl
		err          error
	)
	if pageNum < 0 {
		return shortUrlList, nil
	}
	offset := (pageNum - 1) * pageSize
	if pageSize > 0 && pageNum > 0 {
		err = r.DB(ctx).Where("1=1").Order("id asc").Offset(offset).Limit(pageSize).Find(&shortUrlList).Error
	} else {
		err = r.DB(ctx).Find(&shortUrlList).Error
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return shortUrlList, nil
}

func (r *shorturlRepository) Create(ctx context.Context, shorturl *model.ShortUrl) error {
	if err := r.DB(ctx).Create(shorturl).Error; err != nil {
		return err
	}
	return nil
}

func (r *shorturlRepository) Update(ctx context.Context, shorturl *model.ShortUrl) error {
	if err := r.DB(ctx).Save(shorturl).Error; err != nil {
		return err
	}
	return nil
}

func (r *shorturlRepository) FirstByDesturl(ctx context.Context, url string) (*model.ShortUrl, error) {
	var shorturl model.ShortUrl
	if err := r.DB(ctx).Where("short_url =?", url).First(&shorturl).Error; err != nil {
		return nil, err
	}
	return &shorturl, nil
}

func (r *shorturlRepository) DeleteByurl(ctx context.Context, url string) error {
	if err := r.DB(ctx).Where("short_url =?", url).Delete(&model.ShortUrl{}).Error; err != nil {
		return err
	}
	return nil
}
