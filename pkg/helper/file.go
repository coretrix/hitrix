package helper

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"strings"
)

const (
	png  = "png"
	jpeg = "jpeg"
)

var supportedFileExtensions = []string{
	png, jpeg,
}

func GetFileExtensionFromBase64Header(base64Header string) (string, error) {
	for _, ext := range supportedFileExtensions {
		if strings.Contains(base64Header, ext) {
			return "." + ext, nil
		}
	}

	return "", errors.New("unsupported image type")
}

func ExistsInDir(filename, dir string) bool {
	fileInfo, err := os.ReadDir(dir)
	if err != nil {
		return false
	}

	for _, file := range fileInfo {
		if filename == file.Name() {
			return true
		}
	}

	return false
}

func GetTemplateDirectory() string {
	var dir string

	templateFolder, has := os.LookupEnv("APP_TEMPLATE_FOLDER")
	if has {
		dir = templateFolder
	} else {
		_, filename, _, _ := runtime.Caller(0)
		dir = path.Join(path.Dir(filename), "../../../templates")
	}

	return dir
}

func CSVToMap(reader *bytes.Buffer) ([]map[string]string, error) {
	r := csv.NewReader(reader)
	r.Comma = ';'

	var rows []map[string]string
	var header []string

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		if header == nil {
			header = record
		} else {
			dict := map[string]string{}

			for i := range header {
				dict[header[i]] = record[i]
			}

			rows = append(rows, dict)
		}
	}

	return rows, nil
}

func FileExistsInDir(filename, dir string) bool {
	fileInfo, err := os.ReadDir(dir)
	if err != nil {
		return false
	}

	for _, file := range fileInfo {
		if filename == file.Name() {
			return true
		}
	}

	return false
}

func Copy(src, dst string, bufferSize int64) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return errors.New(src + " is not a regular file.")
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	_, err = os.Stat(dst)
	if err == nil {
		return fmt.Errorf("file %s already exists", dst)
	}

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	buf := make([]byte, bufferSize)

	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}

		if n == 0 {
			break
		}

		if _, err := destination.Write(buf[:n]); err != nil {
			return err
		}
	}

	return err
}
