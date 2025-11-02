import 'dart:math' as math;
import 'package:flutter/material.dart';

extension StringExtensions on String {
  bool get isValidEmail {
    return RegExp(r'^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$').hasMatch(this);
  }

  bool get isValidPhone {
    return RegExp(r'^\+?[\d\s-()]+$').hasMatch(this);
  }

  bool get isValidUrl {
    try {
      final uri = Uri.parse(this);
      return uri.hasAbsolutePath && (uri.scheme == 'http' || uri.scheme == 'https');
    } catch (e) {
      return false;
    }
  }

  String get capitalize {
    if (isEmpty) return this;
    return '${this[0].toUpperCase()}${substring(1).toLowerCase()}';
  }

  String get titleCase {
    return split(' ').map((word) => word.capitalize).join(' ');
  }

  String get truncate {
    if (length <= 50) return this;
    return '${substring(0, 50)}...';
  }

  String truncateTo(int maxLength) {
    if (length <= maxLength) return this;
    return '${substring(0, maxLength)}...';
  }

  String get removeWhitespace {
    return replaceAll(RegExp(r'\s+'), '');
  }

  String get removeSpecialCharacters {
    return replaceAll(RegExp(r'[^\w\s]'), '');
  }

  bool get isNumeric {
    return double.tryParse(this) != null;
  }

  String get currencyFormat {
    final amount = double.tryParse(this);
    if (amount == null) return this;
    return '\$${amount.toStringAsFixed(2)}';
  }
}

extension DateTimeExtensions on DateTime {
  bool get isToday {
    final now = DateTime.now();
    return year == now.year && month == now.month && day == now.day;
  }

  bool get isYesterday {
    final yesterday = DateTime.now().subtract(const Duration(days: 1));
    return year == yesterday.year && month == yesterday.month && day == yesterday.day;
  }

  bool get isTomorrow {
    final tomorrow = DateTime.now().add(const Duration(days: 1));
    return year == tomorrow.year && month == tomorrow.month && day == tomorrow.day;
  }

  String get timeAgo {
    final now = DateTime.now();
    final difference = now.difference(this);

    if (difference.inDays > 0) {
      return '${difference.inDays}d ago';
    } else if (difference.inHours > 0) {
      return '${difference.inHours}h ago';
    } else if (difference.inMinutes > 0) {
      return '${difference.inMinutes}m ago';
    } else {
      return 'Just now';
    }
  }

  String get formattedDate {
    return '${day.toString().padLeft(2, '0')}/${month.toString().padLeft(2, '0')}/$year';
  }

  String get formattedTime {
    return '${hour.toString().padLeft(2, '0')}:${minute.toString().padLeft(2, '0')}';
  }

  String get formattedDateTime {
    return '$formattedDate $formattedTime';
  }

  DateTime get startOfDay {
    return DateTime(year, month, day);
  }

  DateTime get endOfDay {
    return DateTime(year, month, day, 23, 59, 59, 999);
  }

  bool get isWeekend {
    return weekday == DateTime.saturday || weekday == DateTime.sunday;
  }

  int get daysInMonth {
    return DateTime(year, month + 1, 0).day;
  }
}

extension IntExtensions on int {
  String get formatNumber {
    if (this >= 1000000) {
      return '${(this / 1000000).toStringAsFixed(1)}M';
    } else if (this >= 1000) {
      return '${(this / 1000).toStringAsFixed(1)}K';
    } else {
      return toString();
    }
  }

  Duration get milliseconds {
    return Duration(milliseconds: this);
  }

  Duration get seconds {
    return Duration(seconds: this);
  }

  Duration get minutes {
    return Duration(minutes: this);
  }

  Duration get hours {
    return Duration(hours: this);
  }

  Duration get days {
    return Duration(days: this);
  }
}

extension DoubleExtensions on double {
  String get currencyFormat {
    return '\$${toStringAsFixed(2)}';
  }

  String get compactCurrency {
    if (this >= 1000000) {
      return '\$${(this / 1000000).toStringAsFixed(1)}M';
    } else if (this >= 1000) {
      return '\$${(this / 1000).toStringAsFixed(1)}K';
    } else {
      return currencyFormat;
    }
  }

  String get percentageFormat {
    return '${(this * 100).toStringAsFixed(1)}%';
  }

  double roundToDecimalPlaces(int places) {
    final factor = math.pow(10, places);
    return (this * factor).round() / factor;
  }
}

extension ListExtensions<T> on List<T> {
  T? get firstOrNull {
    return isEmpty ? null : first;
  }

  T? get lastOrNull {
    return isEmpty ? null : last;
  }

  List<T> get unique {
    final seen = <T>{};
    return where((element) => seen.add(element)).toList();
  }

  List<T> get reversedList {
    return reversed.toList();
  }

  T? elementAtOrNull(int index) {
    if (index < 0 || index >= length) return null;
    return this[index];
  }

  List<List<T>> chunk(int size) {
    final chunks = <List<T>>[];
    for (var i = 0; i < length; i += size) {
      chunks.add(sublist(i, (i + size).clamp(0, length)));
    }
    return chunks;
  }
}

extension BuildContextExtensions on BuildContext {
  ThemeData get theme => Theme.of(this);
  
  ColorScheme get colorScheme => Theme.of(this).colorScheme;
  
  TextTheme get textTheme => Theme.of(this).textTheme;
  
  MediaQueryData get mediaQuery => MediaQuery.of(this);
  
  Size get size => mediaQuery.size;
  
  double get width => size.width;
  
  double get height => size.height;
  
  bool get isKeyboardOpen => mediaQuery.viewInsets.bottom > 0;
  
  bool get isTablet => width >= 600 && width < 1200;
  
  bool get isDesktop => width >= 1200;
  
  bool get isMobile => width < 600;
  
  void hideKeyboard() {
    FocusScope.of(this).unfocus();
  }
  
  void showSnackBar(String message, {Color? backgroundColor}) {
    ScaffoldMessenger.of(this).showSnackBar(
      SnackBar(
        content: Text(message),
        backgroundColor: backgroundColor,
      ),
    );
  }
  
  Future<T?> push<T>(Widget page) {
    return Navigator.of(this).push<T>(
      MaterialPageRoute(builder: (context) => page),
    );
  }
  
  void pop<T>([T? result]) {
    Navigator.of(this).pop<T>(result);
  }
  
  Future<T?> pushReplacement<T>(Widget page) {
    return Navigator.of(this).pushReplacement<T, T>(
      MaterialPageRoute(builder: (context) => page),
    );
  }
}

extension ColorExtensions on Color {
  Color get lighter {
    final hsl = HSLColor.fromColor(this);
    return hsl.withLightness((hsl.lightness + 0.1).clamp(0.0, 1.0)).toColor();
  }

  Color get darker {
    final hsl = HSLColor.fromColor(this);
    return hsl.withLightness((hsl.lightness - 0.1).clamp(0.0, 1.0)).toColor();
  }

  Color get opposite {
    return Color(0xFFFFFFFF - value);
  }

  String get toHex {
    return '#${value.toRadixString(16).padLeft(8, '0').substring(2)}';
  }
}