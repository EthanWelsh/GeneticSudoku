package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"unicode"
)

type Board struct {
	board *Chromosome
}

// Initializes a blank, unassigned sudoku board.
func Init() (new_board Board) {

	var newChromosome Chromosome

	new_board.board = &newChromosome

	for r := 0; r < NUMBER_OF_ROWS; r++ {
		for c := 0; c < NUMBER_OF_COLS; c++ {
			new_board.Set(r, c, UNASSIGNED)
		}
	}

	return new_board
}

// Reads a board in from file and returns it
func BoardParser(filename string) (board Board, genesThatCanBeMutated []int) {
	board = Init()
	data, err := ioutil.ReadFile(filename)

	if err != nil {
		fmt.Println("Could not find", filename)
		os.Exit(0)
	}

	counter := 0
	row := 0
	col := 0
	for counter < len(data) {
		if unicode.IsDigit(rune(data[counter])) {

			num, _ := strconv.Atoi(string(data[counter]))
			board.Set(row, col, uint8(num))
			col++
		} else if data[counter] == '\n' {
			row++
			col = 0
		} else if data[counter] == '-' {
			board.Set(row, col, UNASSIGNED)
			col++
		}
		counter++
	}

	for r := 0; r < NUMBER_OF_ROWS; r++ {
		for c := 0; c < NUMBER_OF_COLS; c++ {
			possibles := board.PossibleCells(r, c)

			if len(possibles) == 1 {
				board.Set(r, c, possibles[0])
			} else {
				genesThatCanBeMutated = append(genesThatCanBeMutated, col+(row*9))
			}
		}
	}

	return
}

// Grade the given board on its completion
func (b *Board) Grade() (score float64) {

	//    Scoring Criteria | Score    |
	//    -----------------------------
	//    Assigned Square  |        1 |
	//    Unique Row       |        3 |
	//    Unique Column    |        3 |
	//    Unique Box       |        3 |

	_, errorCount := b.IsWrong()

	score += (float64(NUMBER_OF_ROWS+NUMBER_OF_COLS+NUMBER_OF_BOXES) - float64(errorCount)) * ERROR_MODIFIER

	for i := 0; i < NUMBER_OF_ROWS; i++ {

		// award REWARD_FOR_COMPLETE_BOARD_ELEMENT points for each complete row
		if b.isUniqueRow(i) {
			score += REWARD_FOR_COMPLETE_BOARD_ELEMENT
		}

		// award REWARD_FOR_COMPLETE_BOARD_ELEMENT points for each complete column
		if b.isUniqueColumn(i) {
			score += REWARD_FOR_COMPLETE_BOARD_ELEMENT
		}
	}

	for r := 0; r < NUMBER_OF_ROWS; r += int(math.Sqrt(NUMBER_OF_ROWS)) {
		for c := 0; c < NUMBER_OF_COLS; c += int(math.Sqrt(NUMBER_OF_COLS)) {
			// award REWARD_FOR_COMPLETE_BOARD_ELEMENT points for each complete box
			if b.isUniqueBox(r, c) {
				score += REWARD_FOR_COMPLETE_BOARD_ELEMENT
			}
		}
	}

	return
}

// A small utility function for checking if the row of a given board allows that number in it
func (b *Board) uniqueRows(possible_num uint8, row int) bool {

	for c := 0; c < NUMBER_OF_COLS; c++ {
		if b.Get(row, c) == possible_num {
			return false
		}
	}

	return true
}

// checks to see if given row is unique
func (b *Board) isUniqueRow(r int) bool {
	counter := make([]int, NUMBER_OF_ROWS)

	for c := 0; c < NUMBER_OF_COLS; c++ {

		if b.Get(r, c) == 0 || counter[b.Get(r, c)-1] >= 1 {
			return false
		}
		counter[b.Get(r, c)-1]++
	}
	return true
}

// Refer to uniqueRows, except columns
func (b *Board) uniqueColumns(possible_num uint8, col int) bool {

	for r := 0; r < NUMBER_OF_ROWS; r++ {
		if b.Get(col, r) == possible_num {
			return false
		}
	}

	return true
}

