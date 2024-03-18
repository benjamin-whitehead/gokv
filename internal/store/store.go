package store

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

type Store struct {
	sync.RWMutex
	file     *os.File
	snapshot *os.File
	log      map[string]logEntry
}

func NewStore(name string) (*Store, error) {
	file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	snapshot, err := os.OpenFile(snapshotFileFromPath(name), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	log := make(map[string]logEntry)

	return &Store{
		file:     file,
		snapshot: snapshot,
		log:      log,
	}, nil
}

func NewStoreFromFile(path string) (*Store, error) {
	if !fileExists(path) {
		return nil, ErrFileNotFound(path)
	}

	snapshotFilePath := snapshotFileFromPath(path)
	if !fileExists(snapshotFilePath) {
		return nil, ErrSnapshotFileNotFound(snapshotFilePath)
	}

	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	snapshot, err := os.OpenFile(snapshotFilePath, os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	snapshotStat, err := snapshot.Stat()
	if err != nil {
		return nil, err
	}

	log := make(map[string]logEntry)
	buffer := make([]byte, snapshotStat.Size())

	_, err = snapshot.Read(buffer)
	if err != nil {
		return nil, err
	}

	decoder := gob.NewDecoder(bytes.NewReader(buffer))
	err = decoder.Decode(&log)
	if err != nil {
		return nil, err
	}

	return &Store{
		file:     file,
		snapshot: snapshot,
		log:      log,
	}, nil
}

func (s *Store) Write(key string, value string) error {
	s.Lock()
	defer s.Unlock()

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
		Offset: int(offset),
		Length: len(entry),
	}

	return nil
}

func (s *Store) Read(key string) (string, error) {
	s.RLock()
	defer s.RUnlock()

	if _, ok := s.log[key]; !ok {
		return "", ErrKeyNotFound(key)
	}

	length := s.log[key].Length
	offset := s.log[key].Offset

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
	s.Lock()
	defer s.Unlock()

	if _, ok := s.log[key]; !ok {
		return ErrKeyNotFound(key)
	}

	delete(s.log, key)

	return nil
}

func (s *Store) Close() error {
	defer s.file.Close()
	defer s.snapshot.Close()

	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)

	err := encoder.Encode(s.log)
	if err != nil {
		return err
	}

	_, err = s.snapshot.Write(buffer.Bytes())
	if err != nil {
		return err
	}

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

func fileExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}

func snapshotFileFromPath(path string) string {
	return fmt.Sprintf("%v.snapshot", path)
}

type pair struct {
	key   string
	value string
}

type logEntry struct {
	Offset int
	Length int
}
