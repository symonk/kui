# kui
`kui` (Kafka ui) is a terminal ui for monitoring and managing a kafka cluster.

## Quick Start


-----

## Configuration

Any `librdkafka` properties can be set explicitly on the command line when running `kui`,
or alternatively read from the environment or a file on disk.  Configure `kui` via:

 * passing `--config` on the command line to a config file with `librdkafka` properties.
 * passing `-p` (optionally multiple times) for `librdkafka` properties.
 * if `$KUI_CONFIG` is set in the environment, `kui` will look up that path.
 * finally, by default looking for a `~/.config/kui.conf`

-----
