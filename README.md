# Flake Id Generator

FlakeId is distributed unique ID generator inspired by [Twitter's Snowflake](https://blog.twitter.com/engineering/en_us/a/2010/announcing-snowflake)

A FlakeId composed of

```plaintext
| flag  | timeoffset in millisec | machine id | sequence num |
| ----- | ---------------------- | ---------- | ------------ |
| 1 bit | 48 bits                | 8 bits     | 7 bits       |
```

## Installation

```bash
go get github.com/marverlous811/go-flake-id
```

## Usage

The function `NewIdFlakeGenerator` creates a new Flake Id generator instance

``` go
func NewIdFlakeGenerator(st IdFlakeGeneratorSetting) *IdFlakeGenerator
```

You can configure Flake by the struct settings

```go
type IdFlakeGeneratorSetting struct {
    StartTime       time.Time
    MachineIdGetter func() uint8
}
```

- `StartTime` is the time since which the Flake time is defined as elapsed time. If `StartTime` is 0, the start time of the flake generator is set to `2022/06/23T00:00:00Z`
- `MachineIdGetter` will return the unique ID of the Flake Generator instance

In order to get a new unique ID, you just have to call the method NextID

```go
func (gen *IdFlakeGenerator) NextId() (uint64, error)
```

> NextId can continue to generate IDs for about many thousands years with 48 bits allocated. But after that, the function with return errors
