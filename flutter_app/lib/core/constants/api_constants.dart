class ApiConstants {
  static const String baseUrl = 'http://localhost:8080/api';
  
  // Authentication endpoints
  static const String login = '/auth/login';
  static const String register = '/auth/register';
  static const String profile = '/auth/profile';
  static const String refreshToken = '/auth/refresh';
  
  // Vehicle endpoints
  static const String vehicles = '/vehicles';
  static String vehicleById(int id) => '/vehicles/$id';
  static const String vehicleSearch = '/vehicles/search';
  
  // Customer endpoints
  static const String customers = '/customers';
  static String customerById(int id) => '/customers/$id';
  
  // Transaction endpoints
  static const String transactions = '/transactions';
  static String transactionById(int id) => '/transactions/$id';
  static String updateTransactionStatus(int id) => '/transactions/$id/status';
  
  // User management endpoints
  static const String users = '/users';
  static String userById(int id) => '/users/$id';
}