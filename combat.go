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
	DebugLogs      bool
	AmmunitionType AmmunitionType
	// AttackType                AttackType
	Attacker                  CurrentCharacter
	Defender                  CurrentCharacter
	Bodyguard                 CurrentCharacter
	RangeBand                 RangeBand
	TotalAttacks              int
	DamageDone                int
	NumberOfRounds            int
	NotApplicable             bool
	NumberOfReloads           int
	RoundsSpentRunning        int
	InGrapple                 bool
	DistanceBetweenCharacters int
	ScenarioCost              int
	SetupCost                 int
	ArmorJuryRigsRemaining    int
	ShieldJuryRigsRemaining   int
}

func NewCombatScenario(params ScenarioParams) CombatScenario {
	// Create a new random seed to make this scenario different
	rand.Seed(time.Now().UnixNano())

	popupShieldHP := 0
	if params.Defender.HasPopupShield {
		popupShieldHP = 10
	}

	return CombatScenario{
		DebugLogs:      params.DebugLogs,
		AmmunitionType: params.Ammunition,
		Attacker: CurrentCharacter{
			Character: params.Attacker,
			CurrentHP: params.Attacker.MaxHP,
			CurrentWeapon: CurrentWeapon{
				Weapon:          params.Attacker.Weapon,
				CurrentClipSize: params.Attacker.getClipSize(),
			},
			CurrentSP:  params.Attacker.ArmorValue,
			AttackType: params.AttackType,
		},
		Defender: CurrentCharacter{
			Character: params.Defender,
			CurrentHP: params.Defender.MaxHP,
			CurrentWeapon: CurrentWeapon{
				Weapon:          params.Defender.Weapon,
				CurrentClipSize: params.Defender.getClipSize(),
			},
			CurrentSP:     params.Defender.ArmorValue,
			PopupShieldHP: popupShieldHP,
			AttackType:    SingleShot,
		},
		Bodyguard: CurrentCharacter{
			Character: params.Bodyguard,
			CurrentHP: params.Bodyguard.MaxHP,
			CurrentWeapon: CurrentWeapon{
				Weapon:          params.Bodyguard.Weapon,
				CurrentClipSize: params.Bodyguard.getClipSize(),
			},
			CurrentSP:     params.Bodyguard.ArmorValue,
			PopupShieldHP: popupShieldHP,
			AttackType:    SingleShot,
		},
		InGrapple:                 false,
		RangeBand:                 params.RangeBand,
		DistanceBetweenCharacters: params.RangeBand.MaxDistance,

		TotalAttacks:            0,
		DamageDone:              0,
		NumberOfRounds:          0,
		NumberOfReloads:         0,
		ArmorJuryRigsRemaining:  1,
		ShieldJuryRigsRemaining: 1,
	}
}

func (scenario *CombatScenario) Execute() CombatScenario {
	for i := 0; i < MAXIMUMROUNDS; i++ {
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
			if scenario.makeGrappleCheckAgainstEnemy(&scenario.Attacker) {
				scenario.InGrapple = true
				continue
			}
		}

		if scenario.Defender.IsTech {
			if scenario.Defender.HasPopupShield && scenario.ShieldJuryRigsRemaining > 0 {
				if scenario.Defender.PopupShieldHP <= 0 {
					scenario.Defender.PopupShieldHP = 10
					scenario.ShieldJuryRigsRemaining -= 1
				}
			}

			if scenario.ArmorJuryRigsRemaining > 0 {
				halfArmor := int(math.Ceil(float64(scenario.Defender.ArmorValue) / 2))
				if scenario.Defender.CurrentSP < halfArmor+2 {
					// Jury Rig your armor
					scenario.Defender.CurrentSP = scenario.Defender.ArmorValue
					scenario.ArmorJuryRigsRemaining -= 1
				}
			}
		}

		// Attack
		attacksResult := scenario.CalculateAttacks(&scenario.Attacker)

		// Bodyguard Attack
		if scenario.Attacker.HasBodyguardTeammate {
			scenario.CalculateAttacks(&scenario.Bodyguard)
		}

		if scenario.DebugLogs {
			DisplayRound(attacksResult, i+1)
			scenario.DisplayResult()
		}
	}

	scenario.setCost()

	return *scenario
}

