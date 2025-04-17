package game

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"bitcrawler/pkg/entity"
	"bitcrawler/pkg/logging"
	"bitcrawler/pkg/room"
)

type Game struct {
	Room      *room.Room
	Player    *entity.Character
	Enemies   []*entity.Character
	Turn      int
	Logger    *logging.Logger
	StartTime time.Time
}

var (
	ValidCommands = []string{
		InputActionMove,
		InputActionAttack,
		InputActionUse,
		InputActionLook,
		InputActionExamine,
		InputActionOpen,
		InputActionClose,
		InputActionPick,
		InputActionDrop,
		InputActionTalk,
		InputActionRead,
		InputActionCast,
		InputActionEquip,
		InputActionUnequip,
		InputActionDrink,
		InputActionEat,
		InputActionClimb,
		InputActionSwim,
		InputActionJump,
		InputActionSneak,
		InputActionRun,
		InputActionHide,
		InputActionSearch,
		InputActionRest,
		InputActionWait,
		InputActionSleep,
		InputActionSave,
		InputActionLoad,
		InputActionQuit,
		InputActionExit,
		InputActionHelp,
		InputActionInventory,
		InputActionStatus,
		InputActionStats,
		InputActionQuests,
		InputActionJournal,
	}
)

func (g *Game) movePlayerOnInput(input string) error {
	if len(input) == 0 {
		return fmt.Errorf("cannot move without a direction")
	}

	var dirX, dirY int
	switch input {
	case DirectionUp, DirectionNorth:
		dirX, dirY = 0, Up
	case DirectionDown, DirectionSouth:
		dirX, dirY = 0, Down
	case DirectionLeft, DirectionWest:
		dirX, dirY = Left, 0
	case DirectionRight, DirectionEast:
		dirX, dirY = Right, 0
	case DirectionUpRight, DirectionNortheast:
		dirX, dirY = Right, Up
	case DirectionUpLeft, DirectionNorthwest:
		dirX, dirY = Left, Up
	case DirectionDownRight, DirectionSoutheast:
		dirX, dirY = Right, Down
	case DirectionDownLeft, DirectionSouthwest:
		dirX, dirY = Left, Down
	default:
		return fmt.Errorf("You can't go that way.")
	}

	if err := g.Room.Move(g.Player, dirX, dirY); err != nil {
		g.Room.LogView.WriteString(err.Error() + "\n")
	} else {
		g.Room.LogView.WriteString(fmt.Sprintf("%s moves %s\n", g.Player.Name, input))
	}

	return nil
}

func resolveActionObject(input string) (string, string, error) {
	var action, object string
	cmd := strings.Split(input, " ")
	if len(cmd) < 2 {
		return action, object, fmt.Errorf("invalid command")
	}

	// check for a valid command
	endOfInput := len(cmd) - 1
	var validCmdIndex, validObjIndex int
	for i, v := range cmd {
		if slices.Contains(ValidCommands, strings.ToLower(v)) {
			// a valid command needs an object after it
			if i < endOfInput {
				validCmdIndex = i
				validObjIndex = i + 1
				break
			} else {
				return action, object, fmt.Errorf("invalid command")
			}
		}

		// a valid command must be in the input
		if i == endOfInput {
			return action, object, fmt.Errorf("unknown input")
		}
	}

	action = strings.ToLower(cmd[validCmdIndex])
	object = strings.ToLower(cmd[validObjIndex])
	return action, object, nil
}

