package store

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type Store struct {
	file *os.File
	log  map[string]logEntry
}

func NewStore(name string) (*Store, error) {
	file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	log := make(map[string]logEntry)

	return &Store{
		file: file,
		log:  log,
	}, nil
}

func NewStoreFromFile(path string) (*Store, error) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return nil, ErrFileNotFound(path)
	}

	_, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	// TODO: Finish implementing
	// We are going to need to know where each entry we are parsing from is in the file
	// So we can write the offset into the log

	return nil, nil
}

func (s *Store) Write(key string, value string) error {
	entry := []byte(fmt.Sprintf("%s:::%s\n", key, value))
	_, err := s.file.Write(entry)
	if err != nil {
		return err
	}

	offset, err := s.file.Seek(0, io.SeekCurrent)
	if err != nil {
		return err
	}

	s.log[key] = logEntry{
		offset: int(offset),
		length: len(entry),
	}

	return nil
}

func (s *Store) Read(key string) (string, error) {
	if _, ok := s.log[key]; !ok {
		return "", ErrKeyNotFound(key)
	}

	length := s.log[key].length
	offset := s.log[key].offset

	readBuffer := make([]byte, length)

	_, err := s.file.Seek(0, io.SeekStart)
	if err != nil {
		return "", err
	}

	_, err = s.file.ReadAt(readBuffer, int64(offset-(length)))
	if err != nil {
		return "", err
	}

	pair, err := parsePairFromBuffer(readBuffer)
	if err != nil {
		return "", err
	}

	readKey, readValue := pair.key, pair.value
	if readKey != key {
		return "", ErrReadIncorrectKey(key, readKey)
	}

	return readValue, nil
}

func (s *Store) Delete(key string) error {
	if _, ok := s.log[key]; !ok {
		return ErrKeyNotFound(key)
	}

	delete(s.log, key)

	return nil
}

func parsePairFromBuffer(buffer []byte) (pair, error) {
	entry := strings.Split(string(buffer), ":::")
	if len(entry) != 2 {
		return pair{}, ErrDecodeEntry
	}

	key, value := entry[0], entry[1]
	trimmedValue := strings.TrimSuffix(value, "\n")

	return pair{key: key, value: trimmedValue}, nil
}

type pair struct {
	key   string
	value string
}

type logEntry struct {
	offset int
	length int
}
