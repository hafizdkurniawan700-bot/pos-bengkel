import 'package:json_annotation/json_annotation.dart';
import 'package:equatable/equatable.dart';

part 'user_model.g.dart';

@JsonSerializable()
class User extends Equatable {
  final int id;
  final String username;
  final String email;
  final String role;
  @JsonKey(name: 'created_at')
  final DateTime createdAt;
  @JsonKey(name: 'updated_at')
  final DateTime updatedAt;

  const User({
    required this.id,
    required this.username,
    required this.email,
    required this.role,
    required this.createdAt,
    required this.updatedAt,
  });

  factory User.fromJson(Map<String, dynamic> json) => _$UserFromJson(json);
  Map<String, dynamic> toJson() => _$UserToJson(this);

  @override
  List<Object?> get props => [id, username, email, role, createdAt, updatedAt];

  // Helper methods
  bool get isAdmin => role == 'admin';
  bool get isSales => role == 'sales';
  bool get isCashier => role == 'cashier';
  bool get isCustomer => role == 'customer';
  
  bool get canManageVehicles => isAdmin || isSales;
  bool get canManageCustomers => isAdmin || isSales;
  bool get canManageTransactions => isAdmin || isSales || isCashier;
  bool get canViewReports => isAdmin || isSales;

  User copyWith({
    int? id,
    String? username,
    String? email,
    String? role,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) {
    return User(
      id: id ?? this.id,
      username: username ?? this.username,
      email: email ?? this.email,
      role: role ?? this.role,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
    );
  }
}

@JsonSerializable()
class LoginRequest extends Equatable {
  final String username;
  final String password;

  const LoginRequest({
    required this.username,
    required this.password,
  });

  factory LoginRequest.fromJson(Map<String, dynamic> json) => _$LoginRequestFromJson(json);
  Map<String, dynamic> toJson() => _$LoginRequestToJson(this);

  @override
  List<Object?> get props => [username, password];
}

@JsonSerializable()
class LoginResponse extends Equatable {
  final String token;
  final User user;

  const LoginResponse({
    required this.token,
    required this.user,
  });

  factory LoginResponse.fromJson(Map<String, dynamic> json) => _$LoginResponseFromJson(json);
  Map<String, dynamic> toJson() => _$LoginResponseToJson(this);

  @override
  List<Object?> get props => [token, user];
}

@JsonSerializable()
class RegisterRequest extends Equatable {
  final String username;
  final String email;
  final String password;
  final String role;

  const RegisterRequest({
    required this.username,
    required this.email,
    required this.password,
    required this.role,
  });

  factory RegisterRequest.fromJson(Map<String, dynamic> json) => _$RegisterRequestFromJson(json);
  Map<String, dynamic> toJson() => _$RegisterRequestToJson(this);

  @override
  List<Object?> get props => [username, email, password, role];
}