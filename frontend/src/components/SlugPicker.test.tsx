import { describe, it, expect, vi } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import SlugPicker from './SlugPicker'

describe('SlugPicker', () => {
  it('renders candidates as buttons', () => {
    render(<SlugPicker candidates={['abc123', 'def456']} loading={false} onSelect={vi.fn()} />)
    expect(screen.getByText('abc123')).toBeInTheDocument()
    expect(screen.getByText('def456')).toBeInTheDocument()
  })

  it('calls onSelect with slug when button clicked', () => {
    const onSelect = vi.fn()
    render(<SlugPicker candidates={['abc123']} loading={false} onSelect={onSelect} />)
    fireEvent.click(screen.getByText('abc123'))
    expect(onSelect).toHaveBeenCalledWith('abc123')
  })

  it('shows loading indicator when loading=true', () => {
    render(<SlugPicker candidates={[]} loading={true} onSelect={vi.fn()} />)
    expect(screen.getByText('…')).toBeInTheDocument()
  })
})
