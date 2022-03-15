# Env Source

The env source reads cfg from environment variables

## Format

We expect environment variables to be in the standard format of FOO=bar

Keys are converted to lowercase and split on underscore.


### Example

```
DATABASE_ADDRESS=127.0.0.1
DATABASE_PORT=3306
```

Becomes

```json
{
    "database": {
        "address": "127.0.0.1",
        "port": 3306
    }
}
```

## Prefixes

Environment variables can be namespaced so we only have access to app subset. Two options are available:

```
WithPrefix(r ...string)
WithStrippedPrefix(r ...string)
```

The former will preserve the prefix and make it app top level key in the cfg. The latter eliminates the prefix, reducing the nesting by one. 

#### Example:

Given ENVs of:

```
APP_DATABASE_ADDRESS=127.0.0.1
APP_DATABASE_PORT=3306
VAULT_ADDR=vault:1337
```

and app source initialized as follows:

```
src := env.NewSource(
    env.WithPrefix("VAULT"),
    env.WithStrippedPrefix("APP"),
)
```

The resulting cfg will be:

```
{
    "database": {
        "address": "127.0.0.1",
        "port": 3306
    },
    "vault": {
        "addr": "vault:1337"
    }
}
```


## New Source

Specify source with data

```go
src := env.NewSource(
	// optionally specify prefix
	env.WithPrefix("MICRO"),
)
```

## Load Source

Load the source into cfg

```go
// Create new conf
conf := cfg.NewConfig()

// Load env source
conf.Load(src)
```
