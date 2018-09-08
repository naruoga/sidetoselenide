package selenidegen

import (
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"

	"../side"
)

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
		return "", fmt.Errorf("Invalid Locator: %s", target)
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

func GenerateJava(side side.Side, className string) (javacode []string, err error) {
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