func (scenario *CombatScenario) CalculateAttacks(attacker *CurrentCharacter) AttacksResult {
	weapon := &attacker.CurrentWeapon
	attribute := 0
	skill := attacker.GetAttackSkill()

	if weapon.Ranged {
		attribute = attacker.Reflexes
	} else {
		attribute = attacker.Dexterity
	}

	damageDoneThisRound := 0
	armorAblatedThisRound := 0
	criticalInjuriesThisRound := 0
	attacksDoneThisRound := 0

	numberOfAttacks := weapon.RateOfFire
	if attacker.AttackType == Autofire || attacker.AttackType == Headshot || attacker.CurrentWeapon.ShouldChoke {
		numberOfAttacks = 1
	}

	if attacker.CurrentWeapon.Weapon.Name == MilitechPerseus.Name {
		if scenario.NumberOfRounds%2 == 0 {
			numberOfAttacks = 2
		}
	}

	for i := 0; i < numberOfAttacks; i++ {
		if scenario.Defender.CurrentHP <= 0 {
			break
		}

		mustReload := attacker.mustReload(attacker.AttackType)
		if mustReload {
			scenario.CalculateReloadFor(attacker)
			break
		}

		attacksDoneThisRound++

		if attacker.CurrentWeapon.ShouldChoke && scenario.InGrapple {
			scenario.Defender.CurrentHP -= attacker.CurrentWeapon.ChokeDamage
			damageDoneThisRound += attacker.CurrentWeapon.ChokeDamage
			attacker.ConsecutiveChokeRounds++

			if attacker.ConsecutiveChokeRounds >= 3 {
				scenario.Defender.CurrentHP = 0
			}

			continue
		}

		scenario.SubtractAmmo(attacker)

		dv := scenario.GetDV(attacker)

		hitParams := HitParams{
			Attribute:      attribute,
			Skill:          skill,
			AttackModifier: scenario.GetAttackModifiers(attacker),
			DV:             dv,
		}

		toHitResult := CalculateHit(hitParams)
		damageResult := scenario.CalculateDamage(DamageParams{
			ToHitResult:    toHitResult,
			AttackType:     attacker.AttackType,
			AmmunitionType: scenario.AmmunitionType,
			HalvesArmor:    attacker.Weapon.HalvesArmor,
			UsesAmmunition: attacker.Weapon.Ranged,
			AttackNumber:   i + 1,
			Attacker:       attacker,
		})

		damageDoneThisRound += damageResult.TotalDamage
		armorAblatedThisRound += damageResult.ArmorAblated
		criticalInjuriesThisRound += damageResult.NumberOfCriticalInjuries
	}

	scenario.TotalAttacks += attacksDoneThisRound

	// Simulate enemy escaping from their grapple on their turn
	if scenario.InGrapple {
		if !scenario.makeGrappleCheckAgainstEnemy(&scenario.Attacker) {
			scenario.InGrapple = false
			attacker.ConsecutiveChokeRounds = 0
		}
	}

	return AttacksResult{
		TotalDamage:      damageDoneThisRound,
		ArmorAblated:     armorAblatedThisRound,
		CriticalInjuries: criticalInjuriesThisRound,
		AttacksDone:      attacksDoneThisRound,
	}
}

func (scenario *CombatScenario) CalculateDamage(params DamageParams) DamageResult {
	damageTotal := 0

	damageResult := DamageResult{
		TotalDamage:              0,
		NumberOfCriticalInjuries: 0,
		ArmorAblated:             0,
	}

	if !params.ToHitResult.Hit {
		return damageResult
	}

	numberOfDiceToRoll := GetDamageDice(params.Attacker.Weapon, params.AttackType)
	diceResult := RollD6s(numberOfDiceToRoll)

	damageTotal = diceResult.Total

	// Get autofire damage here before armor
	if params.AttackType == Autofire {
		damageTotal = params.Attacker.CurrentWeapon.getAutofireDamage(damageTotal, params.ToHitResult.Difference)
	}

	if params.AttackNumber == 1 {
		damageTotal += params.Attacker.CombatAwarenessSpotWeakness
	}

	if scenario.Defender.HasPopupShield && scenario.Defender.PopupShieldHP > 0 {
		scenario.Defender.PopupShieldHP -= damageTotal
		damageResult.ShieldDamage += damageTotal
		return damageResult
	}

	if diceResult.NumberOf6s >= 2 {
		scenario.Defender.CurrentHP = scenario.Defender.CurrentHP - 5
		damageResult.NumberOfCriticalInjuries++
		damageResult.TotalDamage += 5
	}

	armorValue := scenario.Defender.CurrentSP
	if params.Attacker.CurrentWeapon.IgnoresArmorUnder > 0 {
		if armorValue < params.Attacker.CurrentWeapon.IgnoresArmorUnder {
			armorValue = 0
		}
	} else if armorValue > 0 && params.HalvesArmor {
		armorValue = int(math.Ceil(float64(armorValue) / 2))
	}

	damageDifference := damageTotal - armorValue

	if damageDifference > 0 {
		damageDone := GetDamageDone(params, damageDifference)
		damageResult.TotalDamage += damageDone
		scenario.DamageDone += damageDone

		armorAblated := GetArmorAblated(params)
		damageResult.ArmorAblated += armorAblated

		scenario.Defender.CurrentHP = scenario.Defender.CurrentHP - damageDone

		if scenario.Defender.CurrentSP >= armorAblated {
			scenario.Defender.CurrentSP = scenario.Defender.CurrentSP - armorAblated
		} else {
			scenario.Defender.CurrentSP = 0
		}
	}

	return damageResult
}

