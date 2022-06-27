## konvoy-image completion bash

Generate the autocompletion script for bash

### Synopsis

Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

	source <(konvoy-image completion bash)

To load completions for every new session, execute once:

#### Linux:

	konvoy-image completion bash > /etc/bash_completion.d/konvoy-image

#### macOS:

	konvoy-image completion bash > $(brew --prefix)/etc/bash_completion.d/konvoy-image

You will need to start a new shell for this setup to take effect.


```
konvoy-image completion bash
```

### Options

```
  -h, --help              help for bash
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

