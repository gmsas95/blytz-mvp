import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:cached_network_image/cached_network_image.dart';
import 'package:getwidget/getwidget.dart';
import '../../../data/models/auction_model.dart';
import '../../../../core/utils/formatters.dart';
import '../../../../core/utils/extensions.dart';

class AuctionCard extends ConsumerWidget {
  final AuctionModel auction;
  final VoidCallback? onTap;
  final VoidCallback? onWatch;
  final bool showWatchButton;

  const AuctionCard({
    super.key,
    required this.auction,
    this.onTap,
    this.onWatch,
    this.showWatchButton = true,
  });

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Card(
      elevation: 4,
      shadowColor: Colors.black.withOpacity(0.1),
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(12),
      ),
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
                        icon: Icon(
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
                    Positioned(
                      bottom: 8,
                      right: 8,
                      child: Container(
                        padding: const EdgeInsets.symmetric(
                          horizontal: 8,
                          vertical: 4,
                        ),
                        decoration: BoxDecoration(
                          color: Colors.red,
                          borderRadius: BorderRadius.circular(12),
                        ),
                        child: const Row(
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
                    Text(
                      auction.title,
                      style: Theme.of(context).textTheme.titleSmall?.copyWith(
                        fontWeight: FontWeight.w600,
                      ),
                      maxLines: 2,
                      overflow: TextOverflow.ellipsis,
                    ),
                    
                    const SizedBox(height: 4),
                    
                    // Category
                    if (auction.categories.isNotEmpty)
                      Text(
                        auction.categories.first,
                        style: Theme.of(context).textTheme.bodySmall?.copyWith(
                          color: Theme.of(context).primaryColor,
                        ),
                      ),
                    
                    const Spacer(),
                    
                    // Current Price
                    Row(
                      children: [
                        Text(
                          'Current Bid',
                          style: Theme.of(context).textTheme.bodySmall?.copyWith(
                            color: Theme.of(context).colorScheme.onSurface.withOpacity(0.7),
                          ),
                        ),
                        const Spacer(),
                        Text(
                          auction.currentBidAmount.currencyFormat,
                          style: Theme.of(context).textTheme.titleMedium?.copyWith(
                            fontWeight: FontWeight.bold,
                            color: Theme.of(context).primaryColor,
                          ),
                        ),
                      ],
                    ),
                    
                    const SizedBox(height: 4),
                    
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
                            const SizedBox(width: 4),
                            Text(
                              '${auction.totalBids} bids',
                              style: Theme.of(context).textTheme.bodySmall?.copyWith(
                                color: Theme.of(context).colorScheme.onSurface.withOpacity(0.6),
                              ),
                            ),
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
                            const SizedBox(width: 4),
                            Text(
                              auction.timeLeft,
                              style: Theme.of(context).textTheme.bodySmall?.copyWith(
                                color: auction.isEndingSoon 
                                    ? Colors.orange 
                                    : Theme.of(context).colorScheme.onSurface.withOpacity(0.6),
                                fontWeight: auction.isEndingSoon ? FontWeight.w600 : FontWeight.normal,
                              ),
                            ),
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
        break;
      case 'ending_soon':
        backgroundColor = Colors.orange;
        textColor = Colors.white;
        text = 'ENDING SOON';
        break;
      case 'ended':
        backgroundColor = Colors.grey;
        textColor = Colors.white;
        text = 'ENDED';
        break;
      case 'scheduled':
        backgroundColor = Colors.blue;
        textColor = Colors.white;
        text = 'SCHEDULED';
        break;
      default:
        backgroundColor = Colors.grey;
        textColor = Colors.white;
        text = auction.status.toUpperCase();
    }

    return Container(
      padding: const EdgeInsets.symmetric(
        horizontal: 8,
        vertical: 4,
      ),
      decoration: BoxDecoration(
        color: backgroundColor,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Text(
        text,
        style: TextStyle(
          color: textColor,
          fontSize: 10,
          fontWeight: FontWeight.bold,
        ),
      ),
    );
  }
}