package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"unicode"
)

type Board struct {
	board [9][9]int
}

// initializes a blank, unassigned sudoku board.
func Init() Board {
	new_board := Board{}
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			new_board.board[i][j] = UNASSIGNED
		}
	}

	return new_board
}

// Reads a board in from file and returns it
func BoardParser(filename string) (board Board) {
	board = Init()
	data, _ := ioutil.ReadFile(filename)
	counter := 0
	row := 0
	col := 0
	for counter < len(data) {
		if unicode.IsDigit(rune(data[counter])) {
			board.board[row][col], _ = strconv.Atoi(string(data[counter]))
			col++
		} else if data[counter] == '\n' {
			row++
			col = 0
		} else if data[counter] == '-' {
			board.board[row][col] = UNASSIGNED
			col++
		}
		counter++
	}
	return
}

// Grade the given board on its completion
func (b *Board) Grade() (score int) {

	//    Scoring Criteria | Score    |
	//	  -----------------------------
	//	  Assigned Square  |        1 |
	//    Unique Row 	   |		3 |
	//	  Unique Column    |        3 |
	//	  Unique Box       |        3 |

	for r := 0; r < NUMBER_OF_ROWS; r++ {
		for c := 0; c < NUMBER_OF_COLS; c++ {
			if b.Get(r, c) != UNASSIGNED {
				score++
			}
		}
	}

	for i := 0; i < 9; i++ {
		if b.isUniqueRow(i) {
			score += REWARD_FOR_COMPLETE_BOARD_ELEMENT
		}

		if b.isUniqueColumn(i) {
			score += REWARD_FOR_COMPLETE_BOARD_ELEMENT
		}
	}

	for i := 0; i < 9; i += 3 {
		for j := 0; j < 9; j += 3 {
			if b.isUniqueBox(i, j) {
				score += REWARD_FOR_COMPLETE_BOARD_ELEMENT
			}
		}
	}

	return
}

// A small utility function for checking if the row of a given board allows that number in it
func (b *Board) uniqueRows(possible_num int, row int) bool {
	for _, cell := range b.board[row] {
		if possible_num == cell {
			return false
		}
	}
	return true
}

// checks to see if given row is unique
func (b *Board) isUniqueRow(r int) bool {
	counter := make([]int, 9)

	for c := 0; c < NUMBER_OF_COLS; c++ {

		if b.Get(r, c) == 0 || counter[b.Get(r, c)-1] >= 1 {
			return false
		}
		counter[b.Get(r, c)-1]++
	}
	return true
}

// Refer to uniqueRows, except columns
func (b *Board) uniqueColumns(possible_num int, column int) bool {
	for _, row := range b.board {
		for col, cell := range row {
			if col == column && possible_num == cell {
				return false
			}
		}
	}
	return true
}

// checks to see if given column is unique
func (b *Board) isUniqueColumn(c int) bool {
	counter := make([]int, 9)

	for r := 0; r < NUMBER_OF_ROWS; r++ {

		if b.Get(r, c) == 0 || counter[b.Get(r, c)-1] >= 1 {
			return false
		}
		counter[b.Get(r, c)-1]++
	}
	return true
}

// A small utility function for checking if the box of a cell is unique based on the cells around it.
func (b *Board) uniqueBox(possible_num int, row int, col int) bool {
	// check the box using math!!!
	starting_row := (row / 3) * 3
	starting_col := (col / 3) * 3
	ending_row := starting_row + 3
	ending_col := starting_col + 3

	for i := starting_row; i < ending_row; i++ {
		for j := starting_col; j < ending_col; j++ {
			if b.board[i][j] == possible_num {
				return false
			}
		}
	}
	return true
}

// checks to see if given box is unique
func (b *Board) isUniqueBox(row int, col int) bool {

	counter := make([]int, 9)

	starting_row := row
	starting_col := col
	ending_row := starting_row + 3
	ending_col := starting_col + 3

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
	for i, row := range b.board {
		for j, cell := range row {
			if cell == UNASSIGNED && len(b.PossibleCells(i, j)) == 0 {
				return false
			}
		}
	}
	return true
}

// Returns a slice of possible numbers that are allowed to be assigned in a board at the given position
func (b *Board) PossibleCells(row int, col int) (possibles []int) {
	possibles = []int{}
	for i := 1; i <= 9; i++ {
		if b.uniqueRows(i, row) && b.uniqueColumns(i, col) && b.uniqueBox(i, row, col) {
			possibles = append(possibles, i)
		}
	}

	return possibles
}

// Checks to see if a given board violates any of the rules of Sudoku
func (b Board) IsWrong() bool {

	for i := 0; i < 9; i++ {
		if !b.isUniqueRow(i) {
			return true
		}

		if !b.isUniqueColumn(i) {
			return true
		}
	}

	for i := 0; i < 9; i += 3 {
		for j := 0; j < 9; j += 3 {
			if !b.isUniqueBox(i, j) {
				return true
			}
		}
	}
	return false
}

// Checks to see if all cells have an assigned value. Complete =/= Correct.
func (b *Board) IsComplete() bool {
	for _, row := range b.board {
		for _, cell := range row {
			if cell == UNASSIGNED {
				return false
			}
		}
	}
	return true
}

// Checks to see if the board represents a complete and correct solution
func (b *Board) IsCorrect() bool {
	return b.Grade() == (NUMBER_OF_ROWS*NUMBER_OF_COLS)+(REWARD_FOR_COMPLETE_BOARD_ELEMENT*NUMBER_OF_ROWS)+(REWARD_FOR_COMPLETE_BOARD_ELEMENT*NUMBER_OF_COLS)+(REWARD_FOR_COMPLETE_BOARD_ELEMENT*NUMBER_OF_BOXES)
}

// Return deep copy of given board
func (b *Board) Clone() Board {
	new_board := Board{}
	for i, _ := range b.board {
		for j, _ := range b.board[i] {
			new_board.board[i][j] = b.board[i][j]
		}
	}

	return new_board
}

// Returns the integer at the given location of the board
func (b *Board) Get(r int, c int) int {
	return b.board[r][c]
}

// Sets a given location on the board to a certain integer
func (b *Board) Set(r int, c int, value int) {
	b.board[r][c] = value
}

// Prints the board!
func (b *Board) Print() {
	for _, row := range b.board {
		for _, cell := range row {
			if cell == UNASSIGNED {
				fmt.Print(" - ")
			} else {
				fmt.Print(" ", cell, " ")
			}

		}
		fmt.Println()
	}
}
