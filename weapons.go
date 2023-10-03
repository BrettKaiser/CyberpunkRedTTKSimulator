package main

var WeaponsList = []Weapon{
	BrawlingChoke7Body,
	BrawlingChoke10Body,
	BrawlingChoke12Body,
	BrawlingChoke14Body,
	BrawlingStrike7Body,
	BrawlingStrike11Body,

	HeavyMelee,
	VeryHeavyMelee,
	KendachiMonoThree,
	KendachiMonoWakizashi,
	RostovicKleaver,

	Body7MartialArt,
	Body11MartialArt,

	MediumPistol,
	HeavyPistol,
	VeryHeavyPistol,
	ExoticHeavyPistol,
	MilitechPerseus,

	GrenadeLauncher,
	RocketLauncher,

	AssaultRifle,
	RhinemetallRailgun,

	SniperRifle,

	TsunamiArmsHelix,

	SMG,
	HeavySMG,

	Shotgun,
}

var AmmunitionTypes = map[string]Ammunition{
	"basic":          BasicAmmunition,
	"armor_piercing": ArmorPiercingAmmunition,
}

type Ammunition struct {
	Name string
	Cost int `json:"cost"`
}

var BasicAmmunition = Ammunition{
	Name: "Basic",
	Cost: 10,
}

var ArmorPiercingAmmunition = Ammunition{
	Name: "Armor Piercing",
	Cost: 100,
}

type AmmunitionType string

const (
	Basic         AmmunitionType = "basic"
	ArmorPiercing AmmunitionType = "armor_piercing"
	Incendiary    AmmunitionType = "incendiary"
)

type Weapon struct {
	Name                 string            `json:"name"`
	DamageDice           int               `json:"damage_dice"`
	ClipSize             int               `json:"clip_size"`
	ExtendedClipSize     int               `json:"extended_clip_size"`
	DrumClipSize         int               `json:"drum_clip_size"`
	RangeBandDVs         map[RangeBand]int `json:"range_band_dv_s"`
	AutofireRangeBandDVs map[RangeBand]int `json:"autofire_range_band_dv_s"`
	RequiredHands        int               `json:"required_hands"`
	RateOfFire           int               `json:"rate_of_fire"`
	Skill                Skill             `json:"skill"`
	Ranged               bool              `json:"ranged"`
	AutofireMax          int               `json:"autofire_max"`
	CanAutofire          bool              `json:"can_autofire"`
	CanAimedShot         bool              `json:"can_aimed_shot"`
	CannotSingleShot     bool              `json:"cannot_single_shot"`
	HalvesArmor          bool
	ShouldChoke          bool
	ChokeDamage          int
	Unarmed              bool
	Cost                 int
	Explosive            bool
	AutofireSpent        int
	IgnoresArmorUnder    int
	TurnsToReload        int
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
	AutofireMax:   0,
	HalvesArmor:   true,
	Cost:          100,
}

var VeryHeavyMelee = Weapon{
	Name:          "Very Heavy Melee",
	DamageDice:    4,
	RangeBandDVs:  nil,
	RequiredHands: 2,
	RateOfFire:    1,
	Skill:         Melee,
	Ranged:        false,
	CanAutofire:   false,
	CanAimedShot:  true,
	AutofireMax:   0,
	HalvesArmor:   true,
	Cost:          500,
}

var RostovicKleaver = Weapon{
	Name:              "Rostovic Kleaver",
	DamageDice:        5,
	RangeBandDVs:      nil,
	RequiredHands:     2,
	RateOfFire:        1,
	Skill:             Melee,
	Ranged:            false,
	CanAutofire:       false,
	CanAimedShot:      true,
	AutofireMax:       0,
	HalvesArmor:       true,
	Cost:              500,
	IgnoresArmorUnder: 11,
}

var KendachiMonoWakizashi = Weapon{
	Name:              "Kendachi Mono-Wakizashi",
	DamageDice:        3,
	RangeBandDVs:      nil,
	RequiredHands:     2,
	RateOfFire:        2,
	Skill:             Melee,
	Ranged:            false,
	CanAutofire:       false,
	CanAimedShot:      true,
	AutofireMax:       0,
	HalvesArmor:       true,
	Cost:              100,
	IgnoresArmorUnder: 7,
}

