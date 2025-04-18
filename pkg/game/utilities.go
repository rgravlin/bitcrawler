package game

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

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

func resolveDirection(input string) (int, int) {
	switch input {
	case DirectionUp, DirectionNorth:
		return 0, Up
	case DirectionDown, DirectionSouth:
		return 0, Down
	case DirectionLeft, DirectionWest:
		return Left, 0
	case DirectionRight, DirectionEast:
		return Right, 0
	case DirectionUpRight, DirectionNortheast:
		return Right, Up
	case DirectionUpLeft, DirectionNorthwest:
		return Left, Up
	case DirectionDownRight, DirectionSoutheast:
		return Right, Down
	case DirectionDownLeft, DirectionSouthwest:
		return Left, Down
	default:
		return 0, 0
	}
}

func isValidDirection(input string) bool {
	if len(input) == 0 {
		return false
	}

	switch input {
	case DirectionUp, DirectionNorth:
	case DirectionDown, DirectionSouth:
	case DirectionLeft, DirectionWest:
	case DirectionRight, DirectionEast:
	case DirectionUpRight, DirectionNortheast:
	case DirectionUpLeft, DirectionNorthwest:
	case DirectionDownRight, DirectionSoutheast:
	case DirectionDownLeft, DirectionSouthwest:
	default:
		return false
	}
	return true
}

func getUserInput() (string, error) {
	var input string
	var str string
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		input = scanner.Text()
		break
	}

	if err := scanner.Err(); err != nil {
		return str, err
	}

	// convert input to slice and lowercase
	str = strings.ToLower(input)

	return str, nil
}

func calculateDirectionVector(x1, y1, x2, y2 int) (int, int) {
	return x1 - x2, y1 - y2
}

func normalizeVector(x, y int) (int, int) {
	var dx, dy int
	if x > 0 {
		dx = Right
	} else if x < 0 {
		dx = Left
	}

	if y > 0 {
		dy = Up
	} else if y < 0 {
		dy = Down
	}

	return dx, dy
}
