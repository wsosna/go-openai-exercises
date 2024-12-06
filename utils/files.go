package utils

import (
	"bytes"
	"encoding/base64"
	"net/http"
	"os"
	"strings"
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

func ReadFileToBase64(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	var base64Encoding string
	mimeType := http.DetectContentType(data)

	switch mimeType {
	case "image/jpeg":
		base64Encoding += "data:image/jpeg;base64,"
	case "image/png":
		base64Encoding += "data:image/png;base64,"
	}

	base64Encoding += base64.StdEncoding.EncodeToString(data)
	return base64Encoding, nil
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

func ReadFilesToPrompt(path string, promptTitle ...string) string {
	dir, err := os.ReadDir(path)
	HandleFatalError(err)

	title := "file"
	if len(promptTitle) > 0 {
		title = promptTitle[0]
	}

	result := ""
	for _, entry := range dir {
		if entry.IsDir() {
			continue
		}
		content, err := ReadFileToString(path + "/" + entry.Name())
		HandleFatalError(err)

		if strings.Contains(content, "entry deleted") {
			continue
		}
		content = strings.ReplaceAll(content, "\n", " ")
		if len(content) > 0 {
			result += "<" + title + ">\n"
			result += "<name>\n"
			result += entry.Name()
			result += "\n</name>\n"
			result += "<content>\n"
			result += content
			result += "\n</content>\n"
			result += "<" + title + ">\n"
		}
	}
	return result
}
