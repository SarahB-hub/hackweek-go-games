package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell"
)

const paddleHeight = 6
const initialBallVelocityCol = 2
const initialBallVelocityRow = 1

const paddleSymbol = 0x2588
const ballSymbol = 0x25CF

type GameObject struct {
	col, row, width, height int
	velCol, velRow          int
	symbol                  rune
}

var screen tcell.Screen
var player1Paddle *GameObject
var player2Paddle *GameObject
var ball *GameObject
var gameObjects []*GameObject //initialising empty slice to fill once each object has been constructed
//from gameObject struct

func main() {

	initScreen()
	initGameState()
	inputChan := initUserInput()
	timePause := 75000

	for {
		handleUserInput(readInput(inputChan))
		updateState()
		drawState()

		time.Sleep(time.Duration(timePause) * time.Microsecond)
		if timePause > 25000 {
			timePause = timePause - 100 //keeps getting faster up to a certain point
		}
	}
}

func initScreen() {
	var err error
	screen, err = tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		gameExit()
	}
	if err := screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		gameExit()
	}

	defStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)
	screen.SetStyle(defStyle)

}

func initGameState() {
	width, height := screen.Size()
	paddleStartHeight := height/2 - paddleHeight/2

	player1Paddle = &GameObject{
		col:    0,
		row:    paddleStartHeight,
		width:  2,
		height: paddleHeight,
		velCol: 0,
		velRow: 0,
		symbol: paddleSymbol,
	}

	player2Paddle = &GameObject{
		col:    width - 2,
		row:    paddleStartHeight,
		width:  2,
		height: paddleHeight,
		velCol: 0,
		velRow: 0,
		symbol: paddleSymbol,
	}

	ball = &GameObject{
		col:    width / 2,
		row:    height / 2,
		width:  1,
		height: 1,
		velCol: initialBallVelocityCol,
		velRow: initialBallVelocityRow,
		symbol: ballSymbol,
	}
	gameObjects = []*GameObject{
		player1Paddle, player2Paddle, ball,
	} //filling empty slice now we've set the object structures
}

func printScreen(column, row, width, height int, ch rune) {
	for r := 0; r < height; r++ { //for loop starting at r = 0, running as long as r is less than
		//height, and adding 1 to r each iteration.
		for c := 0; c < width; c++ {
			screen.SetContent(column+c, row+r, ch, nil, tcell.StyleDefault)
			//print our char along the length of the row + the r iteration, and same for width
		}
	}
}

func printString(col, row int, str string) {
	for _, char := range str {
		screen.SetContent(col, row, char, nil, tcell.StyleDefault)
		col += 1
	}
}

func updateState() {
	for i := range gameObjects {
		gameObjects[i].row += gameObjects[i].velRow
		gameObjects[i].col += gameObjects[i].velCol
	}

	if wallCollision(ball) {
		ball.velRow = -ball.velRow //Passes the obj info to wallCollision, which will flip to true
		//and return the bool if the ball's row location matches the top or bottom of the screen. It then inverts the
		//row trajectory.
	}
	if paddleCollision(ball, player1Paddle) || paddleCollision(ball, player2Paddle) {
		ball.velCol = -ball.velCol
	}

	if winCondition(ball, player1Paddle, player2Paddle) == 1 {
		screen.Clear()
		printWinner(2)
		screen.Show()
		time.Sleep(5 * time.Second)
		gameExit()

	}
	if winCondition(ball, player1Paddle, player2Paddle) == 2 {
		screen.Clear()
		printWinner(1)
		screen.Show()
		time.Sleep(5 * time.Second)
		gameExit()
	}
}

func drawState() {
	screen.Clear()
	for _, object := range gameObjects {
		printScreen(object.col, object.row, object.width, object.height, object.symbol)
	}

	screen.Show()
}

func wallCollision(object *GameObject) bool {
	_, screenHeight := screen.Size()
	if object.row+object.velRow >= 0 && object.row+object.velRow < screenHeight {
		return false
	}
	return true
}

func paddleCollision(ball *GameObject, paddle *GameObject) bool {
	var collidesOnColumn bool
	if ball.col < paddle.col {
		collidesOnColumn = ball.col+ball.velCol >= paddle.col
	} else {
		collidesOnColumn = ball.col+ball.velCol <= paddle.col
	}

	return collidesOnColumn &&
		ball.row+ball.velRow >= paddle.row &&
		ball.row+ball.velRow < paddle.row+paddle.height
}

func initUserInput() chan string {
	inputChan := make(chan string) //creating a channel that holds string inputs we can use to return
	go func() {                    //this, note the () closer, is a background process.
		for {
			switch ev := screen.PollEvent().(type) {
			case *tcell.EventResize: // the event type is a resize, do this
				screen.Sync()
				drawState()
			case *tcell.EventKey: // if the event type is a key click, do that
				inputChan <- ev.Name()
			}
		}
	}()

	return inputChan
}

func readInput(inputChan chan string) string {
	var key string //this leaves the key var as the default empty, unless a new command
	//arrives from the input channel created. That way the rest of the program can keep
	//running regardless of waiting on input
	select {
	case key = <-inputChan:
	default:
		key = ""
	}

	return key
}

func handleUserInput(key string) {
	_, screenHeight := screen.Size()
	if key == "Rune[q]" {
		screen.Fini()
		gameExit()
	} else if key == "Escape" {
		screen.Fini()
		gameExit()
	} else if key == "Rune[w]" && player1Paddle.row > 0 {
		player1Paddle.row = player1Paddle.row - 2
	} else if key == "Rune[s]" && player1Paddle.row+player1Paddle.height < screenHeight {
		player1Paddle.row = player1Paddle.row + 2
	} else if key == "Up" && player2Paddle.row > 0 {
		player2Paddle.row = player2Paddle.row - 2
	} else if key == "Down" && player2Paddle.row+player2Paddle.height < screenHeight {
		player2Paddle.row = player2Paddle.row + 2
	}
}

func winCondition(ball *GameObject, player1Paddle *GameObject, player2Paddle *GameObject) int {
	if ball.col < player1Paddle.col {
		return 1
	} else if ball.col > player2Paddle.col {
		return 2
	}
	return 0 //default return meaning no win state yet
}

func gameExit() {
	screen.Fini()
	os.Exit(0)
}

func printWinner(winner int) {
	w, h := screen.Size()
	printString(w/2-5, h/2, "Game Over!")
	printString(w/2-11, h/2+2, fmt.Sprintf("Player %d is the winner", winner))
	screen.Show()
}
