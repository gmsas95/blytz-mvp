import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:getwidget/getwidget.dart';
import 'package:velocity_x/velocity_x.dart';

class CategoriesPage extends ConsumerWidget {
  const CategoriesPage({super.key});

  final List<Map<String, dynamic>> _categories = const [
    {
      'name': 'Electronics',
      'icon': Icons.devices,
      'color': Colors.blue,
      'description': 'Gadgets, computers, smartphones, and more',
      'streamCount': 234,
      'featuredImage': 'https://example.com/electronics.jpg',
    },
    {
      'name': 'Fashion',
      'icon': Icons.checkroom,
      'color': Colors.pink,
      'description': 'Clothing, accessories, designer items',
      'streamCount': 456,
      'featuredImage': 'https://example.com/fashion.jpg',
    },
    {
      'name': 'Collectibles',
      'icon': Icons.collections,
      'color': Colors.amber,
      'description': 'Rare items, memorabilia, unique finds',
      'streamCount': 189,
      'featuredImage': 'https://example.com/collectibles.jpg',
    },
    {
      'name': 'Home & Garden',
      'icon': Icons.home,
      'color': Colors.green,
      'description': 'Furniture, decor, outdoor equipment',
      'streamCount': 167,
      'featuredImage': 'https://example.com/home.jpg',
    },
    {
      'name': 'Sports',
      'icon': Icons.sports_soccer,
      'color': Colors.orange,
      'description': 'Sporting goods, fitness equipment',
      'streamCount': 123,
      'featuredImage': 'https://example.com/sports.jpg',
    },
    {
      'name': 'Toys & Games',
      'icon': Icons.toys,
      'color': Colors.purple,
      'description': 'Collectible toys, board games, video games',
      'streamCount': 98,
      'featuredImage': 'https://example.com/toys.jpg',
    },
    {
      'name': 'Books',
      'icon': Icons.book,
      'color': Colors.brown,
      'description': 'Rare books, comics, literature',
      'streamCount': 76,
      'featuredImage': 'https://example.com/books.jpg',
    },
    {
      'name': 'Art',
      'icon': Icons.palette,
      'color': Colors.indigo,
      'description': 'Paintings, sculptures, digital art',
      'streamCount': 145,
      'featuredImage': 'https://example.com/art.jpg',
    },
    {
      'name': 'Jewelry',
      'icon': Icons.diamond,
      'color': Colors.teal,
      'description': 'Fine jewelry, watches, accessories',
      'streamCount': 89,
      'featuredImage': 'https://example.com/jewelry.jpg',
    },
    {
      'name': 'Vintage',
      'icon': Icons.watch,
      'color': Colors.grey,
      'description': 'Antiques, retro items, historical pieces',
      'streamCount': 234,
      'featuredImage': 'https://example.com/vintage.jpg',
    },
    {
      'name': 'Automotive',
      'icon': Icons.directions_car,
      'color': Colors.red,
      'description': 'Car parts, accessories, collectible vehicles',
      'streamCount': 67,
      'featuredImage': 'https://example.com/automotive.jpg',
    },
    {
      'name': 'Musical Instruments',
      'icon': Icons.music_note,
      'color': Colors.deepPurple,
      'description': 'Guitars, keyboards, vintage instruments',
      'streamCount': 112,
      'featuredImage': 'https://example.com/music.jpg',
    },
  ];

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Categories'),
        backgroundColor: Theme.of(context).primaryColor,
        foregroundColor: Colors.white,
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Header
            'Browse Categories'
                .text
                .xl2
                .bold
                .make()
                .py(8),

            "Find exactly what you're looking for in our curated categories"
                .text
                .lg
                .color(Colors.grey[600] ?? Colors.grey)
                .make()
                .py(8)
                .pOnly(bottom: 16),

