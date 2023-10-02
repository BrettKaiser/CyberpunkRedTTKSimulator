package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

type AttacksResult struct {
	TotalDamage      int
	ArmorAblated     int
	CriticalInjuries int
	AttacksDone      int
}

type AttacksParams struct {
	CombatScenario
	AttackerModifiers []Modifier
	DefenderModifiers []Modifier
}

type CombatScenario struct {
	DebugLogs                 bool
	Ammunition                AmmunitionType
	AttackType                AttackType
	Attacker                  CurrentCharacter
	Defender                  CurrentCharacter
	RangeBand                 RangeBand
	TotalAttacks              int
	DamageDone                int
	NumberOfRounds            int
	NotApplicable             bool
	NumberOfReloads           int
	RoundsSpentRunning        int
	InGrapple                 bool
	DistanceBetweenCharacters int
}

func NewCombatScenario(params ScenarioParams) CombatScenario {
	rand.Seed(time.Now().UnixNano()) // only run once
	return CombatScenario{
		DebugLogs:  params.DebugLogs,
		Ammunition: params.Ammunition,
		AttackType: params.AttackType,
		Attacker: CurrentCharacter{
			Character: params.Attacker,
			CurrentHP: params.Attacker.MaxHP,
			CurrentWeapon: CurrentWeapon{
				Weapon:          params.Attacker.Weapon,
				CurrentClipSize: params.Attacker.Weapon.ClipSize,
			},
		},
		Defender: CurrentCharacter{
			Character: params.Defender,
			CurrentHP: params.Defender.MaxHP,
			CurrentWeapon: CurrentWeapon{
				Weapon:          params.Defender.Weapon,
				CurrentClipSize: params.Defender.Weapon.ClipSize,
			},
		},
		InGrapple:                 false,
		RangeBand:                 params.RangeBand,
		DistanceBetweenCharacters: params.RangeBand.MaxDistance,

		TotalAttacks:    0,
		DamageDone:      0,
		NumberOfRounds:  0,
		NumberOfReloads: 0,
	}
}

func (scenario *CombatScenario) inMeleeRange() bool {
	if scenario.DistanceBetweenCharacters > VeryClose.MaxDistance {
		return false
	}

	return true
}

func (scenario *CombatScenario) Execute() CombatScenario {
	for i := 0; i < MaximumRounds; i++ {
		if scenario.Defender.CurrentHP <= 0 {
			break
		}

		scenario.NumberOfRounds++

		// Move into melee range
		if !scenario.Attacker.CurrentWeapon.Weapon.Ranged && !scenario.inMeleeRange() {
			scenario.moveCharacterCloser(scenario.Attacker)
			if !scenario.inMeleeRange() {
				scenario.moveCharacterCloser(scenario.Attacker)
				scenario.RoundsSpentRunning++
				continue
			}
		}

		// Attempt to grapple if not grappled
		if !scenario.InGrapple && scenario.Attacker.Weapon.ShouldChoke {
			if scenario.grappledSuccessfully() {
				scenario.InGrapple = true
				continue
			}
		}

		// Attack
		attacksResult := scenario.CalculateAttacks()

		if scenario.DebugLogs {
			DisplayRound(attacksResult, i+1)
			scenario.DisplayResult()
		}
	}

	return *scenario
}

