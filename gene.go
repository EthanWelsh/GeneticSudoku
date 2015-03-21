package main

import (
//"fmt"
)

const GENE_SIZE = NUMBER_OF_ROWS * NUMBER_OF_COLS
const PROB_GRANULARITY = 100000
const MATE_CHANCES = 100

var scoreCache map[string]int

type Gene struct {
	gene [GENE_SIZE]string
}

func (g *Gene) Score() int {

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

func (g *Gene) Mutate(chanceAtMutation float64) (mutationsMade int) {

	if chanceAtMutation == 0 {
		return 0
	}

	for i, _ := range g.gene {
		chanceAtMutation = chanceAtMutation * PROB_GRANULARITY
		r := random(0, PROB_GRANULARITY)

		if float64(r) < chanceAtMutation { // time to mutate!

			temp := *g
			rand := random(1, 9+NUMBER_OF_CHANCES_FOR_UNASSIGNED)
			g.gene[i] = numToBitString(rand)

			if !getBoardFromGene(g).PossibleBoard() {
				g = &temp
			} else {
				mutationsMade++
			}

		}
	}

	return
}

func mateGenes(a *Gene, b *Gene) (res Gene) {

	firstIteration := true

	for i := 0; firstIteration || !getBoardFromGene(&res).PossibleBoard(); i++ {

		//To prevent deadlock, after a certain amount of unsuccessful mating attempts, just return the high board
		if i >= MATE_CHANCES {
			if a.Score() > b.Score() {
				return *a
			} else {
				return *b
			}
		}

		firstIteration = false
		r := random(1, GENE_SIZE)

		for i := 0; i < r; i++ {
			res.gene[i] = a.gene[i]
		}
		for i := r; i < GENE_SIZE; i++ {
			res.gene[i] = b.gene[i]
		}
	}

	return
}

// Generates a random gene sequence that represents a possible partial solution to the given board
func getRandomGene(b *Board) (g Gene) {

	firstIteration := true

	for firstIteration || !getBoardFromGene(&g).PossibleBoard() {
		firstIteration = false
		for r := 0; r < NUMBER_OF_ROWS; r++ {
			for c := 0; c < NUMBER_OF_COLS; c++ {
				if b.Get(r, c) == UNASSIGNED {
					rand := random(1, 9+NUMBER_OF_CHANCES_FOR_UNASSIGNED)
					g.gene[c+(r*9)] = numToBitString(rand)
				} else {
					g.gene[c+(r*9)] = numToBitString(b.Get(r, c))
				}
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

func (g *Gene) String() (ret string) {
	for _, s := range g.gene {
		ret += s
	}
	return
}