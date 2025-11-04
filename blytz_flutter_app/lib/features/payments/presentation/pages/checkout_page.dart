import 'package:cached_network_image/cached_network_image.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:getwidget/getwidget.dart';
import 'package:velocity_x/velocity_x.dart';

class CheckoutPage extends ConsumerStatefulWidget {
  const CheckoutPage({
    required this.productTitle,
    required this.productImage,
    required this.price,
    required this.auctionId,
    super.key,
  });

  final String productTitle;
  final String productImage;
  final double price;
  final String auctionId;

  @override
  ConsumerState<CheckoutPage> createState() => _CheckoutPageState();
}

class _CheckoutPageState extends ConsumerState<CheckoutPage> {
  int _selectedPaymentMethod = 0;
  final _shippingController = TextEditingController();
  final _cardNumberController = TextEditingController();
  final _cardNameController = TextEditingController();
  final _cardExpiryController = TextEditingController();
  final _cardCvvController = TextEditingController();
  final _addressController = TextEditingController();
  final _cityController = TextEditingController();
  final _stateController = TextEditingController();
  final _zipController = TextEditingController();

  String _selectedShipping = 'standard';

  final List<Map<String, dynamic>> _paymentMethods = [
    {
      'name': 'Credit/Debit Card',
      'icon': Icons.credit_card,
      'type': 'card',
    },
    {
      'name': 'PayPal',
      'icon': Icons.account_balance_wallet,
      'type': 'paypal',
    },
    {
      'name': 'Apple Pay',
      'icon': Icons.apple,
      'type': 'apple',
    },
    {
      'name': 'Google Pay',
      'icon': Icons.g_mobiledata,
      'type': 'google',
    },
  ];

  final Map<String, Map<String, dynamic>> _shippingOptions = {
    'standard': {
      'name': 'Standard Shipping',
      'price': 9.99,
      'days': '5-7 business days',
    },
    'express': {
      'name': 'Express Shipping',
      'price': 19.99,
      'days': '2-3 business days',
    },
    'overnight': {
      'name': 'Overnight Shipping',
      'price': 39.99,
      'days': '1 business day',
    },
  };

  @override
  void dispose() {
    _shippingController.dispose();
    _cardNumberController.dispose();
    _cardNameController.dispose();
    _cardExpiryController.dispose();
    _cardCvvController.dispose();
    _addressController.dispose();
    _cityController.dispose();
    _stateController.dispose();
    _zipController.dispose();
    super.dispose();
  }

  double get _subtotal => widget.price;
  double get _shippingCost => (_shippingOptions[_selectedShipping]?['price'] as num?)?.toDouble() ?? 0.0;
  double get _tax => _subtotal * 0.08; // 8% tax
  double get _total => _subtotal + _shippingCost + _tax;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Checkout'),
        backgroundColor: Theme.of(context).primaryColor,
        foregroundColor: Colors.white,
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Product Summary
            'Order Summary'
                .text
                .xl
                .bold
                .make()
                .py(8),

