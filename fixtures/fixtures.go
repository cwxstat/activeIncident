package fixtures

import (
	"bytes"
	"encoding/gob"
	"os"
	"path/filepath"
	"runtime"
)

// basepath is the root directory of this package.
var basepath string

func init() {
	_, currentFile, _, _ := runtime.Caller(0)
	basepath = filepath.Dir(currentFile)
}

// Path returns the absolute path the given relative file or directory path,
// relative to the "fixtures" directory.
// If rel is already absolute, it is returned unmodified.
func Path(rel string) string {
	if filepath.IsAbs(rel) {
		return rel
	}

	return filepath.Join(basepath, rel)
}

func Read(filepath string) ([]byte, error) {
	return os.ReadFile(filepath)
}
func Write(filepath string, data []byte) error {
	return os.WriteFile(filepath, data, 0644)
}

func Encode(filePath string, a any) error {
	var network bytes.Buffer

	enc := gob.NewEncoder(&network)
	err := enc.Encode(a)
	if err != nil {
		return err
	}
	return Write(filePath, network.Bytes())

}
