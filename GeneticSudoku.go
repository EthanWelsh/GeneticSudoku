package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"
)

const (
	UNASSIGNED                        = 0
	NUMBER_OF_ROWS                    = 9
	NUMBER_OF_COLS                    = 9
	NUMBER_OF_BOXES                   = 9
	CHROMOSOME_SIZE                   = NUMBER_OF_ROWS * NUMBER_OF_COLS
	POPULATION_SIZE                   = 1000
	NUMBER_OF_CHANCES_FOR_UNASSIGNED  = 5 // When a random gene or mutation occurs, how many chances should there be that the number will be UNASSIGNED?
	MATE_MAX_RETRIES                  = 10
	REWARD_FOR_COMPLETE_BOARD_ELEMENT = 3
)

var boardCache map[string]Board

func main() {

	rand.Seed(int64(time.Now().UnixNano()))

	scoreCache = make(map[string]int)
	boardCache = make(map[string]Board)

	defer un(trace("BASELINE"))

	startBoard := BoardParser("/Users/welshej/github/GeneticSudoku/src/main/boards/board.txt")

	population := make([]Gene, POPULATION_SIZE)

	for i := range population {
		population[i] = getRandomGene(&startBoard)
	}

	for i := 0; i < 100; i++ {
		avg, max, min := getPopulationStats(population)
		fmt.Printf("%d).\t\t\tAVG: %.2f\t\tMAX: %d\t\tMIN: %d\n", i, avg, max, min)
		population = evolve(population, 100, .001)

		popMax := 0
		popMaxInt := 0

		for j, g := range population {

			if g.Score() > popMax {
				popMax = g.Score()
				popMaxInt = j
			}
		}

		b := getBoardFromGene(population[popMaxInt])
		b.Print()

	}
}

// Performs reproduction and mutations for a given number of iterations and returns the resulting population
func evolve(population []Gene, iterations int, chanceAtMutation float64) []Gene {

	for i := 0; i < iterations; i++ {
		population = getNextGeneration(population)
	}

	Mutate(population, chanceAtMutation)

	return population

}

// Performs reproduction and returns the resulting population
func getNextGeneration(oldPopulation []Gene) (newPopulation []Gene) {

	var randomGeneSelector Spinner
	randomGeneSelector.addOptions(oldPopulation)

	newPopulation = make([]Gene, POPULATION_SIZE)

	for i := range newPopulation {

		phenotypeA := randomGeneSelector.Spin()
		phenotypeB := randomGeneSelector.Spin()

		newPopulation[i] = mateGenes(phenotypeA, phenotypeB)

	}

	return newPopulation
}

// Provide the average, maximum, and minimum board scores in the population
func getPopulationStats(population []Gene) (avg float64, max uint64, min uint64) {

	var total uint64 = 0
	var geneScore uint64 = 0

	max = 0
	min = math.MaxUint64

	for _, gene := range population {

		geneScore = uint64(gene.Score())

		total += geneScore

		if geneScore > max {
			max = geneScore
		}

		if geneScore < min {
			min = geneScore
		}
	}

	avg = float64(total) / float64(len(population))

	return

}

// Given a specific gene, will get the board for that gene
func getBoardFromGene(gene Gene) Board {

	gs := gene.String()

	if val, ok := boardCache[gs]; ok {
		return val
	} else {
		board := Init()
		var index int

		for r := 0; r < NUMBER_OF_ROWS; r++ {
			for c := 0; c < NUMBER_OF_COLS; c++ {
				board.Set(r, c, bitStringToNum(gene.gene[index]))
				index++
			}
		}

		boardCache[gs] = board

		return board
	}
}

// Generates a random integer between min and max (inclusive)
func randomInt(min int, max int) int {
	return rand.Intn(max) + min
}

// Generates a random float between min and max (inclusive)
func randomFloat(min float64, max float64) float64 {

	return rand.Float64()*(max-min) + min
}

// temporary timing function
func trace(s string) (string, time.Time) {
	log.Println("START:", s)

	return s, time.Now()
}

// temporary timing function
func un(s string, startTime time.Time) {
	endTime := time.Now()
	log.Println("  END:", s, "ElapsedTime in seconds:", endTime.Sub(startTime))
}
