package holocure

import (
	b64 "encoding/base64"
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
		if char == 0x7D && (*data)[offset+1] == 0x00 {
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

func check(e error) {
	if e != nil {
		panic(e)
	}
}
