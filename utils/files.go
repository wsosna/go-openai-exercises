package utils

import (
	"bytes"
	"os"
)

func ReadFileToString(path string) (string, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(file), nil
}

func ReadFileToBuffer(path string) (*bytes.Buffer, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(file), nil
}

func WriteStringToFile(text string, path string) error {
	out, err := os.Create(path)
	if err != nil {
		return err

	}
	_, err = out.WriteString(text)
	if err != nil {
		return err

	}
	err = out.Close()
	if err != nil {
		return err

	}
	return nil
}
