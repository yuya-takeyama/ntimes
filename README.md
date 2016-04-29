# ntimes

Command to execute command N times

## Usage

```
$ ntimes 3 -- echo foo bar baz
foo bar baz
foo bar baz
foo bar baz
```

### Set parallel degree of execution (-p)

```
$ ntimes 10 -p 3 -- sh -c 'echo "Hi!"; sleep 1; echo "Bye"'
```

## Author

Yuya Takeyama

## License

The MIT License
