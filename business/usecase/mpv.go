package usecase

import (
	"encoding/json"
	"net/url"
	"time"

	"media-box-ui/business/entity"
	"media-box-ui/pkg/logger"
)

type MPVUseCase struct {
	udsClient   UDSClient
	stateSetter StateSetter
	channels    map[string]*entity.Channel
	log         *logger.Zerolog
}

type UDSClient interface {
	Write(msg []byte) error
}

type StateSetter interface {
	Power(in *bool)
	Track(in string)
	Channel(in *entity.Channel)
	Drop()
}

func NewMPVUseCase(udsClient UDSClient, stateSetter StateSetter, channels map[string]*entity.Channel, log *logger.Zerolog) *MPVUseCase {
	mpv := &MPVUseCase{
		udsClient:   udsClient,
		stateSetter: stateSetter,
		channels:    channels,
		log:         log,
	}

	go mpv.handler()

	return mpv
}

func (uc *MPVUseCase) handler() {
	for {
		uc.getMpvCurrentTrackMetadata()
		uc.getMpvCurrentTrackPath()
		time.Sleep(time.Second)
	}
}

func (uc *MPVUseCase) getMpvCurrentTrackMetadata() {
	cmd := &entity.MpvRequest{Payload: []string{"get_property", "filtered-metadata"}}
	msg, _ := json.Marshal(cmd)

	err := uc.udsClient.Write([]byte(string(msg) + "\n"))
	if err != nil {
		uc.stateSetter.Drop()
		return
	}

	v := true
	uc.stateSetter.Power(&v)
}

func (uc *MPVUseCase) getMpvCurrentTrackPath() {
	cmd := &entity.MpvRequest{Payload: []string{"get_property", "path"}}
	msg, _ := json.Marshal(cmd)

	err := uc.udsClient.Write([]byte(string(msg) + "\n"))
	if err != nil {
		uc.stateSetter.Drop()
		return
	}

	v := true
	uc.stateSetter.Power(&v)
}

func (uc *MPVUseCase) MPVResponse(msg []byte) {
	resp := &entity.MpvResponse{}
	if err := json.Unmarshal(msg, resp); err != nil {
		uc.log.Error().Msgf("Unmarshal error: %v (%s)", err, string(msg))
		return
	}

	if resp.Error != "success" {
		return
	}

	if data, ok := resp.Data.(map[string]interface{}); ok {
		uc.stateSetter.Track(data["icy-title"].(string))
	} else if u, ok := resp.Data.(string); ok {
		ch := uc.getChannel(u)
		if ch != nil {
			uc.stateSetter.Channel(ch)
		}
	}
}

func (uc *MPVUseCase) getChannel(in string) *entity.Channel {
	u, err := url.Parse(in)
	if err != nil {
		uc.log.Error().Msgf("failed to parse channel URL: %v", err)
		return nil
	}

	if ch, ok := uc.channels[u.Path]; ok {
		uc.stateSetter.Channel(ch)
	}

	return nil
}
