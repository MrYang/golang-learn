package file

import (
	"io/ioutil"
	"os"
	"strings"
)

func IsExist(f string) bool {
	_, err := os.Stat(f)
	return err == nil || os.IsExist(err)
}

func ToTrimString(filePath string) (string, error) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(b)), nil
}
