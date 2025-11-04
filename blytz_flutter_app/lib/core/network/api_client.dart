import 'package:blytz_flutter_app/core/constants/api_constants.dart';
import 'package:blytz_flutter_app/features/auth/data/models/auth_model.dart';
import 'package:dio/dio.dart';
import 'package:retrofit/retrofit.dart';

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
  Future<Map<String, dynamic>> getAuctions(
    @Query('page') int page,
    @Query('limit') int limit,
    @Query('category') String? category,
    @Query('status') String? status,
  );

  @GET('${ApiConstants.auctions}/{id}')
  Future<Map<String, dynamic>> getAuction(@Path('id') String id);

  @POST('${ApiConstants.auctions}/{id}/bids')
  Future<Map<String, dynamic>> placeBid(
    @Path('id') String auctionId,
    @Body() Map<String, dynamic> request,
  );
}