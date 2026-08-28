# 🌱 seedling

A tiny link shortener — and the playground for
[**tiny**](https://github.com/tiny-systems/tiny), coding-agent sessions on
your own Kubernetes.

## How this repo works

1. [File an issue](../../issues/new) — a bug, a feature, a question.
2. Add the **`tiny`** label.
3. Watch: within seconds a 🌱 comment says a session picked it up. Code
   arrives as a pull request; answers arrive as comments. The agent runs
   in a Kubernetes cluster, not in CI — the five-second workflow jobs
   only carry messages.

## The app itself

```sh
go run .                                  # listens on :8080
curl -X POST -d url=https://example.com localhost:8080/shorten
curl -i localhost:8080/s1                 # 302 → example.com
```

In-memory store, no persistence, no dedup, no custom codes — on purpose.
Every missing feature is an issue waiting for a label.

MIT.
