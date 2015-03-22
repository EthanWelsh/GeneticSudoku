package main

import (
	"fmt"
	"log"
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

	scoreCache = make(map[string]int)
	boardCache = make(map[string]Board)

	defer un(trace("BASELINE"))

	startBoard := BoardParser("/Users/welshej/github/GeneticSudoku/src/main/boards/board.txt")

	population := make([]Gene, POPULATION_SIZE)

	for i := range population {
		population[i] = getRandomGene(&startBoard)
	}

	for i := 0; i < 10; i++ {
		avg, max, min := getPopulationStats(population)
		fmt.Printf("%d).\t\t\tAVG: %.2f\t\tMAX: %d\t\tMIN: %d\n", i, avg, max, min)
		population = evolve(population, 10, 0)

	}

}

func evolve(population []Gene, iterations int, chanceAtMutation float64) []Gene {

	for i := 0; i < iterations; i++ {
		population = getNextGeneration(population)

		for i := range population {
			population[i].Mutate(chanceAtMutation)
		}
	}

	return population

}

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

var boardCache map[string]Board

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

// Generates a random number between min and max (inclusive)
func random(min int, max int) int {
	return rand.Intn(max) + min
}

func trace(s string) (string, time.Time) {
	log.Println("START:", s)

	return s, time.Now()
}

func un(s string, startTime time.Time) {
	endTime := time.Now()
	log.Println("  END:", s, "ElapsedTime in seconds:", endTime.Sub(startTime))
}
