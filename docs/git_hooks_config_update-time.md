## git hooks config update-time

Changes the Githooks update time.

### Synopsis

Changes the Githooks update time used to check for updates.

Resets the last Githooks update time with the `--reset` option,
causing the update check to run next time if it is enabled.
Use 'git hooks update [--enable|--disable]' to change that setting.

```
git hooks config update-time [flags]
```

### Options

```
  -h, --help    help for update-time
      --print   Print the setting.
      --reset   Reset the setting.
```

### SEE ALSO

* [git hooks config](git_hooks_config.md)	 - Manages various Githooks configuration.

###### Auto generated by spf13/cobra on 11-Jan-2021