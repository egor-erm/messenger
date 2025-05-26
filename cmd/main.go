package main

import "messenger/core"

func main() {
	messenger, err := core.NewCore()
	if err != nil {
		panic(err)
	}

	messenger.Run()
}
