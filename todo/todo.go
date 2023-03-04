package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/alexeyco/simpletable"
)

const (
	TodoFile = ".todos.json"
)

type todo struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

type Todos []todo

func (t *Todos) Add(task string) {
	todo := todo{
		Task:        task,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}

	*t = append(*t, todo)

	err := t.Store(TodoFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	t.Print()
}

func (t *Todos) Complete(index int) error {
	ls := *t

	if index <= 0 || index > len(ls) {
		return errors.New("invalid index")
	}

	ls[index-1].CompletedAt = time.Now()
	ls[index-1].Done = true

	err := ls.Store(TodoFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	t.Print()

	return nil
}

func (t *Todos) Delete(index int) error {
	ls := *t

	if (index <= 0) || index > len(ls) {
		return errors.New("invalid index")
	}

	*t = append(ls[:index-1], ls[index:]...)

	err := ls.Store(TodoFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	t.Print()

	return nil
}

func (t *Todos) Load(filename string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	if len(file) == 0 {
		return err
	}

	err = json.Unmarshal(file, t)
	if err != nil {
		return err
	}

	return nil
}

func (t *Todos) Store(filename string) error {
	data, err := json.Marshal(t)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 8644)
}

func (t *Todos) Print() {
	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "Task"},
			{Align: simpletable.AlignCenter, Text: "Complete?"},
			{Align: simpletable.AlignCenter, Text: "CreatedAt"},
			{Align: simpletable.AlignCenter, Text: "CompletedAt"},
		},
	}

	var cells [][]*simpletable.Cell

	completedAt := func(timeArg time.Time) string {
		if timeArg.IsZero() {
			return "--"
		} else {
			return timeArg.Format(time.RFC822)
		}
	}

	isCompleteIcon := func(isDone bool) string {
		if isDone {
			return "✅"
		} else {
			return "❌"
		}
	}

	for i, todos := range *t {
		i++
		cells = append(cells, []*simpletable.Cell{
			{Text: fmt.Sprintf("%d", i)},
			{Text: todos.Task},
			{Text: isCompleteIcon(todos.Done), Align: simpletable.AlignCenter},
			{Text: todos.CreatedAt.Format(time.RFC822)},
			{Text: completedAt(todos.CompletedAt), Align: simpletable.AlignCenter},
		})
	}

	table.Body = &simpletable.Body{Cells: cells}

	table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
		{Align: simpletable.AlignCenter, Span: 5, Text: fmt.Sprintf("Your have %d pending todos", t.CountPending())},
	}}

	table.SetStyle(simpletable.StyleUnicode)

	table.Println()
}

func (t *Todos) CountPending() int {
	total := 0
	for _, todos := range *t {
		if !todos.Done {
			total++
		}
	}

	return total
}
