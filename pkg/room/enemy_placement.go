package room

import (
	"bitcrawler/pkg/entity"
)

func (r *Room) PlaceGoblinPack(goblins int, goblinLeader bool) []*entity.Character {
	emptyX, emptyY := r.FindEmptySpace()
	emptyArea := r.FindEmptySpacesCloseTogether(emptyX, emptyY, 1)
	//logger.LogMessage(logging.LogLevelDebug,
	//	fmt.Sprintf("Empty area found at coordinates: (%d, %d)", emptyX, emptyY))

	var enemies []*entity.Character
	var counter int
	for _, coord := range emptyArea {
		if counter > goblins {
			if goblinLeader {
				if coord.Entity.ID == entity.ObjEmpty {
					goblinLeaderEnemy := entity.NewEnemy(entity.GoblinLeaderTemplate)
					goblinLeaderEnemy.X = coord.X
					goblinLeaderEnemy.Y = coord.Y

					r.AddEntity(&Coordinate{X: coord.X, Y: coord.Y, Entity: goblinLeaderEnemy})
					enemies = append(enemies, goblinLeaderEnemy)

					//logger.LogMessage(logging.LogLevelDebug,
					//	fmt.Sprintf("Enemy added at coordinates: (%d, %d)", coord.X, coord.Y))
				}

				break
			}
		}

		if coord.Entity == nil || coord.Entity.ID == entity.ObjEmpty {
			// create new enemy from the template
			e := entity.NewEnemy(entity.GoblinEnemyTemplate)
			e.X = coord.X
			e.Y = coord.Y

			r.AddEntity(&Coordinate{X: coord.X, Y: coord.Y, Entity: e})
			enemies = append(enemies, e)

			//logger.LogMessage(logging.LogLevelDebug,
			//	fmt.Sprintf("Enemy added at coordinates: (%d, %d)", coord.X, coord.Y))
		}
		counter++
	}

	return enemies
}
