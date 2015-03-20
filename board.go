package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"unicode"
)

var UNASSIGNED int = 0

type Board struct {
	board [9][9]int
}

type Position struct {
	row int
	col int
}

func (pos *Position) Set(r int, c int) {
	pos.row = r
	pos.col = c
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

// Self describing!
func (b *Board) Clone() Board {
	new_board := Board{}
	for i, _ := range b.board {
		for j, _ := range b.board[i] {
			new_board.board[i][j] = b.board[i][j]
		}
	}

	return new_board
}

// Assigns a cell the given number, returns a new board because functional programming.
func (b *Board) Assign(pos Position, num int) Board {
	new_board := b.Clone()
	new_board.board[pos.row][pos.col] = num
	return new_board
}

// Checks to see if all cells have an assigned value. Complete =/= Correct.
// TODO: Make an IsCorrect function dumby.
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

// Finds the first {row, col} that is "empty" or unassigned
func (board *Board) findUnassignedPosition() (position Position) {
	for i, row := range board.board {
		for j, cell := range row {
			if cell == UNASSIGNED {
				return Position{i, j}
			}
		}
	}

	return Position{-1, -1}
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
