package resp

import (
	"reflect"
	"testing"
)

func TestEncode(t *testing.T) {
	t.Run("encode SET", func(t *testing.T) {
		command := []string{"SET", "hello", "world"}

		want := "*3\r\n$3\r\nSET\r\n$5\r\nhello\r\n$5\r\nworld\r\n"
		got := Encode(command)

		assertStringEquals(t, got, want)
	})

	t.Run("encode GET", func(t *testing.T) {
		command := []string{"GET", "hello"}

		want := "*2\r\n$3\r\nGET\r\n$5\r\nhello\r\n"
		got := Encode(command)

		assertStringEquals(t, got, want)
	})

	t.Run("encode DEL", func(t *testing.T) {
		command := []string{"DEL", "hello"}

		want := "*2\r\n$3\r\nDEL\r\n$5\r\nhello\r\n"
		got := Encode(command)

		assertStringEquals(t, got, want)
	})
}

func TestDecode(t *testing.T) {
	t.Run("decode SET command", func(t *testing.T) {
		command := "*3\r\n$3\r\nSET\r\n$5\r\nhello\r\n$5\r\nworld\r\n"

		want := []string{"SET", "hello", "world"}
		got := Decode(command)

		assertSliceEquals(t, got, want)
	})

	t.Run("decode GET command", func(t *testing.T) {
		command := "*2\r\n$3\r\nGET\r\n$5\r\nhello\r\n"

		want := []string{"GET", "hello"}
		got := Decode(command)

		assertSliceEquals(t, got, want)
	})

	t.Run("decode DEL command", func(t *testing.T) {
		command := "*2\r\n$3\r\nDEL\r\n$5\r\nhello\r\n"

		want := []string{"DEL", "hello"}
		got := Decode(command)

		assertSliceEquals(t, got, want)
	})
}

func TestEncodeDecode(t *testing.T) {
	want := "*3\r\n$3\r\nSET\r\n$5\r\nhello\r\n$5\r\nworld\r\n"
	got := Encode(Decode(want))

	assertStringEquals(t, got, want)
}

func assertStringEquals(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Fatalf("got %s want %s", got, want)
	}
}

func assertSliceEquals(t testing.TB, got, want []string) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}
}
