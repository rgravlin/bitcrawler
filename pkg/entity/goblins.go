package entity

var (
	GoblinEnemyTemplate = Character{
		ID:      ObjEnemy,
		Name:    "Goblin",
		HP:      30,
		Attack:  10,
		Defense: 2,
		Visual:  'g',
	}
	GoblinLeaderTemplate = Character{
		ID:      ObjEnemy,
		Name:    "Goblin Leader",
		HP:      30,
		Attack:  15,
		Defense: 5,
		Visual:  'G',
	}
)

func NewEnemy(character Character) *Character {
	return &Character{
		Name:    character.Name,
		ID:      character.ID,
		HP:      character.HP,
		Attack:  character.Attack,
		Defense: character.Defense,
		Visual:  character.Visual,
	}
}
