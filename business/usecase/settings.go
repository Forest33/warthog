// Package usecase provides business logic.
package usecase

import (
	"sync"
	"time"

	"golang.org/x/net/context"

	"github.com/forest33/warthog/business/entity"
	"github.com/forest33/warthog/pkg/logger"
)

const (
	updateSettingsDelay = 3
)

// SettingsUseCase object capable of interacting with SettingsUseCase
type SettingsUseCase struct {
	ctx          context.Context
	wg           *sync.WaitGroup
	settingsRepo SettingsRepo
	log          *logger.Zerolog
	guiSettings  *entity.Settings
	updateCh     chan *entity.Settings
	grpcClient   GrpcClient
}

// SettingsRepo is the common interface implemented SettingsRepository methods
type SettingsRepo interface {
	Get() (*entity.Settings, error)
	Update(in *entity.Settings) (*entity.Settings, error)
}

// NewSettingsUseCase creates a new SettingsUseCase
func NewSettingsUseCase(ctx context.Context, wg *sync.WaitGroup, log *logger.Zerolog, settingsRepo SettingsRepo, grpcClient GrpcClient) *SettingsUseCase {
	uc := &SettingsUseCase{
		ctx:          ctx,
		wg:           wg,
		settingsRepo: settingsRepo,
		grpcClient:   grpcClient,
		log:          log,
		guiSettings:  &entity.Settings{},
		updateCh:     make(chan *entity.Settings, 10),
	}

	wg.Add(1)

	uc.updateHandler()

	return uc
}

// Get reads and returns current Settings from database
func (uc *SettingsUseCase) Get() (*entity.Settings, error) {
	cfg, err := uc.settingsRepo.Get()
	if err != nil {
		uc.log.Error().Msgf("failed to get settings: %v", err)
		return nil, err
	}
	return cfg, nil
}

// Update updates application settings
func (uc *SettingsUseCase) Update(payload map[string]interface{}) *entity.GUIResponse {
	req := &entity.Settings{}
	if err := req.Model(payload); err != nil {
		return entity.ErrorGUIResponse(err)
	}

	resp, err := uc.update(req)
	if err != nil {
		return entity.ErrorGUIResponse(err)
	}

	uc.grpcClient.SetSettings(resp)

	return &entity.GUIResponse{
		Status:  entity.GUIResponseStatusOK,
		Payload: resp,
	}
}

// Set delayed writes Settings to database
func (uc *SettingsUseCase) Set(cfg *entity.Settings) {
	uc.updateCh <- cfg
}

// Stop stops SettingsUseCase and writes current Settings to database
func (uc *SettingsUseCase) Stop() {
	_, _ = uc.update(uc.guiSettings)
	uc.wg.Done()
}

func (uc *SettingsUseCase) updateHandler() {
	go func() {
		for {
			select {
			case <-uc.ctx.Done():
				return
			case cfg := <-uc.updateCh:
				uc.setGUISettings(cfg)
			case <-time.After(time.Second * updateSettingsDelay):
				_, _ = uc.update(uc.guiSettings)
			}
		}
	}()
}

func (uc *SettingsUseCase) setGUISettings(cfg *entity.Settings) {
	if cfg.WindowWidth > 0 {
		uc.guiSettings.WindowWidth = cfg.WindowWidth
	}
	if cfg.WindowHeight > 0 {
		uc.guiSettings.WindowHeight = cfg.WindowHeight
	}
	if cfg.WindowX != nil {
		uc.guiSettings.WindowX = cfg.WindowX
	}
	if cfg.WindowY != nil {
		uc.guiSettings.WindowY = cfg.WindowY
	}
}

func (uc *SettingsUseCase) update(settings *entity.Settings) (*entity.Settings, error) {
	resp, err := uc.settingsRepo.Update(settings)
	if err != nil {
		uc.log.Error().
			Interface("data", settings).
			Msgf("failed to update settings: %v", err)
	}
	return resp, err
}
