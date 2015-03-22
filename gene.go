package main

import (
//"fmt"
)

var scoreCache map[string]int

type Gene struct {
	gene [CHROMOSOME_SIZE]string
}

// Fitness function used to determine the degree of completion of the board
func (g Gene) Score() int {

	gs := g.String()

	if val, ok := scoreCache[gs]; ok {
		return val
	} else {
		b := getBoardFromGene(g)
		grade := b.Grade()
		scoreCache[gs] = grade
		return grade
	}

}

// Will randomly mutate random genes in a populations
func Mutate(population []Gene, chanceToMutateGene float64) {

	if chanceToMutateGene == 0 {
		return
	}

	chanceToModifyChromosome := chanceToMutateGene * 81
	chanceToModifyPopulation := (chanceToModifyChromosome * 1000) / ((chanceToModifyChromosome * 1000) + 1)

	for randomFloat(0, 1) < chanceToModifyPopulation {
		modifiedChromosome := randomInt(0, len(population))
		modifiedGene := randomInt(0, 81)

		b := getBoardFromGene(population[modifiedChromosome])

		modifyRow := modifiedGene / NUMBER_OF_COLS
		modifyCol := modifiedGene % NUMBER_OF_ROWS

		possibilities := b.PossibleCells(modifyRow, modifyCol)

		rand := randomInt(0, len(possibilities)+NUMBER_OF_CHANCES_FOR_UNASSIGNED)

		var valueToAdd int

		if rand < len(possibilities) {

			valueToAdd = possibilities[rand]

		} else {
			valueToAdd = 0
		}

		population[modifiedChromosome].gene[modifiedGene] = numToBitString(valueToAdd)
	}

	return
}

// Will perform a crossover operation between two genes
func mateGenes(a Gene, b Gene) (res Gene) {

	firstIteration := true

	for i := 0; firstIteration || getBoardFromGene(res).IsWrong(); i++ {

		firstIteration = false

		r := randomInt(1, CHROMOSOME_SIZE)

		for i := 0; i < r; i++ {
			res.gene[i] = a.gene[i]
		}
		for i := r; i < CHROMOSOME_SIZE; i++ {
			res.gene[i] = b.gene[i]
		}

		//To prevent deadlock, after a certain amount of unsuccessful mating attempts, just return the high board
		if i >= MATE_MAX_RETRIES {
			if a.Score() > b.Score() {
				return a
			} else {
				return b
			}
		}
	}

	return res
}

// Generates a random gene sequence that represents a possible partial solution to the given board
func getRandomGene(b *Board) (g Gene) {

	cpy := b.Clone()

	for r := 0; r < NUMBER_OF_ROWS; r++ {
		for c := 0; c < NUMBER_OF_COLS; c++ {

			if b.Get(r, c) == UNASSIGNED {

				possibilities := cpy.PossibleCells(r, c)

				rand := randomInt(0, len(possibilities)+NUMBER_OF_CHANCES_FOR_UNASSIGNED)

				var valueToAdd int

				if rand < len(possibilities) {

					valueToAdd = possibilities[rand]

				} else {
					valueToAdd = 0
				}

				g.gene[c+(r*9)] = numToBitString(valueToAdd)
				cpy.Set(r, c, valueToAdd)

			} else {
				g.gene[c+(r*9)] = numToBitString(b.Get(r, c))
			}
		}
	}

	return g
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

// Returns the string representation of a particular gene
func (g *Gene) String() (ret string) {
	for _, s := range g.gene {
		ret += s
	}
	return
}
