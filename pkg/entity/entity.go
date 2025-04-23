package entity

type Character struct {
	ID            ID
	Name          string
	HP            int
	Attack        int
	Defense       int
	Abilities     []Ability
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
	HasExited     bool
}

type Ability struct {
	Name        string
	Description string
	Effect      Effect
}

type Effect struct {
	Attack  int
	Defense int
	HP      int
}

type ID uint8

const (
	ObjEmpty ID = iota
	ObjPlayer
	ObjEnemy
	ObjWall
	ObjExit
)
