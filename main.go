package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		println("Wrong input")
		return
	}
	input := os.Args[1]
	file, err := os.Open(input)
	if err != nil {
		fmt.Println("Dosya açma hatası:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Dosya okuma hatası", err)
	}
}
