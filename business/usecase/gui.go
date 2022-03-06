package usecase

import (
	"encoding/json"
	"regexp"

	"github.com/gen2brain/beeep"

	"media-box-ui/business/entity"
	"media-box-ui/pkg/logger"
)

type GUIUseCase struct {
	broker Broker
	image  Image
	log    *logger.Zerolog
}

var (
	trackRe = regexp.MustCompile("[[:^ascii:]]")
)

func NewGUIUseCase(broker Broker, image Image, log *logger.Zerolog) *GUIUseCase {
	return &GUIUseCase{
		broker: broker,
		image:  image,
		log:    log,
	}
}

func (uc *GUIUseCase) Start() {
	uc.broker.Subscribe(stateTopic, func(topic string, payload []byte) {
		uc.log.Debug().Msgf("%s - %s", topic, string(payload))

		state, err := uc.parseState(payload)
		if err != nil {
			uc.log.Error().Msgf("failed to parse state: %v", err)
			return
		}

		uc.notify(state)
	})
}

func (uc *GUIUseCase) parseState(payload []byte) (*entity.State, error) {
	state := &entity.State{}
	if err := json.Unmarshal(payload, state); err != nil {
		return nil, err
	}
	return state, nil
}

func (uc *GUIUseCase) notify(state *entity.State) {
	if (state.Track != nil && len(*state.Track) == 0) || (state.Channel != nil && len(state.Channel.Title) == 0) {
		return
	}

	img, err := uc.image.Get(state.Channel.Img)
	if err != nil {
		uc.log.Error().Msgf("failed to get channel image: %v", err)
	}

	title := state.Channel.Title
	if state.Power != nil && !*state.Power {
		title += " [OFF]"
	}
	if state.Mute != nil && *state.Mute {
		title += " [MUTE]"
	}
	if state.Pause != nil && *state.Pause {
		title += " [PAUSE]"
	}

	err = beeep.Notify(title, uc.prepareTrackName(*state.Track), img)
	if err != nil {
		uc.log.Error().Msgf("failed to show notification: %v", err)
	}
}

func (uc *GUIUseCase) prepareTrackName(t string) string {
	// David Helpling \u0026 Jon Jenkins - Two Paths
	return trackRe.ReplaceAllLiteralString(t, "")
}
