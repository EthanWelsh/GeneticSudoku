package main

import (
	"math"
)

const SPOTS_ON_WHEEL = 100000

type Spinner struct {
	score       []float64
	wheel       [SPOTS_ON_WHEEL]int
	chromosomes []Chromosome
}

// Add chromosomes to the spinner to be randomly selected from later
func (s *Spinner) addOptions(c []Chromosome) {

	s.chromosomes = c

	score := make([]float64, len(c))

	var total float64

	// Remember the scores of each of the chromosomes and keep track of the over score for the population
	for i, chromosome := range c {
		score[i] = chromosome.Score()
		total += float64(score[i])
	}

	var chance float64
	wheelPos := 0
	j := 0

	// For every chromosome, determine how many spots that chromosome should get on the wheel
	for i, chromosomeScore := range score {
		chance = Round(float64(chromosomeScore)/total, 1, len(string(SPOTS_ON_WHEEL))-1)
		spotsOnWheel := int(chance * SPOTS_ON_WHEEL)

		for j = wheelPos; j < spotsOnWheel+wheelPos; j++ {
			s.wheel[j] = i
		}
		wheelPos = j
	}
}

// Randomly picks a chromosomes to reproduce, giving preference to those chromosomes with a high fitness value
func (s *Spinner) Spin() Chromosome {

	randomIndexInWheel := randomInt(0, len(s.wheel))
	indexOfChromosome := s.wheel[randomIndexInWheel]
	randomChromosome := s.chromosomes[indexOfChromosome]

	return randomChromosome
}

// Self-explanatory!
func Round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}
