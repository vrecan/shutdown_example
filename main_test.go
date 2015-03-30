package main

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"runtime"
	"testing"
)

func TestApp(t *testing.T) {
	runtime.GOMAXPROCS(4)
	Convey("Test creating and sending message with close", t, func() {
		app := NewApp()
		app.Start()
		for i := 1; i <= 10; i++ {
			app.Data <- fmt.Sprintf("TESTING %d", i)
		}
		app.Close()
	})
}
