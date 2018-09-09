package genselenide

import (
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"

	"../side"
)

func locator(target string) (string, error) {
	locatorMap := map[string]string{
		"id":    "\"#%s\"",
		"name":  "byName(\"%s\")",
		"css":   "\"%s\"",
		"xpath": "byXpath(\"%s\")",
		"index": "%s",
	}

	targetArray := strings.Split(target, "=")
	if len(targetArray) != 2 {
		return target, nil
	}
	kind := targetArray[0]
	locator := targetArray[1]

	fmtStr, ok := locatorMap[kind]
	if ok == false {
		return "", fmt.Errorf("Invalid Locator: %s", target)
	}
	return fmt.Sprintf(fmtStr, locator), nil
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

func GenerateJava(side side.Side, className string) (javacode []string, err error) {
	javaCodeMap := map[string]string{
		"open":        "%sopen(\"%s\")%s",
		"selectFrame": "%sswitchTo().frame(%s)",
		"click":       "%s$(%s).click()%s",
		"type":        "%s$(%s).val(\"%s\")",
		"sendKeys":    "%s$(%s).val(\"%s\")",
	}

	javacode = generateJavaHeader()

	javacode = append(javacode, "public class "+className+" {\n")
	h := "    "
	javacode = append(javacode, h+"@BeforeEach\n")
	javacode = append(javacode, h+"public void setup() {\n")
	javacode = append(javacode, h+h+"Configuration.baseUrl = \""+side.Url+"\";\n")
	javacode = append(javacode, h+h+"Configuration.browser = WebDriverRunner.CHROME;\n")
	javacode = append(javacode, h+"}\n\n")
	for _, test := range side.Tests {
		javacode = append(javacode, h+"@Test\n")
		javacode = append(javacode, h+"public void "+test.Name+"() {\n")
		for _, command := range test.Commands {
			locator, errLocator := locator(command.Target)
			if errLocator != nil {
				err = errLocator
				return
			}
			fmtStr, ok := javaCodeMap[command.Command]
			if ok == false {
				err = &ErrUnknownCommand{command: command.Command, value: command.Value}
				return
			}
			var value string
			switch command.Command {
			case "sendKeys":
				value, err = translateSendKeyValue(command.Value)
				if err != nil {
					return
				}
			case "open":
			case "click":
				value = ""
			default:
				value = command.Value
			}
			javacode = append(javacode, fmt.Sprintf(fmtStr+";\n", h+h, locator, value))
		}
		javacode = append(javacode, h+"}\n")

	}
	javacode = append(javacode, "}\n")
	return
}

func generateJavaHeader() []string {
	var javacode []string
	javacode = append(javacode, heredoc.Doc(`
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import com.codeborne.selenide.Configuration;
import com.codeborne.selenide.WebDriverRunner;

import static com.codeborne.selenide.Selenide.*;
import static com.codeborne.selenide.Condition.*;
import static com.codeborne.selenide.Selectors.*;
import static com.codeborne.selenide.WebDriverRunner.*;

`))

	return javacode
}