func (scenario *CombatScenario) CalculateAttacks() AttacksResult {
	weapon := scenario.Attacker.CurrentWeapon
	attribute := 0
	skill := scenario.GetAttackSkill()

	if weapon.Ranged {
		attribute = scenario.Attacker.Reflexes
	} else {
		attribute = scenario.Attacker.Dexterity
	}

	damageDoneThisRound := 0
	armorAblatedThisRound := 0
	criticalInjuriesThisRound := 0
	attacksDoneThisRound := 0

	numberOfAttacks := weapon.RateOfFire
	if scenario.AttackType == Autofire || scenario.AttackType == Headshot || scenario.Attacker.CurrentWeapon.ShouldChoke {
		numberOfAttacks = 1
	}

	for i := 0; i < numberOfAttacks; i++ {
		if scenario.Defender.CurrentHP <= 0 {
			break
		}

		if scenario.Attacker.mustReload(scenario.AttackType) {
			scenario.Reload(scenario.Attacker.CurrentWeapon)
			break
		}

		attacksDoneThisRound++
		scenario.Attacker.CurrentWeapon.subtractAmmo(scenario.AttackType)

		if scenario.Attacker.CurrentWeapon.ShouldChoke && scenario.InGrapple {
			scenario.Defender.CurrentHP -= scenario.Attacker.CurrentWeapon.ChokeDamage
			damageDoneThisRound += scenario.Attacker.CurrentWeapon.ChokeDamage
			continue
		}

		dv := scenario.GetDV()

		hitParams := HitParams{
			Attribute:      attribute,
			Skill:          skill,
			AttackModifier: scenario.GetAttackModifiers(),
			DV:             dv,
		}

		toHitResult := CalculateHit(hitParams)
		damageResult := scenario.CalculateDamage(DamageParams{
			ToHitResult:    toHitResult,
			AttackType:     scenario.AttackType,
			AmmunitionType: scenario.Ammunition,
			HalvesArmor:    scenario.Attacker.Weapon.HalvesArmor,
		})

		damageDoneThisRound += damageResult.TotalDamage
		armorAblatedThisRound += damageResult.ArmorAblated
		criticalInjuriesThisRound += damageResult.NumberOfCriticalInjuries
	}

	scenario.TotalAttacks += attacksDoneThisRound

	return AttacksResult{
		TotalDamage:      damageDoneThisRound,
		ArmorAblated:     armorAblatedThisRound,
		CriticalInjuries: criticalInjuriesThisRound,
		AttacksDone:      attacksDoneThisRound,
	}
}

func (scenario *CombatScenario) Reload(currentWeapon CurrentWeapon) {
	scenario.NumberOfReloads++
	scenario.Attacker.CurrentWeapon.reload()
}

func (weapon *CurrentWeapon) reload() {
	weapon.CurrentClipSize = weapon.ClipSize
}

func (scenario *CombatScenario) grappledSuccessfully() bool {
	attackerRoll := GetD10CheckResult(scenario.Attacker.Dexterity, scenario.Attacker.Brawling, 0)
	defenderRoll := GetD10CheckResult(scenario.Defender.Dexterity, scenario.Defender.Brawling, 0)

	if attackerRoll > defenderRoll {
		return true
	}

	return false
}

func (scenario *CombatScenario) moveCharacterCloser(character CurrentCharacter) {
	scenario.DistanceBetweenCharacters -= character.Movement

	if scenario.DistanceBetweenCharacters < 0 {
		scenario.DistanceBetweenCharacters = 0
	}
}

func (weapon *CurrentWeapon) subtractAmmo(attackType AttackType) {
	if !weapon.Ranged {
		return
	}

	switch attackType {
	case Autofire:
		weapon.CurrentClipSize -= 10
	case SingleShot:
		weapon.CurrentClipSize -= 1
	case Headshot:
		weapon.CurrentClipSize -= 1
	}
}

func (character *CurrentCharacter) mustReload(attackType AttackType) bool {
	if !character.CurrentWeapon.Ranged {
		return false
	}

	switch attackType {
	case Autofire:
		if character.CurrentWeapon.CurrentClipSize < 10 {
			return true
		}
	case SingleShot:
		if character.CurrentWeapon.CurrentClipSize < 1 {
			return true
		}
	case Headshot:
		if character.CurrentWeapon.CurrentClipSize < 1 {
			return true
		}
	}

	return false
}

func (scenario *CombatScenario) GetAttackSkill() int {
	skill := 0
	switch scenario.Attacker.Weapon.Skill {
	case Handguns:
		skill = scenario.Attacker.Handguns
	case ShoulderArms:
		skill = scenario.Attacker.ShoulderArms
	case HeavyWeapons:
		skill = scenario.Attacker.HeavyWeapons
	case AutoFire:
		skill = scenario.Attacker.AutoFire
	case Brawling:
		skill = scenario.Attacker.Brawling
	case Melee:
		skill = scenario.Attacker.Melee
	case MartialArts:
		skill = scenario.Attacker.MartialArts
	default:
		panic(fmt.Sprintf("missing weapon skill %s", scenario.Attacker.Weapon.Skill))
	}

	return skill
}

