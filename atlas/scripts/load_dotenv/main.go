package main

import (
	"encoding/json"
	"fmt"
	"os"

	dotenv "github.com/joho/godotenv"
)

func main() {
	args := os.Args[1:]
	dotenv.Load(args[0])

	file, err := os.Open(args[0])
	if err != nil {
		panic(err)
	}
	defer file.Close()

	m, err := dotenv.Parse(file)
	if err != nil {
		panic(err)
	}
	j, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(j))
}
