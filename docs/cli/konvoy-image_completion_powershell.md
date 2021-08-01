## konvoy-image completion powershell

generate the autocompletion script for powershell

### Synopsis


Generate the autocompletion script for powershell.

To load completions in your current shell session:
PS C:\> konvoy-image completion powershell | Out-String | Invoke-Expression

To load completions for every new session, add the output of the above command
to your powershell profile.


```
konvoy-image completion powershell [flags]
```

### Options

```
  -h, --help              help for powershell
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

