## konvoy-image completion fish

generate the autocompletion script for fish

### Synopsis


Generate the autocompletion script for the fish shell.

To load completions in your current shell session:
$ konvoy-image completion fish | source

To load completions for every new session, execute once:
$ konvoy-image completion fish > ~/.config/fish/completions/konvoy-image.fish

You will need to start a new shell for this setup to take effect.


```
konvoy-image completion fish [flags]
```

### Options

```
  -h, --help              help for fish
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
      --color     enable color output (default true)
  -v, --v int     select verbosity level, should be between 0 and 6
      --verbose   enable debug level logging (same as --v 5)
```

### SEE ALSO

* [konvoy-image completion](konvoy-image_completion.md)	 - generate the autocompletion script for the specified shell

