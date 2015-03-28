package main
import (
	"math/rand"
)

type Chromosome struct {
	genes [CHROMOSOME_SIZE]uint8
}

// Fitness function used to determine the degree of completion of the board
func (c Chromosome) Score() float64 {
	gs := c.String()

	if val, ok := scoreCache[gs]; ok { // cache hit
		return val
	} else { // cache miss
		b := getBoardFromChromosome(c) // get the board that corresponds to this chromosome
		grade := b.Grade()
		scoreCache[gs] = grade // add it to the cache
		return grade
	}

}

// Will randomly mutate random genes in random chromosomes within a given population
func Mutate(original Board, population []Chromosome, chanceToModifyPopulation float64) []Chromosome {

	if chanceToModifyPopulation == 0 {
		return population
	}

	for rand.Float64() < chanceToModifyPopulation { // if you decided to mutate...

		modifiedChromosome := randomInt(0, len(population)) // pick a random chromosome to modify

		modifiedRow := randomInt(0, NUMBER_OF_ROWS-1)

		randomRow := GetRandomRow(modifiedRow)
		startIndex := modifiedRow * NUMBER_OF_ROWS

		for i := 0; i < NUMBER_OF_COLS; i++ {
			population[modifiedChromosome].genes[startIndex+i] = randomRow[i]
		}

		// add the mutation to the chromosome

	}

	return population
}

// Will perform a crossover operation between two chromosomes
func MateChromosome(a Chromosome, b Chromosome) (Chromosome, Chromosome) {

	if rand.Float64() < CROSSOVER_RATE {

		rowToSwap := randomInt(0, NUMBER_OF_ROWS-1) // pick row to swap

		for c := 0; c < NUMBER_OF_COLS; c++ { // and swap it!
			temp := a.genes[(rowToSwap*NUMBER_OF_ROWS)+c]
			a.genes[(rowToSwap*NUMBER_OF_ROWS)+c] = b.genes[(rowToSwap*NUMBER_OF_ROWS)+c]
			b.genes[(rowToSwap*NUMBER_OF_ROWS)+c] = temp
		}

	}
	return a, b
}

// Generates a random gene sequence that represents a possible partial solution to the given board
func GetRandomChromosome(b *Board) (chromosome Chromosome) {

	for r := 0; r < NUMBER_OF_ROWS; r++ {

		randomRow := GetRandomRow(r)
		startIndex := r * NUMBER_OF_ROWS

		for i := 0; i < NUMBER_OF_COLS; i++ {
			chromosome.genes[startIndex+i] = randomRow[i]
		}
	}

	return

}

// Given a number, will provide the gene that maps to that number
func geneToNum(n uint8) uint8 {

	if n < 10 {
		return n
	} else {
		return 0
	}
}

// Returns the string representation of a particular chromosome
func (c *Chromosome) String() (ret string) {
	for _, s := range c.genes {
		ret += string(s) + " "
	}
	return
}
