package main

import (
	"flag"
	"fmt"
)

func main() {
	// Define a string flag with a default value and a short description.
	name := flag.String("name", "World", "a name to say hello to")

	// Parse the flags provided by the user.
	flag.Parse()

	// Use the flag value in the program.
	fmt.Printf("Hello, %s!\n", *name)
}