// checks to see if given column is unique
func (b *Board) isUniqueColumn(c int) bool {
	counter := make([]int, NUMBER_OF_COLS)

	for r := 0; r < NUMBER_OF_ROWS; r++ {

		if b.Get(r, c) == 0 || counter[b.Get(r, c)-1] >= 1 {
			return false
		}
		counter[b.Get(r, c)-1]++
	}
	return true
}

// A small utility function for checking if the box of a cell is unique based on the cells around it.
func (b *Board) uniqueBox(possible_num uint8, row int, col int) bool {

	spaceToNextBox := int(math.Sqrt(NUMBER_OF_ROWS))

	starting_row := (row / spaceToNextBox) * spaceToNextBox
	starting_col := (col / spaceToNextBox) * spaceToNextBox
	ending_row := starting_row + spaceToNextBox
	ending_col := starting_col + spaceToNextBox

	for i := starting_row; i < ending_row; i++ {
		for j := starting_col; j < ending_col; j++ {
			if b.Get(i, j) == uint8(possible_num) {
				return false
			}
		}
	}
	return true
}

// checks to see if given box is unique
func (b *Board) isUniqueBox(row int, col int) bool {

	counter := make([]int, NUMBER_OF_ROWS)

	spaceToNextBox := int(math.Sqrt(NUMBER_OF_ROWS))

	starting_row := row
	starting_col := col
	ending_row := starting_row + spaceToNextBox
	ending_col := starting_col + spaceToNextBox

	for r := starting_row; r < ending_row; r++ {
		for c := starting_col; c < ending_col; c++ {

			if b.Get(r, c) == 0 {
				return false
			}
			if counter[b.Get(r, c)-1] >= 1 {
				return false
			}

			counter[b.Get(r, c)-1]++
		}
	}
	return true
}

// PossibleBoard returns whether the board is solveable
func (b Board) PossibleBoard() bool {

	for r := 0; r < NUMBER_OF_ROWS; r++ {
		for c := 0; c < NUMBER_OF_COLS; c++ {
			if b.Get(r, c) == UNASSIGNED && len(b.PossibleCells(r, c)) == 0 {
				return false
			}
		}
	}

	return true
}

// Returns a slice of possible numbers that are allowed to be assigned in a board at the given position
func (b *Board) PossibleCells(row int, col int) (possibles []uint8) {
	possibles = []uint8{}

	var i uint8

	for i = 1; i <= NUMBER_OF_ROWS; i++ {
		if b.uniqueRows(i, row) && b.uniqueColumns(i, col) && b.uniqueBox(i, row, col) {
			possibles = append(possibles, uint8(i))
		}
	}

	return possibles
}

// Checks to see if a given board violates any of the rules of Sudoku
func (b Board) IsWrong() (ret bool, errorCount int) {

	ret = false

	// If there are duplicated within the row,
	for i := 0; i < NUMBER_OF_ROWS; i++ {
		nums := b.GetNumbersInRow(i)
		if containsDuplicates(nums) {
			ret = true
			errorCount++
		}
	}

	// column,
	for i := 0; i < NUMBER_OF_COLS; i++ {
		nums := b.GetNumbersInCol(i)
		if containsDuplicates(nums) {
			ret = true
			errorCount++
		}
	}

	// or box, then this is an invalid board
	for r := 0; r < NUMBER_OF_ROWS; r += int(math.Sqrt(NUMBER_OF_ROWS)) {
		for c := 0; c < NUMBER_OF_COLS; c += int(math.Sqrt(NUMBER_OF_COLS)) {
			nums := b.GetNumbersInBox(r, c)
			if containsDuplicates(nums) {
				ret = true
				errorCount++
			}
		}
	}

	return
}

// Checks to see if all cells have an assigned value. Complete =/= Correct.
func (b *Board) IsComplete() bool {

	for r := 0; r < NUMBER_OF_ROWS; r++ {
		for c := 0; c < NUMBER_OF_COLS; c++ {
			if b.Get(r, c) == UNASSIGNED {
				return false
			}
		}
	}

	return true
}

