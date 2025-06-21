package checksum

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

func CalculateFile(path string) (string, error) {
	// 1. Open file
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 2. Create SHA256 hasher
	hasher := sha256.New()

	// 3. Copy file to hasher (magic happens here!)
	_, err = io.Copy(hasher, file)
	if err != nil {
		return "", err
	}

	// 4. Get final hash
	hashBytes := hasher.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)

	return hashString, nil
}

func CompareFiles(path1, path2 string) (bool, error) {
	// Fast check: file size
	info1, err := os.Stat(path1)
	if err != nil {
		return false, err
	}
	info2, err := os.Stat(path2)
	if err != nil {
		return false, err
	}

	// Different size = definitely different!
	if info1.Size() != info2.Size() {
		return false, nil
	}

	// Same size -> need checksum
	checksum1, err := CalculateFile(path1)
	if err != nil {
		return false, err
	}
	checksum2, err := CalculateFile(path2)
	if err != nil {
		return false, err
	}

	return checksum1 == checksum2, nil
}
