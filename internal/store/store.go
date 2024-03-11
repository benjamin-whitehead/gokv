package store

import (
	"io"
	"os"
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

func (s *Store) Write(key string, value string) error {
	_, err := s.file.Write([]byte(value))
	if err != nil {
		return err
	}

	offset, err := s.file.Seek(0, io.SeekCurrent)
	if err != nil {
		return err
	}

	s.log[key] = logEntry{
		offset: int(offset),
		length: len(value),
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

	_, err = s.file.ReadAt(readBuffer, int64(offset-length))
	if err != nil {
		return "", err
	}

	return string(readBuffer), nil
}

func (s *Store) Delete(key string) error {
	if _, ok := s.log[key]; !ok {
		return ErrKeyNotFound(key)
	}

	delete(s.log, key)

	return nil
}

type logEntry struct {
	offset int
	length int
}
