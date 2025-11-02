import 'package:dio/dio.dart';
import 'package:retrofit/retrofit.dart';
import '../constants/api_constants.dart';
import '../../features/auth/data/models/auth_model.dart';
import '../../features/auctions/data/models/auction_model.dart';
import '../../features/chat/data/models/chat_model.dart';

part 'api_client.g.dart';

@RestApi(baseUrl: ApiConstants.baseUrl)
abstract class ApiClient {
  factory ApiClient(Dio dio, {String baseUrl}) = _ApiClient;

  // Authentication endpoints
  @POST(ApiConstants.login)
  Future<AuthResponse> login(@Body() LoginRequest request);

  @POST(ApiConstants.register)
  Future<AuthResponse> register(@Body() RegisterRequest request);

  @POST(ApiConstants.refresh)
  Future<AuthResponse> refreshToken(@Body() RefreshTokenRequest request);

  @POST(ApiConstants.logout)
  Future<void> logout();

  // Auction endpoints
  @GET(ApiConstants.auctions)
  Future<List<AuctionModel>> getAuctions(
    @Query('page') int page,
    @Query('limit') int limit,
    @Query('category') String? category,
    @Query('status') String? status,
  );

  @GET('${ApiConstants.auctions}/{id}')
  Future<AuctionModel> getAuction(@Path('id') String id);

  @POST('${ApiConstants.auctions}/{id}/bids')
  Future<BidResponse> placeBid(
    @Path('id') String auctionId,
    @Body() BidRequest request,
  );

  @GET('${ApiConstants.auctions}/{id}/bids')
  Future<List<BidModel>> getBidHistory(
    @Path('id') String auctionId,
    @Query('page') int page,
    @Query('limit') int limit,
  );

  // Product endpoints
  @GET(ApiConstants.products)
  Future<List<ProductModel>> getProducts(
    @Query('page') int page,
    @Query('limit') int limit,
    @Query('category') String? category,
  );

  @GET('${ApiConstants.products}/{id}')
  Future<ProductModel> getProduct(@Path('id') String id);

  // Chat endpoints
  @GET('${ApiConstants.chat}/rooms')
  Future<List<ChatRoomModel>> getChatRooms();

  @GET('${ApiConstants.chat}/rooms/{roomId}/messages')
  Future<List<MessageModel>> getMessages(
    @Path('roomId') String roomId,
    @Query('limit') int limit,
    @Query('before') String? before,
  );

  @POST('${ApiConstants.chat}/rooms/{roomId}/messages')
  Future<MessageModel> sendMessage(
    @Path('roomId') String roomId,
    @Body() SendMessageRequest request,
  );
}