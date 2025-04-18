package main

import (
	"fmt"
	"math/rand"
	"time"

	"bitcrawler/pkg/entity"
	"bitcrawler/pkg/game"
	"bitcrawler/pkg/logging"
	"bitcrawler/pkg/room"
)

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
	logger.LogMessage(logging.LogLevelDebug,
		fmt.Sprintf("Room initialized with dimensions: %d x %d", rm.Width, rm.Height))

	// Find an empty space in the room to place the player
	randX, randY := rm.FindEmptySpace()
	if randX == -1 && randY == -1 {
		panic("Cannot initialize room, no empty space found!")
	}

	// Initialize our player character with our random coordinates
	player := &entity.Character{
		X:       randX,
		Y:       randY,
		ID:      entity.ObjPlayer,
		Name:    "Hero",
		HP:      100,
		Attack:  10,
		Defense: 5,
		Visual:  '@',
	}
	logger.LogMessage(logging.LogLevelDebug,
		fmt.Sprintf("Player initialized at coordinates: (%d, %d)", randX, randY))

	// Add the player to the room
	rm.AddEntity(&room.Coordinate{X: randX, Y: randY, Entity: player})
	logger.LogMessage(logging.LogLevelDebug, "Player added to the room")

	// Find the farthest distance from the player coordinates
	exitX, exitY := rm.FindFarthestDistance(randX, randY, false)
	logger.LogMessage(logging.LogLevelDebug,
		fmt.Sprintf("Farthest exit found at coordinates: (%d, %d)", exitX, exitY))

	// Create the exit
	exit := &entity.Character{
		X: exitX,
		Y: exitY,
		ID: entity.ObjExit,
		Name: "Exit",
	}
	logger.LogMessage(logging.LogLevelDebug, "Exit created")

	// Add the exit to the room with the farthest coordinates
	rm.AddEntity(&room.Coordinate{X: exitX, Y: exitY, Entity: exit})
	logger.LogMessage(logging.LogLevelDebug, "Exit added to the room")

	// Add the enemies
	enemyCount := rand.Intn(3) + 1
	enemyTemplate := entity.Character{
		ID:      entity.ObjEnemy,
		Name:    "Goblin",
		HP:      30,
		Attack:  5,
		Defense: 2,
		Visual:  'G',
	}
	enemies := rm.AddRandomEntities(enemyTemplate, enemyCount)
	logger.LogMessage(logging.LogLevelDebug,
		fmt.Sprintf("%d enemies added to the room", enemyCount))

	gameBoard := &game.Game{StartTime: startTime, Logger: logger, Room: rm, Player: player, Enemies: enemies}
	logger.LogMessage(logging.LogLevelDebug, "Game board initialized")

	// Game loop
	logger.LogMessage(logging.LogLevelDebug, "Game started")
	for {
		gameBoard.ProcessTurn()
	}
}
