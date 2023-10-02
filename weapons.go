package main

var WeaponsList = []Weapon{
	BrawlingChoke7Body,
	BrawlingChoke10Body,
	BrawlingChoke12Body,
	BrawlingChoke14Body,
	BrawlingStrike7Body,
	BrawlingStrike11Body,
	HeavyMelee,
	Body7MartialArt,
	Body11MartialArt,

	HeavyPistol,
	VeryHeavyPistol,
	ExoticHeavyPistol,
	GrenadeLauncher,
	RocketLauncher,
	AssaultRifle,
	SniperRifle,
	SMG,
	HeavySMG,
	Shotgun,
}

type Weapon struct {
	Name                 string            `json:"name"`
	DamageDice           int               `json:"damage_dice"`
	ClipSize             int               `json:"clip_size"`
	RangeBandDVs         map[RangeBand]int `json:"range_band_dv_s"`
	AutofireRangeBandDVs map[RangeBand]int `json:"autofire_range_band_dv_s"`
	RequiredHands        int               `json:"required_hands"`
	RateOfFire           int               `json:"rate_of_fire"`
	Skill                Skill             `json:"skill"`
	Ranged               bool              `json:"ranged"`
	AutofireDice         int               `json:"autofire_dice"`
	CanAutofire          bool              `json:"can_autofire"`
	CanAimedShot         bool              `json:"can_aimed_shot"`
	HalvesArmor          bool
	ShouldChoke          bool
	ChokeDamage          int
}

var HeavyMelee = Weapon{
	Name:          "Heavy Melee",
	DamageDice:    3,
	RangeBandDVs:  nil,
	RequiredHands: 2,
	RateOfFire:    2,
	Skill:         Melee,
	Ranged:        false,
	CanAutofire:   false,
	CanAimedShot:  true,
	AutofireDice:  0,
	HalvesArmor:   true,
}

var Body11MartialArt = Weapon{
	Name:          "11 Body Martial Arts",
	DamageDice:    4,
	RangeBandDVs:  nil,
	RequiredHands: 1,
	RateOfFire:    2,
	Skill:         MartialArts,
	Ranged:        false,
	CanAutofire:   false,
	CanAimedShot:  true,
	AutofireDice:  0,
	HalvesArmor:   true,
}

var Body7MartialArt = Weapon{
	Name:          "7 Body Martial Arts",
	DamageDice:    3,
	RangeBandDVs:  nil,
	RequiredHands: 1,
	RateOfFire:    2,
	Skill:         MartialArts,
	Ranged:        false,
	CanAutofire:   false,
	CanAimedShot:  true,
	AutofireDice:  0,
	HalvesArmor:   true,
}

var BrawlingStrike7Body = Weapon{
	Name:          "Brawling Strike 11 Body",
	DamageDice:    3,
	RangeBandDVs:  nil,
	RequiredHands: 1,
	RateOfFire:    2,
	Skill:         Brawling,
	Ranged:        false,
	CanAutofire:   false,
	CanAimedShot:  true,
	AutofireDice:  0,
	HalvesArmor:   false,
}

var BrawlingStrike11Body = Weapon{
	Name:          "Brawling Strike 11 Body",
	DamageDice:    4,
	RangeBandDVs:  nil,
	RequiredHands: 1,
	RateOfFire:    2,
	Skill:         Brawling,
	Ranged:        false,
	CanAutofire:   false,
	CanAimedShot:  true,
	AutofireDice:  0,
	HalvesArmor:   false,
}

var BrawlingChoke7Body = Weapon{
	Name:          "Brawling Choke 7 Body",
	DamageDice:    4,
	RangeBandDVs:  nil,
	RequiredHands: 1,
	RateOfFire:    1,
	Skill:         Brawling,
	Ranged:        false,
	CanAutofire:   false,
	CanAimedShot:  false,
	AutofireDice:  0,
	HalvesArmor:   false,
	ShouldChoke:   true,
	ChokeDamage:   7,
}

