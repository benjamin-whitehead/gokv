package resp

import (
	"bytes"
	"fmt"
	"strings"
)

// Encode encodes a command into a RESP string for processing on the server
func Encode(command []string) string {
	var buffer bytes.Buffer

	commandLength := len(command)
	buffer.WriteString(fmt.Sprintf("*%d\r\n", commandLength))

	for _, element := range command {
		elementLength := len(element)

		buffer.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", elementLength, element))
	}

	return buffer.String()
}

type DataType int

const (
	SIMPLE_STRINGS DataType = iota
	SIMPLE_ERRORS
	INTEGERS
	BULK_STRINGS
	ARRAYS
)

var commandPrefixs = map[rune]DataType{
	'+': SIMPLE_STRINGS,
	'-': SIMPLE_ERRORS,
	':': INTEGERS,
	'$': BULK_STRINGS,
	'*': ARRAYS,
}

// Decode decodes a RESP string into an array representation of the command
func Decode(command string) []string {
	commands := strings.Split(command, "\r\n")
	runes := make([][]rune, len(commands))

	for index := range commands {
		runes[index] = []rune(commands[index])
	}

	buffer := make([]string, 0)
	for _, command := range runes[2:] {
		if len(command) > 1 {
			prefix := command[0]
			if _, ok := commandPrefixs[prefix]; !ok {
				buffer = append(buffer, string(command))
			}
		}
	}

	return buffer
}
