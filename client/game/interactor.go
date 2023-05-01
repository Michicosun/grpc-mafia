package game

type IInteractor interface {
	Run()

	Signal()
	GetPrefixString() string
	GetCurBuf() string
}
