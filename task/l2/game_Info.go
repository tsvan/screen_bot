package l2

type ActionStatus int

const (
	WaitForAttack    ActionStatus = 1
	Attacking        ActionStatus = 1
	WaitingForTarget ActionStatus = 2
)

type GameInfo struct {
	GameWindow      WindowSize
	Character       Character
	Enemy           Enemy
	isEnemyOnTarget bool
}

type Character struct {
	Hp float64
	Mp int
}

type Enemy struct {
	Hp          float64
	WasDefeated bool
}

type WindowSize struct {
	X int
	Y int
	W int
	H int
}
