package main

import (
	"fmt"
	"os"

	"github.com/DaniruKun/tasukeru/holocure"
	"github.com/Songmu/prompter"
)

const Version = "1.1.0"
const HomePage = "https://danpetrov.xyz/tasukeru"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func waitQuit() {
	prompter.YN("press Enter to quit", true)
}

func printHeader() {
	header := `
💀 🐔 🐙 🔱 🔎 💎 🪐 🪶 ⏳ 🌿 🎲

Tasukeru - HoloCure save file importer v%s
Made by DaniruKun

- https://github.com/DaniruKun
- https://twitter.com/DaniruKun

Tasukeru  Copyright (C) 2022  Daniils Petrovs
This program comes with ABSOLUTELY NO WARRANTY; for details see the Github link.
This is free software, and you are welcome to redistribute it
under certain conditions.

Website: %s

I am not affiliated with Cover Corp. or Kay Yu in any way.
`
	fmt.Printf(header, Version, HomePage)
}

func main() {
	args := os.Args

	if len(args) < 2 {
		// Run the UI
		RunGUI()
	} else {
		CLI(args)
	}
}

func CLI(args []string) {
	printHeader()
	var sourceSaveFilePath, targetSaveFilePath string

	sourceSaveFilePath = args[1]

	if len(args) == 3 {
		targetSaveFilePath = args[2]
	} else {
		targetSaveFilePath = holocure.SaveFilePath()
	}

	targetDec := holocure.MergeSaves(sourceSaveFilePath, targetSaveFilePath)

	fmt.Println()

	var confirmed bool = prompter.YN("import new save file? インポートOK？", true)

	if confirmed {
		err := holocure.WriteSaveFile(targetSaveFilePath, targetDec)
		check(err)
		fmt.Println("save file imported succesfully!")
		waitQuit()
	}
}
