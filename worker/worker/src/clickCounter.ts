import type { Env } from './types'

export class ClickCounter implements DurableObject {
  private slug = ''

  constructor(private state: DurableObjectState, private env: Env) {}

  async fetch(request: Request): Promise<Response> {
    const url = new URL(request.url)
    if (url.pathname === '/increment') {
      this.slug = url.searchParams.get('slug') ?? this.slug
      const count = ((await this.state.storage.get<number>('count')) ?? 0) + 1
      await this.state.storage.put('count', count)
      const alarmTime = await this.state.storage.getAlarm()
      if (alarmTime === null) {
        await this.state.storage.setAlarm(Date.now() + 60_000)
      }
      return new Response(JSON.stringify({ count }), {
        headers: { 'content-type': 'application/json' },
      })
    }
    return new Response('not found', { status: 404 })
  }

  async alarm(): Promise<void> {
    const count = await this.state.storage.get<number>('count')
    if (!count || !this.slug) return
    await this.state.storage.delete('count')
    await fetch(`${this.env.API_BASE_URL}/internal/clicks/flush`, {
      method: 'POST',
      headers: { 'content-type': 'application/json' },
      body: JSON.stringify({ slug: this.slug, count }),
    }).catch(() => {})
  }
}
