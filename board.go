package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"unicode"
)

var UNASSIGNED int = 0
var COMPLETE_REWARD = 3

type Board struct {
	board [9][9]int
}

type Position struct {
	row int
	col int
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

// Fitness function for grading a board
//    Scoring Criteria | Score    |
//	  -----------------------------
//	  Assigned Square  |        1 |
//    Unique row 	   |		3 |
//	  Unique column    |        3 |
//	  Unique box       |        3 |

func (b *Board) Grade() (score int) {

	for r := 0; r < NUMBER_OF_ROWS; r++ {
		for c := 0; c < NUMBER_OF_COLS; c++ {
			if b.Get(r, c) != UNASSIGNED {
				score++
			}
		}
	}

	for i := 0; i < 9; i++ {
		if b.isUniqueRow(i) {
			score += COMPLETE_REWARD
		}

		if b.isUniqueColumn(i) {
			score += COMPLETE_REWARD
		}
	}

	for i := 0; i < 9; i += 3 {
		for j := 0; j < 9; j += 3 {
			if b.isUniqueBox(i, j) {
				score += COMPLETE_REWARD
			}
		}
	}

	return
}

// A small utility function for checking if the row of a given board allows that number in it
// e.g. a potential row configuration:
//         1 2 3 4 - 5 6 7 9   -- if 8 is passed in, then returns true. If 9, false
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
// Look at sudoku rules for more information
func (b *Board) uniqueBox(possible_num int, pos Position) bool {
	// check the box using math!!!
	starting_row := (pos.row / 3) * 3
	starting_col := (pos.col / 3) * 3
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
			if cell == UNASSIGNED && len(b.PossibleCells(Position{i, j})) == 0 {
				return false
			}
		}
	}
	return true
}

// PossibleCells returns a slice of possible numbers that are allowed to be assigned
// in a board at the position input
func (b *Board) PossibleCells(pos Position) (possibles []int) {
	possibles = []int{}
	for i := 1; i <= 9; i++ {
		if b.uniqueRows(i, pos.row) && b.uniqueColumns(i, pos.col) && b.uniqueBox(i, pos) {
			possibles = append(possibles, i)
		}
	}

	return possibles
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

	return b.Grade() == (NUMBER_OF_ROWS*NUMBER_OF_COLS)+(COMPLETE_REWARD*NUMBER_OF_ROWS)+(COMPLETE_REWARD*NUMBER_OF_COLS)+(COMPLETE_REWARD*9)
}

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

func (b *Board) Get(r int, c int) int {
	return b.board[r][c]
}

func (b *Board) Set(r int, c int, value int) {
	b.board[r][c] = value
}
