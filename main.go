package main

import (
	b64 "encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Songmu/prompter"
)

const defaultSaveFileName = "save.dat"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func printHeader() {
	header := `
	Tasukeru - HoloCure save file importer
	Made by DaniruKun

	- https://github.com/DaniruKun
	- https://twitter.com/DaniruKun
	- https://danpetrov.xyz

	Tasukeru  Copyright (C) 2022  Daniils Petrovs
	This program comes with ABSOLUTELY NO WARRANTY; for details see the Github link.
	This is free software, and you are welcome to redistribute it
	under certain conditions.

	Source code: 
	`
	fmt.Print(header)
}

func holoCureSaveFilePath() string {
	dir, err := os.UserCacheDir()
	check(err)
	return filepath.Join(dir, "HoloCure", defaultSaveFileName)
}

func getSettingsStartEnd(srcDec *[]byte) (start, end int) {
	for offset, char := range *srcDec {
		if char == 0x7B && (*srcDec)[offset+1] == 0x20 {
			fmt.Println("found save block start at offset", offset)
			start = offset
		}
		if char == 0x7D && (*srcDec)[offset+1] == 0x00 {
			fmt.Println("found save block end at offset", offset)
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
		log.Fatalf("not enough arguments provided")
	}

	var start, end int

	sourceSaveFilePath := args[1]

	fmt.Println("reading origin save file", sourceSaveFilePath)

	srcDat, err := os.ReadFile(sourceSaveFilePath)
	check(err)
	srcDec, err := b64.URLEncoding.DecodeString(string(srcDat))
	check(err)

	start, end = getSettingsStartEnd(&srcDec)

	srcSettingsJson := srcDec[start : end+1]
	fmt.Println("source settings:", string(srcSettingsJson))

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

	start, _ = getSettingsStartEnd(&targetDec)

	for i, char := range srcSettingsJson {
		targetDec[start+i] = char
	}

	fmt.Println("patched settings", string(targetDec))

	var confirmed bool = prompter.YN("import new save file? インポートOK？", true)

	if confirmed {
		targetEnc := b64.URLEncoding.EncodeToString(targetDec)
		err = os.WriteFile(holoCureSaveFilePath(), []byte(targetEnc), 0644)
		check(err)
		fmt.Println("save file imported succesfully!")
	}

}
