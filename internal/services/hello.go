package services

import (
	"github.com/0xhoang/go-kit/config"
	log "go.uber.org/zap"
	"gorm.io/gorm"
	"time"
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

func (d *HelloService) CreateCompetition() error {
	time.Sleep(7 * time.Second)

	d.logger.Info("Create Competition Successfully")

	return nil
}

func (d *HelloService) CreateTeam() error {
	time.Sleep(7 * time.Second)
	d.logger.Info("Create Team Successfully")
	return nil
}
