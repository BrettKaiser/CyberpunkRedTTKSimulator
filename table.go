package main

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
	"sort"
	"strconv"
)

func displayTableByLowestRTKPerRangeBand(perBandResults []PerBandResult) {
	headerRow := table.Row{}

	for _, rangeBand := range RangeBands {
		headerRow = append(headerRow, fmt.Sprintf("%s (RTK)", string(rangeBand.Name)))
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(headerRow)
	t.AppendRows(getRowsTableByLowestRTKPerRangeBand(perBandResults))
	t.SetStyle(table.StyleColoredBright)
	t.Render()
}

func getRowsTableByLowestRTKPerRangeBand(perBandResults []PerBandResult) []table.Row {
	rows := []table.Row{}

	type WeaponAttackTypeResult struct {
		WeaponAttackType    string
		AverageRoundsToKill string
	}

	rangeBandValues := map[string][]WeaponAttackTypeResult{}

	// Combine all the results for each range band
	for _, bandResult := range perBandResults {
		for _, weaponResults := range bandResult.RunResultsByWeapon {
			for _, attackTypeResults := range weaponResults.RunResults {
				if _, ok := rangeBandValues[bandResult.RangeBandName]; !ok {
					rangeBandValues[bandResult.RangeBandName] = []WeaponAttackTypeResult{}
				}
				rangeBandValues[bandResult.RangeBandName] = append(rangeBandValues[bandResult.RangeBandName], WeaponAttackTypeResult{
					WeaponAttackType:    fmt.Sprintf("%s / %s", weaponResults.WeaponName, attackTypeResults.AttackType),
					AverageRoundsToKill: attackTypeResults.AverageRoundsToKill,
				})
			}
		}
	}

	// For each range band in rangeBandValues, sort the slice of WeaponAttackTypeResults by AverageRoundsToKill, converting to float64 first and setting NA to 1000
	for _, weaponAttackTypeResultsPerRangeBand := range rangeBandValues {
		sort.Slice(weaponAttackTypeResultsPerRangeBand, func(i, j int) bool {
			val1 := weaponAttackTypeResultsPerRangeBand[i].AverageRoundsToKill
			val2 := weaponAttackTypeResultsPerRangeBand[j].AverageRoundsToKill

			if val1 == "NA" {
				val1 = "1000"
			}
			if val2 == "NA" {
				val2 = "1000"
			}

			val1Float, err := strconv.ParseFloat(val1, 64)
			if err != nil {
				panic("tried to parse a non-float value")
			}
			val2Float, err := strconv.ParseFloat(val2, 64)
			if err != nil {
				panic("tried to parse a non-float value")
			}
			return val1Float < val2Float
		})
	}

	for i := 0; i < len(rangeBandValues[VeryClose.Name])-1; i++ {
		newRow := table.Row{}
		for _, rangeBand := range RangeBands {
			newRow = append(newRow, fmt.Sprintf("%s - %s", rangeBandValues[rangeBand.Name][i].AverageRoundsToKill, rangeBandValues[rangeBand.Name][i].WeaponAttackType))
		}
		rows = append(rows, newRow)
	}

	return rows
}

func displayTableByAverageRTKACrossRangeBands(perBandResults []PerBandResult) {
	headerRow := table.Row{"Weapon / Attack Type"}

	for _, rangeBand := range RangeBands {
		headerRow = append(headerRow, fmt.Sprintf("%s (RTK)", string(rangeBand.Name)))
	}
	headerRow = append(headerRow, "Total Average RTK")
	headerRow = append(headerRow, "Average Eddies Spent In Combat")
	headerRow = append(headerRow, "Eddies Spent On Setup")

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(headerRow)
	t.AppendRows(getRowsTableByAverageRTKACrossRangeBands(perBandResults))
	t.SetStyle(table.StyleColoredBright)
	t.Render()
}

func getRowsTableByAverageRTKACrossRangeBands(perBandResults []PerBandResult) []table.Row {
	rows := []table.Row{}

	// Get a map of Weapon Name / Attack Type / Average Rounds To Kill For Each Range Band
	rangeBandsByWeaponAndAttackType := map[string][]string{}
	totalTimeToKillPerWeaponAndAttackType := map[string]float64{}
	eddiesPerScenarioSpentByAttackType := map[string]float64{}
	setupCost := map[string]string{}

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
					eddiesPerScenarioSpentByAttackType[rowName] = 0
					setupCost[rowName] = "NA"
				} else {
					roundsToKillFloat, err := strconv.ParseFloat(attackTypeResults.AverageRoundsToKill, 64)
					if err != nil {
						panic("tried to parse a non-float value")
					}
					totalTimeToKillPerWeaponAndAttackType[rowName] += roundsToKillFloat

					eddiesPerScenarioFloat, err := strconv.ParseFloat(attackTypeResults.AverageEddiesSpentPerScenario, 64)
					if err != nil {
						panic("tried to parse a non-float value")
					}
					eddiesPerScenarioSpentByAttackType[rowName] += eddiesPerScenarioFloat
					setupCost[rowName] = attackTypeResults.SetupCost
				}
			}
		}
	}

	for key, value := range eddiesPerScenarioSpentByAttackType {
		eddiesPerScenarioSpentByAttackType[key] = value / float64(len(RangeBands))
	}

	// Put the map into a slice of slices of strings
	weaponAttackTypeRows := [][]string{}
	for weaponAttackType, rangeBandResults := range rangeBandsByWeaponAndAttackType {
		newRow := []string{weaponAttackType}
		newRow = append(newRow, rangeBandResults...)
		newRow = append(newRow, fmt.Sprintf("%.3f", totalTimeToKillPerWeaponAndAttackType[weaponAttackType]))
		newRow = append(newRow, fmt.Sprintf("$%.2f", eddiesPerScenarioSpentByAttackType[weaponAttackType]))
		newRow = append(newRow, fmt.Sprintf("$%s", setupCost[weaponAttackType]))
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

	return rows
}

// func oldDisplayTable(perBandResults []PerBandResult) {
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
// 		t.AppendRows(oldGetRows(rangeBandResult.RunResultsByWeapon))
// 		t.SetStyle(table.StyleColoredBright)
// 		t.Render()
// 		fmt.Println("\n*************** END ****************\n")
// 	}
// }

// func oldGetRows(weaponRunResults []WeaponRunResult) []table.Row {
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
