import 'package:equatable/equatable.dart';

import '../../../data/models/user_model.dart';

abstract class AuthEvent extends Equatable {
  const AuthEvent();

  @override
  List<Object?> get props => [];
}

class LoginRequested extends AuthEvent {
  final String username;
  final String password;

  const LoginRequested({
    required this.username,
    required this.password,
  });

  @override
  List<Object?> get props => [username, password];
}

class RegisterRequested extends AuthEvent {
  final String username;
  final String email;
  final String password;
  final String role;

  const RegisterRequested({
    required this.username,
    required this.email,
    required this.password,
    required this.role,
  });

  @override
  List<Object?> get props => [username, email, password, role];
}

class LogoutRequested extends AuthEvent {
  const LogoutRequested();
}

class RefreshTokenRequested extends AuthEvent {
  const RefreshTokenRequested();
}

class CheckAuthStatus extends AuthEvent {
  const CheckAuthStatus();
}