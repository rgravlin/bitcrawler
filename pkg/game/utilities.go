package game

import (
	"bufio"
	"os"
	"strings"
)

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