var BrawlingChoke10Body = Weapon{
	Name:          "Brawling Choke 10 Body",
	DamageDice:    4,
	RangeBandDVs:  nil,
	RequiredHands: 1,
	RateOfFire:    1,
	Skill:         Brawling,
	Ranged:        false,
	CanAutofire:   false,
	CanAimedShot:  false,
	AutofireDice:  0,
	HalvesArmor:   false,
	ShouldChoke:   true,
	ChokeDamage:   10,
}

var BrawlingChoke12Body = Weapon{
	Name:          "Brawling Choke 12 Body",
	DamageDice:    4,
	RangeBandDVs:  nil,
	RequiredHands: 1,
	RateOfFire:    1,
	Skill:         Brawling,
	Ranged:        false,
	CanAutofire:   false,
	CanAimedShot:  false,
	AutofireDice:  0,
	HalvesArmor:   false,
	ShouldChoke:   true,
	ChokeDamage:   12,
}

var BrawlingChoke14Body = Weapon{
	Name:          "Brawling Choke 14 Body",
	DamageDice:    4,
	RangeBandDVs:  nil,
	RequiredHands: 1,
	RateOfFire:    1,
	Skill:         Brawling,
	Ranged:        false,
	CanAutofire:   false,
	CanAimedShot:  false,
	AutofireDice:  0,
	HalvesArmor:   false,
	ShouldChoke:   true,
	ChokeDamage:   14,
}

var HeavyPistol = Weapon{
	Name:          "Heavy Pistol",
	DamageDice:    3,
	RangeBandDVs:  PistolRangeBands,
	RequiredHands: 1,
	RateOfFire:    2,
	Skill:         Handguns,
	Ranged:        true,
	CanAutofire:   false,
	CanAimedShot:  true,
	AutofireDice:  0,
	HalvesArmor:   false,
	ClipSize:      8,
}

var VeryHeavyPistol = Weapon{
	Name:          "Very Heavy Pistol",
	DamageDice:    4,
	RangeBandDVs:  PistolRangeBands,
	RequiredHands: 1,
	RateOfFire:    1,
	Skill:         Handguns,
	Ranged:        true,
	CanAutofire:   false,
	CanAimedShot:  true,
	AutofireDice:  0,
	HalvesArmor:   false,
	ClipSize:      8,
}

var ExoticHeavyPistol = Weapon{
	Name:          "Exotic Pistol",
	DamageDice:    5,
	RangeBandDVs:  PistolRangeBands,
	RequiredHands: 1,
	RateOfFire:    1,
	Skill:         Handguns,
	Ranged:        true,
	CanAutofire:   false,
	CanAimedShot:  true,
	AutofireDice:  0,
	HalvesArmor:   false,
	ClipSize:      8,
}

var GrenadeLauncher = Weapon{
	Name:          "Grenade Launcher",
	DamageDice:    6,
	RangeBandDVs:  GrenadeLauncherRangeBands,
	RequiredHands: 2,
	RateOfFire:    1,
	Skill:         HeavyWeapons,
	Ranged:        true,
	CanAutofire:   false,
	CanAimedShot:  false,
	AutofireDice:  0,
	HalvesArmor:   false,
	ClipSize:      2,
}

var RocketLauncher = Weapon{
	Name:          "Rocket Launcher",
	DamageDice:    8,
	RangeBandDVs:  GrenadeLauncherRangeBands,
	RequiredHands: 2,
	RateOfFire:    1,
	Skill:         HeavyWeapons,
	Ranged:        true,
	CanAutofire:   false,
	CanAimedShot:  false,
	AutofireDice:  0,
	HalvesArmor:   false,
	ClipSize:      1,
}

