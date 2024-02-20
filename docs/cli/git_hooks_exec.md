## git hooks exec

Execute namespace paths pointing to an executable or run configuration.

### Synopsis

Execute namespace paths
pointing to an executable or a run configuration
(e.g. `ns:xxx/mypath/a/b/c.sh` or `ns:xxx/mypath/a/b/c.yaml`).
The execution is run the same as Githooks performs during its execution.

Its not meant to execute hooks but rather add-on scripts
inside Githooks repositories.

If containerized hooks are enabled the execution always runs containerized.

```
git hooks exec namespace-path [args...]
```

### Options

```
      --containerized   Force the execution to be containerized.
  -h, --help            help for exec
```

### SEE ALSO

* [git hooks](git_hooks.md)	 - Githooks CLI application

###### Auto generated by spf13/cobra 