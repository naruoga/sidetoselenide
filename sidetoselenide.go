package main

import (
	"log"
	"os"
	"strings"

	"./genselenide"
	"./side"
)

func writeSideToJava(side side.Side) error {
	className := strings.Replace(side.Name, " ", "", -1)
	file, err := os.Create(className + ".java")
	log.Println("generate: " + file.Name())
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer file.Close()

	javacode, err3 := genselenide.GenerateJava(side, className)
	if err3 != nil {
		log.Fatal(err3)
		return err3
	}
	for _, l := range javacode {
		file.WriteString(l)
	}

	return nil
}

func main() {
	if len(os.Args) == 1 {
		log.Fatal(os.Args[0] + ": 引数に *.sideファイルを指定してください")
	}

	for i, filename := range os.Args {
		if i > 0 {
			side, err := side.Read(filename)
			if err != nil {
				log.Fatal(err)
			}
			writeSideToJava(side)
		}
	}
}
