package main

import (
    "sulfurite/utils"
    "sulfurite/src"
)

func main() {
    config := utils.LoadConfig("config.json")

    screen := utils.InitScreen()
    if screen == nil { return }
    defer screen.Fini()

    gui := src.NewGUI(screen)
    gui.StartLoop(config)
}
