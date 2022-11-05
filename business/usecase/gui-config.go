package usecase

import (
	"time"

	"golang.org/x/net/context"

	"github.com/Forest33/warthog/business/entity"
	"github.com/Forest33/warthog/pkg/logger"
)

const (
	updateGUIConfigDelay = 3
)

type GUIConfigUseCase struct {
	ctx           context.Context
	guiConfigRepo GUIConfigRepo
	log           *logger.Zerolog
	guiConfig     *entity.GUIConfig
	updateCh      chan *entity.GUIConfig
}

type GUIConfigRepo interface {
	Get() (*entity.GUIConfig, error)
	Update(in *entity.GUIConfig) (*entity.GUIConfig, error)
}

func NewGUIConfigUseCase(ctx context.Context, log *logger.Zerolog, guiConfigRepo GUIConfigRepo) *GUIConfigUseCase {
	uc := &GUIConfigUseCase{
		ctx:           ctx,
		guiConfigRepo: guiConfigRepo,
		log:           log,
		guiConfig:     &entity.GUIConfig{},
		updateCh:      make(chan *entity.GUIConfig, 10),
	}

	uc.updateHandler()

	return uc
}

func (uc *GUIConfigUseCase) Get() (*entity.GUIConfig, error) {
	cfg, err := uc.guiConfigRepo.Get()
	if err != nil {
		uc.log.Error().Msgf("failed to get GUI config: %v", err)
		return nil, err
	}
	return cfg, nil
}

func (uc *GUIConfigUseCase) Set(cfg *entity.GUIConfig) {
	uc.updateCh <- cfg
}

func (uc *GUIConfigUseCase) updateHandler() {
	go func() {
		for {
			select {
			case <-uc.ctx.Done():
				uc.updateGUIConfig()
				return
			case cfg := <-uc.updateCh:
				uc.setGUIConfig(cfg)
			case <-time.After(time.Second * updateGUIConfigDelay):
				uc.updateGUIConfig()
			}
		}
	}()
}

func (uc *GUIConfigUseCase) setGUIConfig(cfg *entity.GUIConfig) {
	if cfg.WindowWidth > 0 {
		uc.guiConfig.WindowWidth = cfg.WindowWidth
	}
	if cfg.WindowHeight > 0 {
		uc.guiConfig.WindowHeight = cfg.WindowHeight
	}
	if cfg.WindowX != nil {
		uc.guiConfig.WindowX = cfg.WindowX
	}
	if cfg.WindowY != nil {
		uc.guiConfig.WindowY = cfg.WindowY
	}
}

func (uc *GUIConfigUseCase) updateGUIConfig() {
	if _, err := uc.guiConfigRepo.Update(uc.guiConfig); err != nil {
		uc.log.Error().
			Interface("data", uc.guiConfig).
			Msgf("failed to update GUI config: %v", err)
	}
}
