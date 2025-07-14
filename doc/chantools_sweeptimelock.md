## chantools sweeptimelock

Sweep the force-closed state after the time lock has expired

### Synopsis

Use this command to sweep the funds from channels that
you force-closed with the forceclose command. You **MUST** use the result file
that was created with the forceclose command, otherwise it won't work. You also
have to wait until the highest time lock (can be up to 2016 blocks which is more
than two weeks) of all the channels has passed. If you only want to sweep
channels that have the default CSV limit of 1 day, you can set the --maxcsvlimit
parameter to 144.

```
chantools sweeptimelock [flags]
```

### Examples

```
chantools sweeptimelock \
	--fromsummary results/forceclose-xxxx-yyyy.json \
	--sweepaddr bc1q..... \
	--feerate 10 \
  	--publish
```

### Options

```
      --apiurl string            API URL to use (must be esplora compatible) (default "https://api.node-recovery.com")
      --bip39                    read a classic BIP39 seed and passphrase from the terminal instead of asking for lnd seed format or providing the --rootkey flag
      --feerate uint32           fee rate to use for the sweep transaction in sat/vByte (default 30)
      --fromchanneldb string     channel input is in the format of an lnd channel.db file
      --fromchanneldump string   channel input is in the format of a channel dump file
      --fromsummary string       channel input is in the format of chantool's channel summary; specify '-' to read from stdin
  -h, --help                     help for sweeptimelock
      --listchannels string      channel input is in the format of lncli's listchannels format; specify '-' to read from stdin
      --maxcsvlimit uint16       maximum CSV limit to use (default 2016)
      --pendingchannels string   channel input is in the format of lncli's pendingchannels format; specify '-' to read from stdin
      --publish                  publish sweep TX to the chain API instead of just printing the TX
      --rootkey string           BIP32 HD root key of the wallet to use for deriving keys; leave empty to prompt for lnd 24 word aezeed
      --sweepaddr string         address to recover the funds to; specify 'fromseed' to derive a new address from the seed automatically
      --walletdb string          read the seed/master root key to use for deriving keys from an lnd wallet.db file instead of asking for a seed or providing the --rootkey flag
```

### Options inherited from parent commands

```
      --nologfile   If set, no log file will be created. This is useful for testing purposes where we don't want to create a log file.
  -r, --regtest     Indicates if regtest parameters should be used
  -s, --signet      Indicates if the public signet parameters should be used
  -t, --testnet     Indicates if testnet parameters should be used
```

### SEE ALSO

* [chantools](chantools.md)	 - Chantools helps recover funds from lightning channels

