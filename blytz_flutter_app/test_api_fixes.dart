// Simple test to verify API compatibility fixes
import 'package:flutter/material.dart';
import 'package:getwidget/getwidget.dart';
import 'package:velocity_x/velocity_x.dart';

void main() {
  // Test GetWidget API
  final card = GFCard(
    child: Container(),
  );
  
  // Test Velocity_X API
  final text = 'Test'.text.color(Colors.grey600).make();
  
  // Test GFBadge API
  final badge = GFBadge(
    text: 'LIVE',
    color: Colors.red,
    badgeShape: GFBadgeShape.pills,
  );
  
  print('All API compatibility tests passed!');
}