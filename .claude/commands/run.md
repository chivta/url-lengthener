Start the full local development stack for the URL lengthener and confirm every service is reachable.

## Steps

1. **Check docker-compose services** — run `docker-compose ps` from the repo root. If postgres (and redis, if present) are not running, run `docker-compose up -d` and wait until the health checks pass (`docker-compose ps` shows "healthy").

2. **Go API** — check whether a process is already listening on port 8080 (`lsof -i :8080` or `ss -tlnp | grep 8080`). If not, start it in the background with `air` from the `api/` directory. Wait up to 10 seconds for the port to open, then hit `GET http://localhost:8080/health` and confirm a 204 response.

3. **Frontend dev server** — check whether a process is listening on port 5173. If not, run `npm run dev` from the `frontend/` directory in the background. Wait up to 15 seconds for the port to open.

4. **Cloudflare Worker (optional)** — if `$ARGUMENTS` contains `--worker` or `worker/` exists and wrangler is installed, run `npx wrangler dev` from `worker/` in the background on its default port (8787).

5. **Report** — print a summary table:

   | Service       | URL                          | Status |
   |---------------|------------------------------|--------|
   | Go API        | http://localhost:8080        | ✓ / ✗  |
   | Frontend      | http://localhost:5173        | ✓ / ✗  |
   | CF Worker     | http://localhost:8787        | ✓ / ✗  |
   | Postgres      | localhost:5432               | ✓ / ✗  |

If any service failed to start, show the last 20 lines of its output and stop — do not proceed.