            // Search Bar
            Container(
              decoration: BoxDecoration(
                color: Colors.grey[100],
                borderRadius: BorderRadius.circular(25),
              ),
              child: TextField(
                decoration: const InputDecoration(
                  labelText: 'Search categories...',
                  prefixIcon: Icon(Icons.search),
                  border: InputBorder.none,
                  contentPadding: EdgeInsets.symmetric(horizontal: 20, vertical: 15),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildFeaturedCategoryCard(Map<String, dynamic> category) {
    return Container(
      width: 160,
      margin: const EdgeInsets.only(right: 12),
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(16),
        color: (category['color'] as Color).withOpacity(0.1),
        border: Border.all(
          color: (category['color'] as Color).withOpacity(0.3),
        ),
      ),
      child: Stack(
        children: [
          // Background Image
          ClipRRect(
            borderRadius: BorderRadius.circular(16),
            child: Container(
              decoration: BoxDecoration(
                image: DecorationImage(
                  image: NetworkImage(category['featuredImage']),
                  fit: BoxFit.cover,
                  onError: (exception, stackTrace) {
                    // Fallback to solid color
                  },
                ),
              ),
              child: Container(
                decoration: BoxDecoration(
                  borderRadius: BorderRadius.circular(16),
                  gradient: LinearGradient(
                    begin: Alignment.topCenter,
                    end: Alignment.bottomCenter,
                    colors: [
                      Colors.transparent,
                      Colors.black.withOpacity(0.7),
                    ],
                  ),
                ),
              ),
            ),
          ),

          // Content
          Padding(
            padding: const EdgeInsets.all(12),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              mainAxisAlignment: MainAxisAlignment.end,
              children: [
                Container(
                  padding: const EdgeInsets.all(8),
                  decoration: BoxDecoration(
                    color: Colors.white,
                    borderRadius: BorderRadius.circular(8),
                  ),
                  child: Icon(
                    category['icon'],
                    color: category['color'],
                    size: 20,
                  ),
                ),

                const SizedBox(height: 8),

                Text(
                  category['name'],
                  style: const TextStyle(
                    color: Colors.white,
                    fontSize: 16,
                    fontWeight: FontWeight.bold,
                  ),
                ),

                const SizedBox(height: 4),

                Row(
                  children: [
                    const Icon(
                      Icons.live_tv,
                      color: Colors.white,
                      size: 12,
                    ),
                    const SizedBox(width: 4),
                    Text(
                      '${category['streamCount']} live',
                      style: const TextStyle(
                        color: Colors.white,
                        fontSize: 12,
                      ),
                    ),
                  ],
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildCategoryCard(Map<String, dynamic> category, BuildContext context) {
    return GestureDetector(
      onTap: () => _navigateToCategory(category['name'], context),
      child: Container(
        margin: EdgeInsets.zero,
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(
          borderRadius: BorderRadius.circular(12),
          color: (category['color'] as Color).withOpacity(0.1),
        ),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
          // Icon
          Container(
            padding: const EdgeInsets.all(16),
            decoration: BoxDecoration(
              color: Colors.white,
              borderRadius: BorderRadius.circular(16),
              boxShadow: [
                BoxShadow(
                  color: (category['color'] as Color).withOpacity(0.2),
                  blurRadius: 8,
                  offset: const Offset(0, 2),
                ),
              ],
            ),
            child: Icon(
              category['icon'],
              color: category['color'],
              size: 32,
            ),
          ),

          const SizedBox(height: 12),

          // Category Name
          Text(
            category['name'],
            textAlign: TextAlign.center,
            style: TextStyle(
              fontSize: 14,
              fontWeight: FontWeight.bold,
              color: Theme.of(context).brightness == Brightness.dark
                  ? Colors.white
                  : Colors.black87,
            ),
            maxLines: 1,
            overflow: TextOverflow.ellipsis,
          ),

          const SizedBox(height: 4),

          // Stream Count
          Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              Icon(
                Icons.live_tv,
                size: 12,
                color: category['color'],
              ),
              const SizedBox(width: 4),
              Text(
                '${category['streamCount']}',
                style: TextStyle(
                  fontSize: 12,
                  color: category['color'],
                  fontWeight: FontWeight.w500,
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildTrendingItem(String category, String stream, String viewers, BuildContext context) {
    return ListTile(
      contentPadding: EdgeInsets.zero,
      leading: Container(
        width: 40,
        height: 40,
        decoration: BoxDecoration(
          color: Theme.of(context).primaryColor.withOpacity(0.1),
          borderRadius: BorderRadius.circular(8),
        ),
        child: Icon(
          Icons.trending_up,
          color: Theme.of(context).primaryColor,
          size: 20,
        ),
      ),
      title: Text(
        stream,
        style: const TextStyle(fontWeight: FontWeight.bold),
      ),
      subtitle: Text(category),
      trailing: Text(
        viewers,
        style: TextStyle(
          color: Theme.of(context).primaryColor,
          fontWeight: FontWeight.w500,
          fontSize: 12,
        ),
      ),
    );
  }

  void _navigateToCategory(String categoryName, BuildContext context) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text('Browsing $categoryName...')),
    );
  }
}