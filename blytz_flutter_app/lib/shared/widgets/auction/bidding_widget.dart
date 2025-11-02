import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../data/models/auction_model.dart';
import '../../../data/providers/auction_detail_provider.dart';
import '../../../core/utils/formatters.dart';

class BiddingWidget extends ConsumerStatefulWidget {
  final AuctionModel auction;

  const BiddingWidget({
    super.key,
    required this.auction,
  });

  @override
  ConsumerState<BiddingWidget> createState() => _BiddingWidgetState();
}

class _BiddingWidgetState extends ConsumerState<BiddingWidget> {
  final TextEditingController _bidController = TextEditingController();
  final FocusNode _bidFocusNode = FocusNode();
  bool _isPlacingBid = false;

  @override
  void initState() {
    super.initState();
    // Set minimum bid as default
    _bidController.text = (widget.auction.currentBidAmount + 1.0).toStringAsFixed(2);
  }

  @override
  void dispose() {
    _bidController.dispose();
    _bidFocusNode.dispose();
    super.dispose();
  }

  double get _minimumBid => widget.auction.currentBidAmount + 1.0;
  double get _currentBid => double.tryParse(_bidController.text) ?? 0.0;

  bool _isValidBid() {
    return _currentBid >= _minimumBid;
  }

  void _onQuickBid(double amount) {
    setState(() {
      _bidController.text = amount.toStringAsFixed(2);
    });
    _bidFocusNode.requestFocus();
  }

  Future<void> _placeBid() async {
    if (!_isValidBid()) {
      _showError('Bid must be at least ${_minimumBid.currencyFormat}');
      return;
    }

    setState(() {
      _isPlacingBid = true;
    });

    try {
      final bidData = {
        'auctionId': widget.auction.id,
        'amount': _currentBid,
      };

      await ref.read(bidPlacementProvider(bidData).future);

      if (mounted) {
        _showSuccess('Bid placed successfully!');
        // Refresh auction data
        ref.invalidate(auctionDetailProvider(widget.auction.id));
        Navigator.of(context).pop();
      }
    } catch (e) {
      if (mounted) {
        _showError(e.toString());
      }
    } finally {
      if (mounted) {
        setState(() {
          _isPlacingBid = false;
        });
      }
    }
  }

  void _showError(String message) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text(message),
        backgroundColor: Colors.red,
      ),
    );
  }

  void _showSuccess(String message) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text(message),
        backgroundColor: Colors.green,
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: Theme.of(context).colorScheme.surface,
        borderRadius: const BorderRadius.vertical(top: Radius.circular(16)),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.1),
            blurRadius: 4,
            offset: const Offset(0, -2),
          ),
        ],
      ),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Header
          Row(
            children: [
              Icon(
                Icons.gavel,
                color: Theme.of(context).primaryColor,
              ),
              const SizedBox(width: 8),
              Text(
                'Place Your Bid',
                style: Theme.of(context).textTheme.titleLarge?.copyWith(
                  fontWeight: FontWeight.bold,
                ),
              ),
            ],
          ),

          const SizedBox(height: 16),

          // Current Bid Info
          Container(
            padding: const EdgeInsets.all(12),
            decoration: BoxDecoration(
              color: Colors.grey[100],
              borderRadius: BorderRadius.circular(8),
            ),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      'Current Bid',
                      style: Theme.of(context).textTheme.bodySmall?.copyWith(
                        color: Colors.grey[600],
                      ),
                    ),
                    Text(
                      widget.auction.currentBidAmount.currencyFormat,
                      style: Theme.of(context).textTheme.titleMedium?.copyWith(
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ],
                ),
                Column(
                  crossAxisAlignment: CrossAxisAlignment.end,
                  children: [
                    Text(
                      'Min. Next Bid',
                      style: Theme.of(context).textTheme.bodySmall?.copyWith(
                        color: Colors.grey[600],
                      ),
                    ),
                    Text(
                      _minimumBid.currencyFormat,
                      style: Theme.of(context).textTheme.titleMedium?.copyWith(
                        fontWeight: FontWeight.bold,
                        color: Theme.of(context).primaryColor,
                      ),
                    ),
                  ],
                ),
              ],
            ),
          ),

          const SizedBox(height: 16),

          // Quick Bid Buttons
          Text(
            'Quick Bid',
            style: Theme.of(context).textTheme.titleMedium?.copyWith(
              fontWeight: FontWeight.w600,
            ),
          ),
          const SizedBox(height: 8),
          Row(
            children: [
              _buildQuickBidButton(_minimumBid),
              const SizedBox(width: 8),
              _buildQuickBidButton(_minimumBid + 10),
              const SizedBox(width: 8),
              _buildQuickBidButton(_minimumBid + 25),
              const SizedBox(width: 8),
              _buildQuickBidButton(_minimumBid + 50),
            ],
          ),

          const SizedBox(height: 16),

          // Custom Bid Input
          Text(
            'Your Bid Amount',
            style: Theme.of(context).textTheme.titleMedium?.copyWith(
              fontWeight: FontWeight.w600,
            ),
          ),
          const SizedBox(height: 8),
          TextField(
            controller: _bidController,
            focusNode: _bidFocusNode,
            keyboardType: const TextInputType.numberWithOptions(decimal: true),
            decoration: InputDecoration(
              prefixText: '\$',
              border: OutlineInputBorder(
                borderRadius: BorderRadius.circular(8),
              ),
              errorText: _isValidBid() ? null : 'Minimum bid is ${_minimumBid.currencyFormat}',
            ),
            onChanged: (value) {
              setState(() {}); // Trigger rebuild to update validation
            },
          ),

          const SizedBox(height: 16),

          // Bid Button
          SizedBox(
            width: double.infinity,
            child: ElevatedButton(
              onPressed: (_isPlacingBid || !_isValidBid()) ? null : _placeBid,
              style: ElevatedButton.styleFrom(
                backgroundColor: Theme.of(context).primaryColor,
                padding: const EdgeInsets.symmetric(vertical: 16),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(8),
                ),
              ),
              child: _isPlacingBid
                  ? const Row(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        SizedBox(
                          width: 20,
                          height: 20,
                          child: CircularProgressIndicator(
                            strokeWidth: 2,
                            valueColor: AlwaysStoppedAnimation<Color>(Colors.white),
                          ),
                        ),
                        SizedBox(width: 8),
                        Text('Placing Bid...'),
                      ],
                    )
                  : const Text(
                      'Place Bid',
                      style: TextStyle(
                        fontSize: 16,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
            ),
          ),

          const SizedBox(height: 8),

          // Terms
          Text(
            'By placing a bid, you commit to purchase this item if you win the auction.',
            style: Theme.of(context).textTheme.bodySmall?.copyWith(
              color: Colors.grey[600],
            ),
            textAlign: TextAlign.center,
          ),
        ],
      ),
    );
  }

  Widget _buildQuickBidButton(double amount) {
    return Expanded(
      child: ElevatedButton(
        onPressed: () => _onQuickBid(amount),
        style: ElevatedButton.styleFrom(
          backgroundColor: Colors.grey[200],
          foregroundColor: Colors.black87,
          padding: const EdgeInsets.symmetric(vertical: 12),
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(8),
          ),
        ),
        child: Text(
          amount.currencyFormat,
          style: const TextStyle(
            fontWeight: FontWeight.w600,
          ),
        ),
      ),
    );
  }
}