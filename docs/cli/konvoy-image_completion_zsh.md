## konvoy-image completion zsh

Generate the autocompletion script for zsh

### Synopsis

Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:

	echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions for every new session, execute once:

#### Linux:

	konvoy-image completion zsh > "${fpath[1]}/_konvoy-image"

#### macOS:

	konvoy-image completion zsh > /usr/local/share/zsh/site-functions/_konvoy-image

You will need to start a new shell for this setup to take effect.


```
konvoy-image completion zsh [flags]
```

### Options

```
  -h, --help              help for zsh
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
      --color     enable color output (default true)
  -v, --v int     select verbosity level, should be between 0 and 6
      --verbose   enable debug level logging (same as --v 5)
```

### SEE ALSO

* [konvoy-image completion](konvoy-image_completion.md)	 - Generate the autocompletion script for the specified shell

