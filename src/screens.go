package src

import (
    "goco/utils"
    "github.com/gdamore/tcell"
    "os"
)


type GUI struct {
    screen tcell.Screen
}

func NewGUI(screen tcell.Screen) *GUI {
    return &GUI{screen: screen}
}

func (gui *GUI) StartLoop(config *utils.Config) {
    
    for {
        switch gui.Home() {
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

    utils.FancyText(gui.screen, 10, 5, "goco ~ macros", tcell.StyleDefault.Bold(true).Foreground(tcell.Color75))
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

    utils.FancyText(gui.screen, 10, 5, "goco ~ macros", tcell.StyleDefault.Bold(true).Foreground(tcell.Color75))

    return textInput.Show(gui.screen)
}

func (gui *GUI) EditMacro() {

    gui.screen.Clear()

    macroslist, _ := os.ReadDir("macros")

    options := []string{}

    for _, file := range macroslist {
        options = append(options, file.Name())
    }

    if len(options) == 0 {
        utils.FancyText(gui.screen, 10, 5, "goco ~ macros", tcell.StyleDefault.Bold(true).Foreground(tcell.Color75))
        utils.FancyText(gui.screen, 10, 7, "No macros found", tcell.StyleDefault.Foreground(tcell.Color39))
        utils.FancyText(gui.screen, 10, 10, "Press any key to continue", tcell.StyleDefault.Foreground(tcell.Color39))

        gui.screen.Show()

        utils.WaitUntilKey(gui.screen)
        return
    }


    selection := utils.Selection{
        Title: "",
        Options: options,
        Selected: 1,
        Coord: utils.Coord{X: 8, Y: 10},
    }

    utils.FancyText(gui.screen, 10, 5, "goco ~ macros", tcell.StyleDefault.Bold(true).Foreground(tcell.Color75))
    utils.FancyText(gui.screen, 10, 7, "[up & down arrows to move the cursor]", tcell.StyleDefault.Foreground(tcell.Color39))
    utils.FancyText(gui.screen, 10, 8, "[enter to select]", tcell.StyleDefault.Foreground(tcell.Color39))

    selection.Show(gui.screen)
    
}
