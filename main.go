package main

import (
	"fmt"
	DEATH "github.com/vrecan/death"
	"sync"
	SYS "syscall"
)

//App is an example struct that will run goroutines
type App struct {
	wg   *sync.WaitGroup //allow us to wait for close
	done chan bool
	Data chan string
}

func NewApp() (app *App) {
	app = new(App)
	app.wg = new(sync.WaitGroup)
	app.done = make(chan bool, 1)
	app.Data = make(chan string, 100)
	return app
}

//Start go routines
func (a *App) Start() {
	a.wg.Add(1)
	go a.runSelect()
}

//Run with select statement to manage shutdown
func (a *App) runSelect() {
	defer a.wg.Done()
loop:
	for {
		select {

		case <-a.done:
			fmt.Println("Detected close doing cleanup and exiting")
			a.wg.Add(1)
			a.runClose()
			break loop
			//do cleanup if needed
		case s, ok := <-a.Data:
			if ok { //ok will tell us if the channel has been shutdown
				fmt.Println("select: ", s)
			} else {
				//done processing all messages exit
				fmt.Println("Breaking loop")
				break loop
			}
		}
	}
}

//process all messaging on the queue until we get close message
func (a *App) runClose() {
	for s := range a.Data {
		fmt.Println("close: ", s)
	}
}

//Cleanly close background routines
func (a *App) Close() {
	if a != nil {
		a.done <- true
		close(a.Data) //this lets the recv know there is no more data
		a.wg.Wait()
	}
}

//Simple example on managing shutdown
func main() {
	death := DEATH.NewDeath(SYS.SIGINT, SYS.SIGTERM)
	app := NewApp()
	app.Start()
	// app.Start()
	for i := 1; i <= 100; i++ {
		app.Data <- fmt.Sprintf("TESTING %d", i)
	}
	//use death
	death.WaitForDeath(app)
	//or call close manually
	//app.Close()
}
