const strings = {
  app: {
    title: 'URL lengthener',
  },
  form: {
    url: {
      label: 'Long URL',
      placeholder: 'https://example.com/very/long/path',
    },
    customSlug: {
      label: 'Custom slug (optional)',
      placeholder: 'my-slug',
    },
    submit: 'Lengthen',
  },
  result: {
    longUrl: 'Your long URL',
    copy: 'Copy',
    copied: 'Copied!',
  },
  error: {
    invalid_url: 'Please enter a valid URL starting with http:// or https://',
    slug_conflict: 'That slug is already taken. Try another.',
    url_not_found: 'Long URL not found.',
    internal_error: 'Something went wrong. Please try again.',
  },
  redirect: {
    loading: 'Redirecting…',
    notFound: 'This long URL does not exist.',
  },
} as const

export default strings
