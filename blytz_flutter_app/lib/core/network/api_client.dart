import 'package:blytz_flutter_app/core/constants/api_constants.dart';
import 'package:blytz_flutter_app/features/auth/data/models/auth_model.dart';
import 'package:dio/dio.dart';
import 'package:retrofit/retrofit.dart';

part 'api_client.g.dart';

// Auth Service API Client
@RestApi(baseUrl: ApiConstants.authBaseUrl)
abstract class AuthApiClient {
  factory AuthApiClient(Dio dio, {String baseUrl}) = _AuthApiClient;

  @POST(ApiConstants.authServiceLogin)
  Future<AuthResponse> login(@Body() LoginRequest request);

  @POST(ApiConstants.authServiceRegister)
  Future<AuthResponse> register(@Body() RegisterRequest request);

  @POST(ApiConstants.authServiceRefresh)
  Future<AuthResponse> refreshToken(@Body() RefreshTokenRequest request);

  @POST(ApiConstants.authServiceLogout)
  Future<void> logout(@Header('Authorization') String token);

  @GET(ApiConstants.authServiceMe)
  Future<AuthModel> getCurrentUser(@Header('Authorization') String token);

  @PUT(ApiConstants.authServiceProfile)
  Future<AuthModel> updateProfile(
    @Header('Authorization') String token,
    @Body() Map<String, dynamic> request,
  );

  @GET(ApiConstants.authServiceVerify)
  Future<Map<String, dynamic>> verifyToken(@Header('Authorization') String token);
}

// Product Service API Client
@RestApi(baseUrl: ApiConstants.productBaseUrl)
abstract class ProductApiClient {
  factory ProductApiClient(Dio dio, {String baseUrl}) = _ProductApiClient;

  @GET(ApiConstants.products)
  Future<Map<String, dynamic>> getProducts(
    @Query('page') int page,
    @Query('limit') int limit,
    @Query('category') String? category,
    @Query('search') String? search,
    @Query('minPrice') double? minPrice,
    @Query('maxPrice') double? maxPrice,
    @Query('featured') bool? featured,
  );

  @GET('${ApiConstants.products}/{id}')
  Future<Map<String, dynamic>> getProduct(@Path('id') String id);

  @POST(ApiConstants.products)
  Future<Map<String, dynamic>> createProduct(
    @Header('Authorization') String token,
    @Body() Map<String, dynamic> request,
  );

  @PUT('${ApiConstants.products}/{id}')
  Future<Map<String, dynamic>> updateProduct(
    @Header('Authorization') String token,
    @Path('id') String id,
    @Body() Map<String, dynamic> request,
  );

  @DELETE('${ApiConstants.products}/{id}')
  Future<void> deleteProduct(
    @Header('Authorization') String token,
    @Path('id') String id,
  );

  @GET(ApiConstants.productsFeatured)
  Future<Map<String, dynamic>> getFeaturedProducts(
    @Query('limit') int limit,
  );

  @GET(ApiConstants.productsMe)
  Future<Map<String, dynamic>> getMyProducts(
    @Header('Authorization') String token,
    @Query('page') int page,
    @Query('limit') int limit,
  );
}

// Auction Service API Client
@RestApi(baseUrl: ApiConstants.auctionBaseUrl)
abstract class AuctionApiClient {
  factory AuctionApiClient(Dio dio, {String baseUrl}) = _AuctionApiClient;

  @GET(ApiConstants.auctions)
  Future<Map<String, dynamic>> getAuctions(
    @Query('page') int page,
    @Query('limit') int limit,
    @Query('category') String? category,
    @Query('status') String? status,
    @Query('sellerId') String? sellerId,
  );

  @GET(ApiConstants.auctionsActive)
  Future<Map<String, dynamic>> getActiveAuctions(
    @Query('page') int page,
    @Query('limit') int limit,
  );

  @GET('${ApiConstants.auctions}/{id}')
  Future<Map<String, dynamic>> getAuction(@Path('id') String id);

  @POST(ApiConstants.auctions)
  Future<Map<String, dynamic>> createAuction(
    @Header('Authorization') String token,
    @Body() Map<String, dynamic> request,
  );

  @PUT('${ApiConstants.auctions}/{id}')
  Future<Map<String, dynamic>> updateAuction(
    @Header('Authorization') String token,
    @Path('id') String id,
    @Body() Map<String, dynamic> request,
  );

  @DELETE('${ApiConstants.auctions}/{id}')
  Future<void> deleteAuction(
    @Header('Authorization') String token,
    @Path('id') String id,
  );

  @POST('${ApiConstants.auctions}/{id}/bids')
  Future<Map<String, dynamic>> placeBid(
    @Header('Authorization') String token,
    @Path('id') String auctionId,
    @Body() Map<String, dynamic> request,
  );

  @GET('${ApiConstants.auctions}/{id}${ApiConstants.auctionBids}')
  Future<Map<String, dynamic>> getAuctionBids(
    @Path('id') String auctionId,
    @Query('page') int page,
    @Query('limit') int limit,
  );
}

// Order Service API Client
@RestApi(baseUrl: ApiConstants.orderBaseUrl)
abstract class OrderApiClient {
  factory OrderApiClient(Dio dio, {String baseUrl}) = _OrderApiClient;

  @POST(ApiConstants.orders)
  Future<Map<String, dynamic>> createOrder(
    @Header('Authorization') String token,
    @Body() Map<String, dynamic> request,
  );

