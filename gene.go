package main

const GENE_SIZE = NUMBER_OF_ROWS * NUMBER_OF_COLS

type Gene struct {
	gene [GENE_SIZE]string
}

func (g *Gene) Score() int {
	b := getBoardFromGene(*g)
	return b.Grade()
}

func mateGenes(a Gene, b Gene) (res Gene) {

	r := random(1, GENE_SIZE-1)
	aFirst := random(0, 1)

	for i := 0; i < r; i++ {
		if aFirst == 1 {
			res.gene[i] = a.gene[i]
		} else {
			res.gene[i] = b.gene[i]
		}
	}
	for i := r; i < GENE_SIZE; i++ {
		if aFirst == 1 {
			res.gene[i] = b.gene[i]
		} else {
			res.gene[i] = a.gene[i]
		}
	}
	return
}

// Generates a random gene sequence that represents a possible partial solution to the given board
// TODO Non-Recursive Solution
func getRandomGene(b Board) (g Gene) {

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

	if getBoardFromGene(g).PossibleBoard() {
		return g
	} else {
		return getRandomGene(b)
	}
}
