package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"math/rand"
	"time"
)

type Position struct {
	x, y int
}

type Game struct {
	board      [8][8]rune
	player     Position
	emptySpace Position
	score      int
}

func getRandomPosition() Position {
	return Position{
		x: rand.Intn(8),
		y: rand.Intn(8),
	}
}

func getValidRandomPositions() (Position, Position) {
	player := getRandomPosition()

	var emptySpace Position
	for {
		emptySpace = getRandomPosition()
		dx := abs(player.x - emptySpace.x)
		dy := abs(player.y - emptySpace.y)
		if dx+dy >= 3 {
			break
		}
	}

	return player, emptySpace
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func newGame(score int) *Game {
	game := &Game{}
	game.score = score

	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			game.board[i][j] = 'x'
		}
	}

	game.player, game.emptySpace = getValidRandomPositions()

	game.board[game.player.y][game.player.x] = 'o'
	game.board[game.emptySpace.y][game.emptySpace.x] = ' '

	return game
}

func (g *Game) draw() {
	_ = termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	drawString(0, 0, "Use arrow keys or WASD to move (ESC to quit)", termbox.ColorWhite, termbox.ColorDefault)
	drawString(0, 1, "Reach the empty space to win!", termbox.ColorWhite, termbox.ColorDefault)
	drawString(0, 2, fmt.Sprintf("Score: %d", g.score), termbox.ColorYellow, termbox.ColorDefault)

	drawString(0, 3, "------------------------", termbox.ColorWhite, termbox.ColorDefault)

	for i := 0; i < 8; i++ {
		termbox.SetCell(0, i+4, '|', termbox.ColorWhite, termbox.ColorDefault)

		for j := 0; j < 8; j++ {
			char := g.board[i][j]
			fg := termbox.ColorWhite
			if char == 'o' {
				fg = termbox.ColorGreen
			}
			termbox.SetCell(2+j*2, i+4, char, fg, termbox.ColorDefault)
		}

		termbox.SetCell(17, i+4, '|', termbox.ColorWhite, termbox.ColorDefault)
	}

	drawString(0, 12, "------------------------", termbox.ColorWhite, termbox.ColorDefault)

	_ = termbox.Flush()
}

func (g *Game) movePlayer(dx, dy int) bool {
	newX := g.player.x + dx
	newY := g.player.y + dy

	if newX < 0 || newX >= 8 || newY < 0 || newY >= 8 {
		return false
	}

	g.board[g.player.y][g.player.x] = 'x'
	g.player.x = newX
	g.player.y = newY
	g.board[g.player.y][g.player.x] = 'o'

	return true
}

func (g *Game) checkWin() bool {
	return g.player.x == g.emptySpace.x && g.player.y == g.emptySpace.y
}

func drawString(x, y int, str string, fg, bg termbox.Attribute) {
	for i, char := range str {
		termbox.SetCell(x+i, y, char, fg, bg)
	}
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	game := newGame(0)
	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	game.draw()

	for {
		select {
		case ev := <-eventQueue:
			if ev.Type == termbox.EventKey {
				moved := false

				switch {
				case ev.Key == termbox.KeyEsc:
					return
				case ev.Key == termbox.KeyArrowUp || ev.Ch == 'w' || ev.Ch == 'W':
					moved = game.movePlayer(0, -1)
				case ev.Key == termbox.KeyArrowDown || ev.Ch == 's' || ev.Ch == 'S':
					moved = game.movePlayer(0, 1)
				case ev.Key == termbox.KeyArrowLeft || ev.Ch == 'a' || ev.Ch == 'A':
					moved = game.movePlayer(-1, 0)
				case ev.Key == termbox.KeyArrowRight || ev.Ch == 'd' || ev.Ch == 'D':
					moved = game.movePlayer(1, 0)
				}

				game.draw()

				if moved && game.checkWin() {
					drawString(0, 13, "Congratulations! You won!", termbox.ColorGreen, termbox.ColorDefault)
					_ = termbox.Flush()
					time.Sleep(1 * time.Second)
					game = newGame(game.score + 1)
					game.draw()
				}
			}
		}
	}
}
