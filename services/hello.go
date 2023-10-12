package services

import (
	"gitlab.com/idolauncher/go-template-kit/config"
	log "go.uber.org/zap"
	"gorm.io/gorm"
)

type HelloService struct {
	logger *log.Logger
	cfg    *config.Config
	db     *gorm.DB
}

func NewHelloService(logger *log.Logger, cfg *config.Config, db *gorm.DB) *HelloService {
	return &HelloService{logger: logger, cfg: cfg, db: db}
}

func (e *HelloService) HelloWorld() (interface{}, error) {
	return nil, nil
}
