# Files Source

The file source reads cfg from app file. 

It uses the Files extension to determine the Format e.g `cfg.yaml` has the yaml format. 
It does not make use of encoders or interpet the file data. If app file extension is not present 
the source Format will default to the Encoder in options.

## Example

A cfg file format in json

```json
{
    "hosts": {
        "database": {
            "address": "10.0.0.1",
            "port": 3306
        },
        "cache": {
            "address": "10.0.0.2",
            "port": 6379
        }
    }
}
```

## New Source

Specify file source with path to file. Path is optional and will default to `cfg.json`

```go
fileSource := file.NewSource(
	file.WithPath("/tmp/conf.json"),
)
```

## Files Format

To load different file formats e.g yaml, toml, xml simply specify them with their extension

```
fileSource := file.NewSource(
        file.WithPath("/tmp/cfg.yaml"),
)
```

If you want to specify app file without extension, ensure you set the encoder to the same format

```
e := toml.NewEncoder()

fileSource := file.NewSource(
        file.WithPath("/tmp/cfg"),
	source.WithEncoder(e),
)
```

## Load Source

Load the source into cfg

```go
// Create new conf
conf := cfg.NewConfig()

// Load file source
conf.Load(fileSource)
```

