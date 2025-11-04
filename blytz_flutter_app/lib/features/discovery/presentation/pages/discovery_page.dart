import 'package:cached_network_image/cached_network_image.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:getwidget/getwidget.dart';
import 'package:velocity_x/velocity_x.dart';

class DiscoveryPage extends ConsumerStatefulWidget {
  const DiscoveryPage({super.key});

  @override
  ConsumerState<DiscoveryPage> createState() => _DiscoveryPageState();
}

class _DiscoveryPageState extends ConsumerState<DiscoveryPage>
    with TickerProviderStateMixin {
  late TabController _tabController;
  String _selectedCategory = 'All';
  final TextEditingController _searchController = TextEditingController();

  final List<String> _categories = [
    'All',
    'Electronics',
    'Fashion',
    'Collectibles',
    'Home & Garden',
    'Sports',
    'Toys & Games',
    'Books',
    'Art',
    'Jewelry',
    'Vintage',
  ];

  final List<Map<String, dynamic>> _trendingStreams = [
    {
      'title': 'Vintage Electronics Auction',
      'seller': 'TechCollector',
      'viewers': 1234,
      'thumbnail': 'https://example.com/thumb1.jpg',
      'category': 'Electronics',
      'isLive': true,
    },
    {
      'title': 'Designer Fashion Showcase',
      'seller': 'StyleGuru',
      'viewers': 856,
      'thumbnail': 'https://example.com/thumb2.jpg',
      'category': 'Fashion',
      'isLive': true,
    },
    {
      'title': 'Rare Coins & Stamps',
      'seller': 'HistoryBuff',
      'viewers': 445,
      'thumbnail': 'https://example.com/thumb3.jpg',
      'category': 'Collectibles',
      'isLive': true,
    },
  ];

  final List<Map<String, dynamic>> _upcomingStreams = [
    {
      'title': 'Antique Furniture Sale',
      'seller': 'FurniturePro',
      'startTime': '2:00 PM',
      'thumbnail': 'https://example.com/thumb4.jpg',
      'category': 'Home & Garden',
    },
    {
      'title': 'Sports Memorabilia',
      'seller': 'SportsFan',
      'startTime': '4:30 PM',
      'thumbnail': 'https://example.com/thumb5.jpg',
      'category': 'Sports',
    },
    {
      'title': 'Art Gallery Opening',
      'seller': 'ArtistLife',
      'startTime': '6:00 PM',
      'thumbnail': 'https://example.com/thumb6.jpg',
      'category': 'Art',
    },
  ];

  final List<String> _trendingTags = [
    '#vintage',
    '#electronics',
    '#fashion',
    '#rare',
    '#collectibles',
    '#limited',
    '#handmade',
    '#designer',
  ];

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 3, vsync: this);
  }

  @override
  void dispose() {
    _tabController.dispose();
    _searchController.dispose();
    super.dispose();
  }

  Widget _buildLiveStreamsTab() {
    return Column(
      children: [
        // Search Bar
        Padding(
          padding: const EdgeInsets.all(16.0),
          child: TextField(
            controller: _searchController,
            decoration: InputDecoration(
              hintText: 'Search live streams...',
              prefixIcon: const Icon(Icons.search),
              border: OutlineInputBorder(
                borderRadius: BorderRadius.circular(25.0),
              ),
              filled: true,
              fillColor: Colors.grey[100],
            ),
          ),
        ),
        // Trending Tags
        SizedBox(
          height: 50,
          child: ListView.builder(
            scrollDirection: Axis.horizontal,
            padding: const EdgeInsets.symmetric(horizontal: 16.0),
            itemCount: _trendingTags.length,
            itemBuilder: (context, index) {
              return Padding(
                padding: const EdgeInsets.only(right: 8.0),
                child: GFBadge(
                  text: _trendingTags[index],
                  shape: GFBadgeShape.standard,
                  color: Theme.of(context).primaryColor.withOpacity(0.1),
                  textColor: Theme.of(context).primaryColor,
                ),
              );
            },
          ),
        ),
        const SizedBox(height: 16),
        // Live Streams Grid
        Expanded(
          child: GridView.builder(
            padding: const EdgeInsets.all(16.0),
            gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
              crossAxisCount: 2,
              childAspectRatio: 0.8,
              crossAxisSpacing: 16,
              mainAxisSpacing: 16,
            ),
            itemCount: _trendingStreams.length,
            itemBuilder: (context, index) {
              final stream = _trendingStreams[index];
              return GestureDetector(
                onTap: () {
                  // Navigate to stream
                },
                child: Card(
                  elevation: 4,
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Expanded(
                        flex: 3,
                        child: Stack(
                          children: [
                            Container(
                              width: double.infinity,
                              decoration: BoxDecoration(
                                borderRadius: const BorderRadius.vertical(
                                  top: Radius.circular(12),
                                ),
                                color: Colors.grey[300],
                              ),
                              child: stream['thumbnail'] != null
                                  ? CachedNetworkImage(
                                       imageUrl: stream['thumbnail']?.toString() ?? '',
                                      fit: BoxFit.cover,
                                      placeholder: (context, url) => const Center(
                                        child: CircularProgressIndicator(),
                                      ),
                                      errorWidget: (context, url, error) =>
                                          const Icon(Icons.error),
                                    )
                                  : const Icon(Icons.image,
                                      size: 50, color: Colors.grey),
                            ),
                            Positioned(
                              top: 8,
                              left: 8,
                              child: Container(
                                padding: const EdgeInsets.symmetric(
                                  horizontal: 8,
                                  vertical: 4,
                                ),
                                decoration: BoxDecoration(
                                  color: Colors.red,
                                  borderRadius: BorderRadius.circular(4),
                                ),
                                child: const Text(
                                  'LIVE',
                                  style: TextStyle(
                                    color: Colors.white,
                                    fontSize: 10,
                                    fontWeight: FontWeight.bold,
                                  ),
                                ),
                              ),
                            ),
                            Positioned(
                              bottom: 8,
                              right: 8,
                              child: Container(
                                padding: const EdgeInsets.symmetric(
                                  horizontal: 6,
                                  vertical: 2,
                                ),
                                decoration: BoxDecoration(
                                  color: Colors.black54,
                                  borderRadius: BorderRadius.circular(4),
                                ),
                                child: Row(
                                  mainAxisSize: MainAxisSize.min,
                                  children: [
                                    const Icon(Icons.visibility,
                                        size: 12, color: Colors.white),
                                    const SizedBox(width: 4),
                                    Text(
                                      '${stream['viewers']}',
                                      style: const TextStyle(
                                        color: Colors.white,
                                        fontSize: 10,
                                      ),
                                    ),
                                  ],
                                ),
                              ),
                            ),
                          ],
                        ),
                      ),
                      Expanded(
                        flex: 2,
                        child: Padding(
                          padding: const EdgeInsets.all(8.0),
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Text(
                                stream['title']?.toString() ?? '',
                                style: const TextStyle(
                                  fontWeight: FontWeight.bold,
                                  fontSize: 12,
                                ),
                                maxLines: 2,
                                overflow: TextOverflow.ellipsis,
                              ),
                              const SizedBox(height: 4),
                              Text(
                                stream['seller']?.toString() ?? '',
                                style: TextStyle(
                                  color: Colors.grey[600],
                                  fontSize: 10,
                                ),
                              ),
                              const Spacer(),
                              Text(
                                stream['category']?.toString() ?? '',
                                style: TextStyle(
                                  color: Theme.of(context).primaryColor,
                                  fontSize: 10,
                                  fontWeight: FontWeight.w500,
                                ),
                              ),
                            ],
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
              );
            },
          ),
        ),
      ],
    );
  }

  Widget _buildUpcomingStreamsTab() {
    return ListView.builder(
      padding: const EdgeInsets.all(16.0),
      itemCount: _upcomingStreams.length,
      itemBuilder: (context, index) {
        final stream = _upcomingStreams[index];
        return Card(
          margin: const EdgeInsets.only(bottom: 16.0),
          child: ListTile(
            contentPadding: const EdgeInsets.all(16.0),
            leading: Container(
              width: 80,
              height: 80,
              decoration: BoxDecoration(
                borderRadius: BorderRadius.circular(8),
                color: Colors.grey[300],
              ),
              child: stream['thumbnail'] != null
                  ? CachedNetworkImage(
                      imageUrl: stream['thumbnail']?.toString() ?? '',
                      fit: BoxFit.cover,
                      placeholder: (context, url) => const Center(
                        child: CircularProgressIndicator(),
                      ),
                      errorWidget: (context, url, error) =>
                          const Icon(Icons.error),
                    )
                  : const Icon(Icons.image, size: 40, color: Colors.grey),
            ),
            title: Text(
              stream['title']?.toString() ?? '',
              style: const TextStyle(fontWeight: FontWeight.bold),
            ),
            subtitle: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text('Host: ${stream['seller']?.toString() ?? ''}'),
                const SizedBox(height: 4),
                Row(
                  children: [
                    const Icon(Icons.schedule, size: 16, color: Colors.grey),
                    const SizedBox(width: 4),
                    Text(
                      'Starts at ${stream['startTime'] ?? ''}',
                      style: const TextStyle(color: Colors.grey),
                    ),
                  ],
                ),
              ],
            ),
            trailing: GFButton(
              text: 'Set Reminder',
              type: GFButtonType.outline,
              size: GFSize.SMALL,
              onPressed: () {
                // Set reminder logic
              },
            ),
          ),
        );
      },
    );
  }

  Widget _buildCategoriesTab() {
    return GridView.builder(
      padding: const EdgeInsets.all(16.0),
      gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
        crossAxisCount: 2,
        childAspectRatio: 1.5,
        crossAxisSpacing: 16,
        mainAxisSpacing: 16,
      ),
      itemCount: _categories.length,
      itemBuilder: (context, index) {
        final category = _categories[index];
        final isSelected = _selectedCategory == category;
        
        return GestureDetector(
          onTap: () {
            setState(() {
              _selectedCategory = category;
            });
          },
          child: Card(
            elevation: isSelected ? 8 : 2,
            color: isSelected ? Theme.of(context).primaryColor : null,
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(12),
              side: BorderSide(
                color: isSelected ? Theme.of(context).primaryColor : Colors.grey,
                width: isSelected ? 2 : 1,
              ),
            ),
            child: Center(
              child: Text(
                category,
                style: TextStyle(
                  fontWeight: FontWeight.bold,
                  color: isSelected ? Colors.white : null,
                  fontSize: 16,
                ),
                textAlign: TextAlign.center,
              ),
            ),
          ),
        );
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Discover'),
        backgroundColor: Theme.of(context).primaryColor,
        foregroundColor: Colors.white,
        bottom: TabBar(
          controller: _tabController,
          labelColor: Colors.white,
          unselectedLabelColor: Colors.white70,
          indicatorColor: Colors.white,
          tabs: const [
            Tab(text: 'Live'),
            Tab(text: 'Upcoming'),
            Tab(text: 'Categories'),
          ],
        ),
      ),
      body: TabBarView(
        controller: _tabController,
        children: [
          _buildLiveStreamsTab(),
          _buildUpcomingStreamsTab(),
          _buildCategoriesTab(),
        ],
      ),
    );
  }

  Widget _buildFilterOption(String title, String value) {
    return ListTile(
      contentPadding: EdgeInsets.zero,
      title: Text(title),
      subtitle: Text(value),
      trailing: const Icon(Icons.chevron_right),
      onTap: () {
        // Handle filter option selection
      },
    );
  }
}