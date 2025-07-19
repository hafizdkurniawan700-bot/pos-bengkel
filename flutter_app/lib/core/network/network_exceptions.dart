import 'package:dio/dio.dart';

class NetworkExceptions implements Exception {
  final String message;
  final int? statusCode;

  const NetworkExceptions(this.message, [this.statusCode]);

  static NetworkExceptions fromDioError(DioException dioError) {
    switch (dioError.type) {
      case DioExceptionType.connectionTimeout:
        return const NetworkExceptions('Connection timeout. Please check your internet connection.');
      
      case DioExceptionType.sendTimeout:
        return const NetworkExceptions('Send timeout. Please try again.');
      
      case DioExceptionType.receiveTimeout:
        return const NetworkExceptions('Receive timeout. Please try again.');
      
      case DioExceptionType.badResponse:
        return _handleResponse(dioError.response);
      
      case DioExceptionType.cancel:
        return const NetworkExceptions('Request was cancelled.');
      
      case DioExceptionType.connectionError:
        return const NetworkExceptions('No internet connection. Please check your network settings.');
      
      case DioExceptionType.badCertificate:
        return const NetworkExceptions('Certificate error. Please contact support.');
      
      case DioExceptionType.unknown:
      default:
        return const NetworkExceptions('Something went wrong. Please try again.');
    }
  }

  static NetworkExceptions _handleResponse(Response? response) {
    if (response == null) {
      return const NetworkExceptions('No response from server.');
    }

    final statusCode = response.statusCode;
    final data = response.data;

    // Try to extract error message from response
    String errorMessage = 'Unknown error occurred.';
    
    if (data is Map<String, dynamic>) {
      errorMessage = data['error'] ?? 
                    data['message'] ?? 
                    data['detail'] ?? 
                    errorMessage;
    } else if (data is String) {
      errorMessage = data;
    }

    switch (statusCode) {
      case 400:
        return NetworkExceptions('Bad request: $errorMessage', statusCode);
      case 401:
        return NetworkExceptions('Unauthorized access. Please login again.', statusCode);
      case 403:
        return NetworkExceptions('Access forbidden. You don\'t have permission.', statusCode);
      case 404:
        return NetworkExceptions('Resource not found.', statusCode);
      case 409:
        return NetworkExceptions('Conflict: $errorMessage', statusCode);
      case 422:
        return NetworkExceptions('Validation error: $errorMessage', statusCode);
      case 429:
        return NetworkExceptions('Too many requests. Please try again later.', statusCode);
      case 500:
        return NetworkExceptions('Internal server error. Please try again later.', statusCode);
      case 502:
        return NetworkExceptions('Bad gateway. Server is temporarily unavailable.', statusCode);
      case 503:
        return NetworkExceptions('Service unavailable. Please try again later.', statusCode);
      case 504:
        return NetworkExceptions('Gateway timeout. Please try again later.', statusCode);
      default:
        return NetworkExceptions('Error $statusCode: $errorMessage', statusCode);
    }
  }

  @override
  String toString() {
    return message;
  }

  // Helper methods for specific error types
  bool get isUnauthorized => statusCode == 401;
  bool get isForbidden => statusCode == 403;
  bool get isNotFound => statusCode == 404;
  bool get isValidationError => statusCode == 422;
  bool get isServerError => statusCode != null && statusCode! >= 500;
  bool get isClientError => statusCode != null && statusCode! >= 400 && statusCode! < 500;
}