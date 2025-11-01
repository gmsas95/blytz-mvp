import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { axe } from 'jest-axe';
import React from 'react';

import { BidButton } from '../bid-button';

// Mock the utils module
jest.mock('@/lib/utils', () => ({
  formatPrice: (price: number) => `$${price.toFixed(2)}`,
  cn: (...classes: string[]) => classes.filter(Boolean).join(' '),
}));

// Mock UI components
jest.mock('@/components/ui/button', () => ({
  Button: ({ children, onClick, disabled, ...props }: any) => (
    <button onClick={onClick} disabled={disabled} {...props}>
      {children}
    </button>
  ),
}));

jest.mock('@/components/ui/input', () => ({
  Input: ({ onChange, value, disabled, ...props }: any) => (
    <input onChange={onChange} value={value} disabled={disabled} {...props} />
  ),
}));

jest.mock('@/components/ui/card', () => ({
  Card: ({ children, className }: any) => <div className={className}>{children}</div>,
  CardContent: ({ children }: any) => <div>{children}</div>,
}));

jest.mock('@/components/ui/badge', () => ({
  Badge: ({ children }: any) => <span>{children}</span>,
}));

jest.mock('@/components/ui/alert', () => ({
  Alert: ({ children, className }: any) => <div className={className}>{children}</div>,
  AlertDescription: ({ children }: any) => <div>{children}</div>,
}));

describe('BidButton', () => {
  const defaultProps = {
    currentBid: 100,
    minBidIncrement: 10,
    onPlaceBid: jest.fn(),
    auctionId: 'test-auction',
    currentUserId: 'user-123',
  };

  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('renders correctly with default props', () => {
    render(<BidButton {...defaultProps} />);

    expect(screen.getByRole('button', { name: /place bid \$110\.00/i })).toBeInTheDocument();
  });

  it('displays the correct suggested bid amount', () => {
    render(<BidButton {...defaultProps} />);

    const button = screen.getByRole('button', { name: /place bid/i });
    expect(button).toHaveTextContent('$110.00');
  });

  it('handles bid placement successfully', async () => {
    const mockOnPlaceBid = jest
      .fn()
      .mockResolvedValue({ success: true, data: { bidId: 'bid-123' } });
    const user = userEvent.setup();

    render(<BidButton {...defaultProps} onPlaceBid={mockOnPlaceBid} />);

    const button = screen.getByRole('button', { name: /place bid/i });
    await user.click(button);

    expect(mockOnPlaceBid).toHaveBeenCalledWith(110);
  });

  it('shows loading state during bid placement', async () => {
    const mockOnPlaceBid = jest.fn(
      () => new Promise((resolve) => setTimeout(() => resolve({ success: true }), 100))
    );
    const user = userEvent.setup();

    render(<BidButton {...defaultProps} onPlaceBid={mockOnPlaceBid} />);

    const button = screen.getByRole('button', { name: /place bid/i });
    await user.click(button);

    expect(screen.getByRole('button')).toBeDisabled();
    expect(screen.getByText(/placing/i)).toBeInTheDocument();
  });

  it('displays success message after successful bid', async () => {
    const mockOnPlaceBid = jest.fn().mockResolvedValue({ success: true });
    const user = userEvent.setup();

    render(<BidButton {...defaultProps} onPlaceBid={mockOnPlaceBid} />);

    const button = screen.getByRole('button', { name: /place bid/i });
    await user.click(button);

    await waitFor(() => {
      expect(screen.getByText(/bid placed successfully/i)).toBeInTheDocument();
    });
  });

  it('displays error message when bid fails', async () => {
    const mockOnPlaceBid = jest
      .fn()
      .mockResolvedValue({ success: false, error: 'Insufficient funds' });
    const user = userEvent.setup();

    render(<BidButton {...defaultProps} onPlaceBid={mockOnPlaceBid} />);

    const button = screen.getByRole('button', { name: /place bid/i });
    await user.click(button);

    await waitFor(() => {
      expect(screen.getByText(/insufficient funds/i)).toBeInTheDocument();
    });
  });

  it('is disabled when disabled prop is true', () => {
    render(<BidButton {...defaultProps} disabled />);

    const button = screen.getByRole('button', { name: /place bid/i });
    expect(button).toBeDisabled();
  });

  it('shows disabled reason when provided', () => {
    render(<BidButton {...defaultProps} disabled disabledReason="Auction has ended" />);

    expect(screen.getByText(/auction has ended/i)).toBeInTheDocument();
  });

  it('renders quick bid buttons when showQuickBids is true', () => {
    render(<BidButton {...defaultProps} showQuickBids />);

    // Should show main bid button plus quick bid buttons
    expect(screen.getAllByRole('button').length).toBeGreaterThan(1);
  });

  it('renders compact variant correctly', () => {
    render(<BidButton {...defaultProps} variant="compact" />);

    expect(screen.getByRole('button', { name: /bid \$110\.00/i })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /custom/i })).toBeInTheDocument();
  });

  it('renders card variant correctly', () => {
    render(<BidButton {...defaultProps} variant="card" />);

    expect(screen.getByText(/\$100\.00/i)).toBeInTheDocument();
    expect(screen.getByText(/\+\$10\.00/i)).toBeInTheDocument();
  });

  it('handles custom bid input', async () => {
    const mockOnPlaceBid = jest.fn().mockResolvedValue({ success: true });
    const user = userEvent.setup();

    render(<BidButton {...defaultProps} variant="card" />);

    const customButton = screen.getByRole('button', { name: /custom/i });
    await user.click(customButton);

    const input = screen.getByRole('spinbutton');
    await user.type(input, '150');

    const placeBidButton = screen.getByRole('button', { name: /place bid/i });
    await user.click(placeBidButton);

    expect(mockOnPlaceBid).toHaveBeenCalledWith(150);
  });

  it('validates custom bid amount', async () => {
    const mockOnPlaceBid = jest.fn();
    const user = userEvent.setup();

    render(<BidButton {...defaultProps} variant="card" />);

    const customButton = screen.getByRole('button', { name: /custom/i });
    await user.click(customButton);

    const input = screen.getByRole('spinbutton');
    await user.type(input, '50'); // Lower than current bid

    const placeBidButton = screen.getByRole('button', { name: /place bid/i });
    await user.click(placeBidButton);

    expect(screen.getByText(/bid must be greater than \$100\.00/i)).toBeInTheDocument();
    expect(mockOnPlaceBid).not.toHaveBeenCalled();
  });

  it('has no accessibility violations', async () => {
    const { container } = render(<BidButton {...defaultProps} />);
    const results = await axe(container);
    expect(results).toHaveNoViolations();
  });

  it('handles network errors gracefully', async () => {
    const mockOnPlaceBid = jest.fn().mockRejectedValue(new Error('Network error'));
    const user = userEvent.setup();

    render(<BidButton {...defaultProps} onPlaceBid={mockOnPlaceBid} />);

    const button = screen.getByRole('button', { name: /place bid/i });
    await user.click(button);

    await waitFor(() => {
      expect(screen.getByText(/unknown error occurred/i)).toBeInTheDocument();
    });
  });

  it('generates correct quick bid amounts', () => {
    const quickBidAmounts = [125, 150, 175, 200];
    render(<BidButton {...defaultProps} quickBidAmounts={quickBidAmounts} variant="card" />);

    // Should show the provided quick bid amounts that are higher than current bid
    expect(screen.getByText(/\$125\.00/i)).toBeInTheDocument();
    expect(screen.getByText(/\$150\.00/i)).toBeInTheDocument();
  });
});
