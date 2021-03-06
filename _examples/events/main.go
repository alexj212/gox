package main

import (
	"fmt"

	"github.com/alexj212/gox/events"
)

func main() {
	e := events.New[string]()

	unsubscribe := e.Subscribe(func(data string) {
		fmt.Println("Subscriber 1:", data)
	})

	e.Subscribe(func(data string) {
		fmt.Println("Subscriber 2:", data)
	})

	e.SubscribeOnce(func(data string) {
		fmt.Println("Subscriber 3:", data)
	})

	e.Dispatch("Hello")
	unsubscribe()
	e.Dispatch("World")
}
