package main

var scoreCache map[string]int

type Chromosome struct {
	genes [CHROMOSOME_SIZE]string
}

// Fitness function used to determine the degree of completion of the board
func (c Chromosome) Score() int {
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

		var valueToAdd int

		if rand < len(possibilities) {
			valueToAdd = possibilities[rand]

		} else {
			valueToAdd = UNASSIGNED
		}

		// add the mutation
		population[modifiedChromosome].genes[modifiedGene] = NumToBitString(valueToAdd)
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
func GetRandomChromosome(b *Board) (c Chromosome) {

	cpy := b.Clone()

	for r := 0; r < NUMBER_OF_ROWS; r++ {
		for c := 0; c < NUMBER_OF_COLS; c++ {

			if b.Get(r, c) == UNASSIGNED { // for every unassigned cell in the board

				possibilities := cpy.PossibleCells(r, c) // get all the potential numbers that can go in that cell

				rand := randomInt(0, len(possibilities)+NUMBER_OF_CHANCES_FOR_UNASSIGNED) // randomly assign a number or set unassigned

				var valueToAdd int

				if rand < len(possibilities) {

					valueToAdd = possibilities[rand]

				} else {
					valueToAdd = 0
				}

				c.genes[c+(r*9)] = NumToBitString(valueToAdd)
				cpy.Set(r, c, valueToAdd)

			} else {
				c.genes[c+(r*9)] = NumToBitString(b.Get(r, c))
			}
		}
	}

	return c
}

// Given a bit string, will provide the number which maps to that bit string
func BitStringToNum(s string) int {
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
func NumToBitString(num int) string {
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

// Returns the string representation of a particular chromosome
func (c *Chromosome) String() (ret string) {
	for _, s := range c.genes {
		ret += s
	}
	return
}
