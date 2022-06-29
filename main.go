package main

import (
	b64 "encoding/base64"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Songmu/prompter"
)

const defaultSaveFileName = "save.dat"
const version = "0.4"

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
Tasukeru - HoloCure save file importer v%s
Made by DaniruKun

- https://github.com/DaniruKun
- https://twitter.com/DaniruKun

Tasukeru  Copyright (C) 2022  Daniils Petrovs
This program comes with ABSOLUTELY NO WARRANTY; for details see the Github link.
This is free software, and you are welcome to redistribute it
under certain conditions.

Website: https://danpetrov.xyz/tasukeru

I am not affiliated with Cover Corp. or Kay Yu in any way.
`
	fmt.Printf(header, version)
}

func holoCureSaveFilePath() string {
	dir, err := os.UserCacheDir()
	check(err)
	return filepath.Join(dir, "HoloCure", defaultSaveFileName)
}

// generally it seems that the start offset will always be the same
// across machines, but safer to find save block dynamically
func getSaveBlockStartEnd(srcDec *[]byte) (start, end int) {
	for offset, char := range *srcDec {
		if char == 0x7B && (*srcDec)[offset+1] == 0x20 {
			start = offset
		}
		if char == 0x7D && (*srcDec)[offset+1] == 0x00 {
			end = offset
		}
	}
	return
}

func main() {
	printHeader()
	args := os.Args

	if len(args) < 2 {
		// normally when used as drag n drop on windows, will be exactly 2
		fmt.Println("not enough arguments provided")
		fmt.Println("did you Drag n Drop the source save file onto tasukeru.exe ?")
		fmt.Println("do not forget to Drag n Drop the new save.dat onto tasukeru.exe")

		waitQuit()
		os.Exit(1)
	}

	var start, end int

	sourceSaveFilePath := args[1]

	fmt.Println("reading origin save file", sourceSaveFilePath)

	srcDat, err := os.ReadFile(sourceSaveFilePath)
	check(err)
	srcDec, err := b64.URLEncoding.DecodeString(string(srcDat))
	check(err)

	start, end = getSaveBlockStartEnd(&srcDec)

	srcSaveBlock := srcDec[start : end+1]

	var targetFilePath string

	if len(args) == 3 {
		targetFilePath = args[2]
	} else {
		targetFilePath = holoCureSaveFilePath()
	}

	targetDat, err := os.ReadFile(targetFilePath)
	check(err)
	targetDec, err := b64.URLEncoding.DecodeString(string(targetDat))
	check(err)

	start, _ = getSaveBlockStartEnd(&targetDec)

	// iterate over the source save block and overwrite the dst save block with its data
	for i, char := range srcSaveBlock {
		targetOffset := start + i

		if targetOffset >= len(targetDec) {
			targetDec = append(targetDec, char)
		} else {
			targetDec[targetOffset] = char
		}
	}

	fmt.Println("patched save:", string(targetDec))
	fmt.Println()

	var confirmed bool = prompter.YN("import new save file? インポートOK？", true)

	if confirmed {
		targetEnc := b64.URLEncoding.EncodeToString(targetDec)
		err = os.WriteFile(holoCureSaveFilePath(), []byte(targetEnc), 0644)
		check(err)
		fmt.Println("save file imported succesfully!")
		waitQuit()
	}
}
