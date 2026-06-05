import { useState } from 'react'
import strings from '../i18n'
import URLForm from '../components/URLForm'
import ResultCard from '../components/ResultCard'
import { URLRecord } from '../api/urls'

export default function ShortenPage() {
  const [result, setResult] = useState<URLRecord | null>(null)

  return (
    <main style={{ maxWidth: 600, margin: '0 auto', padding: '2rem' }}>
      <h1>{strings.app.title}</h1>
      <URLForm onResult={setResult} />
      {result && <div style={{ marginTop: '1.5rem' }}><ResultCard record={result} /></div>}
    </main>
  )
}
