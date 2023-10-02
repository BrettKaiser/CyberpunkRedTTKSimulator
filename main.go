package main

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
	"sort"
	"strconv"
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
	enemy := BoostGanger
	perBandResults := make([]PerBandResult, 0)
	attacker := Character{
		CharacterStats:  PlayerCharacter,
		Weapon:          HeavyPistol,
		ArmorValue:      11,
		ArmorPenalty:    0,
		HasSmartLink:    true,
		AimedShotBonus:  1,
		CombatAwareness: 3,
		DrumClip:        true,
	}

	fmt.Printf("\nCharacter Name: %s / Has Smartlink: %t / Combat Awareness: %d / Aimed Shot: %d / Extended Clip: %t / Drum Clip: %t\n",
		attacker.Name, attacker.HasSmartLink, attacker.CombatAwareness, attacker.AimedShotBonus, attacker.ExtendedClip, attacker.DrumClip,
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
							Ammunition: Basic,
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

	displayTable(perBandResults)
}

type WeaponRunResult struct {
	RunResults []RunResult
	WeaponName string
}

type RunResult struct {
	AverageAttacksToKill      string
	AverageRoundsToKill       string
	AverageNumberOfReloads    string
	AverageRoundsSpentRunning string
	AttackType                string
	WeaponName                string
}

