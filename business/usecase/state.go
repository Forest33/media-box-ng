package usecase

import (
	"encoding/json"
	"os/exec"
	"sync"

	"media-box-ui/business/entity"
	"media-box-ui/pkg/logger"
)

type StateUseCase struct {
	cfg    *StateConfig
	broker Broker
	log    *logger.Zerolog
	state  *entity.State
	mu     sync.RWMutex
}

type StateConfig struct {
	PowerOnPath string
}

func NewStateUseCase(cfg *StateConfig, broker Broker, log *logger.Zerolog) *StateUseCase {
	return &StateUseCase{
		cfg:    cfg,
		broker: broker,
		log:    log,
		state:  defaultState(),
	}
}

func defaultState() *entity.State {
	bVal := false
	sVal := ""
	chVal := entity.Channel{}

	return &entity.State{
		Power:   &bVal,
		Mute:    &bVal,
		Pause:   &bVal,
		Track:   &sVal,
		Channel: &chVal,
	}
}

func (uc *StateUseCase) Power(in *bool) {
	if in == nil {
		cmd := exec.Command(uc.cfg.PowerOnPath)
		resp, err := cmd.Output()
		if err != nil {
			uc.log.Error().Msg(err.Error())
			return
		}
		uc.log.Debug().Msg(string(resp))
	}

	uc.mu.RLock()
	defer uc.mu.RUnlock()

	if in == nil {
		v := !*uc.state.Power
		uc.state.Power = &v
		uc.publish()
	} else if in != uc.state.Power {
		uc.state.Power = in
	}
}

func (uc *StateUseCase) Pause() {
	uc.mu.RLock()
	defer uc.mu.RUnlock()

	v := !*uc.state.Pause
	uc.state.Pause = &v
	uc.publish()
}

func (uc *StateUseCase) Mute() {
	uc.mu.RLock()
	defer uc.mu.RUnlock()

	v := !*uc.state.Mute
	uc.state.Mute = &v
	uc.publish()
}

func (uc *StateUseCase) Track(in string) {
	uc.mu.RLock()
	defer uc.mu.RUnlock()

	if uc.state.Track == nil || *uc.state.Track != in {
		uc.state.Track = &in
		uc.publish()
	}
}

func (uc *StateUseCase) Channel(in *entity.Channel) {
	uc.mu.RLock()
	defer uc.mu.RUnlock()

	if uc.state.Channel == nil || uc.state.Channel.Title != in.Title {
		uc.state.Channel = in
		uc.publish()
	}
}

func (uc *StateUseCase) Drop() {
	uc.mu.RLock()
	defer uc.mu.RUnlock()

	power := false
	track := ""
	channel := entity.Channel{}
	uc.state.Power = &power
	uc.state.Track = &track
	uc.state.Channel = &channel
}

func (uc *StateUseCase) publish() {
	data, err := json.Marshal(uc.state)
	if err != nil {
		uc.log.Error().Msg(err.Error())
		return
	}
	uc.broker.PublishState(data)
}
