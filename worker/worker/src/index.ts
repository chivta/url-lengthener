import type { Env } from './types'
import { handleRedirect } from './redirect'
import { getOrCreateQR } from './qr'

export { ClickCounter } from './clickCounter'

export default {
  async fetch(request: Request, env: Env): Promise<Response> {
    const url = new URL(request.url)

    if (url.pathname === '/health') {
      return new Response('ok', { status: 200 })
    }

    if (request.method === 'GET' && url.pathname.startsWith('/qr/')) {
      const slug = url.pathname.slice(4)
      if (!slug) return new Response('not found', { status: 404 })
      try {
        const qrURL = await getOrCreateQR(slug, env)
        return Response.redirect(qrURL, 302)
      } catch {
        return new Response('qr generation failed', { status: 500 })
      }
    }

    if (request.method === 'GET' && url.pathname.length > 1) {
      const slug = url.pathname.slice(1)
      return handleRedirect(slug, env)
    }

    return new Response('not found', { status: 404 })
  },
}