            GFCard(
              content: Row(
                children: [
                  Container(
                    width: 80,
                    height: 80,
                    decoration: BoxDecoration(
                      color: Colors.grey[200],
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: ClipRRect(
                      borderRadius: BorderRadius.circular(8),
                      child: CachedNetworkImage(
                        imageUrl: widget.productImage,
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
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        widget.productTitle.text.semiBold.maxLines(2).overflow(TextOverflow.ellipsis).make(),
                        4.heightBox,
                        widget.price.currencyFormat.text.lg.color(Theme.of(context).primaryColor).bold.make(),
                        4.heightBox,
                        'Winning bid'.text.sm.color(Colors.grey[600]).make(),
                      ],
                    ),
                  ),
                ],
              ),
            ),

            24.heightBox,

            // Shipping Address
            'Shipping Address'
                .text
                .xl
                .bold
                .make()
                .py(8),

            GFCard(
              content: Column(
                children: [
                  TextFormField(
                    controller: _addressController,
                    decoration: const InputDecoration(
                      labelText: 'Street Address',
                      border: OutlineInputBorder(),
                    ),
                    validator: (value) {
                      if (value == null || value.isEmpty) {
                        return 'Please enter your address';
                      }
                      return null;
                    },
                  ),
                  12.heightBox,

                  Row(
                    children: [
                      Expanded(
                        child: TextFormField(
                          controller: _cityController,
                          decoration: const InputDecoration(
                            labelText: 'City',
                            border: OutlineInputBorder(),
                          ),
                          validator: (value) {
                            if (value == null || value.isEmpty) {
                              return 'City required';
                            }
                            return null;
                          },
                        ),
                      ),
                      12.widthBox,
                      Expanded(
                        child: TextFormField(
                          controller: _stateController,
                          decoration: const InputDecoration(
                            labelText: 'State',
                            border: OutlineInputBorder(),
                          ),
                          validator: (value) {
                            if (value == null || value.isEmpty) {
                              return 'State required';
                            }
                            return null;
                          },
                        ),
                      ),
                    ],
                  ),
                  12.heightBox,

                  TextFormField(
                    controller: _zipController,
                    keyboardType: TextInputType.number,
                    decoration: const InputDecoration(
                      labelText: 'ZIP Code',
                      border: OutlineInputBorder(),
                    ),
                    validator: (value) {
                      if (value == null || value.isEmpty) {
                        return 'ZIP code required';
                      }
                      return null;
                    },
                  ),
                ],
              ),
            ),

            24.heightBox,

            // Shipping Options
            'Shipping Method'
                .text
                .xl
                .bold
                .make()
                .py(8),

            GFCard(
              content: Column(
                children: _shippingOptions.entries.map((entry) {
                  final key = entry.key;
                  final option = entry.value;
                  final isSelected = _selectedShipping == key;

                  return Container(
                    margin: const EdgeInsets.only(bottom: 8),
                    child: InkWell(
                      onTap: () {
                        setState(() {
                          _selectedShipping = key;
                        });
                      },
                      child: Container(
                        padding: const EdgeInsets.all(12),
                        decoration: BoxDecoration(
                          border: Border.all(
                            color: isSelected ? Theme.of(context).primaryColor : Colors.grey[300]!,
                          ),
                          borderRadius: BorderRadius.circular(8),
                          color: isSelected ? Theme.of(context).primaryColor.withOpacity(0.1) : null,
                        ),
                        child: Row(
                          children: [
                            Icon(
                              isSelected ? Icons.radio_button_checked : Icons.radio_button_unchecked,
                              color: isSelected ? Theme.of(context).primaryColor : Colors.grey[600],
                            ),
                            12.widthBox,
                            Expanded(
                              child: Column(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                children: [
                                  (option['name']?.toString() ?? '').text.semiBold.make(),
                                  (option['days']?.toString() ?? '').text.sm.color(Colors.grey[600]).make(),
                                ],
                              ),
                            ),
                            '\$${(option['price'] as num? ?? 0).toStringAsFixed(2)}'.text.lg.bold.color(Theme.of(context).primaryColor).make(),
                          ],
                        ),
                      ),
                    ),
                  );
                }).toList(),
              ),
            ),

            24.heightBox,

            // Payment Method
            'Payment Method'
                .text
                .xl
                .bold
                .make()
                .py(8),

            GFCard(
              content: Column(
                children: [
                  // Payment Method Selection
                  ...List.generate(_paymentMethods.length, (index) {
                    final method = _paymentMethods[index];
                    final isSelected = _selectedPaymentMethod == index;

                    return Container(
                      margin: const EdgeInsets.only(bottom: 8),
                      child: InkWell(
                        onTap: () {
                          setState(() {
                            _selectedPaymentMethod = index;
                          });
                        },
                        child: Container(
                          padding: const EdgeInsets.all(12),
                          decoration: BoxDecoration(
                            border: Border.all(
                              color: isSelected ? Theme.of(context).primaryColor : Colors.grey[300]!,
                            ),
                            borderRadius: BorderRadius.circular(8),
                            color: isSelected ? Theme.of(context).primaryColor.withOpacity(0.1) : null,
                          ),
                          child: Row(
                            children: [
                              Icon(
                                isSelected ? Icons.radio_button_checked : Icons.radio_button_unchecked,
                                color: isSelected ? Theme.of(context).primaryColor : Colors.grey[600],
                              ),
                              12.widthBox,
                              Icon(method['icon'] as IconData? ?? Icons.payment),
                              12.widthBox,
                              (method['name']?.toString() ?? '').text.semiBold.make(),
                            ],
                          ),
                        ),
                      ),
                    );
                  }),

                  // Card Details (if card selected)
                  if (_selectedPaymentMethod == 0) ...[
                    16.heightBox,
                    TextFormField(
                      controller: _cardNumberController,
                      keyboardType: TextInputType.number,
                      maxLength: 16,
                      decoration: const InputDecoration(
                        labelText: 'Card Number',
                        border: OutlineInputBorder(),
                      ),
                      validator: (value) {
                        if (value == null || value.isEmpty) {
                          return 'Card number required';
                        }
                        if (value.length != 16) {
                          return 'Invalid card number';
                        }
                        return null;
                      },
                    ),
                    12.heightBox,
                    TextFormField(
                      controller: _cardNameController,
                      decoration: const InputDecoration(
                        labelText: 'Cardholder Name',
                        border: OutlineInputBorder(),
                      ),
                      validator: (value) {
                        if (value == null || value.isEmpty) {
                          return 'Cardholder name required';
                        }
                        return null;
                      },
                    ),
                    12.heightBox,
                    Row(
                      children: [
                        Expanded(
                          child: TextFormField(
                            controller: _cardExpiryController,
                            maxLength: 5,
                            decoration: const InputDecoration(
                              labelText: 'MM/YY',
                              border: OutlineInputBorder(),
                            ),
                            validator: (value) {
                              if (value == null || value.isEmpty) {
                                return 'Expiry date required';
                              }
                              return null;
                            },
                          ),
                        ),
                        12.widthBox,
                        Expanded(
                          child: TextFormField(
                            controller: _cardCvvController,
                            keyboardType: TextInputType.number,
                            maxLength: 3,
                            decoration: const InputDecoration(
                              labelText: 'CVV',
                              border: OutlineInputBorder(),
                            ),
                            validator: (value) {
                              if (value == null || value.isEmpty) {
                                return 'CVV required';
                              }
                              return null;
                            },
                          ),
                        ),
                      ],
                    ),
                  ],
                ],
              ),
            ),

            24.heightBox,

            // Order Summary (Cost Breakdown)
            'Order Total'
                .text
                .xl
                .bold
                .make()
                .py(8),

            GFCard(
              content: Column(
                children: [
                  _buildSummaryRow('Subtotal', _subtotal.currencyFormat),
                  _buildSummaryRow('Shipping', _shippingCost.currencyFormat),
                  _buildSummaryRow('Tax', _tax.currencyFormat),
                  const Divider(),
                  _buildSummaryRow(
                    'Total',
                    _total.currencyFormat,
                    isBold: true,
                    color: Theme.of(context).primaryColor,
                  ),
                ],
              ),
            ),

            32.heightBox,

            // Place Order Button
            GFButton(
              onPressed: _placeOrder,
              color: Theme.of(context).primaryColor,
              size: GFSize.LARGE,
              fullWidthButton: true,
              child: Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  const Icon(Icons.lock, color: Colors.white, size: 16),
                  8.widthBox,
                  'Place Order • ${_total.currencyFormat}'
                      .text
                      .white
                      .bold
                      .make(),
                ],
              ),
            ),

            16.heightBox,

            // Security Note
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                const Icon(Icons.security, color: Colors.grey, size: 16),
                4.widthBox,
                'Secure checkout powered by Blytz Payments'.text.sm.color(Colors.grey[600]).make(),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildSummaryRow(String label, String value, {bool isBold = false, Color? color}) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(
            label,
            style: TextStyle(
              fontWeight: isBold ? FontWeight.bold : FontWeight.normal,
            ),
          ),
          Text(
            value,
            style: TextStyle(
              fontWeight: isBold ? FontWeight.bold : FontWeight.normal,
              color: color,
            ),
          ),
        ],
      ),
    );
  }

  void _placeOrder() {
    // Validate form
    if (_addressController.text.isEmpty ||
        _cityController.text.isEmpty ||
        _stateController.text.isEmpty ||
        _zipController.text.isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('Please complete shipping address'),
          backgroundColor: Colors.red,
        ),
      );
      return;
    }

    if (_selectedPaymentMethod == 0 &&
        (_cardNumberController.text.isEmpty ||
            _cardNameController.text.isEmpty ||
            _cardExpiryController.text.isEmpty ||
            _cardCvvController.text.isEmpty)) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('Please complete payment details'),
          backgroundColor: Colors.red,
        ),
      );
      return;
    }

    // Show processing dialog
    showDialog(
      context: context,
      barrierDismissible: false,
      builder: (context) => const AlertDialog(
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            GFLoader(type: GFLoaderType.circle),
            SizedBox(height: 16),
            Text('Processing your order...'),
          ],
        ),
      ),
    );

    // Simulate payment processing
    Future.delayed(const Duration(seconds: 3), () {
      Navigator.pop(context); // Close processing dialog

      // Show success
      showDialog(
        context: context,
        builder: (context) => AlertDialog(
          title: Row(
            children: [
              const Icon(Icons.check_circle, color: Colors.green),
              8.widthBox,
              const Text('Order Confirmed!'),
            ],
          ),
          content: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              'Your order has been successfully placed!'.text.make(),
              16.heightBox,
              'Order ID: #${DateTime.now().millisecondsSinceEpoch}'.text.make(),
              'Estimated delivery: ${_shippingOptions[_selectedShipping]?['days']}'.text.make(),
            ],
          ),
          actions: [
            GFButton(
              onPressed: () {
                Navigator.pop(context); // Close success dialog
                Navigator.pop(context); // Go back to previous screen
              },
              text: 'View Orders',
              color: Theme.of(context).primaryColor,
            ),
          ],
        ),
      );
    });
  }
}

extension CurrencyFormat on double {
  String get currencyFormat => '\$${toStringAsFixed(2)}';
}