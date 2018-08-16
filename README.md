# stdout-heartbeat
Wraps your existing command with interleaved STDOUT heartbeats. Can be used as an alternative to [`travis_wait`](https://docs.travis-ci.com/user/common-build-problems/#build-times-out-because-no-output-was-received).

## Installation

Builds for linux_amd64 and darwin_amd64 are available in the [Releases](https://github.com/goodeggs/stdout-heartbeat/releases).

## Usage

```
stdout-heartbeat <interval> <command> [<arg>,...]
```

- `interval` can be any string compatible with [Go's `time.ParseDuration`](https://golang.org/pkg/time/#ParseDuration), but isn't reliable with values less than a few seconds.
- Your command's exit code will be passed through

## Example

```
$ stdout-heartbeat 10s sh -c 'echo hello; sleep 15; echo world; sleep 5; echo goodbye world' | ts '[%H:%M:%S]'
[10:11:01] hello
[10:11:11] ♥
[10:11:16] world
[10:11:21] goodbye world
```

