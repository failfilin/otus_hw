package main

import (
	"fmt"
	"strings"
)

func main() {
	boardSize := 15
	ChessBoard := strings.Builder{}
	for i := 0; i < boardSize; i++ {
		switch i % 2 {
		case 0:
			ChessBoard.WriteString(chessboardLineGeneretor(boardSize, 0))
		case 1:
			ChessBoard.WriteString(chessboardLineGeneretor(boardSize, 1))
		}
	}
	fmt.Println(ChessBoard.String())
}

func chessboardLineGeneretor(size int, startType int) string { // количество элементов в линии и тип первого элемента 0 - пробел 1 - решетка
	space := " " // подменил на пробел для простоты визуала
	lattice := "#"
	boardline := strings.Builder{}
	if startType == 0 {
		for i := 0; i < size; i++ {
			switch i % 2 {
			case 0:
				boardline.WriteString(space)
			case 1:
				boardline.WriteString(lattice)
			}
		}
	} else {
		for i := 0; i < size; i++ {
			switch i % 2 {
			case 0:
				boardline.WriteString(lattice)
			case 1:
				boardline.WriteString(space)
			}
		}
	}
	boardline.WriteString("\n")
	return boardline.String()
}
