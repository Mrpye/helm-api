## helm-api completion zsh

Generate the autocompletion script for zsh

### Synopsis

Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:

	echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions in your current shell session:

	source <(helm-api completion zsh); compdef _helm-api helm-api

To load completions for every new session, execute once:

#### Linux:

	helm-api completion zsh > "${fpath[1]}/_helm-api"

#### macOS:

	helm-api completion zsh > $(brew --prefix)/share/zsh/site-functions/_helm-api

You will need to start a new shell for this setup to take effect.


```
helm-api completion zsh [flags]
```

### Options

```
  -h, --help              help for zsh
      --no-descriptions   disable completion descriptions
```

### SEE ALSO

* [helm-api completion](helm-api_completion.md)	 - Generate the autocompletion script for the specified shell

###### Auto generated by spf13/cobra on 14-Feb-2023
