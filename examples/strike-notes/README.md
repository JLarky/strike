Based on

- [the original](https://github.com/jazzypants1989/rsc-from-scratch/tree/main/server-components-demo)
- [and the fork](https://github.com/jazzypants1989/rsc-from-scratch/tree/main/server-components-demo)

## How to run

```bash
cd examples/strike-notes
go run .
```

Or

```bash
air
```

Build jsx -> js (optional)

```bash
(cd examples/strike-notes/public/; NODE_ENV=production bun --watch build.ts)
```

## Deploy

- Install Fly

```bash
flyctl launch
flyctl deploy
```

## Profiling

```
plow http://127.0.0.1:8080/ -c 500 -n 200000 -d 10s
go tool pprof -http :8888 http://:8080/debug/pprof/profile?seconds=5
```
