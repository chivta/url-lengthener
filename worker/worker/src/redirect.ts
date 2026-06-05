import type { Env } from './types'

export async function handleRedirect(slug: string, env: Env): Promise<Response> {
  const cached = await env.URL_KV.get(slug, { cacheTtl: 300 })
  if (cached) {
    incrementCounter(slug, env)
    return Response.redirect(cached, 301)
  }

  const apiRes = await fetch(`${env.API_BASE_URL}/api/v1/urls/${slug}`)
  if (!apiRes.ok) {
    return new Response('not found', { status: 404, headers: { 'content-type': 'text/plain' } })
  }

  const data = await apiRes.json<{ original_url: string }>()
  await env.URL_KV.put(slug, data.original_url, { expirationTtl: 3600 })
  incrementCounter(slug, env)
  return Response.redirect(data.original_url, 301)
}

function incrementCounter(slug: string, env: Env): void {
  const id = env.CLICK_COUNTER.idFromName(slug)
  const stub = env.CLICK_COUNTER.get(id)
  void stub.fetch(`https://internal/increment?slug=${encodeURIComponent(slug)}`, { method: 'POST' })
}
