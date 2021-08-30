package fileextractor

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type FileExtractor struct {
	founds map[string]string
}

type ExtractParams struct {
	SearchPath string
	Excludes   []string
	Expression string
}

func (l *FileExtractor) Extract(params ExtractParams) ([]string, error) {
	l.founds = map[string]string{}
	err := filepath.Walk(params.SearchPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if strings.HasSuffix(path, ".go") &&
				!strings.HasSuffix(path, "_test.go") &&
				!strings.HasSuffix(path, "_gen.go") &&
				!strings.HasSuffix(path, "generated.go") {
				l.extractFormFile(path, params.Expression)
			}
			return err
		})
	if err != nil {
		return nil, err
	}

	var foundsSlice []string
	for key := range l.founds {
		foundsSlice = append(foundsSlice, key)
	}

	return foundsSlice, err
}

func (l *FileExtractor) extractFormFile(pathToread string, expression string) {
	fileContent, err := l.readFile(pathToread)
	if err != nil {
		log.Fatal(err)
	}
	reg := *regexp.MustCompile(expression)
	res := reg.FindAllStringSubmatch(fileContent, -1)
	for i := range res {
		l.founds[res[i][1]] = ""
	}
}

func (l *FileExtractor) readFile(pathToread string) (string, error) {
	file, err := os.Open(pathToread)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	defer file.Close()

	srcbuf, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return string(srcbuf), nil
}

func NewFileExtractor() *FileExtractor {
	return &FileExtractor{}
}
