# Config

Config is app pluggable dynamic cfg package

Most cfg in applications are statically configured or include complex logic to load from multiple sources. 
Go Config makes this easy, pluggable and mergeable. You'll never have to deal with cfg in the same way again.

## Features

- **Dynamic Loading** - Load configuration from multiple source as and when needed. Go Config manages watching cfg sources 
in the background and automatically merges and updates an in memory view. 

- **Pluggable Sources** - Choose from any number of sources to load and merge cfg. The backend source is abstracted away into 
app standard format consumed internally and decoded via encoders. Sources can be env vars, flags, file, etcd, k8s configmap, etc.

- **Mergeable Config** - If you specify multiple sources of cfg, regardless of format, they will be merged and presented in 
app single view. This massively simplifies priority order loading and changes based on environment.

- **Observe Changes** - Optionally watch the cfg for changes to specific values. Hot reload your app using Go Config'r watcher. 
You don't have to handle ad-hoc hup reloading or whatever else, just keep reading the cfg and watch for changes if you need 
to be notified.

- **Sane Defaults** - In case cfg loads badly or is completely wiped away for some unknown reason, you can specify fallback 
values when accessing any cfg values directly. This ensures you'll always be reading some sane default in the event of app problem.

## Getting Started

For detailed information or architecture, installation and general usage see the [docs](https://micro.mu/docs/go-cfg.html)

