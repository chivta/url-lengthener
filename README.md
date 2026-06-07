# URL Lengthener

Fully vibe coded URL lengthener project built as a learning project covering Go, React, Cloudflare Workers, AWS/EKS, Terraform, and Claude AI.

## Stack

| Layer | Tech |
|-------|------|
| API | Go, Gin, pgx/v5, golang-migrate |
| Frontend | React, Vite, TypeScript |
| Edge | Cloudflare Worker, KV, Durable Objects |
| Infra | Terraform, EKS, RDS, ElastiCache, ECR |
| CI/CD | GitHub Actions |

## Features

- Lengthen URLs with an auto-generated or custom slug
- AI-powered slug suggestions via Claude (streamed as SSE)
- Edge redirects via Cloudflare Worker with KV caching
- Click counting via Cloudflare Durable Objects, flushed to Postgres every 60s

## TODO

- QR-code generation
- Stripe integration with additional features like:
  - Even longer urls
  - Scary-looking urls
  - Urls with lot's of query parameters to be even longer
  - Freaky qr-codes

## Quick start

```bash
docker-compose up -d
```

Starts api, frontend with hot reload in dev mode

- API: http://localhost:8080
- Frontend: http://localhost:5173
- Worker (local): `cd worker && npx wrangler dev`

## Environment variables

| Variable | Required | Default | Notes |
|---|---|---|---|
| `ANTHROPIC_API_KEY` | no | — | Claude slug suggestions, copy .env.example to .env and fill it if needed |
| `ALLOWED_ORIGINS` | yes | — | Comma-separated CORS origins |
| `DATABASE_URL` | auto | set by compose | Set in docker compose for dev |
| `REDIS_URL` | auto | set by compose | Set in docker compose for dev |
| `PORT` | no | `8080` | |
| `ENVIRONMENT` | no | `development` | `production` enables Gin release mode |

## GitHub Actions secrets

CI/CD workflows (`.github/workflows/`) require the following repository secrets to be configured (Settings → Secrets and variables → Actions):

| Secret | Used by | Notes |
|---|---|---|
| `AWS_ROLE_ARN` | CD | IAM role assumed via OIDC for ECR push and manifest commits — no static AWS credentials |
| `VITE_API_BASE_URL_PROD` | CD (push to `main`) | Baked into the frontend build deployed to the `prod` workspace |
| `VITE_API_BASE_URL_DEV` | CD (push to `dev`) | Baked into the frontend build deployed to the `dev` workspace |

## Monorepo layout

```
api/        Go service
frontend/   React + Vite app
worker/     Cloudflare Worker
infra/      Terraform (EKS, RDS, ElastiCache, ECR)
k8s/        Kubernetes manifests (GitOps target)
.github/    CI (test/lint) + CD (build & deploy)
```
