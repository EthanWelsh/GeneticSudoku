package main

import (
	"math"
)

type Spinner struct {
	score []int
	wheel [1000]int
	genes []Gene
}

// add genes to the spinner to be randomly selected from later
func (s *Spinner) addOptions(g []Gene) {

	s.genes = g

	score := make([]int, len(g))

	var total float64

	for i, gene := range g {

		score[i] = gene.Score()
		total += float64(score[i])
	}

	var chance float64
	wheelPos := 0
	j := 0

	for i, geneScore := range score {
		chance = Round(float64(geneScore)/total, 1, 3)
		spotsOnWheel := int(chance * 1000)

		for j = wheelPos; j < spotsOnWheel+wheelPos; j++ {
			s.wheel[j] = i
		}
		wheelPos = j
	}
}

func (s *Spinner) Spin() Gene {

	randomIndexInWheel := random(0, len(s.wheel))
	indexOfGene := s.wheel[randomIndexInWheel]
	randomGene := s.genes[indexOfGene]

	return randomGene
}

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
