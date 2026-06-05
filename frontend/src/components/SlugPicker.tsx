import strings from '../i18n'

interface Props {
  candidates: string[]
  loading: boolean
  onSelect: (slug: string) => void
}

export default function SlugPicker({ candidates, loading, onSelect }: Props) {
  return (
    <div style={{ display: 'flex', flexWrap: 'wrap', gap: '0.5rem', alignItems: 'center' }}>
      {candidates.map((slug) => (
        <button
          key={slug}
          type="button"
          onClick={() => onSelect(slug)}
          style={{ padding: '0.25rem 0.75rem', cursor: 'pointer' }}
        >
          {slug}
        </button>
      ))}
      {loading && <span style={{ color: '#888' }}>…</span>}
      {!loading && candidates.length === 0 && (
        <span style={{ color: '#888' }}>{strings.form.suggest}</span>
      )}
    </div>
  )
}
