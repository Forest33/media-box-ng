package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/soellman/pidfile"

	"media-box-ui/adapter/broker"
	"media-box-ui/adapter/uds"
	"media-box-ui/business/usecase"
	"media-box-ui/pkg/logger"
)

var (
	log *logger.Zerolog

	brokerClient *broker.Client
	udsServer    *uds.Server
	udsClient    *uds.Client

	stateUseCase *usecase.StateUseCase
	mpvUseCase   *usecase.MPVUseCase
)

const (
	pidFile = "/tmp/media-box.pid"
)

// GOOS=linux GOARCH=arm go build

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

	if err := pidfile.Write(pidFile); err != nil {
		log.Fatal().Msgf("failed to create pid file: %v", err)
	}

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

	udsServer, err = uds.NewUDSServer(&uds.ServerConfig{
		SocketPath:     cfg.UDS.ServerSocket,
		CommandTimeout: 2,
	}, log)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	udsClient = uds.NewUDSClient(&uds.ClientConfig{
		SocketPath: cfg.UDS.ClientSocket,
	}, log)
}

func initUseCases() {
	stateUseCase = usecase.NewStateUseCase(&usecase.StateConfig{
		PowerOnPath: cfg.MPV.PowerOnPath,
	}, brokerClient, log)
	mpvUseCase = usecase.NewMPVUseCase(udsClient, stateUseCase, cfg.Channels, log)

	uds.SetStateUseCase(stateUseCase)
	uds.SetMPVUseCase(mpvUseCase)
}

func shutdown() {
	brokerClient.Close()
	udsServer.Close()
	udsClient.Close()
	_ = pidfile.Remove(pidFile)
}
