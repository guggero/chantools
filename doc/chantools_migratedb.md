## chantools migratedb

Apply all recent lnd channel database migrations

### Synopsis

This command opens an lnd channel database in write mode
and applies all recent database migrations to it. This can be used to update
an old database file to be compatible with the current version that chantools
needs to read the database content.

CAUTION: Running this command will make it impossible to use the channel DB
with an older version of lnd. Downgrading is not possible and you'll need to
run lnd v0.19.0-beta or later after using this command!'

```
chantools migratedb [flags]
```

### Examples

```
chantools migratedb \
	--channeldb ~/.lnd/data/graph/mainnet/channel.db
```

### Options

```
      --channeldb string   lnd channel.db file to migrate
  -h, --help               help for migratedb
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

