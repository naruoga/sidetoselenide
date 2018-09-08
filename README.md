# sidetoselenide: Selenide code generator from *.side

just a mockup, my first golang code.

generate Selenide + JUnit5 java code snippet from Selenium IDE's `*.side` format.

## usage

```command
go run sidetoselenide.go foo.side
```

## todo

- devide one big `.go` file into several module files
- improve packaging
- support more commands
- support more locators
- move supported commands / locators to file (JSON or else)
