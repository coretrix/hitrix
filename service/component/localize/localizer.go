package localize

import (
	"encoding/json"
	"errors"
	errorlogger "github.com/coretrix/hitrix/service/component/error_logger"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	separator = "|-|"
)

type ILocalizer interface {
	T(bucket string, key string) (string, error)
	LoadBucketFromFile(bucket string, path string, append bool)
	LoadBucketFromMap(bucket string, pairs map[string]string, append bool)
	SaveBucketToFile(bucket string, path string)
	PushBucketToSource(bucket string) (err error)
	PullBucketFromSource(bucket string, append bool) (err error)
}

type SimpleLocalizer struct {
	lock        sync.RWMutex
	pairs       map[string]string
	source      Source
	errorLogger errorlogger.ErrorLogger
}

func (l *SimpleLocalizer) T(bucket string, key string) string {
	if bucket == "" {
		panic("localization bucket not provided")
	}
	l.lock.RLock()
	defer l.lock.RUnlock()
	if val, ok := l.pairs[l.genKey(bucket, key)]; ok {
		return val
	}

	l.errorLogger.LogError("missing translations for key " + l.genKey(bucket, key))
	return key
}

func (l *SimpleLocalizer) LoadBucketFromMap(bucket string, pairs map[string]string, append bool) {
	if bucket == "" {
		return
	}
	l.lock.Lock()
	defer l.lock.Unlock()
	if l.pairs == nil {
		l.pairs = map[string]string{}
	}
	if !append {
		l.removeBucket(bucket)
	}
	for k, v := range pairs {
		l.pairs[l.genKey(bucket, k)] = v
	}
}

func (l *SimpleLocalizer) LoadBucketFromFile(bucket string, path string, append bool) {
	jsonBytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic("no such file or directory: " + path)
	}
	var tempParis map[string]string
	err = json.Unmarshal(jsonBytes, &tempParis)
	if err != nil {
		log.Println("translation file not well formated json", err)
	}
	l.LoadBucketFromMap(bucket, tempParis, append)
}

// func (l *SimpleLocalizer) SaveToFile(path string) {
// 	//TODO: not implemented
// 	panic("not implemented")
// }

func (l *SimpleLocalizer) SaveBucketToFile(bucket string, path string) {
	err := l.touchFile(path)
	if err != nil {
		panic(err)
	}
	tempPairs := l.getBucketPairsWithoutPrefix(bucket)
	jsonBytes, _ := json.MarshalIndent(tempPairs, "", " ")
	err = ioutil.WriteFile(path, jsonBytes, 0644)
	if err != nil {
		panic(err)
	}
}

func (l *SimpleLocalizer) PushBucketToSource(bucket string) (err error) {
	if l.source == nil {
		return errors.New("there is no defined sourced")
	}

	terms := l.getBucketTermsWithoutPrefix(bucket)
	err = l.source.Push(terms)
	if err != nil {
		log.Fatal(err)
		return
	}

	return
}

func (l *SimpleLocalizer) PullBucketFromSource(bucket string, append bool) (err error) {
	if l.source == nil {
		return errors.New("there is no defined sourced")
	}

	terms, err := l.source.Pull()
	if err != nil {
		log.Fatal(err)
		return
	}
	l.LoadBucketFromMap(bucket, terms, append)

	return
}

func (l *SimpleLocalizer) removeBucket(bucket string) {
	for k := range l.pairs {
		if strings.HasPrefix(k, l.genKey(bucket, "")) {
			delete(l.pairs, k)
		}
	}
}

func (l *SimpleLocalizer) getBucketTermsWithoutPrefix(bucket string) (terms []string) {
	tempPairs := l.getBucketPairsWithoutPrefix(bucket)
	for k := range tempPairs {
		terms = append(terms, k)
	}

	return
}

func (l *SimpleLocalizer) getBucketPairsWithoutPrefix(bucket string) map[string]string {
	tempPairs := map[string]string{}
	for k, v := range l.pairs {
		if !strings.HasPrefix(k, bucket) {
			continue
		}
		tempPairs[l.removeKeyPrefix(bucket, k)] = v
	}

	return tempPairs
}

func (l *SimpleLocalizer) touchFile(path string) error {
	directoryPath := filepath.Dir(path)
	if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
		err := os.MkdirAll(directoryPath, 0755)
		if err != nil {
			return err
		}
	}
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	return file.Close()
}

func (l *SimpleLocalizer) genKey(bucket string, key string) string {
	return bucket + separator + key
}

func (l *SimpleLocalizer) removeKeyPrefix(bucket string, key string) string {
	return strings.Replace(key, bucket+separator, "", 1)
}

func NewSimpleLocalizer(errorLogger errorlogger.ErrorLogger, source Source, localePath string) *SimpleLocalizer {
	localizerService := &SimpleLocalizer{
		source:      source,
		errorLogger: errorLogger,
	}

	files, err := ioutil.ReadDir(localePath)
	if err != nil {
		panic(err)
	}

	shouldAppend := false
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		localizerService.LoadBucketFromFile(strings.TrimSuffix(file.Name(), ".json"), localePath+"/"+file.Name(), shouldAppend)
		shouldAppend = true
	}

	return localizerService
}
