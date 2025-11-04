import 'dart:async';
import 'dart:convert';
import 'package:blytz_flutter_app/core/constants/api_constants.dart';
import 'package:blytz_flutter_app/features/auction/data/models/auction_model.dart';
import 'package:blytz_flutter_app/features/auction/data/models/bid_model.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:web_socket_channel/web_socket_channel.dart';

// WebSocket event types
class WebSocketEvent {
  const WebSocketEvent({required this.type, required this.data});

  final String type;
  final Map<String, dynamic> data;

  factory WebSocketEvent.fromJson(Map<String, dynamic> json) {
    return WebSocketEvent(
      type: json['type'] as String,
      data: json['data'] as Map<String, dynamic>,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'type': type,
      'data': data,
    };
  }
}

class WebSocketService {
  WebSocketChannel? _channel;
  StreamSubscription? _subscription;
  StreamController<WebSocketEvent>? _eventController;
  bool _isConnected = false;

  // Stream for WebSocket events
  Stream<WebSocketEvent> get events => _eventController?.stream ?? const Stream.empty();

  bool get isConnected => _isConnected;

  // Connect to auction updates WebSocket
  Future<void> connectToAuctionUpdates(String auctionId) async {
    try {
      if (_isConnected) {
        await disconnect();
      }

      final uri = Uri.parse('${ApiConstants.wsAuctionUpdates}/$auctionId');
      _channel = WebSocketChannel.connect(uri);

      _eventController = StreamController<WebSocketEvent>.broadcast();
      _subscription = _channel!.stream.listen(
        _handleMessage,
        onError: _handleError,
        onDone: _handleDisconnection,
      );

      _isConnected = true;
    } catch (e) {
      _isConnected = false;
      throw Exception('Failed to connect to auction updates: $e');
    }
  }

  // Connect to live stream WebSocket
  Future<void> connectToLiveStream(String streamId) async {
    try {
      if (_isConnected) {
        await disconnect();
      }

      final uri = Uri.parse('${ApiConstants.wsLiveStream}/$streamId');
      _channel = WebSocketChannel.connect(uri);

      _eventController = StreamController<WebSocketEvent>.broadcast();
      _subscription = _channel!.stream.listen(
        _handleMessage,
        onError: _handleError,
        onDone: _handleDisconnection,
      );

      _isConnected = true;
    } catch (e) {
      _isConnected = false;
      throw Exception('Failed to connect to live stream: $e');
    }
  }

  // Connect to chat WebSocket
  Future<void> connectToChat(String roomId) async {
    try {
      if (_isConnected) {
        await disconnect();
      }

      final uri = Uri.parse('${ApiConstants.wsChat}/$roomId');
      _channel = WebSocketChannel.connect(uri);

      _eventController = StreamController<WebSocketEvent>.broadcast();
      _subscription = _channel!.stream.listen(
        _handleMessage,
        onError: _handleError,
        onDone: _handleDisconnection,
      );

      _isConnected = true;
    } catch (e) {
      _isConnected = false;
      throw Exception('Failed to connect to chat: $e');
    }
  }

  // Send message through WebSocket
  void sendMessage(WebSocketEvent event) {
    if (_isConnected && _channel != null) {
      final message = jsonEncode(event.toJson());
      _channel!.sink.add(message);
    }
  }

  // Handle incoming WebSocket messages
  void _handleMessage(dynamic data) {
    try {
      final Map<String, dynamic> json = jsonDecode(data as String);
      final event = WebSocketEvent.fromJson(json);
      _eventController?.add(event);
    } catch (e) {
      print('Error parsing WebSocket message: $e');
    }
  }

  // Handle WebSocket errors
  void _handleError(dynamic error) {
    print('WebSocket error: $error');
    _isConnected = false;
    _eventController?.add(WebSocketEvent(
      type: 'error',
      data: {'message': error.toString()},
    ));
  }

  // Handle WebSocket disconnection
  void _handleDisconnection() {
    print('WebSocket disconnected');
    _isConnected = false;
    _eventController?.add(const WebSocketEvent(
      type: 'disconnected',
      data: {'message': 'WebSocket connection closed'},
    ));
  }

  // Disconnect from WebSocket
  Future<void> disconnect() async {
    _isConnected = false;
    await _subscription?.cancel();
    await _channel?.sink.close();
    await _eventController?.close();
    _subscription = null;
    _channel = null;
    _eventController = null;
  }

  // Dispose method
  void dispose() {
    disconnect();
  }
}

// Riverpod providers
final webSocketServiceProvider = Provider<WebSocketService>((ref) {
  final service = WebSocketService();

  // Auto dispose when provider is disposed
  ref.onDispose(() {
    service.dispose();
  });

  return service;
});

// Auction-specific WebSocket provider
final auctionWebSocketProvider = StreamProvider.autoDispose<WebSocketEvent>((ref) {
  final wsService = ref.watch(webSocketServiceProvider);

  return wsService.events.map((event) {
    // Filter only auction-related events
    if (event.type.startsWith('auction_') ||
        event.type.startsWith('bid_') ||
        event.type == 'error' ||
        event.type == 'disconnected') {
      return event;
    }
    return event; // Or filter out unwanted events
  });
});

// Live stream WebSocket provider
final liveStreamWebSocketProvider = StreamProvider.autoDispose<WebSocketEvent>((ref) {
  final wsService = ref.watch(webSocketServiceProvider);

  return wsService.events.map((event) {
    // Filter only live stream related events
    if (event.type.startsWith('stream_') ||
        event.type.startsWith('viewer_') ||
        event.type == 'error' ||
        event.type == 'disconnected') {
      return event;
    }
    return event;
  });
});

// Chat WebSocket provider
final chatWebSocketProvider = StreamProvider.autoDispose<WebSocketEvent>((ref) {
  final wsService = ref.watch(webSocketServiceProvider);

  return wsService.events.map((event) {
    // Filter only chat-related events
    if (event.type.startsWith('chat_') ||
        event.type.startsWith('message_') ||
        event.type == 'error' ||
        event.type == 'disconnected') {
      return event;
    }
    return event;
  });
});