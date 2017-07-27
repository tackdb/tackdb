package tackdb

import "testing"

func TestMain(m *testing.M) {
	done := make(chan error)
	server := NewServer()
	go func() {
		done <- server.Listen().Serve()
	}()
	m.Run()
	done <- nil
	server.listener.Close()
	if err := <-done; err != nil {
		panic(err)
	}
}
