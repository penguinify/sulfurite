package utils

import (
	"fmt"
	"github.com/gdamore/tcell"
)

type Coord struct {
    X int
    Y int
}

type Styling struct {
    Bold bool
    Foreground tcell.Color
    Background tcell.Color
}

type Selection struct {
    Title string
    Options []string
    Selected int
    Coord
}

type TextInput struct {
    Title string
    Value string
    Coord
    Styling
}

func InitScreen() tcell.Screen {
    screen, err := tcell.NewScreen()
    if err != nil {
        fmt.Println("Error initializing screen:", err)
        return nil
    }

    if err = screen.Init(); err != nil {
        fmt.Println("Error initializing screen:", err)
        return nil
    }

    return screen
}

func WaitUntilKey(screen tcell.Screen) {
    for {
        ev := screen.PollEvent()
        switch ev.(type) {
        case *tcell.EventKey:
            return
        }
    }
}

func FancyText(screen tcell.Screen, x int, y int, text string, style tcell.Style) {
    for i, char := range text {
        screen.SetContent(x + i, y, char, nil, style)
    }

    screen.Show()
}


func DrawText(screen tcell.Screen, x int, y int, text string) {
    for i, char := range text {
        screen.SetContent(x + i, y, char, nil, tcell.StyleDefault)
    }

    screen.Show()
}

// PARA SELECCIONAR
func (s *Selection) draw(screen tcell.Screen) {
    DrawText(screen, s.Coord.X + 2, s.Coord.Y, s.Title)
    for i, option := range s.Options {
        i++
        if i == s.Selected {
            FancyText(screen, s.Coord.X, s.Coord.Y + i, "> " + option, tcell.StyleDefault.Bold(true))
        } else {
            DrawText(screen, s.Coord.X, s.Coord.Y + i, "  " + option)
        }
    }
}


func (s *Selection) Show(screen tcell.Screen) int {

    for {
        s.draw(screen)
        screen.Show()

        ev := screen.PollEvent()

        switch ev := ev.(type) {
        case *tcell.EventKey:
            switch ev.Key() {
            case tcell.KeyUp:
                if s.Selected > 1 {
                    s.Selected--
                }
            case tcell.KeyDown:
                if s.Selected < len(s.Options) {
                    s.Selected++
                }
            case tcell.KeyRight, tcell.KeyEnter:
                return s.Selected

            case tcell.KeyLeft:
                return -1
            }
        }
    }
    
}

// PARA INGRESAR TEXTO
func (i *TextInput) draw(screen tcell.Screen) {
    DrawText(screen, i.Coord.X, i.Coord.Y + 1, "                                                               ")
    FancyText(screen, i.Coord.X, i.Coord.Y, i.Title, tcell.StyleDefault.Bold(i.Bold).Foreground(i.Foreground).Background(i.Background))
    FancyText(screen, i.Coord.X, i.Coord.Y + 1, i.Value, tcell.StyleDefault.Foreground(i.Foreground).Background(i.Background))
}

func (i *TextInput) Show(screen tcell.Screen) string {

    defer screen.HideCursor()

    // set default style in case it's not set
    if i.Foreground == 0 {
        i.Foreground = tcell.ColorDefault
    }

    if i.Background == 0 {
        i.Background = tcell.ColorDefault
    }
        
    for {
        screen.ShowCursor(i.Coord.X + len(i.Value), i.Coord.Y + 1)
        i.draw(screen)
        screen.Show()

        ev := screen.PollEvent()

        switch ev := ev.(type) {
        case *tcell.EventKey:
            switch ev.Key() {
            case tcell.KeyRune:
                i.Value += string(ev.Rune())
            case tcell.KeyBackspace, tcell.KeyBackspace2:
                if len(i.Value) > 0 {
                    i.Value = i.Value[:len(i.Value)-1]
                }
            case tcell.KeyRight, tcell.KeyEnter:
                return i.Value
            case tcell.KeyLeft:
                return ""
            }
        }
    }
}
