# quote-engine

A small Go service for freight-style price quotes. Synthetic data only. This is **not** the Frete Gestão codebase.

It shows the patterns I used in production when I moved a quote engine from PHP to Go:

- formulas stored as configuration, evaluated at runtime
- in-memory cache with TTL
- singleflight so a cache miss is loaded once
- the HTTP handler stays thin; the domain stays testable

In production, average latency dropped from 450ms to 70ms at around 400,000 quotes per day. This repo is the public, simplified version of that idea.

## Run

```bash
go test ./...
go run ./cmd/api
```

```bash
curl -s localhost:8080/v1/quotes \
  -H 'content-type: application/json' \
  -d '{"origin":"SAO","destination":"CWB","weight_kg":10}'
```

Expected amount for that lane: `12.5 + 10 * 2.4 = 36.5`.

## What is in here

| Piece | Why |
|---|---|
| `data/tables.json` | Rate tables as config, not `if origin ==` branches |
| `internal/quote` | Deterministic calculation, covered by tests |
| `singleflight` | Avoid a stampede when the cache expires |
| `POST /v1/quotes` | One endpoint, easy to read |

What is **not** in here: company rules, CTe documents, customer data, or the production service.

## Stack

Go · expr-lang · stdlib `net/http`

## Author

[Marcel Peruch](https://github.com/mperuch) · [LinkedIn](https://linkedin.com/in/marcel-peruch)