var KendachiMonoThree = Weapon{
	Name:              "Kendachi Mono Three",
	DamageDice:        4,
	RangeBandDVs:      nil,
	RequiredHands:     2,
	RateOfFire:        1,
	Skill:             Melee,
	Ranged:            false,
	CanAutofire:       false,
	CanAimedShot:      true,
	AutofireMax:       0,
	HalvesArmor:       true,
	Cost:              5000,
	IgnoresArmorUnder: 11,
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
	AutofireMax:   0,
	HalvesArmor:   true,
	Unarmed:       true,
	Cost:          2000,
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
	AutofireMax:   0,
	HalvesArmor:   true,
	Unarmed:       true,
	Cost:          0,
}

var BrawlingStrike7Body = Weapon{
	Name:          "Brawling Strike 7 Body",
	DamageDice:    3,
	RangeBandDVs:  nil,
	RequiredHands: 1,
	RateOfFire:    2,
	Skill:         Brawling,
	Ranged:        false,
	CanAutofire:   false,
	CanAimedShot:  true,
	AutofireMax:   0,
	HalvesArmor:   false,
	Unarmed:       true,
	Cost:          0,
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
	AutofireMax:   0,
	HalvesArmor:   false,
	Unarmed:       true,
	Cost:          2000,
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
	AutofireMax:   0,
	HalvesArmor:   false,
	ShouldChoke:   true,
	ChokeDamage:   7,
	Unarmed:       true,
	Cost:          0,
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
	AutofireMax:   0,
	HalvesArmor:   false,
	ShouldChoke:   true,
	ChokeDamage:   10,
	Unarmed:       true,
	Cost:          1000,
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
	AutofireMax:   0,
	HalvesArmor:   false,
	ShouldChoke:   true,
	ChokeDamage:   12,
	Unarmed:       true,
	Cost:          2000,
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
	AutofireMax:   0,
	HalvesArmor:   false,
	ShouldChoke:   true,
	ChokeDamage:   14,
	Unarmed:       true,
	Cost:          7000,
}

var MediumPistol = Weapon{
	Name:             "Medium Pistol",
	DamageDice:       2,
	RangeBandDVs:     PistolRangeBands,
	RequiredHands:    1,
	RateOfFire:       2,
	Skill:            Handguns,
	Ranged:           true,
	CanAutofire:      false,
	CanAimedShot:     true,
	AutofireMax:      0,
	HalvesArmor:      false,
	ClipSize:         12,
	ExtendedClipSize: 18,
	DrumClipSize:     36,
	Cost:             50,
}

var HeavyPistol = Weapon{
	Name:             "Heavy Pistol",
	DamageDice:       3,
	RangeBandDVs:     PistolRangeBands,
	RequiredHands:    1,
	RateOfFire:       2,
	Skill:            Handguns,
	Ranged:           true,
	CanAutofire:      false,
	CanAimedShot:     true,
	AutofireMax:      0,
	HalvesArmor:      false,
	ClipSize:         8,
	ExtendedClipSize: 14,
	DrumClipSize:     28,
	Cost:             100,
}

var VeryHeavyPistol = Weapon{
	Name:             "Very Heavy Pistol",
	DamageDice:       4,
	RangeBandDVs:     PistolRangeBands,
	RequiredHands:    1,
	RateOfFire:       1,
	Skill:            Handguns,
	Ranged:           true,
	CanAutofire:      false,
	CanAimedShot:     true,
	AutofireMax:      0,
	HalvesArmor:      false,
	ClipSize:         8,
	ExtendedClipSize: 14,
	DrumClipSize:     28,
	Cost:             100,
}

var MilitechPerseus = Weapon{
	Name:             "Militech Perseus",
	DamageDice:       4,
	RangeBandDVs:     PistolRangeBands,
	RequiredHands:    1,
	RateOfFire:       1,
	Skill:            Handguns,
	Ranged:           true,
	CanAutofire:      false,
	CanAimedShot:     true,
	AutofireMax:      0,
	HalvesArmor:      false,
	ClipSize:         11,
	ExtendedClipSize: 11,
	DrumClipSize:     11,
	Cost:             5000,
}

var ExoticHeavyPistol = Weapon{
	Name:             "Exotic Pistol",
	DamageDice:       5,
	RangeBandDVs:     PistolRangeBands,
	RequiredHands:    1,
	RateOfFire:       1,
	Skill:            Handguns,
	Ranged:           true,
	CanAutofire:      false,
	CanAimedShot:     true,
	AutofireMax:      0,
	HalvesArmor:      false,
	ClipSize:         8,
	ExtendedClipSize: 8,
	DrumClipSize:     8,
	Cost:             10000,
}

