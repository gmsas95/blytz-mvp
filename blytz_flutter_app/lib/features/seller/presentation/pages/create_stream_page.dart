import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:getwidget/getwidget.dart';
import 'package:image_picker/image_picker.dart';
import 'package:velocity_x/velocity_x.dart';

class CreateStreamPage extends ConsumerStatefulWidget {
  const CreateStreamPage({super.key});

  @override
  ConsumerState<CreateStreamPage> createState() => _CreateStreamPageState();
}

class _CreateStreamPageState extends ConsumerState<CreateStreamPage> {
  final _formKey = GlobalKey<FormState>();
  final _titleController = TextEditingController();
  final _descriptionController = TextEditingController();
  final _startingPriceController = TextEditingController();
  final _streamTitleController = TextEditingController();

  String? _selectedCategory;
  final List<String> _selectedProducts = [];
  final List<XFile> _selectedImages = [];
  DateTime? _scheduledTime;
  bool _isScheduled = false;

  final List<String> _categories = [
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

  @override
  void dispose() {
    _titleController.dispose();
    _descriptionController.dispose();
    _startingPriceController.dispose();
    _streamTitleController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Create Live Stream'),
        backgroundColor: Theme.of(context).primaryColor,
        foregroundColor: Colors.white,
        actions: [
          TextButton(
            onPressed: _createStream,
            child: const Text(
              'Go Live',
              style: TextStyle(color: Colors.white, fontWeight: FontWeight.bold),
            ),
          ),
        ],
      ),
      body: Form(
        key: _formKey,
        child: SingleChildScrollView(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // Stream Information
              'Stream Information'
                  .text
                  .xl
                  .bold
                  .make()
                  .py(8),

              GFCard(
                content: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    TextFormField(
                      controller: _streamTitleController,
                      decoration: const InputDecoration(
                        labelText: 'Stream Title',
                        border: OutlineInputBorder(),
                      ),
                      validator: (value) {
                        if (value == null || value.isEmpty) {
                          return 'Please enter a stream title';
                        }
                        return null;
                      },
                    ),
                    16.heightBox,

                    TextFormField(
                      controller: _descriptionController,
                      maxLines: 3,
                      decoration: const InputDecoration(
                        labelText: 'Stream Description',
                        border: OutlineInputBorder(),
                      ),
                      validator: (value) {
                        if (value == null || value.isEmpty) {
                          return 'Please enter a stream description';
                        }
                        return null;
                      },
                    ),
                    16.heightBox,

                    // Category Selection
                    'Category'
                        .text
                        .lg
                        .semiBold
                        .make()
                        .py(4),

                    GFDropdown(
                      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
                      borderRadius: BorderRadius.circular(8),
                      border: BorderSide(color: Colors.grey[300]!),
                      value: _selectedCategory,
                      hint: const Text('Select Category'),
                      items: _categories.map((category) {
                        return DropdownMenuItem(
                          value: category,
                          child: Text(category),
                        );
                      }).toList(),
                      onChanged: (value) {
                        setState(() {
                          _selectedCategory = value;
                        });
                      },
                    ),
                  ],
                ),
              ),

              24.heightBox,

              // Featured Product
              'Featured Product'
                  .text
                  .xl
                  .bold
                  .make()
                  .py(8),

              GFCard(
                content: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    TextFormField(
                      controller: _titleController,
                      decoration: const InputDecoration(
                        labelText: 'Product Title',
                      ),
                      validator: (value) {
                        if (value == null || value.isEmpty) {
                          return 'Please enter a product title';
                        }
                        return null;
                      },
                    ),
                    16.heightBox,

                    TextFormField(
                      controller: _startingPriceController,
                      keyboardType: const TextInputType.numberWithOptions(decimal: true),
                      decoration: const InputDecoration(
                        labelText: r'Starting Price ($)',
                      ),
                      validator: (value) {
                        if (value == null || value.isEmpty) {
                          return 'Please enter a starting price';
                        }
                        final price = double.tryParse(value);
                        if (price == null || price <= 0) {
                          return 'Please enter a valid price';
                        }
                        return null;
                      },
                    ),
                    16.heightBox,

                    // Image Upload
                    'Product Images'
                        .text
                        .lg
                        .semiBold
                        .make()
                        .py(4),

                    _buildImageUpload(),
                  ],
                ),
              ),

              24.heightBox,

              // Scheduling Options
              'Scheduling'
                  .text
                  .xl
                  .bold
                  .make()
                  .py(8),

              GFCard(
                content: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      children: [
                        GFCheckbox(
                          size: GFSize.SMALL,
                          type: GFCheckboxType.square,
                          value: _isScheduled,
                          onChanged: (value) {
                            setState(() {
                              _isScheduled = value ?? false;
                              if (!_isScheduled) {
                                _scheduledTime = null;
                              }
                            });
                          },
                        ),
                        12.widthBox,
                        'Schedule for later'
                            .text
                            .lg
                            .make(),
                      ],
                    ),

                    if (_isScheduled) ...[
                      16.heightBox,
                      GFButton(
                        onPressed: _selectDateTime,
                        color: Theme.of(context).primaryColor,
                        child: Text(
                          _scheduledTime != null
                              ? 'Scheduled: ${_scheduledTime!}'
                              : 'Select Date & Time',
                        ),
                      ),
                    ],
                  ],
                ),
              ),

              24.heightBox,

              // Stream Settings
              'Stream Settings'
                  .text
                  .xl
                  .bold
                  .make()
                  .py(8),

              GFCard(
                content: Column(
                  children: [
                    _buildSettingTile(
                      'Enable Chat',
                      'Allow viewers to chat during stream',
                      true,
                      (value) {},
                    ),
                    const Divider(),
                    _buildSettingTile(
                      'Record Stream',
                      'Save recording for later viewing',
                      false,
                      (value) {},
                    ),
                    const Divider(),
                    _buildSettingTile(
                      'Send Notifications',
                      'Notify followers when stream starts',
                      true,
                      (value) {},
                    ),
                  ],
                ),
              ),

              32.heightBox,

              // Preview Section
              'Preview'
                  .text
                  .xl
                  .bold
                  .make()
                  .py(8),

              GFCard(
                color: Colors.grey[100],
                content: Column(
                  children: [
                    // Thumbnail preview
                    Container(
                      height: 180,
                      decoration: BoxDecoration(
                        color: Colors.grey[300],
                        borderRadius: BorderRadius.circular(8),
                      ),
                      child: _selectedImages.isNotEmpty
                          ? ClipRRect(
                              borderRadius: BorderRadius.circular(8),
                              child: Image.network(
                                _selectedImages.first.path,
                                fit: BoxFit.cover,
                                errorBuilder: (context, error, stackTrace) {
                                  return const Center(
                                    child: Icon(Icons.image, color: Colors.grey, size: 48),
                                  );
                                },
                              ),
                            )
                          : const Center(
                              child: Column(
                                mainAxisAlignment: MainAxisAlignment.center,
                                children: [
                                  Icon(Icons.image, color: Colors.grey, size: 48),
                                  8.heightBox,
                                  Text('Stream Thumbnail', style: TextStyle(color: Colors.grey)),
                                ],
                              ),
                            ),
                    ),

                    16.heightBox,

                    // Stream info preview
                    Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        if (_streamTitleController.text.isNotEmpty)
                          _streamTitleController.text.text.xl.bold.make()
                        else
                          'Your Stream Title'.text.xl.bold.color(Colors.grey[600] ?? Colors.grey).make(),

                        4.heightBox,

                        if (_descriptionController.text.isNotEmpty)
                          _descriptionController.text.text.sm.color(Colors.grey[600] ?? Colors.grey).make()
                        else
                          'Stream description will appear here...'.text.sm.color(Colors.grey[600] ?? Colors.grey).make(),

                        8.heightBox,

                        Row(
                          children: [
                            if (_selectedCategory != null) ...[
                              Container(
                                padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                                decoration: BoxDecoration(
                                  color: Theme.of(context).primaryColor,
                                  borderRadius: BorderRadius.circular(12),
                                ),
                                child: _selectedCategory!.text.white.sm.make(),
                              ),
                              8.widthBox,
                            ],
                            const Icon(Icons.visibility, color: Colors.grey, size: 16),
                            4.widthBox,
                            '0 viewers'.text.color(Colors.grey[600] ?? Colors.grey).sm.make(),
                          ],
                        ),
                      ],
                    ),
                  ],
                ),
              ),

              32.heightBox,

              // Create Stream Button
              GFButton(
                onPressed: _createStream,
                color: Colors.red,
                size: GFSize.LARGE,
                fullWidthButton: true,
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    const Icon(Icons.live_tv, color: Colors.white),
                    8.widthBox,
                    (_isScheduled ? 'Schedule Stream' : 'Go Live Now')
                        .text
                        .white
                        .bold
                        .make(),
                  ],
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildImageUpload() {
    return Column(
      children: [
        // Selected images
        if (_selectedImages.isNotEmpty) ...[
          SizedBox(
            height: 100,
            child: ListView.builder(
              scrollDirection: Axis.horizontal,
              itemCount: _selectedImages.length,
              itemBuilder: (context, index) {
                final image = _selectedImages[index];
                return Container(
                  width: 100,
                  margin: const EdgeInsets.only(right: 8),
                  decoration: BoxDecoration(
                    borderRadius: BorderRadius.circular(8),
                    border: Border.all(color: Colors.grey[300]!),
                  ),
                  child: Stack(
                    children: [
                      ClipRRect(
                        borderRadius: BorderRadius.circular(8),
                        child: Image.network(
                          image.path,
                          width: 100,
                          height: 100,
                          fit: BoxFit.cover,
                          errorBuilder: (context, error, stackTrace) {
                            return Container(
                              width: 100,
                              height: 100,
                              color: Colors.grey[200],
                              child: const Icon(Icons.image, color: Colors.grey),
                            );
                          },
                        ),
                      ),
                      Positioned(
                        top: 4,
                        right: 4,
                        child: GestureDetector(
                          onTap: () {
                            setState(() {
                              _selectedImages.removeAt(index);
                            });
                          },
                          child: Container(
                            padding: const EdgeInsets.all(4),
                            decoration: const BoxDecoration(
                              color: Colors.red,
                              shape: BoxShape.circle,
                            ),
                            child: const Icon(Icons.close, color: Colors.white, size: 16),
                          ),
                        ),
                      ),
                    ],
                  ),
                );
              },
            ),
          ),
          12.heightBox,
        ],

        // Add image button
        GFButton(
          onPressed: _addImage,
          type: GFButtonType.outline2x,
          color: Theme.of(context).primaryColor,
          child: Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              const Icon(Icons.add_photo_alternate),
              8.widthBox,
              'Add Image'.text.make(),
            ],
          ),
        ),
      ],
    );
  }

  Widget _buildSettingTile(
    String title,
    String subtitle,
    bool value,
    Function(bool) onChanged,
  ) {
    return ListTile(
      contentPadding: EdgeInsets.zero,
      title: title.text.make(),
      subtitle: subtitle.text.sm.color(Colors.grey[600] ?? Colors.grey).make(),
      trailing: GFToggle(
        value: value,
        onChanged: onChanged,
      ),
    );
  }

  Future<void> _addImage() async {
    final picker = ImagePicker();
    final images = await picker.pickMultiImage();
    if (images.isNotEmpty) {
      setState(() {
        _selectedImages.addAll(images.take(5 - _selectedImages.length)); // Max 5 images
      });
    }
  }

  Future<void> _selectDateTime() async {
    final date = await showDatePicker(
      context: context,
      initialDate: DateTime.now().add(const Duration(hours: 1)),
      firstDate: DateTime.now(),
      lastDate: DateTime.now().add(const Duration(days: 30)),
    );

    if (date != null) {
      final time = await showTimePicker(
        context: context,
        initialTime: TimeOfDay.now(),
      );

      if (time != null) {
        setState(() {
          _scheduledTime = DateTime(
            date.year,
            date.month,
            date.day,
            time.hour,
            time.minute,
          );
        });
      }
    }
  }

  void _createStream() {
    if (!_formKey.currentState!.validate()) {
      return;
    }

    if (_selectedCategory == null) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('Please select a category'),
          backgroundColor: Colors.red,
        ),
      );
      return;
    }

    // Show success message
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text(_isScheduled ? 'Stream scheduled successfully!' : 'Going live now!'),
        backgroundColor: Colors.green,
      ),
    );

    // Navigate back or to live stream page
    Navigator.pop(context);
  }
}