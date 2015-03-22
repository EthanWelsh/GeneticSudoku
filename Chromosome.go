package main

var scoreCache map[string]float64

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
func Mutate(population []Chromosome, chanceToMutateGene float64) []Chromosome {

	if chanceToMutateGene == 0 {
		return population
	}

	chanceToModifyChromosome := chanceToMutateGene * 81
	averageModificationsMadePerPopulation := chanceToModifyChromosome * 1000
	chanceToModifyPopulation := averageModificationsMadePerPopulation / (averageModificationsMadePerPopulation + 1)

	for randomFloat(0, 1) < chanceToModifyPopulation { // if you decided to mutate...
		modifiedChromosome := randomInt(0, len(population)) // pick a random chromosome to modify
		modifiedGene := randomInt(0, 81)                    // pick a random gene within that chromosome to modify

		b := getBoardFromChromosome(population[modifiedChromosome]) // get the board representation of that chromosome

		modifyRow := modifiedGene / NUMBER_OF_COLS
		modifyCol := modifiedGene % NUMBER_OF_ROWS

		possibilities := b.PossibleCells(modifyRow, modifyCol) // get all the valid mutations that could be made

		rand := randomInt(0, len(possibilities)+NUMBER_OF_CHANCES_FOR_UNASSIGNED) // pick one or change cell to unassigned

		var valueToAdd uint8

		if rand < len(possibilities) {
			valueToAdd = possibilities[rand]

		} else {
			valueToAdd = UNASSIGNED
		}

		// add the mutation
		population[modifiedChromosome].genes[modifiedGene] = geneToNum(valueToAdd)
	}

	return population
}

// Will perform a crossover operation between two chromosomes
func MateChromosome(a Chromosome, b Chromosome) (res Chromosome) {

	firstIteration := true

	for i := 0; firstIteration || getBoardFromChromosome(res).IsWrong(); i++ {

		firstIteration = false

		r := randomInt(1, CHROMOSOME_SIZE) // pick a random spot within the chromosomes to crossover

		for i := 0; i < r; i++ { // get genes from a up until crossover point
			res.genes[i] = a.genes[i]
		}
		for i := r; i < CHROMOSOME_SIZE; i++ { // after that, get genes from b
			res.genes[i] = b.genes[i]
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
func GetRandomChromosome(b *Board) (chromosome Chromosome) {

	cpy := b.Clone()

	for r := 0; r < NUMBER_OF_ROWS; r++ {
		for c := 0; c < NUMBER_OF_COLS; c++ {

			if b.Get(r, c) == UNASSIGNED { // for every unassigned cell in the board

				possibilities := cpy.PossibleCells(r, c) // get all the potential numbers that can go in that cell

				rand := randomInt(0, len(possibilities)+NUMBER_OF_CHANCES_FOR_UNASSIGNED) // randomly assign a number or set unassigned

				var valueToAdd uint8

				if rand < len(possibilities) {

					valueToAdd = possibilities[rand]

				} else {
					valueToAdd = 0
				}

				chromosome.genes[c+(r*9)] = geneToNum(valueToAdd)
				cpy.Set(r, c, valueToAdd)

			} else {
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
