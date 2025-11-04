import 'dart:async';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:getwidget/getwidget.dart';
import 'package:velocity_x/velocity_x.dart';
import 'package:web_socket_channel/web_socket_channel.dart';

class LiveStreamPage extends ConsumerStatefulWidget {
  const LiveStreamPage({
    required this.streamId,
    required this.sellerName,
    required this.productTitle,
    super.key,
  });

  final String streamId;
  final String sellerName;
  final String productTitle;

  @override
  ConsumerState<LiveStreamPage> createState() => _LiveStreamPageState();
}

class _LiveStreamPageState extends ConsumerState<LiveStreamPage>
    with TickerProviderStateMixin {
  late TabController _tabController;
  WebSocketChannel? _channel;
  bool _isConnected = false;
  int _viewerCount = 0;
  bool _isFollowing = false;

  // Video state
  final bool _isVideoEnabled = true;
  bool _isAudioEnabled = true;
  bool _isFullscreen = false;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 3, vsync: this);
    _initializeWebSocket();
    _initializeLiveKit();
  }

  @override
  void dispose() {
    _tabController.dispose();
    _channel?.sink.close();
    super.dispose();
  }

  void _initializeWebSocket() {
    try {
      _channel = WebSocketChannel.connect(
        Uri.parse('wss://api.blytz.app/ws/live/${widget.streamId}'),
      );

      _channel.stream.listen(
        (data) {
          try {
            final message = data as Map<String, dynamic>;
            setState(() {
              _viewerCount = (message['viewerCount'] as int?) ?? _viewerCount;
              // Handle other live updates (bids, chat messages, etc.)
            });
          } catch (e) {
            // Handle malformed messages
            print('Error parsing WebSocket message: $e');
          }
        },
        onError: (error) {
          print('WebSocket error: $error');
          // Don't crash the app, just continue without WebSocket updates
        },
        onDone: () {
          print('WebSocket connection closed');
        },
      );
    } catch (e) {
      print('Failed to connect to WebSocket: $e');
      // Don't crash the app, continue without WebSocket updates
      // The app will function without live updates, and can retry connection
      setState(() {
        _isConnected = false;
        _channel = null;
      });
    }
  }

  Future<void> _initializeLiveKit() async {
    try {
      // Connect to LiveKit room for video streaming
      // Implementation would depend on your LiveKit server setup
      setState(() {
        _isConnected = true;
      });
    } catch (e) {
      _showError('Failed to connect to stream: $e');
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.black,
      body: Column(
        children: [
          // Video Stream Section
          Expanded(
            flex: 3,
            child: Stack(
              children: [
                // Video Player
                Container(
                  width: double.infinity,
                  color: Colors.black,
                  child: _buildVideoPlayer(),
                ),

                // Top Overlay
                Positioned(
                  top: 0,
                  left: 0,
                  right: 0,
                  child: _buildTopOverlay(),
                ),

                // Bottom Controls
                Positioned(
                  bottom: 0,
                  left: 0,
                  right: 0,
                  child: _buildBottomControls(),
                ),

                // Side Actions
                Positioned(
                  right: 16,
                  top: 100,
                  bottom: 100,
                  child: _buildSideActions(),
                ),
              ],
            ),
          ),

          // Content Tabs
          Expanded(
            flex: 2,
            child: Container(
              decoration: const BoxDecoration(
                color: Colors.white,
                borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
              ),
              child: Column(
                children: [
                  _buildTabBar(),
                  Expanded(child: _buildTabContent()),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildVideoPlayer() {
    if (!_isConnected) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Theme(
              data: ThemeData.dark(),
              child: GFLoader(
                type: GFLoaderType.circle,
                size: GFSize.LARGE,
              ),
            ),
            SizedBox(height: 16),
            Text(
              'Connecting to stream...',
              style: TextStyle(color: Colors.white),
            ),
          ],
        ),
      );
    }

    // Placeholder for actual video player
    return Container(
      color: Colors.grey[900],
      child: const Center(
        child: Icon(
          Icons.play_circle_filled,
          size: 80,
          color: Colors.white54,
        ),
      ),
    );
  }

  Widget _buildTopOverlay() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        gradient: LinearGradient(
          begin: Alignment.topCenter,
          end: Alignment.bottomCenter,
          colors: [Colors.black.withOpacity(0.7), Colors.transparent],
        ),
      ),
      child: Row(
        children: [
          // Back Button
          GFIconButton(
            icon: const Icon(Icons.arrow_back, color: Colors.white),
            onPressed: () => Navigator.pop(context),
            type: GFButtonType.transparent,
          ),

          // Stream Info
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  widget.sellerName,
                  style: const TextStyle(
                    color: Colors.white,
                    fontSize: 16,
                    fontWeight: FontWeight.bold,
                  ),
                ),
                Text(
                  widget.productTitle,
                  style: const TextStyle(
                    color: Colors.white70,
                    fontSize: 14,
                  ),
                ),
              ],
            ),
          ),

          // Viewer Count
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
            decoration: BoxDecoration(
              color: Colors.black.withOpacity(0.5),
              borderRadius: BorderRadius.circular(20),
            ),
            child: Row(
              mainAxisSize: MainAxisSize.min,
              children: [
                const Icon(Icons.visibility, color: Colors.white, size: 16),
                const SizedBox(width: 4),
                Text(
                  _viewerCount.toString(),
                  style: const TextStyle(color: Colors.white, fontSize: 12),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildBottomControls() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        gradient: LinearGradient(
          begin: Alignment.bottomCenter,
          end: Alignment.topCenter,
          colors: [Colors.black.withOpacity(0.7), Colors.transparent],
        ),
      ),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          // Audio Toggle
          GFIconButton(
            icon: Icon(
              _isAudioEnabled ? Icons.mic : Icons.mic_off,
              color: Colors.white,
            ),
            onPressed: () {
              setState(() {
                _isAudioEnabled = !_isAudioEnabled;
              });
            },
            type: GFButtonType.transparent,
          ),

          const SizedBox(width: 20),

          // Fullscreen Toggle
          GFIconButton(
            icon: Icon(
              _isFullscreen ? Icons.fullscreen_exit : Icons.fullscreen,
              color: Colors.white,
            ),
            onPressed: () {
              setState(() {
                _isFullscreen = !_isFullscreen;
              });
            },
            type: GFButtonType.transparent,
          ),
        ],
      ),
    );
  }

  Widget _buildSideActions() {
    return Column(
      mainAxisAlignment: MainAxisAlignment.spaceEvenly,
      children: [
        // Follow Button
        GFButtonBadge(
          text: _isFollowing ? 'Following' : 'Follow',
          icon: Icon(_isFollowing ? Icons.check : Icons.add),
          color: _isFollowing ? Colors.grey : Colors.pink,
          onPressed: () {
            setState(() {
              _isFollowing = !_isFollowing;
            });
          },
          shape: GFButtonShape.pills,
          type: GFButtonType.outline2x,
        ),

        // Share Button
        Container(
          decoration: BoxDecoration(
            color: Colors.black.withOpacity(0.5),
            shape: BoxShape.circle,
          ),
          child: GFIconButton(
            icon: const Icon(Icons.share, color: Colors.white),
            onPressed: _shareStream,
            type: GFButtonType.transparent,
          ),
        ),

        // Product Info Button
        Container(
          decoration: BoxDecoration(
            color: Colors.black.withOpacity(0.5),
            shape: BoxShape.circle,
          ),
          child: GFIconButton(
            icon: const Icon(Icons.info_outline, color: Colors.white),
            onPressed: _showProductInfo,
            type: GFButtonType.transparent,
          ),
        ),
      ],
    );
  }

  Widget _buildTabBar() {
    return TabBar(
      controller: _tabController,
      labelColor: Theme.of(context).primaryColor,
      unselectedLabelColor: Colors.grey,
      indicatorColor: Theme.of(context).primaryColor,
      tabs: const [
        Tab(text: 'Details'),
        Tab(text: 'Reviews'),
        Tab(text: 'Similar'),
      ],
    );
  }

  Widget _buildTabContent() {
    return TabBarView(
      controller: _tabController,
      children: [
        _buildDetailsTab(),
        _buildReviewsTab(),
        _buildSimilarTab(),
      ],
    );
  }

  Widget _buildDetailsTab() {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Product Details
          Text(
            'Product Details',
            style: Theme.of(context).textTheme.headlineSmall?.copyWith(
              fontWeight: FontWeight.bold,
            ),
          ),
          const SizedBox(height: 8),

          // Description
          const Text(
            'This is an amazing product available for live bidding. High quality, authentic, and ready to ship!',
            style: TextStyle(color: Colors.grey),
          ),

          const SizedBox(height: 16),

          // Specifications
          Text(
            'Specifications',
            style: Theme.of(context).textTheme.titleMedium?.copyWith(
              fontWeight: FontWeight.bold,
            ),
          ),
          const SizedBox(height: 8),

          _buildSpecificationRow('Condition', 'Brand New'),
          _buildSpecificationRow('Brand', 'Premium Brand'),
          _buildSpecificationRow('Shipping', '2-3 business days'),
          _buildSpecificationRow('Returns', '30 day return policy'),

          const SizedBox(height: 20),

          // Bid Now Button
          GFButton(
            onPressed: _placeBid,
            color: Theme.of(context).primaryColor,
            size: GFSize.LARGE,
            fullWidthButton: true,
            child: 'Place Bid'.text.white.bold.make(),
          ),
        ],
      ),
    );
  }

  Widget _buildSpecificationRow(String label, String value) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4),
      child: Row(
        children: [
          Text('$label: ', style: const TextStyle(fontWeight: FontWeight.bold)),
          Text(value, style: const TextStyle(color: Colors.grey)),
        ],
      ),
    );
  }

  Widget _buildReviewsTab() {
    return const Center(
      child: Text(
        'No reviews yet. Be the first to review!',
        style: TextStyle(color: Colors.grey),
      ),
    );
  }

  Widget _buildSimilarTab() {
    return GridView.builder(
      padding: const EdgeInsets.all(16),
      gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
        crossAxisCount: 2,
        childAspectRatio: 0.8,
        crossAxisSpacing: 12,
        mainAxisSpacing: 12,
      ),
      itemCount: 4, // Mock similar products
      itemBuilder: (context, index) {
        return Container(
          decoration: BoxDecoration(
            color: Colors.grey[100],
            borderRadius: BorderRadius.circular(12),
          ),
          child: Column(
            children: [
              Expanded(
                flex: 3,
                child: Container(
                  decoration: BoxDecoration(
                    color: Colors.grey[300],
                    borderRadius: const BorderRadius.vertical(
                      top: Radius.circular(12),
                    ),
                  ),
                  child: const Center(
                    child: Icon(Icons.image, color: Colors.grey),
                  ),
                ),
              ),
              Expanded(
                child: Padding(
                  padding: const EdgeInsets.all(8),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        'Similar Product ${index + 1}',
                        style: const TextStyle(
                          fontSize: 12,
                          fontWeight: FontWeight.bold,
                        ),
                        maxLines: 2,
                        overflow: TextOverflow.ellipsis,
                      ),
                      const Spacer(),
                      Text(
                        '\$${(index + 1) * 25}',
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
      },
    );
  }

  void _shareStream() {
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('Share feature coming soon!')),
    );
  }

  void _showProductInfo() {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Product Information'),
        content: const Text('Detailed product information will be displayed here.'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Close'),
          ),
        ],
      ),
    );
  }

  void _placeBid() {
    // Navigate to bidding interface
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      builder: (context) => Container(
        height: MediaQuery.of(context).size.height * 0.7,
        decoration: const BoxDecoration(
          borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
        ),
        child: Column(
          children: [
            // Handle bar
            Container(
              width: 40,
              height: 4,
              margin: const EdgeInsets.symmetric(vertical: 8),
              decoration: BoxDecoration(
                color: Colors.grey[300],
                borderRadius: BorderRadius.circular(2),
              ),
            ),

            // Bidding content would go here
            Expanded(
              child: Center(
                child: Text(
                  'Bidding Interface\nComing Soon!',
                  textAlign: TextAlign.center,
                  style: Theme.of(context).textTheme.headlineSmall,
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  void _showError(String message) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text(message),
        backgroundColor: Colors.red,
      ),
    );
  }
}