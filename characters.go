package main

var PCStatsAt = 6

var MaxCombatAwarenessPrecision = Modifier{
	Name:  "CombatAwarenessPrecision",
	Value: 3,
}

var MinCombatAwarenessPrecision = Modifier{
	Name:  "CombatAwarenessPrecision",
	Value: 1,
}

var SmartLink = Modifier{
	Name:  "SmartLink",
	Value: 1,
}

var PlayerCharacter = CharacterStats{
	Name:         "Player",
	MaxHP:        40,
	Reflexes:     8,
	Dexterity:    8,
	Movement:     8,
	Evasion:      PCStatsAt,
	Brawling:     PCStatsAt,
	Handguns:     PCStatsAt,
	ShoulderArms: PCStatsAt,
	HeavyWeapons: PCStatsAt,
	AutoFire:     PCStatsAt,
	Melee:        PCStatsAt,
	MartialArts:  PCStatsAt,
	AttackModifiers: []Modifier{
		MaxCombatAwarenessPrecision,
		// SmartLink,
	},
}

var BoostGangerStats = CharacterStats{
	Name:      "Boost Ganger",
	MaxHP:     20,
	Reflexes:  6,
	Dexterity: 5,
	Evasion:   0,
}

var BoostGanger = Character{
	CharacterStats: BoostGangerStats,
	Weapon:         HeavyPistol,
	ArmorValue:     4,
	ArmorPenalty:   0,
	ShouldDodge:    false,
}

var CyberPsychoStats = CharacterStats{
	Name:      "CyberPsycho",
	MaxHP:     55,
	Reflexes:  8,
	Dexterity: 8,
	Evasion:   6,
	Brawling:  6,
}

var CyberPsycho = Character{
	CharacterStats: CyberPsychoStats,
	Weapon:         HeavyPistol,
	ArmorValue:     12,
	ArmorPenalty:   0,
	ShouldDodge:    true,
}
