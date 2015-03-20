package main

import (
	"math/rand"
	"time"
)

var NUMBER_OF_COLS = 9
var NUMBER_OF_ROWS = 9

// Whenever a random gene or mutation occurs, how many chances should there be that the number will UNASSIGNED?
var NUMBER_OF_CHANCES_FOR_UNASSIGNED = 6

func main() {

	rand.Seed(int64(time.Now().Unix()))

	startBoard := BoardParser("/Users/welshej/github/GeneticSudoku/src/main/board.txt")
	gene := getRandomGene(startBoard)

	b := getBoardFromGene(gene)

	b.Print()

}

// Generates a random gene sequence that represents a possible partial solution to the given board
// TODO Non-Recursive Solution
func getRandomGene(b Board) (gene []string) {

	gene = make([]string, NUMBER_OF_COLS*NUMBER_OF_ROWS)

	for r := 0; r < NUMBER_OF_ROWS; r++ {
		for c := 0; c < NUMBER_OF_COLS; c++ {
			if b.Get(r, c) == UNASSIGNED {
				rand := random(1, 9+NUMBER_OF_CHANCES_FOR_UNASSIGNED)
				gene[c+(r*9)] = numToBitString(rand)
			} else {
				gene[c+(r*9)] = numToBitString(b.Get(r, c))
			}
		}
	}

	if getBoardFromGene(gene).PossibleBoard() {
		return gene
	} else {
		return getRandomGene(b)
	}

}

// Given a specific gene, will get the board for that gene
func getBoardFromGene(gene []string) Board {

	board := Init()

	var pos Position
	var index int

	for r := 0; r < NUMBER_OF_ROWS; r++ {
		for c := 0; c < NUMBER_OF_COLS; c++ {
			pos.Set(r, c)
			board = board.Assign(pos, bitStringToNum(gene[index]))
			index++
		}
	}

	return board
}

// Generates a random number between min and max (inclusive)
func random(min int, max int) int {
	return rand.Intn(max) + min
}

// Given a bit string, will provide the number which maps to that bit string
func bitStringToNum(s string) int {
	switch s {
	case "0001":
		return 1
	case "0010":
		return 2
	case "0011":
		return 3
	case "0100":
		return 4
	case "0101":
		return 5
	case "0110":
		return 6
	case "0111":
		return 7
	case "1000":
		return 8
	case "1001":
		return 9
	default:
		return 0
	}
}

// Given a num, will give the bit string that corresponds to that number
func numToBitString(num int) string {
	switch num {
	case 1:
		return "0001"
	case 2:
		return "0010"
	case 3:
		return "0011"
	case 4:
		return "0100"
	case 5:
		return "0101"
	case 6:
		return "0110"
	case 7:
		return "0111"
	case 8:
		return "1000"
	case 9:
		return "1001"
	default:
		return "0000"
	}
}
