package main

import (
	"os"
	"os/signal"
	"syscall"

	"media-box-ui/adapter/broker"
	"media-box-ui/adapter/image"
	"media-box-ui/business/usecase"
	"media-box-ui/pkg/logger"
)

var (
	log *logger.Zerolog

	brokerClient *broker.Client
	imageClient  *image.Client

	guiUseCase *usecase.GUIUseCase
)

func main() {
	defer shutdown()

	log = logger.NewZerolog(logger.ZeroConfig{
		Level:             cfg.Logger.Level,
		TimeFieldFormat:   cfg.Logger.TimeFieldFormat,
		PrettyPrint:       cfg.Logger.PrettyPrint,
		DisableSampling:   cfg.Logger.DisableSampling,
		RedirectStdLogger: cfg.Logger.RedirectStdLogger,
		ErrorStack:        cfg.Logger.ErrorStack,
		ShowCaller:        cfg.Logger.ShowCaller,
	})

	initAdapters()
	initUseCases()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-quit
}

func initAdapters() {
	var err error
	brokerClient, err = broker.NewBrokerClient(&broker.Config{
		Host:       cfg.Broker.Host,
		Port:       cfg.Broker.Port,
		StateTopic: cfg.Broker.StateTopic,
		ClientID:   cfg.Broker.ClientID,
		UserName:   cfg.Broker.UserName,
		Password:   cfg.Broker.Password,
	}, log)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	imageClient = image.NewImageClient(&image.Config{
		OutputPath: cfg.Image.OutputPath,
	}, log)
}

func initUseCases() {
	usecase.SetStateTopic(cfg.Broker.StateTopic)
	guiUseCase = usecase.NewGUIUseCase(brokerClient, imageClient, log)
	guiUseCase.Start()
}

func shutdown() {
	brokerClient.Close()
}
