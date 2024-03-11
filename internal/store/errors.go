package store

import (
	"fmt"
)

func ErrKeyNotFound(key string) error {
	return fmt.Errorf("key %v not found", key)
}