func (scenario *CombatScenario) setCost() {
	scenarioCost := 0

	setupCost := scenario.Attacker.CurrentWeapon.Cost

	if scenario.Attacker.CurrentWeapon.Ranged {
		ammunitionCost := AmmunitionTypes[string(scenario.AmmunitionType)].Cost

		if scenario.Attacker.CurrentWeapon.Explosive {
			scenarioCost = scenarioCost + (ammunitionCost * scenario.TotalAttacks)
		} else {
			TensOfAmmo := int(math.Ceil(float64(scenario.TotalAttacks) / 10))
			scenarioCost = scenarioCost + (TensOfAmmo * ammunitionCost)
		}

		if scenario.Attacker.HasSmartLink {
			setupCost += 1100
		}

		if scenario.Attacker.DrumClip {
			setupCost += 500
		}

		if scenario.Attacker.ExtendedClip {
			setupCost += 100
		}
	}

	if !scenario.Attacker.CurrentWeapon.Unarmed {
		if scenario.Attacker.ExcellentWeapon {
			switch scenario.Attacker.Weapon.Cost {
			case 100:
				setupCost += 400
			case 500:
				setupCost += 500
			}
		}
	}

	scenario.SetupCost = setupCost
	scenario.ScenarioCost = scenarioCost
}

func (character Character) getClipSize() int {
	clipSize := character.Weapon.ClipSize

	if character.ExtendedClip {
		clipSize = character.Weapon.ExtendedClipSize
	}

	if character.DrumClip {
		clipSize = character.Weapon.DrumClipSize
	}

	return clipSize
}

func (scenario *CombatScenario) inMeleeRange() bool {
	if scenario.DistanceBetweenCharacters > VeryClose.MaxDistance {
		return false
	}

	return true
}

func (weapon *CurrentWeapon) getAutofireDamage(damageTotal, overage int) int {
	autofireMultiplier := overage
	if autofireMultiplier > weapon.AutofireMax {
		autofireMultiplier = weapon.AutofireMax
	}

	damageTotal = damageTotal * autofireMultiplier
	return damageTotal
}

func GetDamageDone(params DamageParams, armorAdjustedDamage int) int {
	damageDone := armorAdjustedDamage
	switch params.AttackType {
	case Headshot:
		damageDone = damageDone * 2
	}

	return damageDone
}

func GetArmorAblated(params DamageParams) int {
	armorAblated := 1

	if params.AmmunitionType == ArmorPiercing && params.UsesAmmunition {
		armorAblated++
	}

	return armorAblated
}

func GetDamageDice(weapon Weapon, attackType AttackType) int {
	if attackType == Autofire {
		if !weapon.CanAutofire {
			panic(fmt.Sprintf("trying to autofire a weapon that cant autofire: '%s'", weapon.Name))
		}

		return 2
	}

	return weapon.DamageDice
}

