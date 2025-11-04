import 'package:blytz_flutter_app/core/constants/route_constants.dart';
import 'package:blytz_flutter_app/data/models/auction_model.dart';
import 'package:blytz_flutter_app/shared/navigation/bottom_navigation.dart';
import 'package:blytz_flutter_app/shared/widgets/auction/auction_card.dart';
import 'package:flutter/material.dart';
import 'package:getwidget/getwidget.dart';
import 'package:go_router/go_router.dart';
import 'package:velocity_x/velocity_x.dart';

class HomePage extends StatelessWidget {
  const HomePage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: GFAppBar(
        title: const Text('Blytz Live Auctions'),
        backgroundColor: Theme.of(context).primaryColor,
        actions: [
          GFIconButton(
            icon: const Icon(Icons.notifications),
            onPressed: () {
              ScaffoldMessenger.of(context).showSnackBar(
                const SnackBar(content: Text('Notifications coming soon!')),
              );
            },
            type: GFButtonType.transparent,
          ),
        ],
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Welcome Section
            GFCard(
              margin: const EdgeInsets.only(bottom: 24),
              gradient: LinearGradient(
                colors: [Colors.blue.shade400, Colors.blue.shade600],
                begin: Alignment.topLeft,
                end: Alignment.bottomRight,
              ),
              borderRadius: BorderRadius.circular(16),
              content: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  'Welcome to Blytz!'
                      .text
                      .white
                      .xl3
                      .bold
                      .make(),
                  8.heightBox,
                  'Discover amazing live auctions from around the world.'
                      .text
                      .white
                      .make()
                      .opacity(value: 0.9),
                ],
              ),
            ),
            
            const SizedBox(height: 24),
            
            // Quick Actions
            'Quick Actions'
                .text
                .xl2
                .bold
                .make()
                .py(8),

            // First Row
            Row(
              children: [
                Expanded(
                  child: _buildActionCard(
                    context,
                    Icons.live_tv,
                    'Watch Live',
                    Colors.red,
                    () => context.push(RouteConstants.liveStream.replaceAll(':streamId', '1')),
                  ),
                ),
                12.widthBox,
                Expanded(
                  child: _buildActionCard(
                    context,
                    Icons.explore,
                    'Discover',
                    Colors.purple,
                    () => context.push(RouteConstants.discovery),
                  ),
                ),
                12.widthBox,
                Expanded(
                  child: _buildActionCard(
                    context,
                    Icons.category,
                    'Categories',
                    Colors.orange,
                    () => context.push(RouteConstants.categories),
                  ),
                ),
              ],
            ),

            12.heightBox,

            // Second Row
            Row(
              children: [
                Expanded(
                  child: _buildActionCard(
                    context,
                    Icons.people,
                    'Community',
                    Colors.green,
                    () => context.push(RouteConstants.community),
                  ),
                ),
                12.widthBox,
                Expanded(
                  child: _buildActionCard(
                    context,
                    Icons.store,
                    'Sell',
                    Colors.blue,
                    () => context.push(RouteConstants.sellerDashboard),
                  ),
                ),
                12.widthBox,
                Expanded(
                  child: _buildActionCard(
                    context,
                    Icons.person,
                    'Profile',
                    Colors.indigo,
                    () => context.push(RouteConstants.profile),
                  ),
                ),
              ],
            ),
            
            const SizedBox(height: 24),
            
            // Featured Auctions Preview
            'Featured Auctions'
                .text
                .xl2
                .bold
                .make()
                .py(8),
            
            // Mock auction cards for preview
            GridView.builder(
              shrinkWrap: true,
              physics: const NeverScrollableScrollPhysics(),
              gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                crossAxisCount: 2,
                childAspectRatio: 0.75,
                crossAxisSpacing: 16,
                mainAxisSpacing: 16,
              ),
              itemCount: 4, // Mock data
              itemBuilder: (context, index) {
                return AuctionCard(
                  auction: _getMockAuction(index),
                  onTap: () {
                    context.push(RouteConstants.liveStream.replaceAll(':streamId', 'stream_$index'));
                  },
                );
              },
            ),
          ],
        ),
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () => context.push(RouteConstants.createStream),
        backgroundColor: Theme.of(context).primaryColor,
        child: const Icon(Icons.live_tv, color: Colors.white),
      ),
      bottomNavigationBar: const BottomNavigation(currentIndex: 0),
    );
  }

  Widget _buildActionCard(
    BuildContext context,
    IconData icon,
    String title,
    Color color,
    VoidCallback onTap,
  ) {
    return GestureDetector(
      onTap: onTap,
      child: GFCard(
        margin: EdgeInsets.zero,
        padding: const EdgeInsets.all(20),
        borderRadius: BorderRadius.circular(12),
        color: Theme.of(context).cardColor,
        elevation: 4,
        content: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(icon, size: 40, color: color),
          12.heightBox,
          title.text
              .lg
              .semiBold
              .make()
              .centered(),
        ],
        ),
      ),
    );
  }

  // Mock auction data for preview
  AuctionModel _getMockAuction(int index) {
    final titles = ['Vintage Watch', 'Modern Art', 'Rare Coins', 'Antique Furniture'];
    final prices = [250.0, 1200.0, 450.0, 800.0];
    final bids = [12, 8, 15, 6];
    
    return AuctionModel(
      id: 'auction_$index',
      title: titles[index % titles.length],
      description: 'Amazing item up for auction',
      images: const [],
      categories: const ['Collectibles'],
      currentBidAmount: prices[index % prices.length],
      startingBid: prices[index % prices.length] * 0.8,
      totalBids: bids[index % bids.length],
      status: 'active',
      startTime: DateTime.now().subtract(const Duration(hours: 2)),
      endTime: DateTime.now().add(Duration(hours: index + 2)),
      isActive: true,
      isEndingSoon: index.isEven,
      timeLeft: '${index + 2}h ${30 - (index * 5)}m',
    );
  }
}