var GrenadeLauncher = Weapon{
	Name:             "Grenade Launcher",
	DamageDice:       6,
	RangeBandDVs:     GrenadeLauncherRangeBands,
	RequiredHands:    2,
	RateOfFire:       1,
	Skill:            HeavyWeapons,
	Ranged:           true,
	CanAutofire:      false,
	CanAimedShot:     false,
	AutofireMax:      0,
	HalvesArmor:      false,
	ClipSize:         2,
	ExtendedClipSize: 4,
	DrumClipSize:     6,
	Cost:             500,
	Explosive:        true,
}

var RocketLauncher = Weapon{
	Name:             "Rocket Launcher",
	DamageDice:       8,
	RangeBandDVs:     GrenadeLauncherRangeBands,
	RequiredHands:    2,
	RateOfFire:       1,
	Skill:            HeavyWeapons,
	Ranged:           true,
	CanAutofire:      false,
	CanAimedShot:     false,
	AutofireMax:      0,
	HalvesArmor:      false,
	ClipSize:         1,
	ExtendedClipSize: 2,
	DrumClipSize:     3,
	Cost:             500,
	Explosive:        true,
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
	AutofireMax:          4,
	AutofireRangeBandDVs: AssaultRifleAutofireRangeBands,
	HalvesArmor:          false,
	ClipSize:             25,
	ExtendedClipSize:     35,
	DrumClipSize:         45,
	Cost:                 500,
}

var RhinemetallRailgun = Weapon{
	Name:              "Rhinemetall Railgun",
	DamageDice:        5,
	RangeBandDVs:      AssaultRifleRangeBands,
	RequiredHands:     2,
	RateOfFire:        1,
	Skill:             HeavyWeapons,
	Ranged:            true,
	CanAutofire:       false,
	CanAimedShot:      false,
	AutofireMax:       0,
	HalvesArmor:       false,
	ClipSize:          4,
	ExtendedClipSize:  4,
	DrumClipSize:      4,
	Cost:              6000,
	TurnsToReload:     2,
	IgnoresArmorUnder: 11,
}

var TsunamiArmsHelix = Weapon{
	Name:                 "Tsunami Arms Helix",
	DamageDice:           5,
	RangeBandDVs:         AssaultRifleRangeBands,
	RequiredHands:        2,
	RateOfFire:           1,
	Skill:                HeavyWeapons,
	Ranged:               true,
	CanAutofire:          true,
	CanAimedShot:         false,
	CannotSingleShot:     true,
	AutofireMax:          5,
	AutofireRangeBandDVs: AssaultRifleAutofireRangeBands,
	HalvesArmor:          false,
	AutofireSpent:        20,
	ClipSize:             40,
	ExtendedClipSize:     40,
	DrumClipSize:         40,
	Cost:                 6000,
	TurnsToReload:        2,
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
	AutofireMax:          0,
	AutofireRangeBandDVs: nil,
	HalvesArmor:          false,
	ClipSize:             4,
	ExtendedClipSize:     8,
	DrumClipSize:         12,
	Cost:                 500,
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
	AutofireMax:          3,
	AutofireRangeBandDVs: SMGAutofireRangeBands,
	HalvesArmor:          false,
	ClipSize:             30,
	ExtendedClipSize:     40,
	DrumClipSize:         50,
	Cost:                 100,
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
	AutofireMax:          3,
	AutofireRangeBandDVs: SMGAutofireRangeBands,
	HalvesArmor:          false,
	ClipSize:             40,
	ExtendedClipSize:     50,
	DrumClipSize:         60,
	Cost:                 100,
}

var Shotgun = Weapon{
	Name:             "Shotgun",
	DamageDice:       5,
	RangeBandDVs:     ShotgunRangeBands,
	RequiredHands:    2,
	RateOfFire:       1,
	Skill:            ShoulderArms,
	Ranged:           true,
	CanAutofire:      false,
	CanAimedShot:     true,
	AutofireMax:      0,
	HalvesArmor:      false,
	ClipSize:         4,
	ExtendedClipSize: 8,
	DrumClipSize:     16,
	Cost:             500,
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