type DamageResult struct {
	TotalDamage              int
	NumberOfCriticalInjuries int
	ArmorAblated             int
	ShieldDamage             int
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

func (scenario *CombatScenario) GetAttackModifiers(character *CurrentCharacter) int {
	totalModifier := 0
	if character.AttackType == Headshot {
		totalModifier -= 8
		totalModifier -= character.AimedShotBonus
	}

	if character.CurrentWeapon.Ranged && character.HasSmartLink {
		totalModifier += 1
	}

	if !character.CurrentWeapon.Unarmed && character.ExcellentWeapon {
		totalModifier += 1
	}

	totalModifier += character.CombatAwarenessPrecisionAttack
	totalModifier += character.GetWoundPenalty()
	totalModifier += character.ArmorPenalty

	return totalModifier + GetTotalModifiers(character.AttackModifiers)
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

func (scenario *CombatScenario) GetDV(attacker *CurrentCharacter) int {
	weapon := attacker.CurrentWeapon.Weapon
	dv := 0
	dvModifier := scenario.GetDefenderDVModifier()
	if weapon.Ranged {
		dv = weapon.RangeBandDVs[scenario.RangeBand]
		if attacker.AttackType == Autofire {
			dv = weapon.AutofireRangeBandDVs[scenario.RangeBand]
		}

		if !(scenario.Defender.HasPopupShield && scenario.Defender.PopupShieldHP > 0) && scenario.Defender.Reflexes >= 8 {
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
	totalModifiers += scenario.Defender.GetWoundPenalty()
	return totalModifiers
}

func (character *CurrentCharacter) GetWoundPenalty() int {
	halfHP := 0
	if character.CurrentHP > 0 {
		halfHP = int(math.Ceil((float64(character.MaxHP) / 2)))
	}

	if character.CurrentHP <= halfHP {
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
	UsesAmmunition bool
	AttackNumber   int
	Attacker       *CurrentCharacter
}

func (scenario *CombatScenario) SubtractAmmo(character *CurrentCharacter) {
	character.CurrentWeapon.subtractAmmo(character.AttackType)
}

func (scenario *CombatScenario) CalculateReloadFor(attacker *CurrentCharacter) {
	scenario.NumberOfReloads++
	attacker.TurnsSpentOnThisReload++

	turnsToReload := 1
	if attacker.CurrentWeapon.Weapon.TurnsToReload > 1 {
		turnsToReload = attacker.CurrentWeapon.Weapon.TurnsToReload
	}

	if attacker.TurnsSpentOnThisReload >= turnsToReload {
		attacker.reload()
	}
}

func (character *CurrentCharacter) reload() {
	character.CurrentWeapon.CurrentClipSize = character.getClipSize()
	character.TurnsSpentOnThisReload = 0
}

func (scenario *CombatScenario) makeGrappleCheckAgainstEnemy(character *CurrentCharacter) bool {
	attackerRoll := GetD10CheckResult(character.Dexterity, character.Brawling, scenario.getGrappleModifiers(*character))
	defenderRoll := GetD10CheckResult(scenario.Defender.Dexterity, scenario.Defender.Brawling, scenario.getGrappleModifiers(scenario.Defender))

	if attackerRoll > defenderRoll {
		return true
	}

	return false
}

func (scenario *CombatScenario) getGrappleModifiers(character CurrentCharacter) int {
	modifiers := character.ArmorPenalty
	modifiers += character.GetWoundPenalty()
	return modifiers
}

func (scenario *CombatScenario) moveCharacterCloser(character CurrentCharacter) {
	scenario.DistanceBetweenCharacters -= character.Movement

	if scenario.DistanceBetweenCharacters < 0 {
		scenario.DistanceBetweenCharacters = 0
	}
}

func (weapon *CurrentWeapon) subtractAmmo(attackType AttackType) int {
	if !weapon.Ranged {
		return 0
	}

	amountToSubtract := 0

	switch attackType {
	case Autofire:
		autofireCost := 10
		if weapon.AutofireSpent > 0 {
			autofireCost = weapon.AutofireSpent
		}
		amountToSubtract -= autofireCost
	case SingleShot:
		amountToSubtract -= 1
	case Headshot:
		amountToSubtract -= 1
	}

	weapon.CurrentClipSize += amountToSubtract
	return amountToSubtract * -1
}

func (character *CurrentCharacter) mustReload(attackType AttackType) bool {
	if !character.CurrentWeapon.Ranged {
		return false
	}

	if character.TurnsSpentOnThisReload > 1 && character.CurrentWeapon.Weapon.TurnsToReload > 1 && character.TurnsSpentOnThisReload < character.CurrentWeapon.Weapon.TurnsToReload {
		return true
	}

	switch attackType {
	case Autofire:
		autofireCost := 10
		if character.CurrentWeapon.AutofireSpent > 0 {
			autofireCost = character.CurrentWeapon.AutofireSpent
		}
		if character.CurrentWeapon.CurrentClipSize < autofireCost {
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

func (character *CurrentCharacter) GetAttackSkill() int {
	skill := 0
	switch character.CurrentWeapon.Weapon.Skill {
	case Handguns:
		skill = character.Handguns
	case ShoulderArms:
		skill = character.ShoulderArms
	case HeavyWeapons:
		skill = character.HeavyWeapons
	case AutoFire:
		skill = character.AutoFire
	case Brawling:
		skill = character.Brawling
	case Melee:
		skill = character.Melee
	case MartialArts:
		skill = character.MartialArts
	default:
		panic(fmt.Sprintf("missing weapon skill %s", character.Weapon.Skill))
	}

	return skill
}
