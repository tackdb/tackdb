package main

import (
	"fmt"
)

func put(m map[string]*string) {
	test := "test"
	m["string"] = &test
}

func ret() (string, error) {
	return "hello~", nil
}

type StringError struct {
	String string
	Error  error
}

func main() {
	// m := make(map[string][]byte)
	// m2 := make(map[string][]byte)
	// m["foo"] = nil
	// // bar := "baz"
	// m2["bar"] = m["bar"]
	// fmt.Println(m)
	// fmt.Println(m2)
	// b, ok := m2["bar"]
	// fmt.Println(b, ok)

	// m := make(map[string]*string)
	// // m2 := make(map[string]*string)
	// // // m["foo"]++ // = m["foo"] + 1
	// // m2["foo"] = m["foo"]
	// put(m)
	// m := make([]string, 0)
	// fmt.Println(m[0])

	// store := tackdb.NewStore()
	// fmt.Println(store)

	s := &StringError{ret()}

	fmt.Println(s)
}
