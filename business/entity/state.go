package entity

type State struct {
	Power   *bool    `json:"power"`
	Mute    *bool    `json:"mute"`
	Pause   *bool    `json:"pause"`
	Track   *string  `json:"track"`
	Channel *Channel `json:"channel"`
}
