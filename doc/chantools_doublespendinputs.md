## chantools doublespendinputs

Replace a transaction by double spending its input

### Synopsis

Tries to double spend the given inputs by deriving the
private for the address and sweeping the funds to the given address. This can
only be used with inputs that belong to an lnd wallet.

```
chantools doublespendinputs [flags]
```

### Examples

```
chantools doublespendinputs \
	--inputoutpoints xxxxxxxxx:y,xxxxxxxxx:y \
	--sweepaddr bc1q..... \
	--feerate 10 \
	--publish
```

### Options

```
      --apiurl string            API URL to use (must be esplora compatible) (default "https://api.node-recovery.com")
      --bip39                    read a classic BIP39 seed and passphrase from the terminal instead of asking for lnd seed format or providing the --rootkey flag
      --feerate uint32           fee rate to use for the sweep transaction in sat/vByte (default 30)
  -h, --help                     help for doublespendinputs
      --inputoutpoints strings   list of outpoints to double spend in the format txid:vout
      --publish                  publish replacement TX to the chain API instead of just printing the TX
      --recoverywindow uint32    number of keys to scan per internal/external branch; output will consist of double this amount of keys (default 2500)
      --rootkey string           BIP32 HD root key of the wallet to use for deriving the input keys; leave empty to prompt for lnd 24 word aezeed
      --sweepaddr string         address to recover the funds to; specify 'fromseed' to derive a new address from the seed automatically
      --walletdb string          read the seed/master root key to use for deriving the input keys from an lnd wallet.db file instead of asking for a seed or providing the --rootkey flag
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

