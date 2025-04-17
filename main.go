package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"bitcrawler/pkg/entity"
	"bitcrawler/pkg/game"
	"bitcrawler/pkg/logging"
	"bitcrawler/pkg/room"
)

type Action int

func main() {
	// initialize the logger
	logger, err := logging.NewLogger(logging.LogLevelDebug)
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}

	// Initialize start time
	startTime := time.Now()
	logger.LogMessage(logging.LogLevelInfo, "Game started")

	// Initialize the room for the game
	rm := room.NewRoom(24, 8, 1)
	logger.LogMessage(logging.LogLevelDebug, fmt.Sprintf("Room initialized with dimensions: %d x %d", rm.Width, rm.Height))

	// Find an empty space in the room to place the player
	randX, randY := rm.FindEmptySpace()
	if randX == -1 && randY == -1 {
		panic("Cannot initialize room, no empty space found!")
	}

	// Initialize our player character with our random coordinates
	player := &entity.Character{X: randX, Y: randY, ID: entity.ObjPlayer, Name: "Hero", HP: 100, Attack: 10, Defense: 5, Visual: '@'}
	logger.LogMessage(logging.LogLevelDebug, fmt.Sprintf("Player initialized at coordinates: (%d, %d)", randX, randY))

	// Add the player to the room
	rm.AddEntity(&room.Coordinate{X: randX, Y: randY, Entity: player})
	logger.LogMessage(logging.LogLevelDebug, "Player added to the room")

	// Find the farthest distance from the player coordinates
	exitX, exitY := FindFarthestDistance(rm, randX, randY, false)
	logger.LogMessage(logging.LogLevelDebug, fmt.Sprintf("Farthest exit found at coordinates: (%d, %d)", exitX, exitY))

	// Create the exit
	exit := &entity.Character{X: exitX, Y: exitY, ID: entity.ObjExit, Name: "Exit", HP: 0, Attack: 0, Defense: 0, Visual: '='}
	logger.LogMessage(logging.LogLevelDebug, "Exit created")

	// Add the exit to the room with the farthest coordinates
	rm.AddEntity(&room.Coordinate{X: exitX, Y: exitY, Entity: exit})
	logger.LogMessage(logging.LogLevelDebug, "Exit added to the room")

	// Add the enemies
	enemyCount := rand.Intn(3) + 1
	enemies := addRandomEnemies(rm, enemyCount, "Goblin", 'G', 30, 5, 2)
	logger.LogMessage(logging.LogLevelDebug, fmt.Sprintf("%d enemies added to the room", enemyCount))

	gameBoard := &game.Game{StartTime: startTime, Logger: logger, Room: rm, Player: player, Enemies: enemies}
	logger.LogMessage(logging.LogLevelDebug, "Game board initialized")

	// Game loop
	logger.LogMessage(logging.LogLevelDebug, "Game started")
	for {
		gameBoard.ProcessTurn()
	}
}

func addRandomEnemies(rm *room.Room, enemyCount int, enemyName string, enemyVisual rune, enemyHP, enemyAtk, enemyDef int) []*entity.Character {
	enemies := make([]*entity.Character, enemyCount)
	for i := 0; i < enemyCount; i++ {
		enemyX, enemyY := rm.FindEmptySpace()
		if enemyX == -1 && enemyY == -1 {
			// can't find empty space so don't spawn the enemy
			continue
		}

		enemy := &entity.Character{X: enemyX, Y: enemyY, ID: entity.ObjEnemy, Name: enemyName, HP: enemyHP, Attack: enemyAtk, Defense: enemyDef, Visual: enemyVisual} //nolint:lll
		rm.AddEntity(&room.Coordinate{X: enemyX, Y: enemyY, Entity: enemy})
		enemies[i] = enemy
	}

	return enemies
}

// Distance calculates the Euclidean distance between two points (x1, y1) and (x2, y2).
func Distance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))
}

func FindFarthestDistance(room *room.Room, playerX, playerY int, skipWalls bool) (int, int) {
	maxX, maxY := 0, 0
	maxDistance := 0.0

	for x := 0; x < room.Width; x++ {
		for y := 0; y < room.Height; y++ {
			if skipWalls {
				if room.Grid[x][y].Entity.ID == entity.ObjWall {
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
