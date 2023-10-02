package main

import (
	"math/rand"
)

// prepare the dice
var d6Values = []int{1, 2, 3, 4, 5, 6}
var d10Values = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

func rollDie(possibleValues []int) int {
	return possibleValues[rand.Intn(len(possibleValues))]
}

// RollD6s rolls a number of six sided dice and returns the result as a D6Result
func RollD6s(numberOfDice int) D6Result {
	total := 0
	numberOf6s := 0
	numberOf1s := 0
	resultsList := []int{}

	for i := 0; i < numberOfDice; i++ {
		result := rollDie(d6Values)
		resultsList = append(resultsList, result)

		total += result
		if result == 6 {
			numberOf6s++
		}
		if result == 1 {
			numberOf1s++
		}
	}

	return D6Result{
		Total:           total,
		IndividualRolls: resultsList,
		NumberOf6s:      numberOf6s,
		NumberOf1s:      numberOf1s,
	}
}

type D6Result struct {
	Total           int
	IndividualRolls []int
	NumberOf6s      int
	NumberOf1s      int
}

func RollD10s(numberOfDice int) D10Result {
	total := 0
	numberOf10s := 0
	numberOf1s := 0
	resultsList := []int{}

	for i := 0; i < numberOfDice; i++ {
		result := rollDie(d10Values)
		resultsList = append(resultsList, result)

		total += result
		if result == 10 {
			numberOf10s++
		}
		if result == 1 {
			numberOf1s++
		}
	}

	return D10Result{
		Total:           total,
		IndividualRolls: resultsList,
		NumberOf10s:     numberOf10s,
		NumberOf1s:      numberOf1s,
	}
}

type D10Result struct {
	Total           int
	IndividualRolls []int
	NumberOf10s     int
	NumberOf1s      int
}