func (scenario *CombatScenario) CalculateDamage(params DamageParams) DamageResult {
	damageResult := DamageResult{
		TotalDamage:              0,
		NumberOfCriticalInjuries: 0,
		ArmorAblated:             0,
	}

	if !params.ToHitResult.Hit {
		return damageResult
	}

	numberOfDiceToRoll := GetDamageDice(scenario.Attacker.Weapon, params.AttackType)
	diceResult := RollD6s(numberOfDiceToRoll)
	if diceResult.NumberOf6s >= 2 {
		scenario.Defender.CurrentHP = scenario.Defender.CurrentHP - 5
		damageResult.NumberOfCriticalInjuries++
		damageResult.TotalDamage += 5
	}

	armorValue := scenario.Defender.ArmorValue
	if armorValue > 0 && params.HalvesArmor {
		armorValue = int(math.Ceil(float64(armorValue) / 2))
	}

	damageDifference := diceResult.Total - armorValue

	if damageDifference > 0 {
		damageDone := GetDamageDone(params, damageDifference)
		damageResult.TotalDamage += damageDone
		scenario.DamageDone += damageDone

		armorAblated := GetArmorAblated(params)
		damageResult.ArmorAblated += armorAblated

		scenario.Defender.CurrentHP = scenario.Defender.CurrentHP - damageDone

		if scenario.Defender.ArmorValue >= armorAblated {
			scenario.Defender.ArmorValue = scenario.Defender.ArmorValue - armorAblated
		} else {
			scenario.Defender.ArmorValue = 0
		}
	}

	return damageResult
}

func GetDamageDone(params DamageParams, damageDifference int) int {
	damageDone := 0
	switch params.AttackType {
	case Headshot:
		damageDone = damageDifference * 2
	case Autofire:
		damageDone = damageDifference * params.ToHitResult.Difference
	case SingleShot:
		damageDone = damageDifference
	}

	return damageDone
}

func GetArmorAblated(params DamageParams) int {
	armorAblated := 1

	if params.AmmunitionType == ArmorPiercing {
		armorAblated++
	}

	return armorAblated
}

func GetDamageDice(weapon Weapon, attackType AttackType) int {
	if attackType == Autofire {
		if !weapon.CanAutofire {
			panic(fmt.Sprintf("trying to autofire a weapon that cant autofire: '%s'", weapon.Name))
		}

		return weapon.AutofireDice
	}

	return weapon.DamageDice
}

type DamageResult struct {
	TotalDamage              int
	NumberOfCriticalInjuries int
	ArmorAblated             int
}

type ToHitResult struct {
	Hit         bool
	Difference  int
	TotalRolled int
}

var AttackTypes = []AttackType{
	SingleShot,
	Headshot,
	Autofire,
}

type AttackType string

const (
	Autofire   AttackType = "Autofire"
	SingleShot AttackType = "SingleShot"
	Headshot   AttackType = "Headshot"
)

func GetD10CheckResult(attribute, skill, totalModifiers int) int {
	base := attribute + skill + totalModifiers
	rollResult := RollD10s(1)
	totalValue := base + rollResult.Total

	if rollResult.NumberOf10s > 0 {
		critResult := RollD10s(1)
		totalValue += critResult.Total
	}

	if rollResult.NumberOf1s > 0 {
		critResult := RollD10s(1)
		totalValue -= critResult.Total
	}

	return totalValue
}

func (scenario *CombatScenario) GetAttackModifiers() int {
	totalModifier := 0
	if scenario.AttackType == Headshot {
		totalModifier -= 8
		totalModifier -= scenario.Attacker.AimedShotBonus
	}

	if scenario.Attacker.CurrentWeapon.Ranged && scenario.Attacker.HasSmartLink {
		totalModifier += 1
	}

	totalModifier += scenario.Attacker.CombatAwareness

	return totalModifier + GetTotalModifiers(scenario.Attacker.AttackModifiers)
}