func runScenario(scenarioParams ScenarioParams) RunResult {
	runResult := RunResult{
		AttackType:           string(scenarioParams.AttackType),
		WeaponName:           scenarioParams.Attacker.Weapon.Name,
		AverageAttacksToKill: "NA",
		AverageRoundsToKill:  "NA",
	}

	scenariosRun := []CombatScenario{}

	fmt.Printf("Running Scenario -  Enemy: %s / Range Band: %s / Weapon: %s / Attack Type: %s\n",
		scenarioParams.Defender.Name, scenarioParams.RangeBand.Name, scenarioParams.Attacker.Weapon.Name, scenarioParams.AttackType)

	switch scenarioParams.AttackType {
	case Autofire:
		if !scenarioParams.Attacker.Weapon.CanAutofire {
			return runResult
		}
	case Headshot:
		if !scenarioParams.Attacker.Weapon.CanAimedShot {
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

	runResult.AverageRoundsToKill = fmt.Sprintf("%.3f", getAverageRoundsToKill(scenariosRun))
	runResult.AverageAttacksToKill = fmt.Sprintf("%.3f", getAverageAttacksToKill(scenariosRun))
	runResult.AverageNumberOfReloads = fmt.Sprintf("%.2f", getAverageReloads(scenariosRun))
	runResult.AverageRoundsSpentRunning = fmt.Sprintf("%.2f", getAverageRunningInstances(scenariosRun))

	return runResult
}

func displayTable(perBandResults []PerBandResult) {
	headerRow := table.Row{"Weapon / Attack Type"}

	for _, rangeBand := range RangeBands {
		headerRow = append(headerRow, fmt.Sprintf("%s (RTK)", string(rangeBand.Name)))
	}
	headerRow = append(headerRow, "Total Average RTK")

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(headerRow)
	t.AppendRows(getRows(perBandResults))
	t.SetStyle(table.StyleColoredBright)
	t.Render()
}

func getRows(perBandResults []PerBandResult) []table.Row {
	rows := []table.Row{}

	// Get a map of Weapon Name / Attack Type / Average Rounds To Kill For Each Range Band
	rangeBandsByWeaponAndAttackType := map[string][]string{}
	totalTimeToKillPerWeaponAndAttackType := map[string]float64{}

	for _, bandResult := range perBandResults {
		for _, weaponResults := range bandResult.RunResultsByWeapon {
			for _, attackTypeResults := range weaponResults.RunResults {
				rowName := fmt.Sprintf("%s / %s", weaponResults.WeaponName, attackTypeResults.AttackType)

				if _, ok := rangeBandsByWeaponAndAttackType[rowName]; !ok {
					rangeBandsByWeaponAndAttackType[rowName] = []string{}
					totalTimeToKillPerWeaponAndAttackType[rowName] = 0
				}
				rangeBandsByWeaponAndAttackType[rowName] = append(rangeBandsByWeaponAndAttackType[rowName], attackTypeResults.AverageRoundsToKill)

				if attackTypeResults.AverageRoundsToKill == "NA" {
					totalTimeToKillPerWeaponAndAttackType[rowName] += 1000
				} else {
					roundsToKillFloat, err := strconv.ParseFloat(attackTypeResults.AverageRoundsToKill, 64)
					if err != nil {
						panic("tried to parse a non-float value")
					}
					totalTimeToKillPerWeaponAndAttackType[rowName] += roundsToKillFloat
				}
			}
		}
	}

	// Put the map into a slice of slices of strings
	weaponAttackTypeRows := [][]string{}
	for weaponAttackType, rangeBandResults := range rangeBandsByWeaponAndAttackType {
		newRow := []string{weaponAttackType}
		newRow = append(newRow, rangeBandResults...)
		newRow = append(newRow, fmt.Sprintf("%.3f", totalTimeToKillPerWeaponAndAttackType[weaponAttackType]))
		weaponAttackTypeRows = append(weaponAttackTypeRows, newRow)
	}

	// Sort the slice of slices of strings by the total time to kill per weapon / attack type
	sort.Slice(weaponAttackTypeRows, func(i, j int) bool {
		val1 := totalTimeToKillPerWeaponAndAttackType[weaponAttackTypeRows[i][0]]
		val2 := totalTimeToKillPerWeaponAndAttackType[weaponAttackTypeRows[j][0]]
		return val1 < val2
	})

	// Put the slice of slices of strings into a slice of table.Rows
	for _, weaponAttackTypeRow := range weaponAttackTypeRows {
		newRow := table.Row{}
		for _, value := range weaponAttackTypeRow {
			newRow = append(newRow, value)
		}
		rows = append(rows, newRow)
	}

	// sort.Slice(weaponAttackTypeRows, func(i, j int) bool {
	// 	val1, _ := strconv.ParseFloat(weaponRunResults[i].RunResults[0].AverageRoundsToKill, 64)
	// 	val2, _ := strconv.ParseFloat(weaponRunResults[j].RunResults[0].AverageRoundsToKill, 64)
	// 	return val1 < val2
	// })
	//
	// for _, weaponResult := range weaponRunResults {
	// 	// fmt.Println("STUFF")
	// 	newRow := table.Row{weaponResult.WeaponName}
	// 	for _, result := range weaponResult.RunResults {
	// 		// newStringValue := fmt.Sprintf("%s/%s", result.AverageAttacksToKill, result.AverageRoundsToKill)
	// 		newStringValue := fmt.Sprintf("%s/R%s/M%s", result.AverageRoundsToKill, result.AverageNumberOfReloads, result.AverageRoundsSpentRunning)
	// 		newRow = append(newRow, newStringValue)
	// 	}
	//
	// 	rows = append(rows, newRow)
	// }

	return rows
}

// func displayTable(perBandResults []PerBandResult) {
// 	for _, rangeBandResult := range perBandResults {
// 		fmt.Println("\n*************** Range Band: ", rangeBandResult.RangeBandName, " ****************")
// 		headerRow := table.Row{"Weapon"}
//
// 		// headerNames := []string{"Weapon"}
// 		for _, attackType := range AttackTypes {
// 			headerRow = append(headerRow, fmt.Sprintf("%s (RTK)", string(attackType)))
// 		}
//
// 		t := table.NewWriter()
// 		t.SetOutputMirror(os.Stdout)
// 		t.AppendHeader(headerRow)
// 		t.AppendRows(getRows(rangeBandResult.RunResultsByWeapon))
// 		t.SetStyle(table.StyleColoredBright)
// 		t.Render()
// 		fmt.Println("\n*************** END ****************\n")
// 	}
// }

// func getRows(weaponRunResults []WeaponRunResult) []table.Row {
// 	rows := []table.Row{}
// 	sort.Slice(weaponRunResults, func(i, j int) bool {
// 		val1, _ := strconv.ParseFloat(weaponRunResults[i].RunResults[0].AverageRoundsToKill, 64)
// 		val2, _ := strconv.ParseFloat(weaponRunResults[j].RunResults[0].AverageRoundsToKill, 64)
// 		return val1 < val2
// 	})
//
// 	for _, weaponResult := range weaponRunResults {
// 		// fmt.Println("STUFF")
// 		newRow := table.Row{weaponResult.WeaponName}
// 		for _, result := range weaponResult.RunResults {
// 			// newStringValue := fmt.Sprintf("%s/%s", result.AverageAttacksToKill, result.AverageRoundsToKill)
// 			newStringValue := fmt.Sprintf("%s/R%s/M%s", result.AverageRoundsToKill, result.AverageNumberOfReloads, result.AverageRoundsSpentRunning)
// 			newRow = append(newRow, newStringValue)
// 		}
//
// 		rows = append(rows, newRow)
// 	}
//
// 	return rows
// }

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
