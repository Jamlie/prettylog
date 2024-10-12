# Pretty Log

Implementing the `slog.Handler` in Go

## Example

```go
func main() {
	slog.SetDefault(slog.New(NewHandler(nil)))

	slog.Info("You did it!", slog.Group(
		"user",
		"name", "foo",
		"drink": "water",
	))

	slog.Warn("No! This will destroy you!", slog.Group(
		"user",
		"name", "bar",
		"drink": "vermouth",
	))

	slog.Error("I TOLD YOU NOT TO DO IT", slog.Group(
		"user",
		"name", "baz",
		"drink": "dizzy drinks",
	))
}
```

## Override

```go
func main() {
	slog.SetDefault(slog.New(NewHandler(&slog.HandlerOptions{
		AddSource: true, // to show the source file, line, etc
		Level:     slog.LevelDebug, // could be Info, Warn, Error, etc
	})))

	slog.Info("You did it!", slog.Group(
		"user",
		"name", "foo",
		"drink": "water",
	))

	slog.Warn("No! This will destroy you!", slog.Group(
		"user",
		"name", "bar",
		"drink": "vermouth",
	))

	slog.Error("I TOLD YOU NOT TO DO IT", slog.Group(
		"user",
		"name", "baz",
		"drink": "dizzy drinks",
	))

	slog.Debug("you can try to do anything lol idc", slog.Group(
		"user",
		"name", "qux",
		"drink": "all",
	))
}
```
