import 'package:cached_network_image/cached_network_image.dart';
import 'package:fl_chart/fl_chart.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:getwidget/getwidget.dart';
import 'package:velocity_x/velocity_x.dart';

class EnhancedProfilePage extends ConsumerStatefulWidget {
  const EnhancedProfilePage({super.key});

  @override
  ConsumerState<EnhancedProfilePage> createState() => _EnhancedProfilePageState();
}

class _EnhancedProfilePageState extends ConsumerState<EnhancedProfilePage>
    with TickerProviderStateMixin {
  late TabController _tabController;
  bool _isEditing = false;
  final bool _isSeller = false;

  final TextEditingController _usernameController = TextEditingController(text: 'JohnDoe123');
  final TextEditingController _emailController = TextEditingController(text: 'john.doe@example.com');
  final TextEditingController _bioController = TextEditingController(text: 'Auction enthusiast & collector');
  final TextEditingController _locationController = TextEditingController(text: 'New York, USA');

  // Mock stats
  final Map<String, int> _stats = {
    'auctionsWon': 47,
    'auctionsSold': 23,
    'totalSpent': 12543,
    'totalEarned': 8756,
    'followers': 1234,
    'following': 567,
    'rating': 48, // 4.8 * 10
    'ratingCount': 156,
  };

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 4, vsync: this);
  }

  @override
  void dispose() {
    _tabController.dispose();
    _usernameController.dispose();
    _emailController.dispose();
    _bioController.dispose();
    _locationController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: CustomScrollView(
        slivers: [
          // Custom App Bar with Cover Image
          SliverAppBar(
            expandedHeight: 200,
            pinned: true,
            backgroundColor: Theme.of(context).primaryColor,
            flexibleSpace: FlexibleSpaceBar(
              background: Container(
                decoration: BoxDecoration(
                  gradient: LinearGradient(
                    begin: Alignment.topLeft,
                    end: Alignment.bottomRight,
                    colors: [
                      Theme.of(context).primaryColor,
                      Theme.of(context).primaryColor.withOpacity(0.8),
                    ],
                  ),
                ),
                child: Stack(
                  children: [
                    // Cover Image Placeholder
                    Container(
                      decoration: BoxDecoration(
                        image: DecorationImage(
                          image: const NetworkImage('https://example.com/cover.jpg'),
                          fit: BoxFit.cover,
                          onError: (exception, stackTrace) {
                            // Gradient fallback
                          },
                        ),
                      ),
                    ),
                    // Edit Cover Button
                    Positioned(
                      top: 60,
                      right: 16,
                      child: IconButton(
                        icon: const Icon(Icons.camera_alt, color: Colors.white),
                        style: IconButton.styleFrom(
                          backgroundColor: Colors.black.withOpacity(0.5),
                        ),
                        onPressed: _changeCoverImage,
                      ),
                    ),
                  ],
                ),
              ),
            ),
          ),

          // Profile Info and Stats
          SliverToBoxAdapter(
            child: Container(
              decoration: BoxDecoration(
                color: Theme.of(context).scaffoldBackgroundColor,
                borderRadius: const BorderRadius.vertical(top: Radius.circular(20)),
              ),
              child: Column(
                children: [
                  // Profile Picture and Basic Info
                  Padding(
                    padding: const EdgeInsets.all(16),
                    child: Column(
                      children: [
                        // Profile Picture
                        Transform.translate(
                          offset: const Offset(0, -50),
                          child: Stack(
                            alignment: Alignment.bottomRight,
                            children: [
                              Container(
                                width: 100,
                                height: 100,
                                decoration: BoxDecoration(
                                  color: Theme.of(context).primaryColor,
                                  borderRadius: BorderRadius.circular(50),
                                  border: Border.all(color: Colors.white, width: 4),
                                ),
                                child: ClipRRect(
                                  borderRadius: BorderRadius.circular(50),
                                  child: CachedNetworkImage(
                                    imageUrl: 'https://example.com/avatar.jpg',
                                    fit: BoxFit.cover,
                                    placeholder: (context, url) => const Icon(
                                      Icons.person,
                                      color: Colors.white,
                                      size: 40,
                                    ),
                                    errorWidget: (context, url, error) => const Icon(
                                      Icons.person,
                                      color: Colors.white,
                                      size: 40,
                                    ),
                                  ),
                                ),
                              ),
                              GFIconButton(
                                icon: const Icon(Icons.camera_alt, size: 16),
                                onPressed: _changeProfilePicture,
                                size: GFSize.SMALL,
                                color: Theme.of(context).primaryColor,
                              ),
                            ],
                          ),
                        ),

                        const SizedBox(height: 20),

                        // Username and Bio
                        Row(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            if (_isEditing)
                              SizedBox(
                                width: 200,
                                child: TextField(
                                  controller: _usernameController,
                                  textAlign: TextAlign.center,
                                  decoration: const InputDecoration(
                                    border: UnderlineInputBorder(),
                                  ),
                                ),
                              )
                            else
                              _usernameController.text.text.xl2.bold.make(),
                            IconButton(
                              icon: Icon(_isEditing ? Icons.check : Icons.edit),
                              onPressed: () {
                                setState(() {
                                  _isEditing = !_isEditing;
                                });
                              },
                            ),
                          ],
                        ),

                        8.heightBox,

                        if (_isEditing)
                          SizedBox(
                            width: 300,
                            child: TextField(
                              controller: _bioController,
                              textAlign: TextAlign.center,
                              maxLines: 3,
                              decoration: const InputDecoration(
                                
                        border: OutlineInputBorder(),
                              ),
                            ),
                          )
                        else
                          _bioController.text.text.center.color(Colors.grey[600]).make(),

                        16.heightBox,

                        // Location
                        Row(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            const Icon(Icons.location_on, size: 16, color: Colors.grey),
                            4.widthBox,
                            if (_isEditing)
                              SizedBox(
                                width: 150,
                                child: TextField(
                                  controller: _locationController,
                                  textAlign: TextAlign.center,
                                  decoration: const InputDecoration(
                                    border: UnderlineInputBorder(),
                                  ),
                                ),
                              )
                            else
                              _locationController.text.text.sm.color(Colors.grey[600]).make(),
                          ],
                        ),

                        16.heightBox,

                        // Action Buttons
                        Row(
                          mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                          children: [
                            GFButton(
                              onPressed: _shareProfile,
                              text: 'Share Profile',
                              type: GFButtonType.outline2x,
                              color: Theme.of(context).primaryColor,
                              size: GFSize.SMALL,
                            ),
                            GFButton(
                              onPressed: _viewPublicProfile,
                              text: 'View Public',
                              type: GFButtonType.outline2x,
                              color: Theme.of(context).primaryColor,
                              size: GFSize.SMALL,
                            ),
                          ],
                        ),

                        16.heightBox,

                        // Stats Cards
                        Container(
                          padding: const EdgeInsets.all(16),
                          decoration: BoxDecoration(
                            color: Colors.grey[100],
                            borderRadius: BorderRadius.circular(12),
                          ),
                          child: Column(
                            children: [
                              Padding(
                                    padding: const EdgeInsets.symmetric(vertical: 8),
                                    child: Text(
                                      'Statistics',
                                      style: Theme.of(context).textTheme.titleMedium?.copyWith(
                                        fontWeight: FontWeight.bold,
                                      ),
                                    ),
                                  ),
                              GridView.count(
                                shrinkWrap: true,
                                physics: const NeverScrollableScrollPhysics(),
                                crossAxisCount: 3,
                                crossAxisSpacing: 12,
                                mainAxisSpacing: 12,
                                childAspectRatio: 1.5,
                                children: [
                                  _buildStatCard('Auctions Won', '${_stats['auctionsWon']}'),
                                  _buildStatCard('Auctions Sold', '${_stats['auctionsSold']}'),
                                  _buildStatCard('Followers', '${_stats['followers']}'),
                                  _buildStatCard('Following', '${_stats['following']}'),
                                  _buildStatCard('Rating', '${(_stats['rating']! / 10).toStringAsFixed(1)}⭐'),
                                  _buildStatCard('Reviews', '${_stats['ratingCount']}'),
                                ],
                              ),
                            ],
                          ),
                        ),
                      ],
                    ),
                  ),

                  // Tabs
                  TabBar(
                    controller: _tabController,
                    labelColor: Theme.of(context).primaryColor,
                    unselectedLabelColor: Colors.grey,
                    indicatorColor: Theme.of(context).primaryColor,
                    tabs: const [
                      Tab(text: 'Activity'),
                      Tab(text: 'History'),
                      Tab(text: 'Wishlist'),
                      Tab(text: 'Settings'),
                    ],
                  ),
                ],
              ),
            ),
          ),

          // Tab Content
          SliverFillRemaining(
            child: TabBarView(
              controller: _tabController,
              children: [
                _buildActivityTab(),
                _buildHistoryTab(),
                _buildWishlistTab(),
                _buildSettingsTab(),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildStatCard(String label, String value) {
    return Container(
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(8),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.05),
            blurRadius: 4,
            offset: const Offset(0, 2),
          ),
        ],
      ),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          value.text.lg.bold.color(Theme.of(context).primaryColor).make(),
          4.heightBox,
          Text(
            label,
            textAlign: TextAlign.center,
            style: TextStyle(
              fontSize: 12,
              color: Colors.grey[600],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildActivityTab() {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          'Recent Activity'
              .text
              .xl
              .bold
              .make()
              .py(8),

          // Activity Chart
          Container(
            height: 200,
            padding: const EdgeInsets.all(16),
            decoration: BoxDecoration(
              color: Colors.grey[100],
              borderRadius: BorderRadius.circular(12),
            ),
            child: Column(
              children: [
                'Activity Overview (Last 30 Days)'.text.semiBold.make(),
                16.heightBox,
                Expanded(
                  child: LineChart(
                    LineChartData(
                      gridData: const FlGridData(show: false),
                      titlesData: const FlTitlesData(show: false),
                      borderData: FlBorderData(show: false),
                      lineBarsData: [
                        LineChartBarData(
                          spots: const [
                            FlSpot(0, 3),
                            FlSpot(1, 7),
                            FlSpot(2, 5),
                            FlSpot(3, 12),
                            FlSpot(4, 8),
                            FlSpot(5, 15),
                            FlSpot(6, 10),
                          ],
                          isCurved: true,
                          color: Theme.of(context).primaryColor,
                          barWidth: 3,
                        ),
                      ],
                    ),
                  ),
                ),
              ],
            ),
          ),

          24.heightBox,

          // Recent Activities List
          _buildActivityItem('Won auction', r'Vintage Watch - $250', '2 hours ago', Icons.gavel, Colors.green),
          _buildActivityItem('Started following', 'TechCollector', '5 hours ago', Icons.person_add, Colors.blue),
          _buildActivityItem('Left review', '★★★★★ - Great seller!', '1 day ago', Icons.star, Colors.amber),
          _buildActivityItem('Sold item', r'Rare Coins - $450', '2 days ago', Icons.attach_money, Colors.purple),
          _buildActivityItem('Added to wishlist', 'Modern Art Piece', '3 days ago', Icons.favorite, Colors.red),
        ],
      ),
    );
  }

  Widget _buildHistoryTab() {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          'Purchase History'
              .text
              .xl
              .bold
              .make()
              .py(8),

          // Won Auctions
          'Won Auctions'.text.lg.semiBold.make().py(4),
          _buildHistoryItem(
            'Vintage Camera',
            r'$1,250',
            'March 10, 2024',
            'Delivered',
            Colors.green,
          ),
          _buildHistoryItem(
            'Antique Vase',
            r'$875',
            'March 5, 2024',
            'In Transit',
            Colors.blue,
          ),

          16.heightBox,

          // Sold Items
          'Sold Items'.text.lg.semiBold.make().py(4),
          _buildHistoryItem(
            'Collectible Stamps',
            r'$450',
            'March 8, 2024',
            'Completed',
            Colors.green,
          ),
          _buildHistoryItem(
            'Designer Watch',
            r'$2,100',
            'March 1, 2024',
            'Completed',
            Colors.green,
          ),
        ],
      ),
    );
  }

  Widget _buildWishlistTab() {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          'My Wishlist'
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
            itemCount: 6, // Mock wishlist items
            itemBuilder: (context, index) {
              return _buildWishlistItem(index);
            },
          ),
        ],
      ),
    );
  }

  Widget _buildSettingsTab() {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          'Account Settings'
              .text
              .xl
              .bold
              .make()
              .py(8),

          _buildSettingTile('Email', _emailController.text, Icons.email, () {}),
          _buildSettingTile('Phone', '+1 (555) 123-4567', Icons.phone, () {}),
          _buildSettingTile('Language', 'English', Icons.language, () {}),
          _buildSettingTile('Currency', r'USD ($)', Icons.attach_money, () {}),

          24.heightBox,

          'Privacy Settings'
              .text
              .xl
              .bold
              .make()
              .py(8),

          _buildToggleTile('Public Profile', true, (value) {}),
          _buildToggleTile('Show Activity', true, (value) {}),
          _buildToggleTile('Email Notifications', true, (value) {}),
          _buildToggleTile('Push Notifications', false, (value) {}),

          24.heightBox,

          'Security'
              .text
              .xl
              .bold
              .make()
              .py(8),

          _buildSettingTile('Change Password', '••••••••', Icons.lock, () {}),
          _buildSettingTile('Two-Factor Auth', 'Enabled', Icons.security, () {}),
          _buildSettingTile('Login History', 'View', Icons.history, () {}),

          24.heightBox,

          GFButton(
            onPressed: _logout,
            text: 'Sign Out',
            color: Colors.red,
            fullWidthButton: true,
          ),
        ],
      ),
    );
  }

  Widget _buildActivityItem(String title, String description, String time, IconData icon, Color color) {
    return Container(
      margin: const EdgeInsets.only(bottom: 12),
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(8),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.05),
            blurRadius: 4,
            offset: const Offset(0, 2),
          ),
        ],
      ),
      child: Row(
        children: [
          Container(
            padding: const EdgeInsets.all(8),
            decoration: BoxDecoration(
              color: color.withOpacity(0.1),
              borderRadius: BorderRadius.circular(8),
            ),
            child: Icon(icon, color: color, size: 20),
          ),
          12.widthBox,
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                title.text.semiBold.make(),
                description.text.sm.color(Colors.grey[600]).make(),
                time.text.xs.color(Colors.grey[500]).make(),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildHistoryItem(String title, String price, String date, String status, Color statusColor) {
    return Container(
      margin: const EdgeInsets.only(bottom: 12),
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(8),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.05),
            blurRadius: 4,
            offset: const Offset(0, 2),
          ),
        ],
      ),
      child: Row(
        children: [
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                title.text.semiBold.make(),
                date.text.sm.color(Colors.grey[600]).make(),
              ],
            ),
          ),
          Column(
            crossAxisAlignment: CrossAxisAlignment.end,
            children: [
              price.text.lg.bold.color(Theme.of(context).primaryColor).make(),
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                decoration: BoxDecoration(
                  color: statusColor.withOpacity(0.1),
                  borderRadius: BorderRadius.circular(12),
                ),
                child: Text(
                  status,
                  style: TextStyle(
                    color: statusColor,
                    fontSize: 12,
                    fontWeight: FontWeight.w500,
                  ),
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildWishlistItem(int index) {
    final titles = ['Vintage Watch', 'Modern Art', 'Rare Book', 'Antique Lamp', 'Collectible Coin', 'Designer Bag'];
    final prices = [1200, 2500, 450, 680, 320, 1800];

    return GFCard(
      margin: EdgeInsets.zero,
      content: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Expanded(
            flex: 3,
            child: Container(
              decoration: BoxDecoration(
                color: Colors.grey[200],
                borderRadius: const BorderRadius.vertical(top: Radius.circular(8)),
              ),
              child: const Center(
                child: Icon(Icons.image, color: Colors.grey),
              ),
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
                    titles[index],
                    style: const TextStyle(
                      fontSize: 12,
                      fontWeight: FontWeight.bold,
                    ),
                    maxLines: 2,
                    overflow: TextOverflow.ellipsis,
                  ),
                  const Spacer(),
                  Text(
                    '\$${prices[index]}',
                    style: TextStyle(
                      fontSize: 12,
                      color: Theme.of(context).primaryColor,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildSettingTile(String title, String value, IconData icon, VoidCallback onTap) {
    return ListTile(
      contentPadding: EdgeInsets.zero,
      leading: Icon(icon),
      title: Text(title),
      subtitle: Text(value),
      trailing: const Icon(Icons.chevron_right),
      onTap: onTap,
    );
  }

  Widget _buildToggleTile(String title, bool value, Function(bool) onChanged) {
    return ListTile(
      contentPadding: EdgeInsets.zero,
      title: Text(title),
      trailing: GFToggle(
        value: value,
        onChanged: (bool? value) => onChanged(value ?? false),
      ),
    );
  }

  void _changeProfilePicture() {
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('Change profile picture feature coming soon!')),
    );
  }

  void _changeCoverImage() {
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('Change cover image feature coming soon!')),
    );
  }

  void _shareProfile() {
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('Share profile feature coming soon!')),
    );
  }

  void _viewPublicProfile() {
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('Public profile view coming soon!')),
    );
  }

  void _logout() {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Sign Out'),
        content: const Text('Are you sure you want to sign out?'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Cancel'),
          ),
          GFButton(
            onPressed: () {
              Navigator.pop(context);
              // Handle logout logic
              ScaffoldMessenger.of(context).showSnackBar(
                const SnackBar(content: Text('Signed out successfully')),
              );
            },
            text: 'Sign Out',
            color: Colors.red,
          ),
        ],
      ),
    );
  }
}