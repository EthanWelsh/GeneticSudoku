package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

const NUMBER_OF_COLS = 9
const NUMBER_OF_ROWS = 9

// Whenever a random gene or mutation occurs, how many chances should there be that the number will UNASSIGNED?
var NUMBER_OF_CHANCES_FOR_UNASSIGNED = 6

var POPULATION_SIZE = 1000

func main() {

	rand.Seed(int64(time.Now().Unix()))

	startBoard := BoardParser("/Users/welshej/github/GeneticSudoku/src/main/boards/board.txt")

	population := make([]Gene, POPULATION_SIZE)

	for i := range population {
		population[i] = getRandomGene(startBoard)
	}

	for i := 0; i < 10; i++ {
		avg, max, min := getPopulationStats(population)
		fmt.Printf("%d).\t\t\tAVG: %.2f\t\tMAX: %d\t\tMIN: %d\n", i, avg, max, min)
		population = evolve(population, 10, .001)
	}

}

func evolve(population []Gene, iterations int, chanceAtMutation float64) []Gene {

	mutationsMade := 0

	for i := 0; i < iterations; i++ {
		population = getNextGeneration(population)

		for i := range population {
			population[i].Mutate(chanceAtMutation)
		}
	}

	fmt.Println("--", mutationsMade, "--")

	return population
}

func getNextGeneration(oldPopulation []Gene) (newPopulation []Gene) {
	var randomGeneSelector Spinner

	randomGeneSelector.addOptions(oldPopulation)

	newPopulation = make([]Gene, POPULATION_SIZE)

	for i := range newPopulation {

		phenotypeA := randomGeneSelector.Spin()
		phenotypeB := randomGeneSelector.Spin()

		newPopulation[i] = mateGenes(&phenotypeA, &phenotypeB)
	}

	return
}

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
func getBoardFromGene(gene *Gene) Board {

	board := Init()
	var index int

	for r := 0; r < NUMBER_OF_ROWS; r++ {
		for c := 0; c < NUMBER_OF_COLS; c++ {
			board.Set(r, c, bitStringToNum(gene.gene[index]))
			index++
		}
	}

	return board
}

// Generates a random number between min and max (inclusive)
func random(min int, max int) int {
	return rand.Intn(max) + min
}
