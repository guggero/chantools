## chantools filterbackup

Filter an lnd channel.backup file and remove certain channels

### Synopsis

Filter an lnd channel.backup file by removing certain 
channels (identified by their funding transaction outpoints).

```
chantools filterbackup [flags]
```

### Examples

```
chantools filterbackup \
	--multi_file ~/.lnd/data/chain/bitcoin/mainnet/channel.backup \
	--discard 2abcdef2b2bffaaa...db0abadd:1,4abcdef2b2bffaaa...db8abadd:0
```

### Options

```
      --bip39               read a classic BIP39 seed and passphrase from the terminal instead of asking for lnd seed format or providing the --rootkey flag
      --discard string      comma separated list of channel funding outpoints (format <fundingTXID>:<index>) to remove from the backup file
  -h, --help                help for filterbackup
      --multi_file string   lnd channel.backup file to filter
      --rootkey string      BIP32 HD root key of the wallet to use for decrypting the backup; leave empty to prompt for lnd 24 word aezeed
      --walletdb string     read the seed/master root key to use for decrypting the backup from an lnd wallet.db file instead of asking for a seed or providing the --rootkey flag
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

