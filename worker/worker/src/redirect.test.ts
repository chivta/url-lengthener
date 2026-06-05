import { describe, it, expect, vi } from 'vitest'
import { handleRedirect } from './redirect'
import type { Env } from './types'

function makeEnv(overrides: Partial<Env> = {}): Env {
  const kv = new Map<string, string>()
  return {
    URL_KV: {
      get: vi.fn(async (key: string) => kv.get(key) ?? null),
      put: vi.fn(async (key: string, value: string) => { kv.set(key, value) }),
      delete: vi.fn(),
      list: vi.fn(),
      getWithMetadata: vi.fn(),
    } as unknown as KVNamespace,
    CLICK_COUNTER: {
      idFromName: vi.fn(() => ({ toString: () => 'id' })),
      get: vi.fn(() => ({ fetch: vi.fn().mockResolvedValue(new Response('{}')) })),
      newUniqueId: vi.fn(),
      idFromString: vi.fn(),
      jurisdiction: vi.fn(),
    } as unknown as DurableObjectNamespace,
    QR_BUCKET: {} as R2Bucket,
    API_BASE_URL: 'http://localhost:8080',
    ...overrides,
  }
}

describe('handleRedirect', () => {
  it('returns 301 with Location from KV cache on hit', async () => {
    const env = makeEnv()
    await env.URL_KV.put('abc123', 'https://example.com')
    ;(env.URL_KV.get as ReturnType<typeof vi.fn>).mockResolvedValueOnce('https://example.com')

    const res = await handleRedirect('abc123', env)
    expect(res.status).toBe(301)
    expect(res.headers.get('Location')).toBe('https://example.com')
  })

  it('returns 404 when KV miss and API returns 404', async () => {
    const env = makeEnv()
    ;(env.URL_KV.get as ReturnType<typeof vi.fn>).mockResolvedValueOnce(null)
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(new Response('not found', { status: 404 })))

    const res = await handleRedirect('missing', env)
    expect(res.status).toBe(404)

    vi.unstubAllGlobals()
  })

  it('returns 301 with Location from API on KV miss', async () => {
    const env = makeEnv()
    ;(env.URL_KV.get as ReturnType<typeof vi.fn>).mockResolvedValueOnce(null)
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(
      new Response(JSON.stringify({ original_url: 'https://example.com' }), {
        status: 200,
        headers: { 'content-type': 'application/json' },
      })
    ))

    const res = await handleRedirect('abc123', env)
    expect(res.status).toBe(301)
    expect(res.headers.get('Location')).toBe('https://example.com')

    vi.unstubAllGlobals()
  })
})
