# Flag Sender

A scheduler that sends flags to jury.

## Plugins

The flag sender allows you to use a plugin (protocol) to send flags in the desired format.

The example for the forcad_http plugin is located [here](../../workers/flag_sender/plugins/forcad_http/client.go)

## Requirements for writing a plugin

- For every plugin should be created own subdirectory such as `workers/flag_sender/plugins/<plugin_name>`
- Every plugin subdirectory should contain file `client.go` with all code in one file
- Plugin should contain variable `NewClient`, which type is `NewClientFunc` from `workers/flag_sender/pkg/plugins/interfaces.go`
- Function should return plugins.IClient, which should implement method SendFlags

During the build of the flag_sender container, plugins are built using the following command:

```docker
RUN find ./plugins -type f -name "client.go" -exec sh -c \
    'for f; do \
        dir=$(dirname "$f"); \
        name=$(basename "$dir"); \
        CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -buildmode=plugin -o "/build/plugin_binaries/${name}.so" "$f"; \
    done' sh {} +
```

The example for the forcad_http plugin is located [here](../../workers/flag_sender/plugins/forcad_http/client.go)

## Configuration

- `FLAG_SENDER_PLUGIN` - type of plugin (protocol) which will be used for sending flags
- `FLAG_SENDER_SUBMIT_PERIOD` - duration of flag sending
- `FLAG_SENDER_MAX_BATCH_SIZE` - maximum number of flags that can be sent at once
- `FLAG_SENDER_JURY_FLAG_URL_OR_HOST` - jury url or host. For http should be like `http://1.2.3.4:5678/api/send/flag/endpoint`. For tcp/grpc should be `1.2.3.4:5678`
- `FLAG_SENDER_SUBMIT_TIMEOUT` - timeout at which the flag sender stops sending flags
- `FLAG_SENDER_SUBMIT_LIMIT` - maximum number of flags that can be sent in a single request
- `FLAG_SENDER_FLAG_TTL` - lifetime of a flag, after which it is marked as OLD and is no longer sent to the verification system


The *Flag Sender* sends flags every *FLAG_SENDER_SUBMIT_PERIOD*, up to *FLAG_SENDER_MAX_BATCH_SIZE* flags at a time to jury by *FLAG_SENDER_PLUGIN* protocol.