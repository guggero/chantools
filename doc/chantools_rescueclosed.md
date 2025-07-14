## chantools rescueclosed

Try finding the private keys for funds that are in outputs of remotely force-closed channels

### Synopsis

If channels have already been force-closed by the remote
peer, this command tries to find the private keys to sweep the funds from the
output that belongs to our side. This can only be used if we have a channel DB
that contains the latest commit point. Normally you would use SCB to get the
funds from those channels. But this method can help if the other node doesn't
know about the channels any more but we still have the channel.db from the
moment they force-closed.

NOTE: Unless your channel was opened before 2019, you very likely don't need to
use this command as things were simplified. Use 'chantools sweepremoteclosed'
instead if the remote party has already closed the channel.

The alternative use case for this command is if you got the commit point by
running the fund-recovery branch of my guggero/lnd fork (see 
https://github.com/guggero/lnd/releases for a binary release) in combination
with the fakechanbackup command. Then you need to specify the --commit_point and 
--force_close_addr flags instead of the --channeldb and --fromsummary flags.

If you need to rescue a whole bunch of channels all at once, you can also
specify the --fromsummary and --lnd_log flags to automatically look for force
close addresses in the summary and the corresponding commit points in the
lnd log file. This only works if lnd is running the fund-recovery branch of my
guggero/lnd (https://github.com/guggero/lnd/releases) fork and only if the
debuglevel is set to debug (lnd.conf, set 'debuglevel=debug').

```
chantools rescueclosed [flags]
```

### Examples

```
chantools rescueclosed \
	--fromsummary results/summary-xxxxxx.json \
	--channeldb ~/.lnd/data/graph/mainnet/channel.db

chantools rescueclosed --force_close_addr bc1q... --commit_point 03xxxx

chantools rescueclosed --fromsummary results/summary-xxxxxx.json \
	--lnd_log ~/.lnd/logs/bitcoin/mainnet/lnd.log
```

### Options

```
      --bip39                     read a classic BIP39 seed and passphrase from the terminal instead of asking for lnd seed format or providing the --rootkey flag
      --channeldb string          lnd channel.db file to use for rescuing force-closed channels
      --commit_point string       the commit point that was obtained from the logs after running the fund-recovery branch of guggero/lnd
      --force_close_addr string   the address the channel was force closed to, look up in block explorer by following funding txid
      --fromchanneldb string      channel input is in the format of an lnd channel.db file
      --fromchanneldump string    channel input is in the format of a channel dump file
      --fromsummary string        channel input is in the format of chantool's channel summary; specify '-' to read from stdin
  -h, --help                      help for rescueclosed
      --listchannels string       channel input is in the format of lncli's listchannels format; specify '-' to read from stdin
      --lnd_log string            the lnd log file to read to get the commit_point values when rescuing multiple channels at the same time
      --num_keys uint32           the number of keys to derive for the brute force attack (default 5000)
      --pendingchannels string    channel input is in the format of lncli's pendingchannels format; specify '-' to read from stdin
      --rootkey string            BIP32 HD root key of the wallet to use for decrypting the backup; leave empty to prompt for lnd 24 word aezeed
      --walletdb string           read the seed/master root key to use for decrypting the backup from an lnd wallet.db file instead of asking for a seed or providing the --rootkey flag
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

