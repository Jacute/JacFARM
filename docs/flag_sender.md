# Flag Sender

A scheduler that sends flags to jury.

## Plugins

The flag sender allows you to use a plugin (protocol) to send flags in the desired format.

TODO: write example

## Configuration

- `FLAG_SENDER_PLUGIN` - type of plugin (protocol) which will be used for sending flags
- `FLAG_SENDER_SEND_DURATION` - duration of flag sending
- `FLAG_SENDER_MAX_BATCH_SIZE` - the maximum number of flags that can be sent at once

The *Flag Sender* sends flags every *FLAG_SENDER_SEND_DURATION* seconds, up to *FLAG_SENDER_MAX_BATCH_SIZE* flags at a time to jury by *FLAG_SENDER_PLUGIN* protocol.