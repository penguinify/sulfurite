package src

import (
	"sulfurite/utils"
	"os"
	"strings"
    "sulfurite/src/editor"
    "context"
	"github.com/gdamore/tcell"
    "github.com/ncruces/zenity"
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


func (gui *GUI) StartLoop(config *utils.ConfigJSON) {
    
    for {
        switch gui.Home() {
        case 1:
            macroName := gui.NewMacro()
            if macroName == "" { continue }

            if _, err := os.Stat(config.MacrosPath); os.IsNotExist(err){
                os.Mkdir(config.MacrosPath, 0755)
            }

            file, _ := os.Create(config.MacrosPath + macroName + ".sulfur")

            file.Close()
        case 2:
            gui.MacroList(config)

        case 3:
            gui.ImportMacro(config)
        case 4:
            return
        }
    }

}

/*
These are functions that define the screens users will see
*/
func (gui *GUI) Home() int  {

    gui.resetScreen()
 
    selection := utils.Selection{
        Title: "",
        Options: []string{"+  New Macro", "✎Macros", "≡Import", "× Exit"},
        Selected: 4,
        Coord: utils.Coord{X: 8, Y: 10},
    }

    return selection.Show(gui.screen)
    
}

func (gui *GUI) NewMacro() string {

    gui.resetScreen()

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

func (gui *GUI) MacroList(config *utils.ConfigJSON) {

    gui.resetScreen()
    
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

    gui.resetScreen()

    /* 
    Option 1: runs the macro using the RunMacro function
    Option 2: opens the editor in the web for editing
    Option 3: renames the macro
    Option 4: blank so people don't accidentally press it
    Option 5: deletes the macro
    */
    switch options.Show(gui.screen) {
        case -1, 4:
            defer gui.MacroList(config)
            return

        case 1:
            gui.runMacro(config.MacrosPath + selectedMacro)

        case 2:
            gui.editMacro(config.MacrosPath + selectedMacro)

        case 3:
            gui.resetScreen()
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
            } 

            os.Rename(config.MacrosPath + selectedMacro, config.MacrosPath + newName)

        case 5:
            os.Remove(config.MacrosPath + selectedMacro)
            return

    }

}

func (gui *GUI) ImportMacro(config *utils.ConfigJSON) {

    path, err := zenity.SelectFile(zenity.Title("Select a macro to import"), zenity.FileFilters{{Name: "Sulfurite Macro", Patterns: []string{"*.sulfur"}}})

    if err != nil || path == "" {
        return
    }

    newPath := config.MacrosPath + strings.Split(path, "/")[len(strings.Split(path, "/")) - 1]

    os.Rename(path, newPath)


    gui.resetScreen()

    utils.FancyText(gui.screen, 10, 10, "Macro imported successfully", colorBright)
    utils.FancyText(gui.screen, 10, 11, "Press any key to continue", colorBright)

    gui.screen.Show()

    utils.WaitUntilKey(gui.screen)
}


// macroSelection returns the name of the selected macro
func (gui *GUI) macroSelection(path string) string {

    gui.resetScreen()

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


func (gui *GUI) runMacro(macroPath string) {
    gui.resetScreen()
    utils.FancyText(gui.screen, 10, 10, "Compiling macro...", colorBright)

    file, _ := os.ReadFile(macroPath)
    parser := NewParser(string(file))

    ast := parser.Parse()

    gui.resetScreen()
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


func (gui *GUI) editMacro(selectedMacro string) {
    s := &editor.Server{
        Addr: "8080",
        File: selectedMacro,
    }
    Server := s.Start()
    defer Server.Server.Shutdown(context.Background())

    gui.resetScreen()

    utils.FancyText(gui.screen, 10, 10, "Opening editor...", colorBright)
    utils.FancyText(gui.screen, 10, 11, "Press any key to quit the editor", colorBright)

    utils.WaitUntilKey(gui.screen)
}


func (gui *GUI) resetScreen() {
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
