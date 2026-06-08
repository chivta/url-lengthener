Smoke-test the running URL lengthener stack end-to-end. Assumes `/run` has already been executed and all services are up.

If `$ARGUMENTS` names a specific area (e.g. `redirect`, `ai`, `analytics`, `frontend`), run only that section. Otherwise run all sections in order.

---

## 1. Health checks

```
GET http://localhost:8080/health        → 204
GET http://localhost:8080/metrics       → 200
```

Fail fast if either is unreachable.

---

## 2. Create a short URL

```
POST http://localhost:8080/api/v1/urls
Content-Type: application/json

{ "url": "https://example.com/some/long/path" }
```

Assert:
- Status 201
- Response body contains `slug` (non-empty string) and `short_url`
- Save the returned `slug` for subsequent steps

---

## 3. Retrieve the URL record

```
GET http://localhost:8080/api/v1/urls/{slug}
```

Assert:
- Status 200
- `original_url` matches the one submitted
- `click_count` is 0

---

## 4. Redirect

```
GET http://localhost:8080/{slug}   (follow redirects disabled)
```

Assert:
- Status 301 or 302
- `Location` header equals the original URL

---

## 5. Click count increments

Hit the redirect endpoint 3 more times (follow redirects disabled), then:

```
GET http://localhost:8080/api/v1/urls/{slug}
```

Assert `click_count` is 3 (or ≥ 3 if counts are eventually consistent).

---

## 6. Cloudflare Worker redirect (if running)

```
GET http://localhost:8787/{slug}   (follow redirects disabled)
```

Assert:
- Status 301 or 302
- `Location` header equals the original URL

---

## 7. Frontend smoke test

Use the Bash tool to fetch the frontend root and check it returns HTML:

```
curl -s -o /dev/null -w "%{http_code}" http://localhost:5173
```

Assert status 200.

---

## 8. Not found handling

```
GET http://localhost:8080/api/v1/urls/doesnotexist999
```

Assert:
- Status 404
- Response body contains a machine-readable `code` field (not a human message)

---

## Report

Print a results table — one row per section, PASS/FAIL, with a one-line note on any failure. If any section failed, exit with a non-zero status so hooks can detect it.
