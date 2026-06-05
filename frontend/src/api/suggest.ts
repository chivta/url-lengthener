const base = (import.meta.env.VITE_API_BASE_URL as string) || ''

export function streamSuggestions(
  originalURL: string,
  onCandidate: (slug: string) => void,
  onDone: () => void,
): () => void {
  const controller = new AbortController()

  fetch(`${base}/api/v1/urls/suggest`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ url: originalURL }),
    signal: controller.signal,
  })
    .then(async (res) => {
      if (!res.body) return
      const reader = res.body.getReader()
      const decoder = new TextDecoder()
      let buf = ''
      while (true) {
        const { done, value } = await reader.read()
        if (done) break
        buf += decoder.decode(value, { stream: true })
        const lines = buf.split('\n')
        buf = lines.pop() ?? ''
        for (const line of lines) {
          if (line.startsWith('event: done')) {
            onDone()
            return
          }
          if (line.startsWith('data: ')) {
            const slug = line.slice(6).trim()
            if (slug) onCandidate(slug)
          }
        }
      }
      onDone()
    })
    .catch(() => onDone())

  return () => controller.abort()
}
