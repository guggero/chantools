## chantools signrescuefunding

Rescue funds locked in a funding multisig output that never resulted in a proper channel; this is the command the remote node (the non-initiator) of the channel needs to run

### Synopsis

This is part 2 of a two phase process to rescue a channel
funding output that was created on chain by accident but never resulted in a
proper channel and no commitment transactions exist to spend the funds locked in
the 2-of-2 multisig.

If successful, this will create a final on-chain transaction that can be
broadcast by any Bitcoin node.

```
chantools signrescuefunding [flags]
```

### Examples

```
chantools signrescuefunding \
	--psbt <the_base64_encoded_psbt_from_step_1>
```

### Options

```
      --bip39             read a classic BIP39 seed and passphrase from the terminal instead of asking for lnd seed format or providing the --rootkey flag
  -h, --help              help for signrescuefunding
      --psbt string       Partially Signed Bitcoin Transaction that was provided by the initiator of the channel to rescue
      --rootkey string    BIP32 HD root key of the wallet to use for deriving keys; leave empty to prompt for lnd 24 word aezeed
      --walletdb string   read the seed/master root key to use for deriving keys from an lnd wallet.db file instead of asking for a seed or providing the --rootkey flag
```

### Options inherited from parent commands

```
      --nologfile           If set, no log file will be created. This is useful for testing purposes where we don't want to create a log file.
  -r, --regtest             Indicates if regtest parameters should be used
      --resultsdir string   Directory where results should be stored (default "./results")
  -s, --signet              Indicates if the public signet parameters should be used
  -t, --testnet             Indicates if testnet parameters should be used
```

### SEE ALSO

* [chantools](chantools.md)	 - Chantools helps recover funds from lightning channels

