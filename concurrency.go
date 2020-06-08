package main

import (
	"fmt"
	"time"
)

func main() {
	/*var wg sync.WaitGroup // Creates a wait group for the goroutine to execute and for the main function to wait
	wg.Add(1)
	go func() {
		count("sheep") // Creates goroutine
		wg.Done()
	}()
	wg.Wait()*/
	// We can make a buffet so we can send values to the channel without worrying about receiving by specifying the size of the channel
	c := make(chan string)
	go count("sheep", c)

	/*for {
		msg, open := <- c // It is a blocking statement
		if !open {
			break
		}
		fmt.Println(msg)
	}*/

	for msg := range c {
		fmt.Println(msg)
	}

}

func count(thing string, c chan string) {
	for i := 1; i<=5 ; i++ {
		c <- thing
		time.Sleep(time.Millisecond * 500)
	}
	close(c)
}
