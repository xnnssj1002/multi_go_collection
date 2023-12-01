package multi_go_collection

import (
	"fmt"
	"testing"
	"time"
)

type CPM struct {
	Name string
	Age  int
}

func TestTask(t *testing.T) {
	handler := NewHandler()

	handler.Add(NewTask(
		func() interface{} { return nil },
		WithTaskName("【nil task】"),
	))

	handler.Add(NewTask(
		func() interface{} { return 1 },
		WithTaskName("【int task】"),
		WithTimeout(1),
	))

	handler.Add(NewTask(
		func() interface{} {
			time.Sleep(2 * time.Second)
			return CPM{
				Name: "CPM1",
				Age:  10,
			}
		},
		WithTaskName("【CPM struct task】"),
		WithTimeout(1),
	))

	handler.Add(NewTask(
		func() interface{} {
			a := 10
			return &a
		},
		WithTaskName("【point task】"),
		WithTimeout(1),
	))

	handler.Run()

	//handler.Range(func(v interface{}) {
	handler.Visit(func(v interface{}) {
		if v == nil {
			fmt.Printf("res is nil\n")
			return
		}

		switch v.(type) {
		case int:
			fmt.Printf("res type is int, value is %d\n", v)
		case string:
			fmt.Printf("res type is string, value is %s\n", v)
		case CPM:
			cpm, ok := v.(CPM)
			fmt.Printf("res type is CPM, assertion is %t, name is %s, age is %d\n", ok, cpm.Name, cpm.Age)
		default:
			fmt.Printf("res type is %T\n", v)
		}
	})

}
