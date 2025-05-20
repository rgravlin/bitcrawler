package gear

import "bitcrawler/pkg/entity"

var (
	AbilityMightStrength = entity.Ability{
		Name: "Mighty Strength",
		Description: "Increases your strength by 5 by a divine force",
		Effect: entity.Effect{
			Attack:  5,
		},
	}
)

// when we attack, iterate through ability effects and add to each field
