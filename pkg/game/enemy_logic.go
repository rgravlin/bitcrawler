package game

import (
	"fmt"

	"bitcrawler/pkg/entity"
	"bitcrawler/pkg/logging"
)

func goblinMoveOrAttack(g *Game, enemy *entity.Character) {
	if g.Turn%2 != 0 {
		g.Logger.LogMessage(logging.LogLevelDebug,
			fmt.Sprintf("Enemy %s is waiting for their turn", enemy.Name))
		return
	}

	g.Logger.LogMessage(logging.LogLevelDebug,
		fmt.Sprintf("Enemy %s takes its turn", enemy.Name))

	// calculate the direction vector
	dx, dy := calculateDirectionVector(g.Player.X, g.Player.Y, enemy.X, enemy.Y)

	// Normalize the direction vector
	dxx, dxy := normalizeVector(dx, dy)
	g.Logger.LogMessage(logging.LogLevelDebug,
		fmt.Sprintf("Enemy %s direction vector: (%d, %d)", enemy.Name, dxx, dxy))

	// Move enemy towards player
	if err := g.Room.Move(enemy, dxx, dxy); err != nil {
		g.Logger.LogMessage(logging.LogLevelDebug, err.Error())
	} else {
		g.Room.LogView.WriteString(fmt.Sprintf("%s moves towards the player\n", enemy.Name))
	}
}
