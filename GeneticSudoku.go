package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"
)

const (
	CHANCE_TO_MODIFY_A_GENE           = .01
	POPULATION_SIZE                   = 1000
	NUMBER_OF_CHANCES_FOR_UNASSIGNED  = 5 // When a random chromosome is generated or a mutation occurs, how many chances should there be that the number will be UNASSIGNED?
	MATE_MAX_RETRIES                  = 10
	REWARD_FOR_COMPLETE_BOARD_ELEMENT = 3
	ITERATIONS                        = 100
	STEPS_PER_ITERATION               = 100
	AVAILABLE_MODIFIER                = 3

	CHROMOSOME_SIZE = NUMBER_OF_ROWS * NUMBER_OF_COLS

	NUMBER_OF_ROWS  = 9
	NUMBER_OF_COLS  = 9
	NUMBER_OF_BOXES = 9

	UNASSIGNED = 0
)

var boardCache map[string]Board

func main() {

	rand.Seed(int64(time.Now().UnixNano()))

	scoreCache = make(map[string]float64)
	boardCache = make(map[string]Board)

	defer un(trace("BASELINE"))

	startBoard := BoardParser("src/main/boards/board.txt")

	population := make([]Chromosome, POPULATION_SIZE)

	// Generate random partial solutions to the given board
	for i := range population {
		population[i] = GetRandomChromosome(&startBoard)
	}

	for i := 0; i < ITERATIONS; i++ {

		avg, max, min := getPopulationStats(population)
		fmt.Printf("%d).\t\t\tAVG: %.2f\t\tMAX: %d\t\tMIN: %d\n", i*STEPS_PER_ITERATION, avg, max, min)
		population = evolve(population, STEPS_PER_ITERATION, CHANCE_TO_MODIFY_A_GENE)

		popMax := 0.0
		popMaxInt := 0

		for j, c := range population {

			if c.Score() > popMax {
				popMax = c.Score()
				popMaxInt = j
			}
		}

		b := getBoardFromChromosome(population[popMaxInt])
		b.Print()

	}
}

// Performs reproduction and mutations for a given number of iterations and returns the resulting population
func evolve(population []Chromosome, iterations int, chanceAtMutation float64) []Chromosome {

	for i := 0; i < iterations; i++ {
		population = getNextGeneration(population)
	}

	population = Mutate(population, chanceAtMutation)

	return population

}

// Performs reproduction and returns the resulting population
func getNextGeneration(oldPopulation []Chromosome) (newPopulation []Chromosome) {

	var randomChromosomeSelector Spinner
	randomChromosomeSelector.addOptions(oldPopulation)

	newPopulation = make([]Chromosome, POPULATION_SIZE)

	for i := range newPopulation {

		// Get mating partner A & B
		phenotypeA := randomChromosomeSelector.Spin()
		phenotypeB := randomChromosomeSelector.Spin()

		// Mate them and add their child to the new population
		newPopulation[i] = MateChromosome(phenotypeA, phenotypeB)

	}

	return newPopulation
}

// Given a specific chromosome, will get the board for that chromosome
func getBoardFromChromosome(chromosome Chromosome) Board {

	gs := chromosome.String()

	if val, ok := boardCache[gs]; ok {
		return val
	} else {
		board := Init()
		var index int

		for r := 0; r < NUMBER_OF_ROWS; r++ {
			for c := 0; c < NUMBER_OF_COLS; c++ {
				board.Set(r, c, geneToNum(chromosome.genes[index]))
				index++
			}
		}

		boardCache[gs] = board

		return board
	}
}

// Provide the average, maximum, and minimum board scores in the population
func getPopulationStats(population []Chromosome) (avg float64, max uint64, min uint64) {

	var total uint64 = 0
	var chromosomeScore uint64 = 0

	max = 0
	min = math.MaxUint64

	for _, chromosome := range population {

		chromosomeScore = uint64(chromosome.Score())

		total += chromosomeScore

		if chromosomeScore > max {
			max = chromosomeScore
		}

		if chromosomeScore < min {
			min = chromosomeScore
		}
	}

	avg = float64(total) / float64(len(population))

	return

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
