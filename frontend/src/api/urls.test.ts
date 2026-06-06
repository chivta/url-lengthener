import { describe, it, expect, vi, afterEach } from 'vitest'
import { lengthenURL, getURL, APIError } from './urls'

function mockFetch(status: number, body: unknown) {
  return vi.fn().mockResolvedValue({
    ok: status >= 200 && status < 300,
    status,
    json: () => Promise.resolve(body),
  })
}

describe('lengthenURL', () => {
  afterEach(() => vi.unstubAllGlobals())

  it('returns URLRecord on 201', async () => {
    const record = { id: '1', slug: 'abc', original_url: 'https://example.com', user_id: null, click_count: 0, expires_at: null, created_at: '' }
    vi.stubGlobal('fetch', mockFetch(201, record))
    const result = await lengthenURL({ url: 'https://example.com' })
    expect(result.slug).toBe('abc')
  })

  it('throws APIError on 409', async () => {
    vi.stubGlobal('fetch', mockFetch(409, { code: 'slug_conflict' }))
    await expect(lengthenURL({ url: 'https://example.com', custom_slug: 'taken' }))
      .rejects.toBeInstanceOf(APIError)
  })

  it('throws APIError with correct code', async () => {
    vi.stubGlobal('fetch', mockFetch(422, { code: 'invalid_url' }))
    try {
      await lengthenURL({ url: 'not-a-url' })
    } catch (e) {
      expect(e).toBeInstanceOf(APIError)
      expect((e as APIError).code).toBe('invalid_url')
    }
  })
})

describe('getURL', () => {
  afterEach(() => vi.unstubAllGlobals())

  it('returns URLRecord on 200', async () => {
    const record = { id: '1', slug: 'abc', original_url: 'https://example.com', user_id: null, click_count: 3, expires_at: null, created_at: '' }
    vi.stubGlobal('fetch', mockFetch(200, record))
    const result = await getURL('abc')
    expect(result.click_count).toBe(3)
  })

  it('throws APIError on 404', async () => {
    vi.stubGlobal('fetch', mockFetch(404, { code: 'url_not_found' }))
    await expect(getURL('missing')).rejects.toBeInstanceOf(APIError)
  })
})