var AssaultRifle = Weapon{
	Name:                 "Assault Rifle",
	DamageDice:           5,
	RangeBandDVs:         AssaultRifleRangeBands,
	RequiredHands:        2,
	RateOfFire:           1,
	Skill:                ShoulderArms,
	Ranged:               true,
	CanAutofire:          true,
	CanAimedShot:         true,
	AutofireDice:         3,
	AutofireRangeBandDVs: AssaultRifleAutofireRangeBands,
	HalvesArmor:          false,
	ClipSize:             25,
}

var SniperRifle = Weapon{
	Name:                 "Sniper Rifle",
	DamageDice:           5,
	RangeBandDVs:         SniperRifleRangeBands,
	RequiredHands:        2,
	RateOfFire:           1,
	Skill:                ShoulderArms,
	Ranged:               true,
	CanAutofire:          false,
	CanAimedShot:         true,
	AutofireDice:         0,
	AutofireRangeBandDVs: nil,
	HalvesArmor:          false,
	ClipSize:             4,
}

var SMG = Weapon{
	Name:                 "SMG",
	DamageDice:           2,
	RangeBandDVs:         SMGRangeBands,
	RequiredHands:        2,
	RateOfFire:           1,
	Skill:                Handguns,
	Ranged:               true,
	CanAutofire:          true,
	CanAimedShot:         true,
	AutofireDice:         2,
	AutofireRangeBandDVs: SMGAutofireRangeBands,
	HalvesArmor:          false,
	ClipSize:             30,
}

var HeavySMG = Weapon{
	Name:                 "Heavy SMG",
	DamageDice:           3,
	RangeBandDVs:         SMGRangeBands,
	RequiredHands:        2,
	RateOfFire:           1,
	Skill:                Handguns,
	Ranged:               true,
	CanAutofire:          true,
	CanAimedShot:         true,
	AutofireDice:         2,
	AutofireRangeBandDVs: SMGAutofireRangeBands,
	HalvesArmor:          false,
	ClipSize:             40,
}

var Shotgun = Weapon{
	Name:          "Shotgun",
	DamageDice:    5,
	RangeBandDVs:  ShotgunRangeBands,
	RequiredHands: 2,
	RateOfFire:    1,
	Skill:         ShoulderArms,
	Ranged:        true,
	CanAutofire:   false,
	CanAimedShot:  true,
	AutofireDice:  0,
	HalvesArmor:   false,
	ClipSize:      4,
}

var PistolRangeBands = map[RangeBand]int{
	VeryClose: 13,
	Close:     15,
	Medium:    20,
	MediumFar: 25,
	Far:       30,
	VeryFar:   30,
}

var SMGRangeBands = map[RangeBand]int{
	VeryClose: 15,
	Close:     13,
	Medium:    15,
	MediumFar: 20,
	Far:       25,
	VeryFar:   25,
}

var ShotgunRangeBands = map[RangeBand]int{
	VeryClose: 13,
	Close:     15,
	Medium:    20,
	MediumFar: 25,
	Far:       30,
	VeryFar:   35,
}

var AssaultRifleRangeBands = map[RangeBand]int{
	VeryClose: 17,
	Close:     16,
	Medium:    15,
	MediumFar: 13,
	Far:       15,
	VeryFar:   20,
}

var SniperRifleRangeBands = map[RangeBand]int{
	VeryClose: 30,
	Close:     25,
	Medium:    25,
	MediumFar: 20,
	Far:       15,
	VeryFar:   16,
}

var GrenadeLauncherRangeBands = map[RangeBand]int{
	VeryClose: 16,
	Close:     15,
	Medium:    15,
	MediumFar: 17,
	Far:       20,
	VeryFar:   22,
}

var SMGAutofireRangeBands = map[RangeBand]int{
	VeryClose: 20,
	Close:     17,
	Medium:    20,
	MediumFar: 25,
	Far:       30,
	VeryFar:   35,
}

var AssaultRifleAutofireRangeBands = map[RangeBand]int{
	VeryClose: 22,
	Close:     20,
	Medium:    17,
	MediumFar: 20,
	Far:       25,
	VeryFar:   30,
}