// Checks to see if the board represents a complete and correct solution
func (b *Board) IsCorrect() bool {
	return b.Grade() == (REWARD_FOR_COMPLETE_BOARD_ELEMENT*NUMBER_OF_ROWS)+ // complete rows
		(REWARD_FOR_COMPLETE_BOARD_ELEMENT*NUMBER_OF_COLS)+ // complete cols
		(REWARD_FOR_COMPLETE_BOARD_ELEMENT*NUMBER_OF_BOXES)+ // complete boxes
		((NUMBER_OF_ROWS+NUMBER_OF_COLS+NUMBER_OF_BOXES)*ERROR_MODIFIER) // lack of errors
}

// Get all assigned numbers in a given row
func (b *Board) GetNumbersInRow(rowNum int) (row []uint8) {

	for i := 0; i < NUMBER_OF_ROWS; i++ {
		x := b.Get(rowNum, i)
		if x != UNASSIGNED {
			row = append(row, x)
		}
	}

	return
}

// Get all assigned numbers in a given column
func (b *Board) GetNumbersInCol(colNum int) (col []uint8) {

	for i := 0; i < NUMBER_OF_COLS; i++ {

		x := b.Get(i, colNum)

		if x != UNASSIGNED {
			col = append(col, x)
		}
	}

	return
}

// Get all assigned numbers in a given box
func (b *Board) GetNumbersInBox(r int, c int) (box []uint8) {

	spaceToNextBox := int(math.Sqrt(NUMBER_OF_ROWS))

	starting_row := r
	starting_col := c
	ending_row := starting_row + spaceToNextBox
	ending_col := starting_col + spaceToNextBox

	for r := starting_row; r < ending_row; r++ {
		for c := starting_col; c < ending_col; c++ {

			x := b.Get(r, c)

			if x != UNASSIGNED {
				box = append(box, x)
			}
		}
	}

	return
}

// Given an index, will return the R and C values for the index element in the board, counting by row
func (b *Board) GetCellByRow(index int) (r, c uint8) {

	r = uint8(index / NUMBER_OF_ROWS)
	c = uint8(index % NUMBER_OF_COLS)

	return
}

// Given an index, will return the R and C values for the index element in the board, counting by col
func (b *Board) GetCellByCol(index int) (r, c uint8) {
	c = uint8(index / NUMBER_OF_COLS)
	r = uint8(index % NUMBER_OF_ROWS)

	return
}

// Given an index, will return the R and C values for the index element in the board, counting by col
// Boxes are read likeso:
// 0  1  2   9  10 11
// 3  4  5   12 13 14
// 6  7  8   15 16 17
func (b *Board) GetCellByBox(index uint8) (r, c uint8) {

	boxNum := uint8(index / NUMBER_OF_ROWS) //0-8

	sizeOfBox := uint8(math.Sqrt(NUMBER_OF_ROWS))

	startR := (boxNum * sizeOfBox) / NUMBER_OF_ROWS
	startC := uint8((boxNum % sizeOfBox) * sizeOfBox)

	index = index % NUMBER_OF_ROWS

	r = startR + uint8(index/sizeOfBox)
	c = startC + uint8(index%sizeOfBox)

	return

}

// Prints the board!
func (b *Board) Print() {

	for r := 0; r < NUMBER_OF_ROWS; r++ {
		for c := 0; c < NUMBER_OF_COLS; c++ {
			cell := b.Get(r, c)

			if cell == UNASSIGNED {
				fmt.Print(" - ")
			} else {
				fmt.Print(" ", cell, " ")
			}
		}
		fmt.Println()
	}
}

// Checks to see if an array contains any duplicate values
func containsDuplicates(arr []uint8) bool {

	nums := make(map[uint8]bool)

	for i := 0; i < len(arr); i++ {
		if _, ok := nums[arr[i]]; ok {
			return true
		} else {
			nums[arr[i]] = true
		}
	}

	return false

}

// Returns the integer at the given location of the board
func (b *Board) Get(r int, c int) uint8 {
	return b.board.genes[c+(r*NUMBER_OF_ROWS)]
}

// Sets a given location on the board to a certain integer
func (b *Board) Set(r int, c int, value uint8) {
	b.board.genes[c+(r*NUMBER_OF_ROWS)] = value
}
