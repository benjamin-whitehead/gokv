package store

import (
	"errors"
	"fmt"
)

func ErrKeyNotFound(key string) error {
	return fmt.Errorf("key %v not found", key)
}

func ErrReadIncorrectKey(key, readKey string) error {
	return fmt.Errorf("can't read key %v, read key %v instead", key, readKey)
}

func ErrFileNotFound(name string) error {
	return fmt.Errorf("can't open file %v, file not found", name)
}

var ErrDecodeEntry = errors.New("can't decode entry from file")
