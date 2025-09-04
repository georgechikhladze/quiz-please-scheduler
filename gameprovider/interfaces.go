package gameprovider

type GameProvider interface {
	GetGamesList() map[int][]Game
}
