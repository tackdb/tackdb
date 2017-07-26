package tackdb

import "testing"

func TestMain(m *testing.M) {
	done := make(chan error)
	go func() {
		done <- NewServer().Listen().Serve()
	}()
	m.Run()
	if err := <-done; err != nil {
		panic(err)
	}
}
