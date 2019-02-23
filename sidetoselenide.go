package main

/*
 * This file is part of sidetoselenide.
 * sidetoselenide is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * sidetoselenide is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with sidetoselenide.  If not, see <http://www.gnu.org/licenses/>.
 */

import (
	"log"
	"os"
	"strings"

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
