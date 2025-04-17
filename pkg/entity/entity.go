package entity

type Character struct {
	ID            ID
	Name          string
	HP            int
	Attack        int
	Defense       int
	Abilities     []string
	Visual        rune
	PreHook       func(*Character)
	PostHook      func(*Character)
	PreviousX     int
	PreviousY     int
	X             int
	Y             int
	Description   string
	HealthyText   string
	DamagedText   string
	WoundedText   string
	DeadText      string
	DeathMessage  string
	SeenMessage   string
	BattleMessage string
	HasDied       bool
}

type ID uint8

const (
	ObjEmpty ID = iota
	ObjPlayer
	ObjEnemy
	ObjWall
	ObjExit
)
