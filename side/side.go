package side

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
	"encoding/json"
	"io/ioutil"
)

type Test struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Commands []Command `json:"commands"`
}

type Command struct {
	ID      string     `json:"id"`
	Comment string     `json:"comment"`
	Command string     `json:"command"`
	Target  string     `json:"target"`
	Targets [][]string `json:"targets"`
	Value   string     `json:"value"`
}

type Side struct {
	ID      string   `json:"id"`
	Version string   `json:"version"`
	Name    string   `json:"name"`
	Url     string   `json:"url"`
	Tests   []Test   `json:"tests"`
	Urls    []string `json:"urls"`
	Plugins []string `json:"plugins"`
}

func Read(filename string) (Side, error) {
	var side Side
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return side, err
	}

	err = json.Unmarshal(bytes, &side)
	return side, err
}
