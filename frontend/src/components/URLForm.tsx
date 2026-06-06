import React, { useState, useTransition } from 'react'
import strings from '../i18n'
import { APIError, LengthenRequest, URLRecord, lengthenURL } from '../api/urls'
import { streamSuggestions } from '../api/suggest'
import SlugPicker from './SlugPicker'

interface Props {
  onResult: (record: URLRecord) => void
}

export default function URLForm({ onResult }: Props) {
  const [url, setUrl] = useState('')
  const [customSlug, setCustomSlug] = useState('')
  const [suggestions, setSuggestions] = useState<string[]>([])
  const [suggesting, setSuggesting] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [isPending, startTransition] = useTransition()

  function handleSuggest() {
    if (!url) return
    setSuggestions([])
    setSuggesting(true)
    streamSuggestions(
      url,
      (slug) => setSuggestions((prev) => [...prev, slug]),
      () => setSuggesting(false),
    )
  }

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setError(null)
    const req: LengthenRequest = { url, custom_slug: customSlug || undefined }
    startTransition(() => {
      lengthenURL(req)
        .then(onResult)
        .catch((err) => {
          if (err instanceof APIError) {
            setError(err.code)
          } else {
            setError('internal_error')
          }
        })
    })
  }

  return (
    <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
      <div>
        <label style={{ display: 'block', marginBottom: '0.25rem' }}>{strings.form.url.label}</label>
        <input
          type="url"
          value={url}
          onChange={(e) => setUrl(e.target.value)}
          placeholder={strings.form.url.placeholder}
          required
          style={{ width: '100%', padding: '0.5rem', fontSize: '1rem', boxSizing: 'border-box' }}
        />
      </div>
      <div>
        <label style={{ display: 'block', marginBottom: '0.25rem' }}>{strings.form.customSlug.label}</label>
        <div style={{ display: 'flex', gap: '0.5rem' }}>
          <input
            type="text"
            value={customSlug}
            onChange={(e) => setCustomSlug(e.target.value)}
            placeholder={strings.form.customSlug.placeholder}
            style={{ flex: 1, padding: '0.5rem', fontSize: '1rem' }}
          />
          <button type="button" onClick={handleSuggest} disabled={!url || suggesting}>
            {strings.form.suggest}
          </button>
        </div>
      </div>
      {(suggestions.length > 0 || suggesting) && (
        <SlugPicker
          candidates={suggestions}
          loading={suggesting}
          onSelect={(slug) => setCustomSlug(slug)}
        />
      )}
      {error && (
        <p style={{ color: 'red', margin: 0 }}>
          {strings.error[error as keyof typeof strings.error] ?? strings.error.internal_error}
        </p>
      )}
      <button type="submit" disabled={isPending}>
        {strings.form.submit}
      </button>
    </form>
  )
}
