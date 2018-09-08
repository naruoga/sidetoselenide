package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
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

type ErrInvalidLocator string

func (e ErrInvalidLocator) Error() string {
	return "Invalid Locator: " + string(e)
}

func locator(target string) (string, error) {
	switch {
	case strings.HasPrefix(target, "id"):
		return fmt.Sprintf("\"#%s\"", strings.TrimPrefix(target, "id=")), nil
	case strings.HasPrefix(target, "name"):
		return fmt.Sprintf("byName(\"%s\")", strings.TrimPrefix(target, "name=")), nil
	case strings.HasPrefix(target, "css"):
		return fmt.Sprintf("\"%s\"", strings.TrimPrefix(target, "css=")), nil
	case strings.HasPrefix(target, "xpath"):
		return fmt.Sprintf("byXpath(\"%s\")", strings.TrimPrefix(target, "xpath=")), nil
	case strings.HasPrefix(target, "index"):
		return strings.TrimPrefix(target, "index="), nil
	default:
		return "", ErrInvalidLocator(target)
	}
}

type ErrUnknownCommand struct {
	command, value string
}

func (e ErrUnknownCommand) Error() string {
	return "Error: unknown command: " + e.command + ":" + e.value
}

func translateSendKeyValue(value string) (string, error) {
	switch value {
	case "${KEY_ENTER}":
		return "Keys.ENTER", nil
	default:
		return "", fmt.Errorf("Error: unknown keycode: %s", value)
	}
}

func generateJava(side Side, className string) (javacode []string, err error) {
	javacode = generateJavaHeader()

	javacode = append(javacode, "public class "+className+" {\n")
	h := "    "
	javacode = append(javacode, h+"@BeforeEach\n")
	javacode = append(javacode, h+"public void setup() {\n")
	javacode = append(javacode, h+h+"Configuration.baseUrl=\""+side.Url+"\";\n")
	javacode = append(javacode, h+h+"Configuration.browser=WebDriverRunner.CHROME;\n")
	javacode = append(javacode, h+"}\n\n")
	for _, test := range side.Tests {
		javacode = append(javacode, h+"@Test\n")
		javacode = append(javacode, h+"public void "+test.Name+"() {\n")
		for _, command := range test.Commands {
			locator, errLocator := locator(command.Target)
			if command.Command != "open" && errLocator != nil {
				err = errLocator
				return
			}
			switch command.Command {
			case "open":
				javacode = append(javacode, fmt.Sprintf("%sopen(\"%s\");\n", h+h, command.Target))
			case "selectFrame":
				javacode = append(javacode, fmt.Sprintf("%sswitchTo().frame(%s);\n", h+h, locator))
			case "click":
				javacode = append(javacode, fmt.Sprintf("%s$(%s).click();\n", h+h, locator))
			case "type":
				javacode = append(javacode, fmt.Sprintf("%s$(%s).val(\"%s\");\n", h+h, locator, command.Value))
			case "sendKeys":
				var keycode string
				keycode, err = translateSendKeyValue(command.Value)
				if err != nil {
					return
				}
				javacode = append(javacode, fmt.Sprintf("%s$(%s).val(\"%s\");\n", h+h, locator, keycode))
			default:
				err = &ErrUnknownCommand{command: command.Command, value: command.Value}
				return
			}
		}
		javacode = append(javacode, h+"}\n")

	}
	javacode = append(javacode, "}\n")
	return
}

func generateJavaHeader() []string {
	var javacode []string
	javacode = append(javacode,
		"import org.junit.jupiter.api.Test;\n",
		"\n",
		"import static com.codeborne.selenide.Selenide.*;\n",
		"import static com.codeborne.selenide.Condition.*;\n",
		"import static com.codeborne.selenide.Selectors.*;\n",
		"\n")

	return javacode
}

func readSide(filename string) (Side, error) {
	var side Side
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return side, err
	}

	err = json.Unmarshal(bytes, &side)
	return side, err
}

func writeSideToJava(side Side) error {
	className := strings.Replace(side.Name, " ", "", -1)
	file, err := os.Create(className + ".java")
	fmt.Println("generate: " + file.Name())
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer file.Close()

	javacode, err3 := generateJava(side, className)
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
			side, err := readSide(filename)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(filename, side)
			writeSideToJava(side)
		}
	}
}
