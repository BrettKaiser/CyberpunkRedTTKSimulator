package main

var VeryClose = RangeBand{
	MinDistance: 0,
	MaxDistance: 6,
	Name:        "VeryClose 0-6",
}

var Close = RangeBand{
	MinDistance: 7,
	MaxDistance: 12,
	Name:        "Close 7-12",
}

var Medium = RangeBand{
	MinDistance: 13,
	MaxDistance: 25,
	Name:        "Medium 13-25",
}

var MediumFar = RangeBand{
	MinDistance: 26,
	MaxDistance: 50,
	Name:        "Medium Far 26-50",
}

var Far = RangeBand{
	MinDistance: 51,
	MaxDistance: 100,
	Name:        "Far 51-100",
}

var VeryFar = RangeBand{
	MinDistance: 101,
	MaxDistance: 200,
	Name:        "Very Far 101-200",
}

// var VeryVeryFar = RangeBand{
// 	MinDistance: 201,
// 	MaxDistance: 400,
// 	Name:        "Very Very Far 201-400",
// }
//
// var InsanelyFar = RangeBand{
// 	MinDistance: 401,
// 	MaxDistance: 800,
// 	Name:        "Insanely Far",
// }

var RangeBands = []RangeBand{
	VeryClose,
	Close,
	Medium,
	MediumFar,
	Far,
	VeryFar,
	// VeryVeryFar,
	// InsanelyFar,
}

type Ammunition struct {
	Name string
	Cost int `json:"cost"`
}

type AmmunitionType string

const (
	Basic         AmmunitionType = "basic"
	ArmorPiercing AmmunitionType = "armor_piercing"
	Incendiary    AmmunitionType = "incendiary"
)

type WeaponType string
