package main

var PCSkillsAt = 6

var CyberPsycho = Character{
	CharacterStats: CyberPsychoStats,
	Weapon:         HeavyPistol,
	ArmorValue:     12,
	ArmorPenalty:   0,
	ShouldDodge:    true,
}

var BoostGanger = Character{
	CharacterStats: BoostGangerStats,
	Weapon:         HeavyPistol,
	ArmorValue:     4,
	ArmorPenalty:   0,
	ShouldDodge:    false,
}

var PlayerCharacter = CharacterStats{
	Name:            "Player",
	MaxHP:           40,
	Reflexes:        8,
	Dexterity:       8,
	Movement:        8,
	Evasion:         PCSkillsAt,
	Brawling:        PCSkillsAt,
	Handguns:        PCSkillsAt,
	ShoulderArms:    PCSkillsAt,
	HeavyWeapons:    PCSkillsAt,
	AutoFire:        PCSkillsAt,
	Melee:           PCSkillsAt,
	MartialArts:     PCSkillsAt,
	AttackModifiers: []Modifier{},
}

var BoostGangerStats = CharacterStats{
	Name:      "Boost Ganger",
	MaxHP:     20,
	Reflexes:  6,
	Dexterity: 5,
	Brawling:  4,
	Movement:  4,
	Evasion:   0,
}

var CyberPsychoStats = CharacterStats{
	Name:      "CyberPsycho",
	MaxHP:     55,
	Reflexes:  8,
	Dexterity: 8,
	Evasion:   6,
	Brawling:  6,
}