func (g *Game) resolveUserInput(input string) error {
	action, object, err := resolveActionObject(input)
	if err != nil {
		return fmt.Errorf("invalid command")
	}

	switch action {
	case InputActionMove:
		if !isValidDirection(object) {
			return fmt.Errorf("You can't go that way.")
		}

		if err := g.movePlayerOnInput(object); err != nil {
			g.Room.LogView.WriteString(err.Error() + "\n")
		}
	case InputActionAttack:
		if !isValidDirection(object) {
			return fmt.Errorf("You can't attack that way.")
		}

		playerX, playerY := g.Player.X, g.Player.Y
		attackX, attackY := resolveDirection(object)

		if err := g.Room.AttackDirection(playerX, playerY, playerX+attackX, playerY+attackY); err != nil {
			g.Room.LogView.WriteString(err.Error() + "\n")
		}
	case InputActionUse:
	case InputActionLook:
	case InputActionExamine:
	case InputActionOpen:
	case InputActionClose:
	case InputActionPick:
	case InputActionDrop:
	case InputActionTalk:
	case InputActionRead:
	case InputActionCast:
	case InputActionEquip:
	case InputActionUnequip:
	case InputActionDrink:
	case InputActionEat:
	case InputActionClimb:
	case InputActionSwim:
	case InputActionJump:
	case InputActionSneak:
	case InputActionRun:
	case InputActionHide:
	case InputActionSearch:
	case InputActionRest:
	case InputActionWait:
	case InputActionSleep:
	case InputActionSave:
	case InputActionLoad:
	case InputActionQuit:
	case InputActionExit:
	case InputActionHelp:
	case InputActionInventory:
	case InputActionStatus:
	case InputActionStats:
	case InputActionQuests:
	case InputActionJournal:
	default:
		return fmt.Errorf("unknown action")
	}

	return nil
}

func (g *Game) ProcessTurn() {
	g.Turn = (g.Turn + 1) % 256
	g.Logger.LogMessage(logging.LogLevelDebug, fmt.Sprintf("Game turn %d", g.Turn))

	g.Room.DrawRoom()
	g.Logger.LogMessage(logging.LogLevelDebug, "Room drawn")

	fmt.Println("Player's turn. Enter a command (e.g., move, attack):")
	command, err := getUserInput()
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}
	g.Logger.LogMessage(logging.LogLevelDebug, fmt.Sprintf("Input: %s", command))

	if err := g.resolveUserInput(command); err != nil {
		g.Room.LogView.WriteString(err.Error() + "\n")
		return
	}
	g.Logger.LogMessage(logging.LogLevelDebug, fmt.Sprintf("Player action resolved: %s", command))

	if g.Player.HasExited {
		g.Room.LogView.WriteString("You have exited the game.\n")
		os.Exit(0)
	}

	// Example: Enemy turn
	for _, enemy := range g.Enemies {
		if enemy.HP <= 0 && enemy.HasDied {
			g.Logger.LogMessage(logging.LogLevelDebug, fmt.Sprintf("Enemy %s is dead and cannot take its turn", enemy.Name))
			continue
		} else if enemy.HP <= 0 && !enemy.HasDied {
			g.Logger.LogMessage(logging.LogLevelDebug, fmt.Sprintf("Enemy %s has died", enemy.Name))
			enemy.HasDied = true
			g.Room.LogView.WriteString(enemy.DeathMessage)
			continue
		}

		if g.Turn%2 != 0 {
			continue
		}
		// enemy.PreHook()
		// enemyChoice := enemy.MoveOrAttack()
		// enemy.Attack()
		// enemy.ProcessMovement()
		// enemy.PostHook()
		// Pre Hooks

		// Check if the enemy is alive
		//if enemy.HP <= 0 {
		//	g.Room.LogView.WriteString(enemy.DeathMessage)
		//	g.Enemies[i] = nil
		//	continue
		//}

		// g.Room.LogView.WriteString(fmt.Sprintf("Enemy %s takes its turn\n", enemy.Name))

		// calculate the direction vector
		dx := g.Player.X - enemy.X
		dy := g.Player.Y - enemy.Y

		var dxx, dxy int
		// Normalize the direction vector
		if dx > 0 {
			dxx = Right
		}

		if dx < 0 {
			dxx = Left
		}

		if dy > 0 {
			dxy = Up
		}

		if dy < 0 {
			dxy = Down
		}

		// g.Room.LogView.WriteString("Enemy direction vector: " + fmt.Sprintf("(%d, %d)\n", dxx, dxy))
		// Move enemy towards player
		if err := g.Room.Move(enemy, dxx, dxy); err != nil {
			g.Room.LogView.WriteString(err.Error() + "\n")
		} else {
			g.Room.LogView.WriteString(fmt.Sprintf("%s moves towards the player\n", enemy.Name))
		}
	}
}
