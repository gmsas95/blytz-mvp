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
    );
  }
  }

  Widget _buildLiveStreamsTab() {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Trending Tags
          'Trending Tags'
              .text
              .lg
              .bold
              .make()
              .py(8),

          Wrap(
            spacing: 8,
            runSpacing: 8,
            children: _trendingTags.map((tag) {
              return Container(
                padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
                decoration: BoxDecoration(
                  color: Theme.of(context).primaryColor.withOpacity(0.1),
                  borderRadius: BorderRadius.circular(16),
                  border: Border.all(
                    color: Theme.of(context).primaryColor.withOpacity(0.3),
                  ),
                ),
                child: Text(
                  tag,
                  style: TextStyle(
                    color: Theme.of(context).primaryColor,
                    fontSize: 12,
                  ),
                ),
              );
            }).toList(),
          ),

          24.heightBox,

          // Live Streams
          'Live Now 🔴'
              .text
              .xl
              .bold
              .make()
              .py(8),

          GridView.builder(
            shrinkWrap: true,
            physics: const NeverScrollableScrollPhysics(),
            gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
              crossAxisCount: 2,
              childAspectRatio: 0.8,
              crossAxisSpacing: 12,
              mainAxisSpacing: 12,
            ),
            itemCount: _trendingStreams.length,
            itemBuilder: (context, index) {
              return _buildStreamCard(_trendingStreams[index], isLive: true);
            },
          ),
        ],
      ),
    );
  }

  Widget _buildUpcomingStreamsTab() {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          'Upcoming Streams'
              .text
              .xl
              .bold
              .make()
              .py(8),

          ListView.builder(
            shrinkWrap: true,
            physics: const NeverScrollableScrollPhysics(),
            itemCount: _upcomingStreams.length,
            itemBuilder: (context, index) {
              return Container(
                margin: const EdgeInsets.only(bottom: 12),
                child: _buildUpcomingStreamCard(_upcomingStreams[index]),
              );
            },
          ),
        ],
      ),
    );
  }

  Widget _buildCategoriesTab() {
    return GridView.builder(
      padding: const EdgeInsets.all(16),
      gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
        crossAxisCount: 2,
        childAspectRatio: 1.2,
        crossAxisSpacing: 12,
        mainAxisSpacing: 12,
      ),
      itemCount: _categories.length - 1, // Exclude "All"
      itemBuilder: (context, index) {
        final category = _categories[index + 1];
        return _buildCategoryCard(category);
      },
    );
  }

  Widget _buildStreamCard(Map<String, dynamic> stream, {required bool isLive}) {
    return Card(
      margin: EdgeInsets.zero,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      elevation: 4,
      child: InkWell(
        onTap: () => _navigateToStream(stream),
        borderRadius: BorderRadius.circular(12),
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
                      color: Colors.grey[300],
                      borderRadius: const BorderRadius.vertical(
                        top: Radius.circular(12),
                      ),
                    ),
                    child: ClipRRect(
                      borderRadius: const BorderRadius.vertical(
                        top: Radius.circular(12),
                      ),
                      child: CachedNetworkImage(
                        imageUrl: stream['thumbnail'],
                        fit: BoxFit.cover,
                        placeholder: (context, url) => Container(
                          color: Colors.grey[300],
                          child: const Center(
                            child: Icon(Icons.image, color: Colors.grey),
                          ),
                        ),
                        errorWidget: (context, url, error) => Container(
                          color: Colors.grey[300],
                          child: const Center(
                            child: Icon(Icons.broken_image, color: Colors.grey),
                          ),
                        ),
                      ),
                    ),
                  ),

                  // Live Badge
                  if (isLive)
                    Positioned(
                      top: 8,
                      left: 8,
                      child: Container(
                        padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
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

                  // Viewer Count
                  Positioned(
                    bottom: 8,
                    right: 8,
                    child: Container(
                      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                      decoration: BoxDecoration(
                        color: Colors.black.withOpacity(0.7),
                        borderRadius: BorderRadius.circular(12),
                      ),
                      child: Row(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          const Icon(
                            Icons.visibility,
                            color: Colors.white,
                            size: 12,
                          ),
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
                padding: const EdgeInsets.all(8),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      stream['title'],
                      style: const TextStyle(
                        fontSize: 12,
                        fontWeight: FontWeight.bold,
                      ),
                      maxLines: 2,
                      overflow: TextOverflow.ellipsis,
                    ),
                    const SizedBox(height: 4),
                    Text(
                      stream['seller'],
                      style: const TextStyle(
                        fontSize: 10,
                        color: Colors.grey,
                      ),
                    ),
                    const Spacer(),
                    Text(
                      stream['category'],
                      style: TextStyle(
                        fontSize: 10,
                        color: Theme.of(context).primaryColor,
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
  }
  }

  Widget _buildUpcomingStreamCard(Map<String, dynamic> stream) {
    return Card(
      margin: EdgeInsets.zero,
      child: Padding(
        padding: const EdgeInsets.all(12),
        child: Row(
          children: [
          // Thumbnail
          Container(
            width: 80,
            height: 80,
            decoration: BoxDecoration(
              color: Colors.grey[300],
              borderRadius: BorderRadius.circular(8),
            ),
            child: ClipRRect(
              borderRadius: BorderRadius.circular(8),
              child: CachedNetworkImage(
                imageUrl: stream['thumbnail'],
                fit: BoxFit.cover,
                placeholder: (context, url) => const Center(
                  child: Icon(Icons.image, color: Colors.grey),
                ),
                errorWidget: (context, url, error) => const Center(
                  child: Icon(Icons.broken_image, color: Colors.grey),
                ),
              ),
            ),
          ),

          12.widthBox,

          // Content
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  stream['title'],
                  style: const TextStyle(
                    fontSize: 14,
                    fontWeight: FontWeight.bold,
                  ),
                  maxLines: 2,
                  overflow: TextOverflow.ellipsis,
                ),
                const SizedBox(height: 4),
                Text(
                  stream['seller'],
                  style: const TextStyle(
                    fontSize: 12,
                    color: Colors.grey,
                  ),
                ),
                const SizedBox(height: 4),
                Row(
                  children: [
                    const Icon(
                      Icons.schedule,
                      size: 12,
                      color: Colors.grey,
                    ),
                    const SizedBox(width: 4),
                    Text(
                      stream['startTime'],
                      style: const TextStyle(
                        fontSize: 12,
                        color: Colors.grey,
                        fontWeight: FontWeight.w500,
                      ),
                    ),
                  ],
                ),
              ],
            ),
          ),

          // Category Badge
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
            decoration: BoxDecoration(
              color: Theme.of(context).primaryColor.withOpacity(0.1),
              borderRadius: BorderRadius.circular(12),
            ),
            child: Text(
              stream['category'],
              style: TextStyle(
                color: Theme.of(context).primaryColor,
                fontSize: 10,
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildCategoryCard(String category) {
    final icons = {
      'Electronics': Icons.devices,
      'Fashion': Icons.checkroom,
      'Collectibles': Icons.collections,
      'Home & Garden': Icons.home,
      'Sports': Icons.sports_soccer,
      'Toys & Games': Icons.toys,
      'Books': Icons.book,
      'Art': Icons.palette,
      'Jewelry': Icons.diamond,
      'Vintage': Icons.watch,
    };

    final colors = {
      'Electronics': Colors.blue,
      'Fashion': Colors.pink,
      'Collectibles': Colors.amber,
      'Home & Garden': Colors.green,
      'Sports': Colors.orange,
      'Toys & Games': Colors.purple,
      'Books': Colors.brown,
      'Art': Colors.indigo,
      'Jewelry': Colors.teal,
      'Vintage': Colors.grey,
    };

    return GFCard(
      margin: EdgeInsets.zero,
      padding: const EdgeInsets.all(16),
      borderRadius: BorderRadius.circular(12),
      color: (colors[category] ?? Colors.grey).withOpacity(0.1),
      onTap: () => _navigateToCategory(category),
      content: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(
            icons[category] ?? Icons.category,
            size: 40,
            color: colors[category] ?? Colors.grey,
          ),
          const SizedBox(height: 8),
          Text(
            category,
            textAlign: TextAlign.center,
            style: TextStyle(
              fontSize: 14,
              fontWeight: FontWeight.bold,
              color: colors[category] ?? Colors.grey,
            ),
          ),
        ],
      ),
    );
  }

  void _navigateToStream(Map<String, dynamic> stream) {
    // Navigate to live stream page
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text('Opening ${stream['title']}...')),
    );
  }

  void _navigateToCategory(String category) {
    // Navigate to category page
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text('Browsing $category...')),
    );
  }

  void _showFilterDialog() {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Filter Streams'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            _buildFilterOption('Price Range', 'Min - Max'),
            _buildFilterOption('Duration', 'Any'),
            _buildFilterOption('Language', 'English'),
            _buildFilterOption('Rating', '4+ Stars'),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Reset'),
          ),
          GFButton(
            onPressed: () => Navigator.pop(context),
            text: 'Apply Filters',
            color: Theme.of(context).primaryColor,
          ),
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