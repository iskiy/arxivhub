package service

import (
	"arxivhub/internal/repository"
)

type Config struct {
	UserRepository  repository.UserRepository
	PaperRepository repository.PaperRepository
}

type Service struct {
	Papers *PaperService
	Users  *UserService
}

//
//func NewService(conf *Config) *Service {
//	var service Service
//
//	service.Papers = NewPaperService(conf.PaperRepository)
//	service.Users = NewUserService(conf.UserRepository)
//
//	return &service
//}
