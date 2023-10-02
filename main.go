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

const ITERATIONS = 50000

const MaximumRounds = 300

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

func Run(params SimulationParams) {
	completedRuns := []WeaponRunResult{}
	// weaponRunResultChan := make(chan WeaponRunResult, len(WeaponsList))
	var weaponGroup sync.WaitGroup

	for _, weapon := range WeaponsList {
		weaponGroup.Add(1)
		go func(weapon Weapon) {
			defer weaponGroup.Done()
			weaponRunResult := WeaponRunResult{
				WeaponName: weapon.Name,
				RunResults: []RunResult{},
			}

			var weaponAttackTypeGroup sync.WaitGroup

			for _, attackType := range AttackTypes {
				weaponAttackTypeGroup.Add(1)
				go func(attackType AttackType) {
					defer weaponAttackTypeGroup.Done()
					scenarioParams := ScenarioParams{
						Ammunition: Basic,
						AttackType: attackType,
						Attacker: Character{
							CharacterStats: PlayerCharacter,
							Weapon:         weapon,
							ArmorValue:     11,
							ArmorPenalty:   0,
						},
						Defender:   CyberPsycho,
						RangeBand:  VeryClose,
						DebugLogs:  params.DebugLogs,
						Iterations: params.Iterations,
					}

					runResult := runScenario(scenarioParams)
					weaponRunResult.RunResults = append(weaponRunResult.RunResults, runResult)
				}(attackType)
			}
			weaponAttackTypeGroup.Wait()

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

			completedRuns = append(completedRuns, weaponRunResult)
		}(weapon)

	}

	weaponGroup.Wait()
	displayTable(completedRuns)
}

type WeaponRunResult struct {
	RunResults []RunResult
	WeaponName string
}

type RunResult struct {
	AverageAttacksToKill   string
	AverageRoundsToKill    string
	AverageNumberOfReloads string
	AttackType             string
	WeaponName             string
}

func runScenario(scenarioParams ScenarioParams) RunResult {
	runResult := RunResult{
		AttackType:           string(scenarioParams.AttackType),
		WeaponName:           scenarioParams.Attacker.Weapon.Name,
		AverageAttacksToKill: "NA",
		AverageRoundsToKill:  "NA",
	}

	scenariosRun := []CombatScenario{}

	fmt.Printf("Running Scenario - Weapon: %s, Enemy: %s, Attack Type: %s, Range Band: %s\n",
		scenarioParams.Attacker.Weapon.Name, scenarioParams.Defender.Name, scenarioParams.AttackType, scenarioParams.RangeBand.Name)

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

	if !scenarioParams.Attacker.Weapon.Ranged && scenarioParams.RangeBand != VeryClose {
		return runResult
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

	averageAttacksToKill := getAverageAttacksToKill(scenariosRun)
	averageRoundsToKill := getAverageRoundsToKill(scenariosRun)

	runResult.AverageRoundsToKill = fmt.Sprintf("%.3f", averageRoundsToKill)
	runResult.AverageAttacksToKill = fmt.Sprintf("%.3f", averageAttacksToKill)
	runResult.AverageNumberOfReloads = fmt.Sprintf("%.2f", getAverageReloads(scenariosRun))

	return runResult
}

func displayTable(weaponRunResults []WeaponRunResult) {
	headerRow := table.Row{"Weapon"}

	// headerNames := []string{"Weapon"}
	for _, attackType := range AttackTypes {
		headerRow = append(headerRow, fmt.Sprintf("%s (RTK)", string(attackType)))
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(headerRow)
	t.AppendRows(getRows(weaponRunResults))
	t.SetStyle(table.StyleColoredBright)
	t.Render()
}

func getRows(weaponRunResults []WeaponRunResult) []table.Row {
	rows := []table.Row{}
	sort.Slice(weaponRunResults, func(i, j int) bool {
		val1, _ := strconv.ParseFloat(weaponRunResults[i].RunResults[0].AverageRoundsToKill, 64)
		val2, _ := strconv.ParseFloat(weaponRunResults[j].RunResults[0].AverageRoundsToKill, 64)
		return val1 < val2
	})

	for _, weaponResult := range weaponRunResults {
		// fmt.Println("STUFF")
		newRow := table.Row{weaponResult.WeaponName}
		for _, result := range weaponResult.RunResults {
			// newStringValue := fmt.Sprintf("%s/%s", result.AverageAttacksToKill, result.AverageRoundsToKill)
			newStringValue := fmt.Sprintf("%s/R%s", result.AverageRoundsToKill, result.AverageNumberOfReloads)
			newRow = append(newRow, newStringValue)
		}

		rows = append(rows, newRow)
	}

	return rows
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
