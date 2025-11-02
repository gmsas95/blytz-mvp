import 'package:form_validator/form_validator.dart';
import '../constants/app_constants.dart';
import '../errors/exceptions.dart';

class Validators {
  // Email validation
  static String? validateEmail(String? value) {
    if (value == null || value.isEmpty) {
      return 'Email is required';
    }
    
    if (!ValidationBuilder().email().build()(value)) {
      return 'Please enter a valid email address';
    }
    
    return null;
  }

  // Password validation
  static String? validatePassword(String? value) {
    if (value == null || value.isEmpty) {
      return 'Password is required';
    }
    
    if (value.length < 8) {
      return 'Password must be at least 8 characters long';
    }
    
    if (!value.contains(RegExp(r'[A-Z]'))) {
      return 'Password must contain at least one uppercase letter';
    }
    
    if (!value.contains(RegExp(r'[a-z]'))) {
      return 'Password must contain at least one lowercase letter';
    }
    
    if (!value.contains(RegExp(r'[0-9]'))) {
      return 'Password must contain at least one number';
    }
    
    if (!value.contains(RegExp(r'[!@#$%^&*(),.?":{}|<>]'))) {
      return 'Password must contain at least one special character';
    }
    
    return null;
  }

  // Confirm password validation
  static String? validateConfirmPassword(String? value, String password) {
    if (value == null || value.isEmpty) {
      return 'Please confirm your password';
    }
    
    if (value != password) {
      return 'Passwords do not match';
    }
    
    return null;
  }

  // Name validation
  static String? validateName(String? value) {
    if (value == null || value.isEmpty) {
      return 'Name is required';
    }
    
    if (value.length < 2) {
      return 'Name must be at least 2 characters long';
    }
    
    if (value.length > 50) {
      return 'Name must be less than 50 characters';
    }
    
    if (!RegExp(r'^[a-zA-Z\s]+$').hasMatch(value)) {
      return 'Name can only contain letters and spaces';
    }
    
    return null;
  }

  // Phone number validation
  static String? validatePhone(String? value) {
    if (value == null || value.isEmpty) {
      return 'Phone number is required';
    }
    
    if (!RegExp(r'^\+?[\d\s-()]+$').hasMatch(value)) {
      return 'Please enter a valid phone number';
    }
    
    if (value.replaceAll(RegExp(r'[^\d]'), '').length < 10) {
      return 'Phone number must be at least 10 digits';
    }
    
    return null;
  }

  // Bid amount validation
  static String? validateBidAmount(String? value, double currentBid) {
    if (value == null || value.isEmpty) {
      return 'Bid amount is required';
    }
    
    final amount = double.tryParse(value);
    if (amount == null) {
      return 'Please enter a valid amount';
    }
    
    if (amount <= currentBid) {
      return 'Bid must be higher than current bid of \$${currentBid.toStringAsFixed(2)}';
    }
    
    if (amount < AppConstants.minBidAmount) {
      return 'Minimum bid amount is \$${AppConstants.minBidAmount.toStringAsFixed(2)}';
    }
    
    if (amount > AppConstants.maxBidAmount) {
      return 'Maximum bid amount is \$${AppConstants.maxBidAmount.toStringAsFixed(2)}';
    }
    
    return null;
  }

  // Auction title validation
  static String? validateAuctionTitle(String? value) {
    if (value == null || value.isEmpty) {
      return 'Auction title is required';
    }
    
    if (value.length < 5) {
      return 'Title must be at least 5 characters long';
    }
    
    if (value.length > 100) {
      return 'Title must be less than 100 characters';
    }
    
    return null;
  }

  // Auction description validation
  static String? validateAuctionDescription(String? value) {
    if (value == null || value.isEmpty) {
      return 'Description is required';
    }
    
    if (value.length < 20) {
      return 'Description must be at least 20 characters long';
    }
    
    if (value.length > 2000) {
      return 'Description must be less than 2000 characters';
    }
    
    return null;
  }

  // Starting price validation
  static String? validateStartingPrice(String? value) {
    if (value == null || value.isEmpty) {
      return 'Starting price is required';
    }
    
    final price = double.tryParse(value);
    if (price == null) {
      return 'Please enter a valid price';
    }
    
    if (price < AppConstants.minBidAmount) {
      return 'Starting price must be at least \$${AppConstants.minBidAmount.toStringAsFixed(2)}';
    }
    
    if (price > AppConstants.maxBidAmount) {
      return 'Starting price must be less than \$${AppConstants.maxBidAmount.toStringAsFixed(2)}';
    }
    
    return null;
  }

  // Category validation
  static String? validateCategory(String? value) {
    if (value == null || value.isEmpty) {
      return 'Please select a category';
    }
    
    return null;
  }

  // Required field validation
  static String? validateRequired(String? value, String fieldName) {
    if (value == null || value.isEmpty) {
      return '$fieldName is required';
    }
    
    return null;
  }

  // Optional field validation
  static String? validateOptional(String? value, {int? minLength, int? maxLength}) {
    if (value == null || value.isEmpty) {
      return null; // Optional field can be empty
    }
    
    if (minLength != null && value.length < minLength) {
      return 'Must be at least $minLength characters long';
    }
    
    if (maxLength != null && value.length > maxLength) {
      return 'Must be less than $maxLength characters';
    }
    
    return null;
  }

  // URL validation
  static String? validateUrl(String? value) {
    if (value == null || value.isEmpty) {
      return null; // Optional
    }
    
    if (!Uri.tryParse(value)?.hasAbsolutePath ?? false) {
      return 'Please enter a valid URL';
    }
    
    return null;
  }

  // Price validation (general)
  static String? validatePrice(String? value, {double? min, double? max}) {
    if (value == null || value.isEmpty) {
      return 'Price is required';
    }
    
    final price = double.tryParse(value);
    if (price == null) {
      return 'Please enter a valid price';
    }
    
    if (price <= 0) {
      return 'Price must be greater than 0';
    }
    
    if (min != null && price < min) {
      return 'Price must be at least \$${min.toStringAsFixed(2)}';
    }
    
    if (max != null && price > max) {
      return 'Price must be less than \$${max.toStringAsFixed(2)}';
    }
    
    return null;
  }
}