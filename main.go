package main

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

const ITERATIONS = 1000

const MAXIMUMROUNDS = 100

type SimulationParams struct {
	Iterations int
	DebugLogs  bool
}

type ScenarioParams struct {
	Ammunition AmmunitionType
	AttackType AttackType
	Attacker   Character
	Defender   Character
	RangeBand  RangeBand
	DebugLogs  bool
	Iterations int
}

func main() {
	simulationParams := SimulationParams{
		Iterations: ITERATIONS,
		DebugLogs:  false,
	}

	fmt.Printf("Starting Simulation, Iterations: %d\n", simulationParams.Iterations)

	startTime := time.Now()
	Run(simulationParams)
	endTime := time.Now()
	totalTime := endTime.Sub(startTime)
	fmt.Printf("\nTook %f seconds to run\n", totalTime.Seconds())
}

type PerBandResult struct {
	RangeBandName      string
	RunResultsByWeapon []WeaponRunResult
}

func Run(params SimulationParams) {
	enemies := []Character{
		BoostGanger,
		CyberPsycho,
	}
	ammunitionTypes := []AmmunitionType{
		// Basic,
		ArmorPiercing,
	}
	for _, ammunitionType := range ammunitionTypes {
		for _, enemy := range enemies {
			perBandResults := make([]PerBandResult, 0)
			attacker := Character{
				CharacterStats:                 PlayerCharacter,
				Weapon:                         HeavyPistol, // this is filler, it gets overwritten in the scenario
				ArmorValue:                     11,
				ArmorPenalty:                   0,
				HasSmartLink:                   true,
				AimedShotBonus:                 1,
				CombatAwarenessPrecisionAttack: 0,
				CombatAwarenessSpotWeakness:    10,
				ExtendedClip:                   false,
				DrumClip:                       true,
				ExcellentWeapon:                true,
			}

			fmt.Printf("\nCharacter Name: %s / Ammunition Type: %s / Has Smartlink: %t / Combat Awareness: %d / Aimed Shot: %d / Extended Clip: %t / Drum Clip: %t / Enemy: %s / Enemy AP: %d\n",
				attacker.Name, ammunitionType, attacker.HasSmartLink, attacker.CombatAwarenessPrecisionAttack, attacker.AimedShotBonus, attacker.ExtendedClip, attacker.DrumClip, enemy.Name, enemy.ArmorValue,
			)

			for _, rangeBand := range RangeBands {
				newRangeBandResult := PerBandResult{
					RangeBandName:      rangeBand.Name,
					RunResultsByWeapon: []WeaponRunResult{},
				}

				var perBandWeaponGroup sync.WaitGroup

				for _, weapon := range WeaponsList {
					perBandWeaponGroup.Add(1)
					go func(weapon Weapon) {
						defer perBandWeaponGroup.Done()
						weaponRunResult := WeaponRunResult{
							WeaponName: weapon.Name,
							RunResults: []RunResult{},
						}

						var perWeaponAttackTypeGroup sync.WaitGroup

						for _, attackType := range AttackTypes {
							perWeaponAttackTypeGroup.Add(1)
							go func(attackType AttackType) {
								defer perWeaponAttackTypeGroup.Done()
								scenarioParams := ScenarioParams{
									Ammunition: ammunitionType,
									AttackType: attackType,
									Attacker:   attacker,
									Defender:   enemy,
									RangeBand:  rangeBand,
									DebugLogs:  params.DebugLogs,
									Iterations: params.Iterations,
								}

								scenarioParams.Attacker.Weapon = weapon

								runResult := runScenario(scenarioParams)
								weaponRunResult.RunResults = append(weaponRunResult.RunResults, runResult)
							}(attackType)
						}
						perWeaponAttackTypeGroup.Wait()

						sort.Slice(weaponRunResult.RunResults, func(i, j int) bool {
							// 	Put SingleShot attack type first, then Headshot and finally Autofire
							if weaponRunResult.RunResults[i].AttackType == string(SingleShot) {
								return true
							} else if weaponRunResult.RunResults[i].AttackType == string(Headshot) && weaponRunResult.RunResults[j].AttackType == string(Autofire) {
								return true
							} else {
								return false
							}
						})

						newRangeBandResult.RunResultsByWeapon = append(newRangeBandResult.RunResultsByWeapon, weaponRunResult)
					}(weapon)
				}

				perBandWeaponGroup.Wait()
				perBandResults = append(perBandResults, newRangeBandResult)
			}

			fmt.Println("\n*************** Results By Average RTK ****************")
			displayTableByAverageRTKACrossRangeBands(perBandResults)
			fmt.Println("\n*************** Results By Lowest RTK Per Range Band ****************")
			displayTableByLowestRTKPerRangeBand(perBandResults)
		}
	}

}

