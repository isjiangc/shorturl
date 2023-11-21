package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"log"
	v1 "nunu_ginblog/api/v1"
	"nunu_ginblog/internal/model"
	"nunu_ginblog/internal/repository"
	"nunu_ginblog/pkg/helper/md5"
	"regexp"
	"time"
)

type ShorturlService interface {
	GetShorturlByUrl(ctx context.Context, url string) (*model.ShortUrl, error)
	GenerateShortUrl(ctx context.Context, req *v1.GenerateShortUrlRequest) (*v1.GenerateShortUrlData, error)
	Search4ShortUrl(ctx context.Context, shortUrl string) (model.MemShortUrl, error)
	ChangeState(ctx context.Context, shortUrl string, req *v1.UpdateShortUrlStateRequest) (*v1.UpdateShortUrlStateData, error)
	DeleteShortUrl(ctx context.Context, shortUrl string) (*v1.DeleteShortUrlStateData, error)
	GetShortUrlList(ctx context.Context, req *v1.GetShortUrlListRequest) ([]model.ShortUrl, error)
}

func NewShorturlService(service *Service, conf *viper.Viper, shorturlRepository repository.ShorturlRepository, rdb *redis.Client) ShorturlService {
	return &shorturlService{
		Service:            service,
		conf:               conf,
		shorturlRepository: shorturlRepository,
		rdb:                rdb,
	}
}

type shorturlService struct {
	conf *viper.Viper
	*Service
	shorturlRepository repository.ShorturlRepository
	rdb                *redis.Client
}

// GenerateShortUrl 生成短链接
func (s *shorturlService) GenerateShortUrl(ctx context.Context, req *v1.GenerateShortUrlRequest) (*v1.GenerateShortUrlData, error) {
	// 判断非法的链接
	regex := regexp.MustCompile(`^https?://.*$`)
	if !regex.MatchString(req.DestUrl) {
		s.Service.logger.Error("The dest_url is illegal")
		return nil, v1.ErrDestUrlIllegal
	}
	if req.OpenType < 0 || req.OpenType > 8 {
		s.Service.logger.Error("The opentype is illegal")
		return nil, v1.ErrOpenTypeIllegal
	}
	if md5.EmptyString(req.Memo) {
		s.Service.logger.Error("The memo is Empty")
		return nil, v1.ErrMemoIsEmpty
	}
	shortUrl, err := generateShortLink(req.DestUrl)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if shorturl, err := s.shorturlRepository.FirstByDesturl(ctx, shortUrl); err == nil && shorturl != nil {
		return nil, v1.ErrDestUrlAlreadyExist
	}
	url := model.ShortUrl{
		ShortUrl: shortUrl,
		DestUrl:  req.DestUrl,
		Valid:    true,
		Memo:     req.Memo,
		OpenType: model.OpenType(req.OpenType),
	}
	err = s.shorturlRepository.Create(ctx, &url)
	if err != nil {
		return nil, v1.ErrInternalServerError
	}

	mu := model.MemShortUrl{DestUrl: url.DestUrl, OpenType: url.OpenType}
	res, err := json.Marshal(mu)
	if err != nil {
		return nil, v1.ErrInternalServerError
	}
	if err := s.redisSet(shortUrl, res, redis.KeepTTL); err != nil {
		log.Println(err)
		return nil, v1.ErrInternalServerError
	}
	return &v1.GenerateShortUrlData{
		ShortUrl: fmt.Sprintf("%s%s", s.conf.GetString("urlprefix"), shortUrl),
	}, nil

}

// GetShorturlByUrl 查找短链接
func (s *shorturlService) GetShorturlByUrl(ctx context.Context, url string) (*model.ShortUrl, error) {
	return s.shorturlRepository.FirstByDesturl(ctx, url)
}

// Search4ShortUrl redis中查询
func (s *shorturlService) Search4ShortUrl(ctx context.Context, shortUrl string) (model.MemShortUrl, error) {
	mu := model.MemShortUrl{}
	found, err := s.redisGetString(shortUrl)
	if err != nil {
		log.Println(err)
		return mu, v1.ErrInternalServerError
	}
	if md5.EmptyString(found) {
		return mu, nil
	}
	return mu, json.Unmarshal([]byte(found), &mu)

}

