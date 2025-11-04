import 'package:blytz_flutter_app/core/utils/extensions.dart';
import 'package:blytz_flutter_app/data/models/auction_model.dart';
import 'package:cached_network_image/cached_network_image.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:getwidget/getwidget.dart';
import 'package:velocity_x/velocity_x.dart';

class AuctionCard extends ConsumerWidget {

  const AuctionCard({
    required this.auction, super.key,
    this.onTap,
    this.onWatch,
    this.showWatchButton = true,
  });
  final AuctionModel auction;
  final VoidCallback? onTap;
  final VoidCallback? onWatch;
  final bool showWatchButton;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Card(
      elevation: 4,
      shadowColor: Colors.black.withOpacity(0.1),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(12),
        child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Image Section
          Expanded(
            flex: 3,
            child: Stack(
              children: [
                // Product Image
                Container(
                  width: double.infinity,
                  decoration: const BoxDecoration(
                    borderRadius: BorderRadius.vertical(
                      top: Radius.circular(12),
                    ),
                  ),
                  child: ClipRRect(
                    borderRadius: const BorderRadius.vertical(
                      top: Radius.circular(12),
                    ),
                    child: auction.images.isNotEmpty
                        ? CachedNetworkImage(
                            imageUrl: auction.images.first,
                            fit: BoxFit.cover,
                            placeholder: (context, url) => Container(
                              color: Colors.grey[300],
                              child: const Center(
                                child: Icon(
                                  Icons.image,
                                  color: Colors.grey,
                                  size: 40,
                                ),
                              ),
                            ),
                            errorWidget: (context, url, error) => Container(
                              color: Colors.grey[300],
                              child: const Center(
                                child: Icon(
                                  Icons.broken_image,
                                  color: Colors.grey,
                                  size: 40,
                                ),
                              ),
                            ),
                          )
                        : Container(
                            color: Colors.grey[300],
                            child: const Center(
                              child: Icon(
                                Icons.image,
                                color: Colors.grey,
                                size: 40,
                              ),
                            ),
                          ),
                  ),
                ),
                
                // Status Badge
                Positioned(
                  top: 8,
                  left: 8,
                  child: _buildStatusBadge(),
                ),
                
                // Watch Button
                if (showWatchButton)
                  Positioned(
                    top: 8,
                    right: 8,
                    child: GFIconButton(
                      type: GFButtonType.transparent,
                      icon: const Icon(
                        Icons.favorite_border,
                        color: Colors.white,
                        size: 20,
                      ),
                      onPressed: onWatch,
                      size: GFSize.SMALL,
                      color: Colors.black.withOpacity(0.5),
                    ),
                  ),
                
                // Live Badge for Active Auctions
                if (auction.isActive)
                  const Positioned(
                    bottom: 8,
                    right: 8,
                    child: GFBadge(
                      text: 'LIVE',
                      color: Colors.red,
                      shape: GFBadgeShape.standard,
                      child: Row(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          Icon(
                            Icons.fiber_manual_record,
                            color: Colors.white,
                            size: 8,
                          ),
                          SizedBox(width: 4),
                          Text(
                            'LIVE',
                            style: TextStyle(
                              color: Colors.white,
                              fontSize: 10,
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                        ],
                      ),
                    ),
                  ),
              ],
            ),
          ),
          
          // Content Section
          Expanded(
            flex: 2,
            child: Padding(
              padding: const EdgeInsets.all(12),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  // Title
                  auction.title
                      .text
                      .semiBold
                      .maxLines(2)
                      .overflow(TextOverflow.ellipsis)
                      .make(),
                  
                  4.heightBox,
                  
                  // Category
                  if (auction.categories.isNotEmpty)
                    auction.categories.first
                        .text
                        .color(Theme.of(context).primaryColor)
                        .sm
                        .make(),
                  
                  const Spacer(),
                  
                  // Current Price
                  Row(
                    children: [
                      'Current Bid'
                          .text
                          .color(Theme.of(context).colorScheme.onSurface.withOpacity(0.7))
                          .sm
                          .make(),
                      const Spacer(),
                      auction.currentBidAmount.currencyFormat
                          .text
                          .color(Theme.of(context).primaryColor)
                          .semiBold
                          .lg
                          .make(),
                    ],
                  ),
                  
                  4.heightBox,
                  
                  // Bid Count and Time Left
                  Row(
                    children: [
                      Row(
                        children: [
                          Icon(
                            Icons.gavel,
                            size: 16,
                            color: Theme.of(context).colorScheme.onSurface.withOpacity(0.6),
                          ),
                          4.widthBox,
                          '${auction.totalBids} bids'
                              .text
                              .color(Theme.of(context).colorScheme.onSurface.withOpacity(0.6))
                              .sm
                              .make(),
                        ],
                      ),
                      const Spacer(),
                      Row(
                        children: [
                          Icon(
                            Icons.access_time,
                            size: 16,
                            color: auction.isEndingSoon 
                                ? Colors.orange 
                                : Theme.of(context).colorScheme.onSurface.withOpacity(0.6),
                          ),
                          4.widthBox,
                          auction.timeLeft
                              .text
                              .color(auction.isEndingSoon 
                                  ? Colors.orange 
                                  : Theme.of(context).colorScheme.onSurface.withOpacity(0.6),)
                              .sm
                              .semiBold
                              .make(),
                        ],
                      ),
                    ],
                  ),
                ],
              ),
            ),
          ),
        ],
      ),
      ),
    );
  }

  Widget _buildStatusBadge() {
    Color backgroundColor;
    Color textColor;
    String text;

    switch (auction.status) {
      case 'active':
        backgroundColor = Colors.green;
        textColor = Colors.white;
        text = 'ACTIVE';
      case 'ending_soon':
        backgroundColor = Colors.orange;
        textColor = Colors.white;
        text = 'ENDING SOON';
      case 'ended':
        backgroundColor = Colors.grey;
        textColor = Colors.white;
        text = 'ENDED';
      case 'scheduled':
        backgroundColor = Colors.blue;
        textColor = Colors.white;
        text = 'SCHEDULED';
      default:
        backgroundColor = Colors.grey;
        textColor = Colors.white;
        text = auction.status.toUpperCase();
    }

    return GFBadge(
      text: text,
      color: backgroundColor,
      textColor: textColor,
      shape: GFBadgeShape.standard,
    );
  }
}