import type { Env } from './types'

export async function getOrCreateQR(slug: string, env: Env): Promise<string> {
  const existing = await env.QR_BUCKET.head(slug)
  if (existing) {
    return `https://qr.lengthener.example.com/${slug}`
  }

  const qrRes = await fetch(
    `https://api.qrserver.com/v1/create-qr-code/?size=256x256&data=${encodeURIComponent(slug)}`,
  )
  if (!qrRes.ok) throw new Error('qr fetch failed')

  await env.QR_BUCKET.put(slug, qrRes.body!, {
    httpMetadata: { contentType: 'image/png' },
  })
  return `https://qr.lengthener.example.com/${slug}`
}
