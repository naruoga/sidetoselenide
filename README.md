# sidetoselenide: Selenide code generator from `*.side`

This is a toy project to learn Go language.

Generate Selenide + JUnit5 java code snippet from Selenium IDE's `*.side` format.

> Note: this project does not to aim whole commands of Selenium IDE,
> just focused on it's capture-and-replay features.
> On the other words, I have no plan to support:
> 
> - variables
> - flow control commands (if, loop, ...)

## usage

```command
go run sidetoselenide.go foo.side
```

## License

see [LICENSE](LICENSE)

## todo

- support more commands
- support more locators
- move supported commands / locators to file (JSON or else)
