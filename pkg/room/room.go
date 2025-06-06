package room

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"

	"bitcrawler/pkg/entity"
)

type Room struct {
	Level            int
	Width, Height    int
	CameraX, CameraY int
	Grid             [][]*Coordinate
	DungeonView      *DungeonView
	LogView          strings.Builder
}

type DungeonView string

type Coordinate struct {
	Entity *entity.Character
	X      int
	Y      int
}

var (
	emptyObject = &entity.Character{
		ID:      entity.ObjEmpty,
		Name:    "empty",
		HP:      0,
		Attack:  0,
		Defense: 0,
		Visual:  '.',
	}
)

func NewRoom(width, height, level int) *Room {
	grid := make([][]*Coordinate, width)
	for i := range grid {
		grid[i] = make([]*Coordinate, height)
		for j := range grid[i] {
			defaultID := entity.ObjEmpty
			isWall := (i == 0 || i == width-1) || (j == 0 || j == height-1)
			if isWall {
				defaultID = entity.ObjWall
			}
			grid[i][j] = &Coordinate{X: i, Y: j, Entity: &entity.Character{ID: defaultID}}
		}
	}
	return &Room{Width: width, Height: height, Grid: grid, Level: level}
}

func (r *Room) AddEntity(c *Coordinate) {
	if c.X >= 0 && c.X < r.Width && c.Y >= 0 && c.Y < r.Height {
		r.Grid[c.X][c.Y] = c
	}
}

func (r *Room) DrawRoom() {
	var builder strings.Builder
	for y := r.Height - 1; y >= 0; y-- { // Start from the top row
		for x := 0; x < r.Width; x++ {
			switch r.Grid[x][y].Entity.ID {
			case entity.ObjEmpty:
				builder.WriteString(". ")
			case entity.ObjPlayer, entity.ObjEnemy:
				if r.Grid[x][y].Entity.ID == entity.ObjEnemy && r.Grid[x][y].Entity.HasDied {
					builder.WriteString("x ")
				} else {
					builder.WriteString(string(r.Grid[x][y].Entity.Visual) + " ")
				}
			case entity.ObjWall:
				builder.WriteString("# ")
			default:
				builder.WriteString("? ") // Unknown entity
			}
		}
		builder.WriteString("\n") // Move to the next row
	}

	builder.WriteString(r.LogView.String())
	r.LogView.Reset()
	fmt.Print("\033[H\033[2J") // Clear the screen and move the cursor to the top-left
	fmt.Printf("%s", builder.String())
}

func (r *Room) Move(character *entity.Character, x, y int) error {
	// Calculate new position
	newX := character.X + x
	newY := character.Y + y

	// Check if the new position is within bounds
	if newX < 0 || newX >= r.Width || newY < 0 || newY >= r.Height {
		return errors.New("you cannot escape into the void")
	}

	// Check if the new position is not a wall
	if r.Grid[newX][newY].Entity != nil && r.Grid[newX][newY].Entity.ID == entity.ObjWall {
		return errors.New("you run into a wall")
	}

	// Check if the new position is an exit
	if r.Grid[newX][newY].Entity != nil && r.Grid[newX][newY].Entity.ID == entity.ObjExit {
		r.LogView.WriteString("You have found the exit!\n")
		character.HasExited = true
		return nil
	}

	// Check if we're both enemies
	if r.Grid[newX][newY].Entity != nil && r.Grid[newX][newY].Entity.ID == entity.ObjEnemy && character.ID == entity.ObjEnemy {
		return errors.New("Enemies cannot move into each other")
	}

	// Check if the new position is an enemy or player
	if r.Grid[newX][newY].Entity != nil && (r.Grid[newX][newY].Entity.ID == entity.ObjEnemy || r.Grid[newX][newY].Entity.ID == entity.ObjPlayer) {
		var shouldAttack bool
		// handle player attacking enemy
		if r.Grid[newX][newY].Entity.ID == entity.ObjEnemy && character.ID == entity.ObjPlayer {
			shouldAttack = true
		}

		// handle enemy attacking player
		if r.Grid[newX][newY].Entity.ID == entity.ObjPlayer && character.ID == entity.ObjEnemy {
			shouldAttack = true
		}

		if shouldAttack {
			r.AttackEntity(character, r.Grid[newX][newY].Entity)
			return nil
		}
	}

	// Update the room Grid
	r.Grid[character.X][character.Y].Entity = emptyObject
	r.Grid[newX][newY].Entity = character

	// Update the character's positions
	character.PreviousX = character.X
	character.PreviousY = character.Y
	character.X = newX
	character.Y = newY

	return nil
}

// AddWallsWithDoorway adds walls to a room and ensures there is a doorway.
func (r *Room) AddWallsWithDoorway(horizontal bool, position int) {
	if horizontal {
		// Add a horizontal wall at the given y position
		if position < 0 || position >= r.Height {
			return // Invalid position
		}
		doorway := rand.Intn(r.Width) // Random doorway position
		for x := 0; x < r.Width; x++ {
			if x == doorway {
				continue // Leave a doorway
			}
			r.Grid[x][position].Entity = &entity.Character{ID: entity.ObjWall, Visual: '#'}
		}
	} else {
		// Add a vertical wall at the given x position
		if position < 0 || position >= r.Width {
			return // Invalid position
		}
		doorway := rand.Intn(r.Height) // Random doorway position
		for y := 0; y < r.Height; y++ {
			if y == doorway {
				continue // Leave a doorway
			}
			r.Grid[position][y].Entity = &entity.Character{ID: entity.ObjWall, Visual: '#'}
		}
	}
}

