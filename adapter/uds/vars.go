package uds

type StateUseCase interface {
	Drop()
	Power(in *bool)
	Pause()
	Mute()
}

type MPVUseCase interface {
	MPVResponse(msg []byte)
}

var (
	stateUseCase StateUseCase
	mpvUseCase   MPVUseCase
)

func SetStateUseCase(uc StateUseCase) {
	stateUseCase = uc
}

func SetMPVUseCase(uc MPVUseCase) {
	mpvUseCase = uc
}
