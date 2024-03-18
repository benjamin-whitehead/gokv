package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/benjamin-whitehead/gokv/internal/resp"
)

type application struct {
	address string
	conn    net.Conn
}

func (a *application) repl() error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("'.exit' to quit")
	for {
		fmt.Print("> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		trimmed := strings.TrimSpace(text)
		cleaned := strings.TrimSuffix(trimmed, "\n")

		if cleaned == ".exit" {
			break
		}

		commands := strings.Split(cleaned, " ")
		encoded := resp.Encode(commands)

		err = a.sendCommandSync(encoded)
		if err != nil {
			log.Println(err)
		}
	}
	return nil
}

func (a *application) sendCommandSync(command string) error {
	err := a.connect()
	if err != nil {
		return err
	}
	defer a.conn.Close()

	fmt.Fprint(a.conn, command)

	status, err := bufio.NewReader(a.conn).ReadString('\n')
	if err != nil && err != io.EOF {
		return err
	}

	fmt.Println(status + "\n")

	return nil
}

func (a *application) connect() error {
	conn, err := net.Dial("tcp", a.address)
	if err != nil {
		return err
	}

	a.conn = conn
	return nil
}

func (a *application) parseFlags() {
	address := flag.String("address", "127.0.0.1:6380", "the address to connect to a Boxer server")
	flag.Parse()

	a.address = *address
}

func main() {
	app := application{}
	app.parseFlags()

	if err := app.repl(); err != nil {
		log.Fatal(err)
	}
}