func DisplayRound(attacksResult AttacksResult, roundNumber int) {
	fmt.Printf("- Round Number: %d\n", roundNumber)
	fmt.Printf("Attacks Done: %d\n", attacksResult.AttacksDone)
	fmt.Printf("Damage Done: %d\n", attacksResult.TotalDamage)
	fmt.Printf("Critical Injuries: %d\n", attacksResult.CriticalInjuries)
	fmt.Printf("Armor Ablated: %d\n", attacksResult.ArmorAblated)
}

func (scenario *CombatScenario) DisplayResult() {
	fmt.Printf("\nTotal Rounds: %d\n", scenario.NumberOfRounds)
	fmt.Printf("Attacks to Kill: %d\n\n", scenario.TotalAttacks)
	fmt.Printf("Damage Done: %d\n", scenario.DamageDone)
}

func (scenario *CombatScenario) GetDV() int {
	weapon := scenario.Attacker.CurrentWeapon.Weapon
	dv := 0
	dvModifier := scenario.GetDefenderDVModifier()
	if weapon.Ranged {
		dv = weapon.RangeBandDVs[scenario.RangeBand]
		if scenario.AttackType == Autofire {
			dv = weapon.AutofireRangeBandDVs[scenario.RangeBand]
		}

		if scenario.Defender.Reflexes >= 8 {
			if scenario.ShouldDodge(dv, dvModifier) {
				dv = GetD10CheckResult(scenario.Defender.Evasion, scenario.Defender.Dexterity, dvModifier)
			}
		}
	} else {
		dv = GetD10CheckResult(scenario.Defender.Evasion, scenario.Defender.Dexterity, dvModifier)
	}

	return dv
}

func (scenario *CombatScenario) GetDefenderDVModifier() int {
	totalModifiers := 0
	totalModifiers += scenario.Defender.ArmorPenalty
	totalModifiers += scenario.GetDefenderWoundPenalty()
	return totalModifiers
}

func (scenario *CombatScenario) GetDefenderWoundPenalty() int {
	halfHP := 0
	if scenario.Defender.CurrentHP > 0 {
		halfHP = int(math.Ceil((float64(scenario.Defender.MaxHP) / 2)))
	}

	if scenario.Defender.CurrentHP <= halfHP {
		return -2
	}

	return 0
}

func GetTotalModifiers(modifiers []Modifier) int {
	total := 0
	for _, modifier := range modifiers {
		total += modifier.Value
	}

	return total
}

func (scenario *CombatScenario) ShouldDodge(rangeDV int, totalModifiers int) bool {
	// baseDodge := scenario.Defender.Evasion + scenario.Defender.Reflexes + 1 - totalModifiers
	// if baseDodge > rangeDV {
	// 	return true
	// }
	//
	// if baseDodge+5 > rangeDV {
	// 	return true
	// }
	//
	// return false

	if scenario.Defender.ShouldDodge && scenario.Defender.Reflexes >= 8 {
		maxDodge := scenario.Defender.Evasion + scenario.Defender.Reflexes + 10 - totalModifiers
		if maxDodge < rangeDV {
			return false
		}

		return true
	}

	return false
}

type HitParams struct {
	Attribute      int
	Skill          int
	AttackModifier int
	DV             int
}

type RoundResult struct {
	NumberOfAttacksRolled int
	HPDamage              int
}

func CalculateHit(params HitParams) ToHitResult {
	d10Result := GetD10CheckResult(params.Attribute, params.Skill, params.AttackModifier)
	difference := d10Result - params.DV

	return ToHitResult{
		Hit:         difference > 0,
		Difference:  difference,
		TotalRolled: d10Result,
	}
}

type DamageParams struct {
	ToHitResult    ToHitResult
	AttackType     AttackType
	AmmunitionType AmmunitionType
	HalvesArmor    bool
}
