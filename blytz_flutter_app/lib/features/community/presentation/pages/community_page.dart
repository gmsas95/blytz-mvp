import 'package:cached_network_image/cached_network_image.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:getwidget/getwidget.dart';
import 'package:velocity_x/velocity_x.dart';

class CommunityPage extends ConsumerStatefulWidget {
  const CommunityPage({super.key});

  @override
  ConsumerState<CommunityPage> createState() => _CommunityPageState();
}

class _CommunityPageState extends ConsumerState<CommunityPage>
    with TickerProviderStateMixin {
  late TabController _tabController;

  final List<Map<String, dynamic>> _topSellers = [
    {
      'name': 'TechCollector',
      'avatar': 'https://example.com/avatar1.jpg',
      'followers': 12543,
      'rating': 4.9,
      'category': 'Electronics',
      'isLive': true,
    },
    {
      'name': 'StyleGuru',
      'avatar': 'https://example.com/avatar2.jpg',
      'followers': 8932,
      'rating': 4.8,
      'category': 'Fashion',
      'isLive': false,
    },
    {
      'name': 'VintageHunter',
      'avatar': 'https://example.com/avatar3.jpg',
      'followers': 6745,
      'rating': 4.9,
      'category': 'Collectibles',
      'isLive': true,
    },
    {
      'name': 'SportsFanatic',
      'avatar': 'https://example.com/avatar4.jpg',
      'followers': 5234,
      'rating': 4.7,
      'category': 'Sports',
      'isLive': false,
    },
  ];

  final List<Map<String, dynamic>> _communityPosts = [
    {
      'author': 'TechCollector',
      'avatar': 'https://example.com/avatar1.jpg',
      'content': 'Just got my hands on a rare vintage gaming console! Going live in 30 minutes to showcase it.',
      'timestamp': '2 hours ago',
      'likes': 234,
      'comments': 45,
      'image': 'https://example.com/post1.jpg',
    },
    {
      'author': 'StyleGuru',
      'avatar': 'https://example.com/avatar2.jpg',
      'content': "New collection dropping tomorrow! Limited edition designer bags. Who's excited? 👜",
      'timestamp': '4 hours ago',
      'likes': 567,
      'comments': 89,
      'image': null,
    },
    {
      'author': 'VintageHunter',
      'avatar': 'https://example.com/avatar3.jpg',
      'content': "Amazing find at today's auction! This 1960s watch is in perfect condition.",
      'timestamp': '6 hours ago',
      'likes': 123,
      'comments': 23,
      'image': 'https://example.com/post2.jpg',
    },
  ];

  final List<Map<String, dynamic>> _trendingDiscussions = [
    {
      'title': 'Tips for beginners in vintage collecting',
      'author': 'CollectiblesExpert',
      'replies': 234,
      'views': 5678,
      'lastActivity': '15 minutes ago',
      'pinned': true,
    },
    {
      'title': 'Best practices for live streaming auctions',
      'author': 'StreamMaster',
      'replies': 156,
      'views': 3421,
      'lastActivity': '1 hour ago',
      'pinned': false,
    },
    {
      'title': 'Authentication guide for luxury items',
      'author': 'LuxuryDealer',
      'replies': 89,
      'views': 2103,
      'lastActivity': '3 hours ago',
      'pinned': false,
    },
  ];

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 4, vsync: this);
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Community'),
        backgroundColor: Theme.of(context).primaryColor,
        foregroundColor: Colors.white,
        actions: [
          GFIconButton(
            icon: const Icon(Icons.search),
            onPressed: _showSearch,
            type: GFButtonType.transparent,
          ),
          GFIconButton(
            icon: const Icon(Icons.add),
            onPressed: _createPost,
            type: GFButtonType.transparent,
          ),
        ],
      ),
      body: Column(
        children: [
          TabBar(
            controller: _tabController,
            labelColor: Theme.of(context).primaryColor,
            unselectedLabelColor: Colors.grey,
            indicatorColor: Theme.of(context).primaryColor,
            tabs: const [
              Tab(text: 'Feed'),
              Tab(text: 'Sellers'),
              Tab(text: 'Forums'),
              Tab(text: 'Events'),
            ],
          ),
          Expanded(
            child: TabBarView(
              controller: _tabController,
              children: [
                _buildFeedTab(),
                _buildSellersTab(),
                _buildForumsTab(),
                _buildEventsTab(),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildFeedTab() {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        children: [
          // Create Post Card
          GFCard(
            color: Colors.grey[100],
            content: Row(
              children: [
                Container(
                  width: 40,
                  height: 40,
                  decoration: BoxDecoration(
                    color: Theme.of(context).primaryColor,
                    borderRadius: BorderRadius.circular(20),
                  ),
                  child: const Icon(
                    Icons.person,
                    color: Colors.white,
                  ),
                ),
                12.widthBox,
                Expanded(
                  child: Text(
                    "What's on your mind?",
                    style: TextStyle(color: Colors.grey[600]),
                  ),
                ),
                IconButton(
                  icon: const Icon(Icons.photo_library),
                  onPressed: _createPost,
                ),
                IconButton(
                  icon: const Icon(Icons.send),
                  onPressed: _createPost,
                ),
              ],
            ),
          ),

          16.heightBox,

          // Community Posts
          ListView.builder(
            shrinkWrap: true,
            physics: const NeverScrollableScrollPhysics(),
            itemCount: _communityPosts.length,
            itemBuilder: (context, index) {
              final post = _communityPosts[index];
              return _buildPostCard(post);
            },
          ),
        ],
      ),
    );
  }

  Widget _buildSellersTab() {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Padding(
            padding: const EdgeInsets.symmetric(vertical: 8),
            child: Text(
              'Top Sellers This Week',
              style: Theme.of(context).textTheme.titleLarge?.copyWith(
                fontWeight: FontWeight.bold,
              ),
            ),
          ),

          GridView.builder(
            shrinkWrap: true,
            physics: const NeverScrollableScrollPhysics(),
            gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
              crossAxisCount: 2,
              childAspectRatio: 1.2,
              crossAxisSpacing: 12,
              mainAxisSpacing: 12,
            ),
            itemCount: _topSellers.length,
            itemBuilder: (context, index) {
              final seller = _topSellers[index];
              return _buildSellerCard(seller);
            },
          ),
        ],
      ),
    );
  }

  Widget _buildForumsTab() {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Padding(
            padding: const EdgeInsets.symmetric(vertical: 8),
            child: Text(
              'Trending Discussions',
              style: Theme.of(context).textTheme.titleLarge?.copyWith(
                fontWeight: FontWeight.bold,
              ),
            ),
          ),

          ListView.builder(
            shrinkWrap: true,
            physics: const NeverScrollableScrollPhysics(),
            itemCount: _trendingDiscussions.length,
            itemBuilder: (context, index) {
              final discussion = _trendingDiscussions[index];
              return _buildDiscussionCard(discussion);
            },
          ),
        ],
      ),
    );
  }

  Widget _buildEventsTab() {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Padding(
            padding: const EdgeInsets.symmetric(vertical: 8),
            child: Text(
              'Upcoming Events',
              style: Theme.of(context).textTheme.titleLarge?.copyWith(
                fontWeight: FontWeight.bold,
              ),
            ),
          ),

          _buildEventCard(
            'Vintage Electronics Fair',
            'March 15, 2024 • 2:00 PM',
            'Join us for the biggest vintage electronics auction of the year!',
            'assets/events/electronics.jpg',
            Colors.blue,
          ),

          12.heightBox,

          _buildEventCard(
            'Fashion Week Live Auction',
            'March 18, 2024 • 7:00 PM',
            'Exclusive designer pieces from top fashion houses.',
            'assets/events/fashion.jpg',
            Colors.pink,
          ),

          12.heightBox,

          _buildEventCard(
            'Collectors Convention',
            'March 22, 2024 • 10:00 AM',
            'Meet fellow collectors and find rare treasures.',
            'assets/events/collectibles.jpg',
            Colors.amber,
          ),
        ],
      ),
    );
  }

  Widget _buildPostCard(Map<String, dynamic> post) {
    return GFCard(
      margin: const EdgeInsets.only(bottom: 16),
      content: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Header
          Row(
            children: [
              Container(
                width: 40,
                height: 40,
                decoration: BoxDecoration(
                  color: Theme.of(context).primaryColor,
                  borderRadius: BorderRadius.circular(20),
                ),
                child: ClipRRect(
                  borderRadius: BorderRadius.circular(20),
                  child: CachedNetworkImage(
                    imageUrl: post['avatar'],
                    fit: BoxFit.cover,
                    placeholder: (context, url) => const Icon(
                      Icons.person,
                      color: Colors.white,
                    ),
                    errorWidget: (context, url, error) => const Icon(
                      Icons.person,
                      color: Colors.white,
                    ),
                  ),
                ),
              ),
              12.widthBox,
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    post['author'].text.semiBold.make(),
                    post['timestamp'].text.xs.color(Colors.grey[600]).make(),
                  ],
                ),
              ),
              GFIconButton(
                icon: const Icon(Icons.more_vert),
                onPressed: () {},
                type: GFButtonType.transparent,
              ),
            ],
          ),

          12.heightBox,

          // Content
          Text(
            post['content']?.toString() ?? '',
            style: Theme.of(context).textTheme.bodyMedium,
          ),

          // Image (if exists)
          if (post['image'] != null) ...[
            12.heightBox,
            ClipRRect(
              borderRadius: BorderRadius.circular(12),
              child: CachedNetworkImage(
                imageUrl: post['image'],
                fit: BoxFit.cover,
                placeholder: (context, url) => Container(
                  height: 200,
                  color: Colors.grey[200],
                  child: const Center(
                    child: Icon(Icons.image, color: Colors.grey),
                  ),
                ),
                errorWidget: (context, url, error) => Container(
                  height: 200,
                  color: Colors.grey[200],
                  child: const Center(
                    child: Icon(Icons.broken_image, color: Colors.grey),
                  ),
                ),
              ),
            ),
          ],

          12.heightBox,

          // Actions
          Row(
            children: [
              _buildActionButton(
                Icons.favorite_border,
                '${post['likes']}',
                () {},
              ),
              16.widthBox,
              _buildActionButton(
                Icons.comment_outlined,
                '${post['comments']}',
                () {},
              ),
              16.widthBox,
              _buildActionButton(
                Icons.share_outlined,
                'Share',
                () {},
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildSellerCard(Map<String, dynamic> seller) {
    return GFCard(
      margin: EdgeInsets.zero,
      padding: const EdgeInsets.all(12),
      content: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          // Avatar with Live Indicator
          Stack(
            alignment: Alignment.bottomRight,
            children: [
              Container(
                width: 60,
                height: 60,
                decoration: BoxDecoration(
                  color: Theme.of(context).primaryColor,
                  borderRadius: BorderRadius.circular(30),
                ),
                child: ClipRRect(
                  borderRadius: BorderRadius.circular(30),
                  child: CachedNetworkImage(
                    imageUrl: seller['avatar'],
                    fit: BoxFit.cover,
                    placeholder: (context, url) => const Icon(
                      Icons.person,
                      color: Colors.white,
                    ),
                    errorWidget: (context, url, error) => const Icon(
                      Icons.person,
                      color: Colors.white,
                    ),
                  ),
                ),
              ),
              if (seller['isLive'])
                Container(
                  width: 16,
                  height: 16,
                  decoration: const BoxDecoration(
                    color: Colors.red,
                    shape: BoxShape.circle,
                  ),
                  child: const Center(
                    child: Icon(
                      Icons.fiber_manual_record,
                      color: Colors.white,
                      size: 10,
                    ),
                  ),
                ),
            ],
          ),

          8.heightBox,

          // Seller Info
          seller['name'].text.semiBold.maxLines(1).overflow(TextOverflow.ellipsis).make(),
          2.heightBox,
          seller['category'].text.xs.color(Colors.grey[600]).make(),
          4.heightBox,

          // Rating
          Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              ...List.generate(5, (index) {
                return Icon(
                  index < seller['rating'].floor()
                      ? Icons.star
                      : index < seller['rating']
                          ? Icons.star_half
                          : Icons.star_border,
                  color: Colors.amber,
                  size: 12,
                );
              }),
              4.widthBox,
              '${seller['rating']}'.text.xs.make(),
            ],
          ),

          4.heightBox,

          // Followers
          '${(seller['followers'] / 1000).toStringAsFixed(1)}k followers'
              .text
              .xs
              .color(Theme.of(context).primaryColor)
              .make(),

          8.heightBox,

          // Follow Button
          GFButton(
            onPressed: () {},
            text: seller['isLive'] ? 'Watch' : 'Follow',
            color: seller['isLive'] ? Colors.red : Theme.of(context).primaryColor,
            size: GFSize.SMALL,
            fullWidthButton: true,
          ),
        ],
      ),
    );
  }

  Widget _buildDiscussionCard(Map<String, dynamic> discussion) {
    return GFCard(
      margin: const EdgeInsets.only(bottom: 12),
      content: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              if (discussion['pinned'])
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
                  decoration: BoxDecoration(
                    color: Colors.orange,
                    borderRadius: BorderRadius.circular(8),
                  ),
                  child: const Text(
                    'PINNED',
                    style: TextStyle(
                      color: Colors.white,
                      fontSize: 10,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                ),
              if (discussion['pinned']) 8.widthBox,
              Expanded(
                child: discussion['title'].text.semiBold.make(),
              ),
            ],
          ),

          4.heightBox,

          Row(
            children: [
              discussion['author'].text.sm.color(Colors.grey[600]).make(),
              8.widthBox,
              '•'.text.color(Colors.grey[400]).make(),
              8.widthBox,
              discussion['lastActivity'].text.sm.color(Colors.grey[600]).make(),
            ],
          ),

          8.heightBox,

          Row(
            children: [
              Icon(Icons.comment_outlined, size: 16, color: Colors.grey[600]),
              4.widthBox,
              '${discussion['replies']} replies'.text.sm.color(Colors.grey[600]).make(),
              16.widthBox,
              Icon(Icons.visibility_outlined, size: 16, color: Colors.grey[600]),
              4.widthBox,
              '${discussion['views']} views'.text.sm.color(Colors.grey[600]).make(),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildEventCard(
    String title,
    String date,
    String description,
    String image,
    Color color,
  ) {
    return GFCard(
      content: Row(
        children: [
          // Event Image
          Container(
            width: 80,
            height: 80,
            decoration: BoxDecoration(
              color: color.withOpacity(0.1),
              borderRadius: BorderRadius.circular(12),
            ),
            child: Icon(
              Icons.event,
              color: color,
              size: 32,
            ),
          ),

          12.widthBox,

          // Event Info
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                title.text.semiBold.make(),
                4.heightBox,
                date.text.sm.color(color).make(),
                8.heightBox,
                description.text.sm.color(Colors.grey[600]).maxLines(2).overflow(TextOverflow.ellipsis).make(),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildActionButton(IconData icon, String label, VoidCallback onPressed) {
    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        IconButton(
          icon: Icon(icon, size: 20),
          onPressed: onPressed,
          visualDensity: VisualDensity.compact,
        ),
        if (label.isNotEmpty) label.text.sm.make(),
      ],
    );
  }

  void _showSearch() {
    showSearch(
      context: context,
      delegate: CommunitySearchDelegate(),
    );
  }

  void _createPost() {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      builder: (context) => Container(
        height: MediaQuery.of(context).size.height * 0.7,
        padding: const EdgeInsets.all(16),
        child: Column(
          children: [
            const Text(
              'Create Post',
              style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
            ),
            const SizedBox(height: 16),
            const TextField(
              maxLines: 8,
              decoration: InputDecoration(
                labelText: "What's on your mind?",
                border: OutlineInputBorder(),
              ),
            ),
            const SizedBox(height: 16),
            Row(
              children: [
                IconButton(
                  icon: const Icon(Icons.photo_library),
                  onPressed: () {},
                ),
                IconButton(
                  icon: const Icon(Icons.camera_alt),
                  onPressed: () {},
                ),
                const Spacer(),
                GFButton(
                  onPressed: () => Navigator.pop(context),
                  text: 'Post',
                  color: Theme.of(context).primaryColor,
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}

class CommunitySearchDelegate extends SearchDelegate<String> {
  @override
  List<Widget> buildActions(BuildContext context) {
    return [
      IconButton(
        icon: const Icon(Icons.clear),
        onPressed: () => query = '',
      ),
    ];
  }

  @override
  Widget buildLeading(BuildContext context) {
    return IconButton(
      icon: const Icon(Icons.arrow_back),
      onPressed: () => close(context, ''),
    );
  }

  @override
  Widget buildResults(BuildContext context) {
    return const Center(
      child: Text('Search results'),
    );
  }

  @override
  Widget buildSuggestions(BuildContext context) {
    return const Center(
      child: Text('Search suggestions'),
    );
  }
}