package main

type AttackParams struct {
	Attribute      int    `json:"attribute"`
	Skill          int    `json:"skill"`
	Modifiers      []int  `json:"modifiers"`
	DV             int    `json:"dv"`
	DamageDice     int    `json:"damage_dice"`
	ClipSize       int    `json:"clip_size"`
	AmmunitionType string `json:"ammunition_type"`
}

type Character struct {
	CharacterStats
	Weapon       Weapon `json:"weapon"`
	ArmorValue   int    `json:"armor"`
	ArmorPenalty int    `json:"armor_penalty"`
	ShouldDodge  bool
}

type CharacterStats struct {
	Name            string     `json:"name"`
	MaxHP           int        `json:"hp"`
	Reflexes        int        `json:"reflexes"`
	Dexterity       int        `json:"dexterity"`
	Movement        int        `json:"movement"`
	Evasion         int        `json:"evasion"`
	Brawling        int        `json:"brawling"`
	Handguns        int        `json:"handguns"`
	ShoulderArms    int        `json:"shoulder_arms"`
	HeavyWeapons    int        `json:"heavy_weapons"`
	MartialArts     int        `json:"martial_arts"`
	Melee           int        `json:"melee"`
	AutoFire        int        `json:"auto_fire"`
	AttackModifiers []Modifier `json:"attack_modifiers"`
}

type CurrentCharacter struct {
	Character
	CurrentHP     int           `json:"current_hp"`
	Modifiers     []Modifier    `json:"modifiers"`
	CurrentWeapon CurrentWeapon `json:"current_weapon"`
}

type Modifier struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type RangeBand struct {
	MinDistance int    `json:"min_distance"`
	MaxDistance int    `json:"max_distance"`
	Name        string `json:"name"`
}

type CurrentWeapon struct {
	Weapon
	CurrentClipSize int `json:"current_clip_size"`
}

type RangeBandDV struct {
	RangeBand RangeBand `json:"range_band"`
	DV        int       `json:"dv"`
}

type Skill string

const (
	Brawling     Skill = "brawling"
	Handguns     Skill = "handguns"
	ShoulderArms Skill = "shoulder_arms"
	HeavyWeapons Skill = "heavy_weapons"
	AutoFire     Skill = "auto_fire"
	Melee        Skill = "melee"
	MartialArts  Skill = "martial_arts"
)