// ChangeState 禁用/启用短链接
func (s *shorturlService) ChangeState(ctx context.Context, shortUrl string, req *v1.UpdateShortUrlStateRequest) (*v1.UpdateShortUrlStateData, error) {
	found, err := s.shorturlRepository.FirstByDesturl(ctx, shortUrl)
	if err != nil {
		return &v1.UpdateShortUrlStateData{
			Result: false,
		}, v1.ErrInternalServerError
	}
	if found.IsEmpty() {
		return &v1.UpdateShortUrlStateData{
			Result: false,
		}, v1.ErrDestUrlNotExist
	}
	found.Valid = req.Enable
	e := s.shorturlRepository.Update(ctx, found)
	if e != nil {
		return &v1.UpdateShortUrlStateData{
			Result: false,
		}, v1.ErrInternalServerError
	}
	if req.Enable {
		mu := model.MemShortUrl{
			DestUrl:  found.DestUrl,
			OpenType: found.OpenType,
		}
		res, err := json.Marshal(mu)
		if err != nil {
			return &v1.UpdateShortUrlStateData{
				Result: false,
			}, v1.ErrInternalServerError
		}
		s.redisSet(found.ShortUrl, res, redis.KeepTTL)
	} else {
		s.redisDelete(found.ShortUrl)
	}
	return &v1.UpdateShortUrlStateData{
		Result: true,
	}, nil
}

// DeleteShortUrl 删除短链
func (s *shorturlService) DeleteShortUrl(ctx context.Context, shortUrl string) (*v1.DeleteShortUrlStateData, error) {
	found, err := s.shorturlRepository.FirstByDesturl(ctx, shortUrl)
	if err != nil {
		return &v1.DeleteShortUrlStateData{
			Result: false,
		}, v1.ErrDestUrlNotExist
	}
	if found.IsEmpty() {
		return &v1.DeleteShortUrlStateData{
			Result: false,
		}, v1.ErrDestUrlNotExist
	}
	err = s.shorturlRepository.DeleteByurl(ctx, shortUrl)
	if err != nil {
		return &v1.DeleteShortUrlStateData{
			Result: false,
		}, v1.ErrInternalServerError
	}
	err = s.redisDelete(found.ShortUrl)
	if err != nil {
		return &v1.DeleteShortUrlStateData{
			Result: false,
		}, v1.ErrInternalServerError
	}
	return &v1.DeleteShortUrlStateData{
		Result: true,
	}, nil
}

func (s *shorturlService) GetShortUrlList(ctx context.Context, req *v1.GetShortUrlListRequest) ([]model.ShortUrl, error) {
	if req.PageNum < 1 || req.PageSize < 1 {
		return nil, nil
	}
	list, err := s.shorturlRepository.FindShortUrlList(ctx, req.PageNum, req.PageSize)
	if err != nil {
		s.Service.logger.Error("get shorturlList error")
		return list, v1.ErrInternalServerError
	}
	return list, nil
}

func generateShortLink(initialLink string) (string, error) {
	if md5.EmptyString(initialLink) {
		return "", fmt.Errorf("empty stirng")
	}
	urlHash, err := md5.Sha256Of(initialLink)
	if err != nil {
		return "", err
	}
	str := md5.Base58Encode(urlHash)
	return str[:8], nil
}

// redis设置
func (s *shorturlService) redisSet(key string, value interface{}, ttl time.Duration) error {
	return s.rdb.Set(context.Background(), key, value, ttl).Err()
}

// redis查询
func (s *shorturlService) redisGetString(key string) (string, error) {
	var result string
	var err error
	result, err = s.rdb.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return result, nil
	}
	return result, err
}

// redis删除
func (s *shorturlService) redisDelete(key ...string) error {
	return s.rdb.Del(context.Background(), key...).Err()
}
