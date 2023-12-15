## The issue
The connector sends requests in the wrong order when `Concurrency` != 1.

## How to run the reproducer

1. `tarantool init.lua`
2. `make build`
3. `./build/gorepro -t 1`


### Arguments
```
> gorepro -h
Usage of gorepro:
  -c int
        tarantool connector concurrency
  -t int
        GOMAXPROCS
```
