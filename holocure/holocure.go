package holocure

import (
	b64 "encoding/base64"
	"fmt"
	"os"
	"path/filepath"
)

const defaultSaveFileName = "save.dat"

// Returns the default HoloCure save file location for the current platform.
func SaveFilePath() string {
	dir, err := os.UserCacheDir()
	check(err)
	return filepath.Join(dir, "HoloCure", defaultSaveFileName)
}

// It seems that the start offset will always be the same
// across machines, but safer to find save block dynamically
func FindSaveBlockStartEnd(data *[]byte) (start, end int) {
	for offset, char := range *data {
		if char == 0x7B && (*data)[offset+1] == 0x20 {
			start = offset
		}
		if char == 0x7D && (*data)[offset-1] == 0x20 {
			end = offset
		}
	}
	return
}

func Decode(data []byte) ([]byte, error) {
	return b64.URLEncoding.DecodeString(string(data))
}

func DecodeSaveFile(filePath string) ([]byte, error) {
	srcDat, err := os.ReadFile(filePath)
	check(err)
	srcDec, err := Decode(srcDat)
	return srcDec, err
}

func WriteSaveFile(filePath string, data []byte) error {
	targetEnc := b64.URLEncoding.EncodeToString(data)
	return os.WriteFile(filePath, []byte(targetEnc), 0644)
}

// Merges the source save file path into the target one.
// Returns the raw bytes of the patches save.
func MergeSaves(sourceSaveFilePath, targetSaveFilePath string) []byte {
	fmt.Println("reading origin save file", sourceSaveFilePath)
	srcDec, err := DecodeSaveFile(sourceSaveFilePath)
	check(err)
	var start, end int

	start, end = FindSaveBlockStartEnd(&srcDec)
	srcSaveBlock := srcDec[start : end+1]

	targetDec, err := DecodeSaveFile(targetSaveFilePath)
	check(err)

	start, _ = FindSaveBlockStartEnd(&targetDec)

	// iterate over the source save block and overwrite the dst save block with its data
	for i, char := range srcSaveBlock {
		targetOffset := start + i

		if targetOffset >= len(targetDec) {
			targetDec = append(targetDec, char)
		} else {
			targetDec[targetOffset] = char
		}
	}
	return targetDec
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
