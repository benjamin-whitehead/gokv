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

func TestCreateStoreFromFile(t *testing.T) {

	t.Run("missing file gives error", func(t *testing.T) {
		store, err := NewStoreFromFile("hello")

		assertError(t, err, ErrFileNotFound("hello"))

		if store != nil {
			t.Error("expected store to be nil, but got value")
		}
	})

	t.Run("existing store is recreated", func(t *testing.T) {
		store := storeWithPairs(t, pair{key: "hello", value: "world"})

		err := store.Close()
		assertNoError(t, err)

		storeFile := store.file.Name()

		recreatedStore, err := NewStoreFromFile(storeFile)
		assertNoError(t, err)

		value, err := recreatedStore.Read("hello")
		assertNoError(t, err)
		if value != "world" {
			t.Errorf("expected to read value %v, got %v", "world", value)
		}
	})

	t.Run("recreate store without closing gives error", func(t *testing.T) {
		// TODO
	})
}

func TestWrite(t *testing.T) {
	name := createFilePath(t)

	store, err := NewStore(name)
	assertNoError(t, err)

	err = store.Write("hello", "world")
	assertNoError(t, err)
}

func TestRead(t *testing.T) {

	t.Run("read value for key", func(t *testing.T) {
		store := storeWithPairs(t, pair{key: "hello", value: "world"})

		value, err := store.Read("hello")
		assertNoError(t, err)

		if value != "world" {
			t.Errorf("expected %v but got %v", "world", value)
		}
	})

	t.Run("missing key gives error", func(t *testing.T) {
		store := emptyStore(t)
		value, err := store.Read("hello")

		if value != "" {
			t.Errorf("expected \"\" but got %v", value)
		}

		assertError(t, err, ErrKeyNotFound("hello"))
	})

	t.Run("supports concurrent reads", func(t *testing.T) {
		store := storeWithPairs(t, pair{key: "hello", value: "world"})

		concurrentIteration(t, 1000, func() {
			value, err := store.Read("hello")
			assertNoError(t, err)

			if value != "world" {
				t.Errorf("expected %v but got %v", "world", value)
			}
		})
	})

}

func TestDelete(t *testing.T) {

	t.Run("delete value for key", func(t *testing.T) {
		store := storeWithPairs(t, pair{key: "hello", value: "world"})

		err := store.Delete("hello")
		assertNoError(t, err)
	})

	t.Run("missing key gives error", func(t *testing.T) {
		store := emptyStore(t)
		err := store.Delete("hello")

		assertError(t, err, ErrKeyNotFound("hello"))
	})

}

func assertNoError(t testing.TB, got error) {
	t.Helper()

	if got != nil {
		t.Error("didn't want an error, but got one")
	}
}

func assertError(t testing.TB, got error, want error) {
	t.Helper()

	if got == nil {
		t.Fatal("wanted an error but didn't get one")
	}

	if got.Error() != want.Error() {
		t.Errorf("expected %v got %q", want, got)
	}
}

func concurrentIteration(t testing.TB, iterations int, toExecute func()) {
	t.Helper()

	for i := 0; i < iterations; i++ {
		go toExecute()
	}
}

func createFilePath(t testing.TB) string {
	t.Helper()

	return path.Join(t.TempDir(), "test.gokv")
}

func emptyStore(t testing.TB) *Store {
	t.Helper()

	store, err := NewStore(createFilePath(t))
	assertNoError(t, err)

	return store
}

func storeWithPairs(t testing.TB, pairs ...pair) *Store {
	t.Helper()

	store := emptyStore(t)
	for _, pair := range pairs {
		err := store.Write(pair.key, pair.value)
		assertNoError(t, err)
	}

	return store
}
