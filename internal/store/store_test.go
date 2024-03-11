package store

import (
	"path"
	"testing"
)

func TestCreateStore(t *testing.T) {
	name := createFilePath(t)

	_, err := NewStore(name)
	assertNoError(t, err)
}

func TestWrite(t *testing.T) {
	name := createFilePath(t)

	store, err := NewStore(name)
	assertNoError(t, err)

	err = store.Write("hello", "world")
	assertNoError(t, err)

	assertMapSize(t, "offsets", store.log, 1)
}

func TestRead(t *testing.T) {
	name := createFilePath(t)

	store, err := NewStore(name)
	assertNoError(t, err)

	err = store.Write("hello", "world")
	assertNoError(t, err)

	value, err := store.Read("hello")
	assertNoError(t, err)

	if value != "world" {
		t.Errorf("expected %v got %v", "world", value)
	}
}

func assertMapSize(t testing.TB, name string, m map[string]logEntry, want int) {
	t.Helper()

	if len(m) != want {
		t.Errorf("expected %v to have %d keys, got %d", name, want, len(m))
	}
}

func assertNoError(t testing.TB, got error) {
	t.Helper()

	if got != nil {
		t.Error("didn't want an error, but got one")
	}
}

func createFilePath(t testing.TB) string {
	t.Helper()

	return path.Join(t.TempDir(), "test.gokv")
}
