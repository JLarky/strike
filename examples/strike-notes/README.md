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
while sleep 1; do (cd examples/strike-notes/public/ ;NODE_ENV=production bun build.ts); done
```
