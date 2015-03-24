package main

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
func Mutate(population []Chromosome, chanceToModifyPopulation float64) []Chromosome {

	if chanceToModifyPopulation == 0 {
		return population
	}

	for randomFloat(0, 1) < chanceToModifyPopulation { // if you decided to mutate...

		modifiedChromosome := randomInt(0, len(population)) // pick a random chromosome to modify
		modifiedGene := randomInt(0, 81)                    // pick a random gene within that chromosome to modify

		rand := randomInt(0, 9+NUMBER_OF_CHANCES_FOR_UNASSIGNED) // randomly assign a number or set unassigned

		// add the mutation to the chromosome
		population[modifiedChromosome].genes[modifiedGene] = geneToNum(uint8(rand))
	}

	return population
}

// Will perform a crossover operation between two chromosomes
func MateChromosome(a Chromosome, b Chromosome) (resa Chromosome, resb Chromosome) {

	if randomFloat(0, 1) < CROSSOVER_RATE {

		r := randomInt(1, CHROMOSOME_SIZE) // pick a random spot within the chromosomes to crossover

		for i := 0; i < r; i++ { // get genes from a up until crossover point
			resa.genes[i] = a.genes[i]
			resb.genes[i] = b.genes[i]
		}
		for i := r; i < CHROMOSOME_SIZE; i++ { // after that, get genes from b
			resa.genes[i] = b.genes[i]
			resb.genes[i] = a.genes[i]
		}

		return resa, resb
	} else {
		return a, b
	}
}

// Generates a random gene sequence that represents a possible partial solution to the given board
func GetRandomChromosome(b *Board) (chromosome Chromosome) {

	for r := 0; r < NUMBER_OF_ROWS; r++ {
		for c := 0; c < NUMBER_OF_COLS; c++ {
			if b.Get(r, c) == UNASSIGNED { // for every unassigned cell in the board

				rand := randomInt(0, 9+NUMBER_OF_CHANCES_FOR_UNASSIGNED) // randomly assign a number or set unassigned
				chromosome.genes[c+(r*9)] = geneToNum(uint8(rand))

			} else { // TODO consider copying all genes over first...
				chromosome.genes[c+(r*9)] = geneToNum(b.Get(r, c))
			}
		}
	}

	return chromosome

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
