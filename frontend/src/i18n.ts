const strings = {
  app: {
    title: 'URL Shortener',
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
    submit: 'Shorten',
    suggest: 'Suggest slugs',
  },
  result: {
    shortUrl: 'Your short URL',
    copy: 'Copy',
    copied: 'Copied!',
  },
  error: {
    invalid_url: 'Please enter a valid URL starting with http:// or https://',
    slug_conflict: 'That slug is already taken. Try another.',
    url_not_found: 'Short URL not found.',
    internal_error: 'Something went wrong. Please try again.',
  },
  redirect: {
    loading: 'Redirecting…',
    notFound: 'This short URL does not exist.',
  },
} as const

export default strings
