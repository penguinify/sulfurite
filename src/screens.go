package src

import (
	"goco/utils"
	"os"
	"strings"

	"github.com/gdamore/tcell"
)


type GUI struct {
    screen tcell.Screen
}

func NewGUI(screen tcell.Screen) *GUI {
    return &GUI{screen: screen}
}

func (gui *GUI) ResetScreen() {
    gui.screen.Clear()
    utils.FancyText(gui.screen, 10, 5, "goco -- macros/automations cli", tcell.StyleDefault.Bold(true).Foreground(tcell.Color81))
    utils.FancyText(gui.screen, 10, 6, "BACK←UP↑DOWN↓SELECT→", tcell.StyleDefault.Foreground(tcell.Color117))

    utils.FancyText(gui.screen, 10, 20, "by @penguinify", tcell.StyleDefault.Foreground(tcell.Color195))
    gui.screen.Show()
}


func (gui *GUI) StartLoop(config *utils.ConfigJSON) {
    
    for {
        switch gui.Home() {
        case 1:
            macroName := gui.NewMacro()
            if macroName == "" { continue }

            if _, err := os.Stat(config.MacrosPath); os.IsNotExist(err){
                os.Mkdir(config.MacrosPath, 0755)
            }

            file, _ := os.Create(config.MacrosPath + macroName + ".goco")

            file.Close()
        case 2:
            gui.EditMacro(config)

        case 4:
            return
        }
    }

}

func (gui *GUI) Home() int  {

    gui.ResetScreen()
 
    selection := utils.Selection{
        Title: "",
        Options: []string{"+  New Macro", "✎Macros", "≡Settings", "× Exit"},
        Selected: 4,
        Coord: utils.Coord{X: 8, Y: 10},
    }

    return selection.Show(gui.screen)
    
}

func (gui *GUI) NewMacro() string {

    gui.ResetScreen()

    textInput := utils.TextInput{
        Title: "Enter the name of the new macro",
        Value: "",
        Coord: utils.Coord{X: 10, Y: 11},
        Styling: utils.Styling{
            Bold: true,
            Foreground: tcell.Color117,
        },
    }

    return textInput.Show(gui.screen)
}

func (gui *GUI) macroSelection(path string) string {

    gui.ResetScreen()

    options := GetMacroList(path)

    if len(options) == 0 {
        utils.FancyText(gui.screen, 10, 11, "No macros found", tcell.StyleDefault.Foreground(tcell.Color117).Bold(true))
        utils.FancyText(gui.screen, 10, 12, "Press any key to continue", tcell.StyleDefault.Foreground(tcell.Color117))

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

    selectedMacro:= selection.Show(gui.screen)

    if selectedMacro == -1 {
        return ""
    }

    return options[selectedMacro - 1]
}

func (gui *GUI) EditMacro(config *utils.ConfigJSON) {

    gui.ResetScreen()
    
    selectedMacro := gui.macroSelection(config.MacrosPath)

    if selectedMacro == "" {
        return
    }
    
    options := utils.Selection{
        Title: selectedMacro,
        Options: []string{"ᐅ Run", "✎ Edit", "𝚃Rename", "⃠⃠⃠⃠⃠⃠⃠⃠", "⦸ Delete"},
        Selected: 1,
        Coord: utils.Coord{X: 8, Y: 11},
    }

    gui.ResetScreen()
    switch options.Show(gui.screen) {
        case -1, 4:
            return
        case 1:
            gui.RunMacro(config.MacrosPath + selectedMacro)
        case 3:
            gui.ResetScreen()
            newNameInput := utils.TextInput{
                Title: "Enter the new name for the macro",
                Value: selectedMacro,
                Coord: utils.Coord{X: 10, Y: 11},
                Styling: utils.Styling{
                    Bold: true,
                    Foreground: tcell.Color117,
                },
            }
            newName := newNameInput.Show(gui.screen)

            if strings.Trim(newName, " ") == "" {
                return
            } else {
                os.Rename(config.MacrosPath + selectedMacro, config.MacrosPath + newName)
            }
        case 5:
            os.Remove(config.MacrosPath + selectedMacro)
            return

    }

}
func (gui *GUI) RunMacro(macroPath string) {
    gui.ResetScreen()
    utils.FancyText(gui.screen, 10, 10, "Compiling macro...", tcell.StyleDefault.Foreground(tcell.Color117))

    file, _ := os.ReadFile(macroPath)
    parser := NewParser(string(file))

    ast := parser.Parse()

    gui.ResetScreen()
    utils.FancyText(gui.screen, 10, 10, "Running macro...", tcell.StyleDefault.Foreground(tcell.Color117))
    utils.FancyText(gui.screen, 10, 11, "Press Ctrl+C to exit", tcell.StyleDefault.Foreground(tcell.Color117))

    interupted := make(chan bool)
    go Interpret(ast, interupted)

    for {
        ev := gui.screen.PollEvent()
        switch ev := ev.(type) {
        case *tcell.EventKey:
            if ev.Key() == tcell.KeyCtrlC {
                utils.FancyText(gui.screen, 10, 12, "Waiting until next instruction to quit...", tcell.StyleDefault.Foreground(tcell.Color117))
                interupted <- true
                return
            }
        }
    }
}
