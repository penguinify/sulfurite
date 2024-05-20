package src

import (
	"goco/utils"
	"os"
	"github.com/gdamore/tcell"
)

func showTitle(screen tcell.Screen) {
    utils.FancyText(screen, 10, 5, "goco ~ macros", tcell.StyleDefault.Bold(true).Foreground(tcell.Color75))
}

type GUI struct {
    screen tcell.Screen
}

func NewGUI(screen tcell.Screen) *GUI {
    return &GUI{screen: screen}
}

func (gui *GUI) StartLoop(config *utils.Config) {
    
    for {
        switch gui.Home() {
        case 1:
            gui.RunMacro()
        case 2:
            macroName := gui.NewMacro()
            if macroName == "" { continue }

            file, _ := os.Create(config.MacrosPath + macroName + ".goco")
            file.WriteString(config.MacroInterpreterVersion + "1.0|Full| goco macro file starts below this line!")

            file.Close()
        case 3:
            gui.EditMacro()

        case 5:
            return
        }
    }

}

func (gui *GUI) Home() int  {

    gui.screen.Clear()
 
    selection := utils.Selection{
        Title: "",
        Options: []string{"ᐅRun Macro", "+  New Macro", "✎Edit Macro", "≡Settings", "× Exit"},
        Selected: 5,
        Coord: utils.Coord{X: 8, Y: 10},
    }

    showTitle(gui.screen)
    utils.FancyText(gui.screen, 10, 7, "[up & down arrows to move the cursor]", tcell.StyleDefault.Foreground(tcell.Color39))
    utils.FancyText(gui.screen, 10, 8, "[enter to select]", tcell.StyleDefault.Foreground(tcell.Color39))

    return selection.Show(gui.screen)
    
}

func (gui *GUI) NewMacro() string {

    gui.screen.Clear()

    textInput := utils.TextInput{
        Title: "Enter the name of the new macro",
        Value: "",
        Coord: utils.Coord{X: 10, Y: 10},
        Styling: utils.Styling{
            Bold: true,
            Foreground: tcell.Color75,
        },
    }

    showTitle(gui.screen)

    return textInput.Show(gui.screen)
}

func (gui *GUI) macroSelection() string {
    gui.screen.Clear()

    options := GetMacroList()

    if len(options) == 0 {
        showTitle(gui.screen)
        utils.FancyText(gui.screen, 10, 7, "No macros found", tcell.StyleDefault.Foreground(tcell.Color39))
        utils.FancyText(gui.screen, 10, 10, "Press any key to continue", tcell.StyleDefault.Foreground(tcell.Color39))

        gui.screen.Show()

        utils.WaitUntilKey(gui.screen)
        return ""
    }

    selection := utils.Selection{
        Title: "",
        Options: options,
        Selected: 1,
        Coord: utils.Coord{X: 8, Y: 10},
    }

    showTitle(gui.screen)
    utils.FancyText(gui.screen, 10, 7, "[up & down arrows to move the cursor]", tcell.StyleDefault.Foreground(tcell.Color39))
    utils.FancyText(gui.screen, 10, 8, "[enter to select]", tcell.StyleDefault.Foreground(tcell.Color39))

    return options[selection.Show(gui.screen) - 1]
}

func (gui *GUI) EditMacro() {

    gui.screen.Clear()
    
    gui.macroSelection()
}

func (gui *GUI) RunMacro() {
    interpreter := Interpreter{
        Macro: "macros/" + gui.macroSelection(),
    }


    gui.screen.Clear()
    showTitle(gui.screen)
    utils.FancyText(gui.screen, 10, 7, "Compiling macro...", tcell.StyleDefault.Foreground(tcell.Color39))
    gui.screen.Show()
    interpreter.CompileToActions()
    gui.screen.Clear()
    showTitle(gui.screen)
    utils.FancyText(gui.screen, 10, 7, "Running macro...", tcell.StyleDefault.Foreground(tcell.Color39))
    gui.screen.Show()

    ch := make(chan bool)
    go func() {
        for {
            select {
            case <-ch:
                return
            default:
                interpreter.Run()
            }
        }
    }()
    
    for {
        ev := gui.screen.PollEvent()

        switch ev := ev.(type) {
        case *tcell.EventKey:
            if ev.Key() == tcell.KeyCtrlC {
                ch <- true
                return
            }
        }

    }

}
