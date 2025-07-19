import 'package:flutter_bloc/flutter_bloc.dart';

import '../../../data/repositories/auth_repository.dart';
import '../../../core/network/network_exceptions.dart';
import 'auth_event.dart';
import 'auth_state.dart';

class AuthBloc extends Bloc<AuthEvent, AuthState> {
  final AuthRepository _authRepository;

  AuthBloc(this._authRepository) : super(const AuthInitial()) {
    on<LoginRequested>(_onLoginRequested);
    on<RegisterRequested>(_onRegisterRequested);
    on<LogoutRequested>(_onLogoutRequested);
    on<RefreshTokenRequested>(_onRefreshTokenRequested);
    on<CheckAuthStatus>(_onCheckAuthStatus);
  }

  Future<void> _onLoginRequested(
    LoginRequested event,
    Emitter<AuthState> emit,
  ) async {
    emit(const AuthLoading());
    
    try {
      final loginResponse = await _authRepository.login(
        username: event.username,
        password: event.password,
      );
      
      await _authRepository.saveToken(loginResponse.token);
      await _authRepository.saveUser(loginResponse.user);
      
      emit(AuthSuccess(
        user: loginResponse.user,
        token: loginResponse.token,
      ));
    } on NetworkExceptions catch (e) {
      emit(AuthFailure(message: e.message));
    } catch (e) {
      emit(AuthFailure(message: 'An unexpected error occurred: $e'));
    }
  }

  Future<void> _onRegisterRequested(
    RegisterRequested event,
    Emitter<AuthState> emit,
  ) async {
    emit(const AuthLoading());
    
    try {
      final user = await _authRepository.register(
        username: event.username,
        email: event.email,
        password: event.password,
        role: event.role,
      );
      
      // Auto-login after successful registration
      add(LoginRequested(
        username: event.username,
        password: event.password,
      ));
    } on NetworkExceptions catch (e) {
      emit(AuthFailure(message: e.message));
    } catch (e) {
      emit(AuthFailure(message: 'Registration failed: $e'));
    }
  }

  Future<void> _onLogoutRequested(
    LogoutRequested event,
    Emitter<AuthState> emit,
  ) async {
    try {
      await _authRepository.logout();
      emit(const AuthUnauthenticated());
    } catch (e) {
      emit(const AuthUnauthenticated());
    }
  }

  Future<void> _onRefreshTokenRequested(
    RefreshTokenRequested event,
    Emitter<AuthState> emit,
  ) async {
    try {
      final newToken = await _authRepository.refreshToken();
      final user = await _authRepository.getCurrentUser();
      
      if (user != null) {
        emit(AuthSuccess(
          user: user,
          token: newToken,
        ));
      } else {
        emit(const AuthUnauthenticated());
      }
    } on NetworkExceptions catch (e) {
      if (e.isUnauthorized) {
        emit(const AuthUnauthenticated());
      } else {
        emit(AuthFailure(message: e.message));
      }
    } catch (e) {
      emit(const AuthUnauthenticated());
    }
  }

  Future<void> _onCheckAuthStatus(
    CheckAuthStatus event,
    Emitter<AuthState> emit,
  ) async {
    try {
      final token = await _authRepository.getToken();
      final user = await _authRepository.getCurrentUser();
      
      if (token != null && user != null) {
        emit(AuthSuccess(
          user: user,
          token: token,
        ));
      } else {
        emit(const AuthUnauthenticated());
      }
    } catch (e) {
      emit(const AuthUnauthenticated());
    }
  }
}