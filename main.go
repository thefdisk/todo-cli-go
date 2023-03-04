package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/thefdisk/todo-cli-go/todo"
)

func main() {
	add := flag.Bool("add", false, "add a new todo")
	complete := flag.Int("complete", 0, "mark a todo as completed")
	del := flag.Int("delete", 0, "delete a todo")
	list := flag.Bool("list", false, "list all todo")

	flag.Parse()

	todos := &todo.Todos{}

	if err := todos.Load(todo.TodoFile); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	switch {
	case *add:
		task, err := filterInput(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		todos.Add(task)
	case *complete > 0:
		err := todos.Complete(*complete)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case *del > 0:
		err := todos.Delete(*del)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case *list:
		todos.Print()
	default:
		fmt.Println("invalid command")
		os.Exit(0)
	}
}

func filterInput(r io.Reader, args ...string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}

	scanner := bufio.NewScanner(r)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return "", nil
	}

	text := scanner.Text()

	if len(text) == 0 {
		return "", errors.New("empty todo is not allowed")
	}

	return text, nil
}
