package src

import (
	"sulfurite/utils"
	"os"
	"strings"
    "sulfurite/src/editor"

	"github.com/gdamore/tcell"
)


type GUI struct {
    screen tcell.Screen
}

var (
    colorBright = tcell.StyleDefault.Foreground(tcell.Color117)
    colorDark = tcell.StyleDefault.Foreground(tcell.Color45)
)

func NewGUI(screen tcell.Screen) *GUI {
    return &GUI{screen: screen}
}

func (gui *GUI) ResetScreen() {
    gui.screen.Clear()
    utils.ASCIIArt(gui.screen, 9, 4, []string{
        "            _  __            _ _       ",
        "  ___ _   _| |/ _|_   _ _ __(_) |_ ___ ",
        " / __| | | | | |_| | | | '__| | __/ _ \\",
        " \\__ \\ |_| | |  _| |_| | |  | | ||  __/",
        " |___/\\__,_|_|_|  \\__,_|_|  |_|\\__\\___|",
    }, colorBright)

    utils.FancyText(gui.screen, 10, 19, "left arrow to go back", colorDark)
    utils.FancyText(gui.screen, 10, 20, "right arrow to select", colorDark)
    utils.FancyText(gui.screen, 10, 21, "up and down arrows to move", colorDark)
    utils.FancyText(gui.screen, 10, 23, "by @penguinify", colorDark)
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

            file, _ := os.Create(config.MacrosPath + macroName + ".sulfurite")

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
        Options: []string{"+  New Macro", "✎Macros", "≡Import", "× Exit"},
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
        utils.FancyText(gui.screen, 10, 11, "No macros found", colorBright.Bold(true))
        utils.FancyText(gui.screen, 10, 12, "Press any key to continue", colorBright)

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
            defer gui.EditMacro(config)
            return
        case 1:
            gui.RunMacro(config.MacrosPath + selectedMacro)
        case 2:
            s := &editor.Server{
                Addr: "8080",
                File: selectedMacro,
            }
            Server := s.Start()
            defer Server.Server.Shutdown(nil)

            gui.ResetScreen()

            utils.FancyText(gui.screen, 10, 10, "Opening editor...", colorBright)
            utils.FancyText(gui.screen, 10, 11, "Press any key to quit the editor", colorBright)

            utils.WaitUntilKey(gui.screen)
            return            

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
    utils.FancyText(gui.screen, 10, 10, "Compiling macro...", colorBright)

    file, _ := os.ReadFile(macroPath)
    parser := NewParser(string(file))

    ast := parser.Parse()

    gui.ResetScreen()
    utils.FancyText(gui.screen, 10, 10, "Running macro...", colorBright)
    utils.FancyText(gui.screen, 10, 11, "Press Ctrl+C to exit", colorBright)

    interupted := make(chan bool)
    go func() {
        _, err := Interpret(ast, interupted)

        if err != nil {
            utils.FancyText(gui.screen, 10, 12, err.Error(), colorBright)
        }

        <-interupted

    }()

    for {
        ev := gui.screen.PollEvent()
        switch ev := ev.(type) {
        case *tcell.EventKey:
            if ev.Key() == tcell.KeyCtrlC {
                utils.FancyText(gui.screen, 10, 12, "Waiting until next instruction to quit...", colorBright)
                interupted <- true
                return
            }
        }
    }
}
