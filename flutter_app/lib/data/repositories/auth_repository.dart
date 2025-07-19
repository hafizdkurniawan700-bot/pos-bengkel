import '../models/user_model.dart';

abstract class AuthRepository {
  Future<LoginResponse> login({
    required String username,
    required String password,
  });

  Future<User> register({
    required String username,
    required String email,
    required String password,
    required String role,
  });

  Future<void> logout();

  Future<String> refreshToken();

  Future<User?> getCurrentUser();

  Future<String?> getToken();

  Future<void> saveToken(String token);

  Future<void> saveUser(User user);

  Future<void> clearAuthData();
}

// Placeholder implementation
class AuthRepositoryImpl implements AuthRepository {
  final dynamic _dataSource;

  const AuthRepositoryImpl(this._dataSource);

  @override
  Future<LoginResponse> login({
    required String username,
    required String password,
  }) async {
    // This would implement actual API call
    throw UnimplementedError('Auth repository not implemented yet');
  }

  @override
  Future<User> register({
    required String username,
    required String email,
    required String password,
    required String role,
  }) async {
    throw UnimplementedError('Auth repository not implemented yet');
  }

  @override
  Future<void> logout() async {
    throw UnimplementedError('Auth repository not implemented yet');
  }

  @override
  Future<String> refreshToken() async {
    throw UnimplementedError('Auth repository not implemented yet');
  }

  @override
  Future<User?> getCurrentUser() async {
    throw UnimplementedError('Auth repository not implemented yet');
  }

  @override
  Future<String?> getToken() async {
    throw UnimplementedError('Auth repository not implemented yet');
  }

  @override
  Future<void> saveToken(String token) async {
    throw UnimplementedError('Auth repository not implemented yet');
  }

  @override
  Future<void> saveUser(User user) async {
    throw UnimplementedError('Auth repository not implemented yet');
  }

  @override
  Future<void> clearAuthData() async {
    throw UnimplementedError('Auth repository not implemented yet');
  }
}