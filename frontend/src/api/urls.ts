const base = (import.meta.env.VITE_API_BASE_URL as string) || ''

export interface URLRecord {
  id: string
  slug: string
  original_url: string
  user_id: string | null
  click_count: number
  expires_at: string | null
  created_at: string
}

export interface LengthenRequest {
  url: string
  custom_slug?: string
}

export async function lengthenURL(req: LengthenRequest): Promise<URLRecord> {
  const res = await fetch(`${base}/api/v1/urls`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(req),
  })
  if (!res.ok) {
    const body = await res.json().catch(() => ({}))
    throw new APIError(res.status, (body as { code?: string }).code ?? 'internal_error')
  }
  return res.json()
}

export async function getURL(slug: string): Promise<URLRecord> {
  const res = await fetch(`${base}/api/v1/urls/${slug}`)
  if (!res.ok) {
    const body = await res.json().catch(() => ({}))
    throw new APIError(res.status, (body as { code?: string }).code ?? 'internal_error')
  }
  return res.json()
}

export async function deleteURL(slug: string): Promise<void> {
  const res = await fetch(`${base}/api/v1/urls/${slug}`, { method: 'DELETE' })
  if (!res.ok) {
    const body = await res.json().catch(() => ({}))
    throw new APIError(res.status, (body as { code?: string }).code ?? 'internal_error')
  }
}

export class APIError extends Error {
  constructor(public status: number, public code: string) {
    super(code)
  }
}
