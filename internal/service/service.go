package service

import (
	"nunu_ginblog/internal/repository"
	"nunu_ginblog/pkg/helper/sid"
	"nunu_ginblog/pkg/jwt"
	"nunu_ginblog/pkg/log"
)

type Service struct {
	logger *log.Logger
	sid    *sid.Sid
	jwt    *jwt.JWT
	tm     repository.Transaction
}

func NewService(tm repository.Transaction, logger *log.Logger, sid *sid.Sid, jwt *jwt.JWT) *Service {
	return &Service{
		logger: logger,
		sid:    sid,
		jwt:    jwt,
		tm:     tm,
	}
}
