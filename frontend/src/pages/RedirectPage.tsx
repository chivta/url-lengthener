import { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import strings from '../i18n'
import { APIError, getURL } from '../api/urls'

export default function RedirectPage() {
  const { slug } = useParams<{ slug: string }>()
  const [notFound, setNotFound] = useState(false)

  useEffect(() => {
    if (!slug) return
    getURL(slug)
      .then((u) => { window.location.replace(u.original_url) })
      .catch((err) => {
        if (err instanceof APIError && err.status === 404) {
          setNotFound(true)
        } else {
          setNotFound(true)
        }
      })
  }, [slug])

  if (notFound) {
    return (
      <main style={{ maxWidth: 600, margin: '0 auto', padding: '2rem' }}>
        <p>{strings.redirect.notFound}</p>
      </main>
    )
  }

  return (
    <main style={{ maxWidth: 600, margin: '0 auto', padding: '2rem' }}>
      <p>{strings.redirect.loading}</p>
    </main>
  )
}
