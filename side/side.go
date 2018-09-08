package side

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
