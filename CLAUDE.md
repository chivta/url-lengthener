# URL lengthener

Learning project covering Go, React+Vite, Cloudflare Workers, AWS/EKS, Terraform, and Claude AI.

## Monorepo layout

```
api/          Go service (gin, pgx/v5, golang-migrate)
frontend/     React + Vite + TypeScript
worker/       Cloudflare Worker (TypeScript, KV, Durable Objects, R2)
infra/        Terraform (EKS, RDS, ElastiCache, ECR, S3)
manifests/    Kubernetes YAML — GitOps target, updated by CD pipeline
.github/      CI (test/lint) and CD (build images, update manifests) workflows
```

## Quick start

```bash
cp .env.example .env          # fill in real values
make dev                      # starts postgres + redis + api (air) + frontend (vite)
```

API: http://localhost:8080
Frontend: http://localhost:5173
Worker (local): `cd worker && npx wrangler dev`

## Key architectural decisions

**Slug format** — up to 8192 bytes (8KB, the "lengthener" bit), `[0-9A-Za-z\-_.~]`, `crypto/rand`. Custom slug if provided, otherwise randomly generated. Auto-retry once on conflict when no custom slug given. Identity/lookup is keyed on `sha256(slug)` (stored as `slug_hash`, unique-indexed) rather than the raw string — Postgres btree index entries cap at ~2.7KB and the planned Cloudflare KV edge cache caps keys at 512 bytes, so a fixed-size digest is the only thing that can be indexed/keyed regardless of slug length.

**Click counting** — Cloudflare Durable Object is the hot counter (per-slug). Alarm flushes to `POST /internal/clicks/flush` on the Go API every 60s. Go API increments postgres `click_count` asynchronously (fire-and-forget goroutine) on every redirect.

**Redirect path** — Cloudflare Worker handles `GET /:slug` at the edge: KV lookup (300s cacheTtl) → fallback to Go API → cache result for 1h. The Go API also has a `GET /:slug` fallback for non-Cloudflare environments.

**Migrations** — `golang-migrate` with `embed.FS`. Run automatically on `server.Run()` before the HTTP listener starts.

**Interfaces at consumer** — `URLRepository`, `ClickRepository`, `URLService` interfaces are defined in `domain/interfaces.go` (consumed by the service/handler layers), not in the packages that implement them.

**Management vs redirect** — `GET /api/v1/urls/:slug` calls `svc.Get()` (no click tracking). `GET /:slug` calls `svc.Resolve()` (fires async click record + increment).

## Environment variables

| Variable | Required | Default | Notes |
|---|---|---|---|
| `DATABASE_URL` | yes | — | postgres DSN |
| `REDIS_URL` | yes | — | `redis://host:port` |
| `ALLOWED_ORIGINS` | yes | — | comma-separated CORS origins |
| `PORT` | no | `8080` | |
| `ENVIRONMENT` | no | `development` | `production` sets gin release mode |
| `LOG_LEVEL` | no | `info` | |

## Go service

```bash
cd api
go test ./...                    # unit tests
go test -tags integration ./...  # + integration tests (needs postgres)
golangci-lint run                # lint
go build ./cmd/server/main.go    # build binary
```

## Frontend

```bash
cd frontend
npm run dev          # Vite dev server (proxies /api to localhost:8080)
npm run typecheck    # tsc --noEmit
npm run lint         # eslint
npm test -- --run    # vitest (non-watch)
```

## Infrastructure

```bash
cd infra
terraform init -backend-config=backend.hcl
terraform plan -var-file=prod.tfvars
terraform apply -var-file=prod.tfvars
```

## CI/CD

- **CI** (`.github/workflows/ci.yml`): runs on every push. Three parallel jobs: `test-api`, `test-frontend`, `lint`.
- **CD** (`.github/workflows/cd.yml`): triggers after CI passes on `main`. Builds and pushes images to ECR (tagged with commit SHA), then commits updated `IMAGE_TAG` in manifests. Commit prefix `deploy:` distinguishes automated commits.

Required GitHub secrets: `AWS_ROLE_ARN`, `VITE_API_BASE_URL_PROD`.
