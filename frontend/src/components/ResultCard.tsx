import { useState } from 'react'
import strings from '../i18n'
import { URLRecord } from '../api/urls'

interface Props {
  record: URLRecord
  baseURL?: string
}

export default function ResultCard({ record, baseURL = window.location.origin }: Props) {
  const [copied, setCopied] = useState(false)
  const longURL = `${baseURL}/${record.slug}`

  function handleCopy() {
    navigator.clipboard.writeText(longURL).then(() => {
      setCopied(true)
      setTimeout(() => setCopied(false), 2000)
    })
  }

  return (
    <div style={{ padding: '1rem', border: '1px solid #ccc', borderRadius: '4px' }}>
      <p style={{ margin: '0 0 0.5rem', fontWeight: 'bold' }}>{strings.result.longUrl}</p>
      <div style={{ display: 'flex', gap: '0.5rem', alignItems: 'center' }}>
        <a href={longURL} target="_blank" rel="noopener noreferrer" style={{ flex: 1, wordBreak: 'break-all' }}>
          {longURL}
        </a>
        <button onClick={handleCopy} type="button">
          {copied ? strings.result.copied : strings.result.copy}
        </button>
      </div>
    </div>
  )
}
