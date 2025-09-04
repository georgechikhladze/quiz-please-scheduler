package gameprovider

type Provider interface {
	GetGamesList() map[int][]Game
}
