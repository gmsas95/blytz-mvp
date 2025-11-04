import 'package:blytz_flutter_app/core/utils/extensions.dart';
import 'package:blytz_flutter_app/data/models/auction_model.dart';
import 'package:blytz_flutter_app/data/providers/auction_detail_provider.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:getwidget/getwidget.dart';
import 'package:velocity_x/velocity_x.dart';

class BiddingWidget extends ConsumerStatefulWidget {

  const BiddingWidget({
    required this.auction, super.key,
  });
  final AuctionModel auction;

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
    return GFCard(
      margin: EdgeInsets.zero,
      padding: const EdgeInsets.all(16),
      borderRadius: const BorderRadius.vertical(top: Radius.circular(16)),
      elevation: 8,
      content: Column(
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
              8.widthBox,
              'Place Your Bid'
                  .text
                  .xl2
                  .bold
                  .make(),
            ],
          ),
          TextFormField(
            controller: _bidController,
            focusNode: _bidFocusNode,
            keyboardType: const TextInputType.numberWithOptions(decimal: true),
            decoration: InputDecoration(
              prefixText: r'$',
              border: const OutlineInputBorder(
                borderRadius: BorderRadius.all(Radius.circular(8)),
              ),
              errorText: _isValidBid() ? null : 'Minimum bid is ${_minimumBid.currencyFormat}',
            ),
            onChanged: (value) {
              setState(() {}); // Trigger rebuild to update validation
            },
          ),

          16.heightBox,

          // Bid Button
          GFButton(
            onPressed: (_isPlacingBid || !_isValidBid()) ? null : _placeBid,
            color: Theme.of(context).primaryColor,
            size: GFSize.LARGE,
            fullWidthButton: true,
            child: _isPlacingBid
                ? Row(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      const GFLoader(
                        type: GFLoaderType.circle,
                        size: GFSize.SMALL,
                      ),
                      8.widthBox,
                      'Placing Bid...'.text.white.make(),
                    ],
                  )
                : 'Place Bid'
                    .text
                    .white
                    .bold
                    .make(),
          ),

          8.heightBox,

          // Terms
          'By placing a bid, you commit to purchase this item if you win the auction.'
              .text
              .color(Colors.grey[600])
              .sm
              .center
              .make(),
        ],
      ),
    );
  }

  Widget _buildQuickBidButton(double amount) {
    return Expanded(
      child: GFButton(
        onPressed: () => _onQuickBid(amount),
        color: Colors.grey[200] ?? Colors.grey,
        textColor: Colors.black87,
        child: Text(
          '\$${amount.toStringAsFixed(2)}',
          style: const TextStyle(
            fontWeight: FontWeight.w600,
            color: Colors.black87,
          ),
        ),
      ),
    );
  }
}