type WeaponRunResult struct {
	RunResults []RunResult
	WeaponName string
}

type RunResult struct {
	AverageAttacksToKill          string
	AverageRoundsToKill           string
	AverageNumberOfReloads        string
	AverageRoundsSpentRunning     string
	AttackType                    string
	WeaponName                    string
	AverageEddiesSpentPerScenario string
	SetupCost                     string
}

func runScenario(scenarioParams ScenarioParams) RunResult {
	runResult := RunResult{
		AttackType:                    string(scenarioParams.AttackType),
		WeaponName:                    scenarioParams.Attacker.Weapon.Name,
		AverageAttacksToKill:          "NA",
		AverageRoundsToKill:           "NA",
		AverageEddiesSpentPerScenario: "NA",
		SetupCost:                     "NA",
	}

	scenariosRun := []CombatScenario{}

	if scenarioParams.DebugLogs {
		fmt.Printf("Running Scenario -  Enemy: %s / Range Band: %s / Weapon: %s / Attack Type: %s\n",
			scenarioParams.Defender.Name, scenarioParams.RangeBand.Name, scenarioParams.Attacker.Weapon.Name, scenarioParams.AttackType)
	}

	switch scenarioParams.AttackType {
	case Autofire:
		if !scenarioParams.Attacker.Weapon.CanAutofire {
			return runResult
		}
	case Headshot:
		if !scenarioParams.Attacker.Weapon.CanAimedShot {
			return runResult
		}
	case SingleShot:
		if scenarioParams.Attacker.Weapon.CannotSingleShot {
			return runResult
		}
	}

	switch scenarioParams.Ammunition {
	case Basic:
		if scenarioParams.Attacker.Weapon.Explosive {
			return runResult
		}
	}

	for i := 0; i < scenarioParams.Iterations; i++ {
		if scenarioParams.DebugLogs {
			fmt.Printf("\n***************\nBegin Scenario %d\n\n", i+1)
		}
		scenario := NewCombatScenario(scenarioParams)
		scenarioResult := scenario.Execute()
		scenariosRun = append(scenariosRun, scenarioResult)

		if scenarioParams.DebugLogs {
			scenario.DisplayResult()
		}
	}

	runResult.AverageRoundsToKill = fmt.Sprintf("%.1f", getAverageRoundsToKill(scenariosRun))
	runResult.AverageAttacksToKill = fmt.Sprintf("%.1f", getAverageAttacksToKill(scenariosRun))
	runResult.AverageNumberOfReloads = fmt.Sprintf("%.1f", getAverageReloads(scenariosRun))
	runResult.AverageRoundsSpentRunning = fmt.Sprintf("%.1f", getAverageRunningInstances(scenariosRun))
	runResult.AverageEddiesSpentPerScenario = fmt.Sprintf("%.1f", getAverageEddiesSpent(scenariosRun))
	runResult.SetupCost = fmt.Sprintf("%d", scenariosRun[0].SetupCost)

	return runResult
}

func getAverageAttacksToKill(scenarios []CombatScenario) float64 {
	total := 0
	numberOfScenarios := 0

	for i, _ := range scenarios {
		numberOfScenarios++
		total += scenarios[i].TotalAttacks
	}

	return float64(total) / float64(numberOfScenarios)
}

func getAverageRoundsToKill(scenarios []CombatScenario) float64 {
	total := 0
	numberOfScenarios := 0

	for i, _ := range scenarios {
		numberOfScenarios++
		total += scenarios[i].NumberOfRounds
	}

	return float64(total) / float64(numberOfScenarios)
}

func getAverageReloads(scenarios []CombatScenario) float64 {
	total := 0
	numberOfScenarios := 0

	for i, _ := range scenarios {
		numberOfScenarios++
		total += scenarios[i].NumberOfReloads
	}

	return float64(total) / float64(numberOfScenarios)
}

func getAverageRunningInstances(scenarios []CombatScenario) float64 {
	total := 0
	numberOfScenarios := 0

	for i, _ := range scenarios {
		numberOfScenarios++
		total += scenarios[i].RoundsSpentRunning
	}

	return float64(total) / float64(numberOfScenarios)
}

func getAverageEddiesSpent(scenarios []CombatScenario) float64 {
	total := 0
	numberOfScenarios := 0

	for i, _ := range scenarios {
		numberOfScenarios++
		total += scenarios[i].ScenarioCost
	}

	return float64(total) / float64(numberOfScenarios)
}
