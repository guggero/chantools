## chantools removechannel

Remove a single channel from the given channel DB

```
chantools removechannel [flags]
```

### Options

```
      --channel string     channel to remove from the DB file, identified by its channel point (<txid>:<txindex>)
      --channeldb string   lnd channel.backup file to remove the channel from
  -h, --help               help for removechannel
```

### Options inherited from parent commands

```
  -r, --regtest   Indicates if regtest parameters should be used
  -t, --testnet   Indicates if testnet parameters should be used
```

### SEE ALSO

* [chantools](chantools.md)	 - Chantools helps recover funds from lightning channels