  @GET('${ApiConstants.orders}/{id}')
  Future<Map<String, dynamic>> getOrder(
    @Header('Authorization') String token,
    @Path('id') String id,
  );

  @GET(ApiConstants.orders)
  Future<Map<String, dynamic>> getMyOrders(
    @Header('Authorization') String token,
    @Query('page') int page,
    @Query('limit') int limit,
    @Query('status') String? status,
  );

  @PUT('${ApiConstants.orders}/{id}${ApiConstants.ordersStatus}')
  Future<Map<String, dynamic>> updateOrderStatus(
    @Header('Authorization') String token,
    @Path('id') String id,
    @Body() Map<String, dynamic> request,
  );

  @POST('${ApiConstants.orders}/{id}${ApiConstants.ordersCancel}')
  Future<Map<String, dynamic>> cancelOrder(
    @Header('Authorization') String token,
    @Path('id') String id,
  );
}

// Payment Service API Client
@RestApi(baseUrl: ApiConstants.paymentBaseUrl)
abstract class PaymentApiClient {
  factory PaymentApiClient(Dio dio, {String baseUrl}) = _PaymentApiClient;

  @POST('${ApiConstants.payments}${ApiConstants.paymentsProcess}')
  Future<Map<String, dynamic>> processPayment(
    @Header('Authorization') String token,
    @Body() Map<String, dynamic> request,
  );

  @GET('${ApiConstants.payments}/{id}')
  Future<Map<String, dynamic>> getPaymentStatus(
    @Header('Authorization') String token,
    @Path('id') String paymentId,
  );

  @GET(ApiConstants.paymentsMethods)
  Future<Map<String, dynamic>> getPaymentMethods(
    @Header('Authorization') String token,
  );

  @POST('${ApiConstants.payments}/{id}${ApiConstants.paymentsRefund}')
  Future<Map<String, dynamic>> processRefund(
    @Header('Authorization') String token,
    @Path('id') String paymentId,
    @Body() Map<String, dynamic> request,
  );
}

// Main API Client that combines all services
class ApiClient {
  final AuthApiClient authClient;
  final ProductApiClient productClient;
  final AuctionApiClient auctionClient;
  final OrderApiClient orderClient;
  final PaymentApiClient paymentClient;

  ApiClient(Dio dio) :
    authClient = AuthApiClient(dio),
    productClient = ProductApiClient(dio),
    auctionClient = AuctionApiClient(dio),
    orderClient = OrderApiClient(dio),
    paymentClient = PaymentApiClient(dio);

  // Auth methods
  Future<AuthResponse> login(LoginRequest request) => authClient.login(request);
  Future<AuthResponse> register(RegisterRequest request) => authClient.register(request);
  Future<AuthResponse> refreshToken(RefreshTokenRequest request) => authClient.refreshToken(request);
  Future<void> logout(String token) => authClient.logout(token);
  Future<AuthModel> getCurrentUser(String token) => authClient.getCurrentUser(token);
  Future<AuthModel> updateProfile(String token, Map<String, dynamic> request) =>
    authClient.updateProfile(token, request);

  // Product methods
  Future<Map<String, dynamic>> getProducts({
    int page = 1,
    int limit = 20,
    String? category,
    String? search,
    double? minPrice,
    double? maxPrice,
    bool? featured,
  }) => productClient.getProducts(
    page, limit, category, search, minPrice, maxPrice, featured,
  );

  Future<Map<String, dynamic>> getProduct(String id) => productClient.getProduct(id);
  Future<Map<String, dynamic>> createProduct(String token, Map<String, dynamic> request) =>
    productClient.createProduct(token, request);
  Future<Map<String, dynamic>> getFeaturedProducts({int limit = 10}) =>
    productClient.getFeaturedProducts(limit);

  // Auction methods
  Future<Map<String, dynamic>> getAuctions({
    int page = 1,
    int limit = 20,
    String? category,
    String? status,
    String? sellerId,
  }) => auctionClient.getAuctions(page, limit, category, status, sellerId);

  Future<Map<String, dynamic>> getActiveAuctions({int page = 1, int limit = 20}) =>
    auctionClient.getActiveAuctions(page, limit);
  Future<Map<String, dynamic>> getAuction(String id) => auctionClient.getAuction(id);
  Future<Map<String, dynamic>> placeBid(String token, String auctionId, Map<String, dynamic> request) =>
    auctionClient.placeBid(token, auctionId, request);

  // Order methods
  Future<Map<String, dynamic>> createOrder(String token, Map<String, dynamic> request) =>
    orderClient.createOrder(token, request);
  Future<Map<String, dynamic>> getOrder(String token, String id) => orderClient.getOrder(token, id);
  Future<Map<String, dynamic>> getMyOrders(String token, {int page = 1, int limit = 20, String? status}) =>
    orderClient.getMyOrders(token, page, limit, status);

  // Payment methods
  Future<Map<String, dynamic>> processPayment(String token, Map<String, dynamic> request) =>
    paymentClient.processPayment(token, request);
  Future<Map<String, dynamic>> getPaymentStatus(String token, String paymentId) =>
    paymentClient.getPaymentStatus(token, paymentId);

  // Auth verification method
  Future<Map<String, dynamic>> verifyToken(String token) =>
    authClient.verifyToken(token);
}