func (r *Room) AttackEntity(attacker, defender *entity.Character) {
	if defender.ID == entity.ObjEmpty {
		r.LogView.WriteString("You attack into the air and almost hit yourself!\n")
		return
	}

	if defender.ID == entity.ObjWall {
		r.LogView.WriteString("You attack and hit a wall!\n")
		return
	}

	if defender.HP > 0 {
		r.LogView.WriteString(fmt.Sprintf("%s attacks %s!\n", attacker.Name, defender.Name))

		// calculate abilities
		var attackerAttackIncrease int
		var defenderDefenseIncrease int

		if len(attacker.Abilities) > 0 {
			for _, ability := range attacker.Abilities {
				attackerAttackIncrease += ability.Effect.Attack
			}
		}

		if len(defender.Abilities) > 0 {
			for _, ability := range defender.Abilities {
				defenderDefenseIncrease += ability.Effect.Defense
			}
		}

		defender.HP -= (attacker.Attack + attackerAttackIncrease) - (defender.Defense - defenderDefenseIncrease)
		if defender.HP <= 0 {
			r.LogView.WriteString(fmt.Sprintf("%s is defeated!\n", defender.Name))
		} else {
			r.LogView.WriteString(fmt.Sprintf("%s has %d HP left.\n", defender.Name, defender.HP))
		}
	} else {
		r.LogView.WriteString(fmt.Sprintf("%s is already defeated!\n", defender.Name))
	}
}

func (r *Room) AttackDirection(x1, y1, x2, y2 int) error {
	// Check if the coordinates are within bounds
	if x1 < 0 || x1 >= r.Width || y1 < 0 || y1 >= r.Height ||
		x2 < 0 || x2 >= r.Width || y2 < 0 || y2 >= r.Height {
		return errors.New("coordinates out of bounds")
	}

	// Check if the attacker and defender are in the same position
	if x1 == x2 && y1 == y2 {
		return errors.New("cannot attack yourself")
	}

	attacker := r.Grid[x1][y1].Entity
	defender := r.Grid[x2][y2].Entity

	if attacker == nil || defender == nil {
		return errors.New("no entity at the given coordinates")
	}

	r.AttackEntity(attacker, defender)
	return nil
}

func (r *Room) FindEmptySpace() (int, int) {
	// put all empty space into a list
	emptySpaces := make([]*Coordinate, 0)
	for x := 0; x < r.Width; x++ {
		for y := 0; y < r.Height; y++ {
			if r.Grid[x][y].Entity.ID == entity.ObjEmpty {
				emptySpaces = append(emptySpaces, r.Grid[x][y])
			}
		}
	}

	// if there are no empty spaces, return -1, -1
	if len(emptySpaces) == 0 {
		return -1, -1
	}

	// pick a random empty space
	randIndex := rand.Intn(len(emptySpaces))
	return emptySpaces[randIndex].X, emptySpaces[randIndex].Y
}

func (r *Room) FindEmptySpacesCloseTogether(x, y, distance int) []*Coordinate {
	emptySpaces := make([]*Coordinate, 0)
	for i := x - distance; i <= x+distance; i++ {
		for j := y - distance; j <= y+distance; j++ {
			// skip walls and out of bounds
			if i < 0 || i >= r.Width-1 || j < 0 || j >= r.Height-1 {
				continue
			}
			if r.Grid[i][j].Entity.ID == entity.ObjEmpty {
				emptySpaces = append(emptySpaces, r.Grid[i][j])
			}
		}
	}
	return emptySpaces
}

func (r *Room) FindFarthestDistance(playerX, playerY int, skipWalls bool) (int, int) {
	maxX, maxY := 0, 0
	maxDistance := 0.0

	for x := 0; x < r.Width; x++ {
		for y := 0; y < r.Height; y++ {
			if skipWalls {
				if r.Grid[x][y].Entity.ID == entity.ObjWall {
					continue
				}
			}
			distance := Distance(float64(playerX), float64(playerY), float64(x), float64(y))
			if distance > maxDistance {
				maxDistance = distance
				maxX, maxY = x, y
			}
		}
	}

	return maxX, maxY
}

func (r *Room) AddRandomEntities(template entity.Character, entityCount int) []*entity.Character {
	entities := make([]*entity.Character, entityCount)
	for i := 0; i < entityCount; i++ {
		entityX, entityY := r.FindEmptySpace()
		if entityX == -1 && entityY == -1 {
			// can't find empty space so don't spawn the e
			continue
		}

		e := &entity.Character{}
		*e = template
		e.X = entityX
		e.Y = entityY
		r.AddEntity(&Coordinate{X: entityX, Y: entityY, Entity: e})
		entities[i] = e
	}

	return entities
}
