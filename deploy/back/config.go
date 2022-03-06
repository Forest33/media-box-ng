package main

import (
	"encoding/json"
	"os"

	"media-box-ui/business/entity"
	"media-box-ui/pkg/logger"
)

const (
	defaultConfigPath = "media-box-back.json"
)

type Config struct {
	Broker   *BrokerConfig              `json:"broker"`
	UDS      *UDSConfig                 `json:"uds"`
	MPV      *MPVConfig                 `json:"mpv"`
	Channels map[string]*entity.Channel `json:"channels"`
	Logger   *LoggerConfig              `json:"logger"`
}

type LoggerConfig struct {
	Level             string `json:"level"`
	TimeFieldFormat   string `json:"time_field_format"`
	PrettyPrint       bool   `json:"pretty_print"`
	DisableSampling   bool   `json:"disable_sampling"`
	RedirectStdLogger bool   `json:"redirect_std_logger"`
	ErrorStack        bool   `json:"error_stack"`
	ShowCaller        bool   `json:"show_caller"`
}

type BrokerConfig struct {
	Host       string `json:"host"`
	Port       int    `json:"port"`
	StateTopic string `json:"state_topic"`
	ClientID   string `json:"client_id"`
	UserName   string `json:"user_name"`
	Password   string `json:"password"`
}

type UDSConfig struct {
	ServerSocket string `json:"server_socket"`
	ClientSocket string `json:"client_socket"`
}

type MPVConfig struct {
	PowerOnPath string `json:"power_on_path"`
}

var (
	cfg = &Config{}
)

func init() {
	log := logger.NewDefaultZerolog()

	path, ok := os.LookupEnv("MEDIA_BOX_BACK_CONFIG")
	if !ok {
		path = defaultConfigPath
	}

	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		log.Fatal().Msg(err.Error())
	}
}
