import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import '../../core/theme/app_theme.dart';
import '../../core/constants/app_constants.dart';
import '../bloc/auth/auth_bloc.dart';
import '../bloc/auth/auth_event.dart';
import '../bloc/auth/auth_state.dart';
import '../widgets/common/custom_button.dart';
import '../widgets/common/loading_widget.dart';

class LoginPage extends StatefulWidget {
  const LoginPage({super.key});

  @override
  State<LoginPage> createState() => _LoginPageState();
}

class _LoginPageState extends State<LoginPage> {
  final _formKey = GlobalKey<FormState>();
  final _usernameController = TextEditingController();
  final _passwordController = TextEditingController();
  bool _obscurePassword = true;

  @override
  void dispose() {
    _usernameController.dispose();
    _passwordController.dispose();
    super.dispose();
  }

  void _onLoginPressed() {
    if (_formKey.currentState!.validate()) {
      context.read<AuthBloc>().add(
            LoginRequested(
              username: _usernameController.text.trim(),
              password: _passwordController.text,
            ),
          );
    }
  }

  void _navigateBasedOnRole(String role) {
    switch (role) {
      case AppConstants.adminRole:
        Navigator.of(context).pushReplacementNamed('/admin-dashboard');
        break;
      case AppConstants.salesRole:
        Navigator.of(context).pushReplacementNamed('/sales-dashboard');
        break;
      case AppConstants.cashierRole:
        Navigator.of(context).pushReplacementNamed('/sales-dashboard');
        break;
      case AppConstants.customerRole:
        Navigator.of(context).pushReplacementNamed('/customer-dashboard');
        break;
      default:
        Navigator.of(context).pushReplacementNamed('/customer-dashboard');
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: BlocListener<AuthBloc, AuthState>(
        listener: (context, state) {
          if (state is AuthSuccess) {
            _navigateBasedOnRole(state.user.role);
          } else if (state is AuthFailure) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                content: Text(state.message),
                backgroundColor: AppTheme.errorColor,
              ),
            );
          }
        },
        child: Container(
          decoration: const BoxDecoration(
            gradient: LinearGradient(
              begin: Alignment.topCenter,
              end: Alignment.bottomCenter,
              colors: [
                AppTheme.primaryBlue,
                Color(0xFF0D47A1),
              ],
            ),
          ),
          child: SafeArea(
            child: Center(
              child: SingleChildScrollView(
                padding: const EdgeInsets.all(24.0),
                child: Card(
                  elevation: 8,
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(16),
                  ),
                  child: Padding(
                    padding: const EdgeInsets.all(32.0),
                    child: Form(
                      key: _formKey,
                      child: Column(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          // Logo and Title
                          const Icon(
                            Icons.directions_car,
                            size: 64,
                            color: AppTheme.primaryBlue,
                          ),
                          const SizedBox(height: 16),
                          Text(
                            AppConstants.appName,
                            style: AppTheme.headlineMedium.copyWith(
                              color: AppTheme.primaryBlue,
                            ),
                          ),
                          const SizedBox(height: 8),
                          Text(
                            'Vehicle Sales Management System',
                            style: AppTheme.bodyMedium.copyWith(
                              color: Colors.grey[600],
                            ),
                            textAlign: TextAlign.center,
                          ),
                          const SizedBox(height: 32),

                          // Username Field
                          TextFormField(
                            controller: _usernameController,
                            decoration: const InputDecoration(
                              labelText: 'Username or Email',
                              prefixIcon: Icon(Icons.person_outline),
                            ),
                            validator: (value) {
                              if (value == null || value.isEmpty) {
                                return 'Please enter your username or email';
                              }
                              return null;
                            },
                            textInputAction: TextInputAction.next,
                          ),
                          const SizedBox(height: 16),

                          // Password Field
                          TextFormField(
                            controller: _passwordController,
                            obscureText: _obscurePassword,
                            decoration: InputDecoration(
                              labelText: 'Password',
                              prefixIcon: const Icon(Icons.lock_outline),
                              suffixIcon: IconButton(
                                icon: Icon(
                                  _obscurePassword
                                      ? Icons.visibility_off
                                      : Icons.visibility,
                                ),
                                onPressed: () {
                                  setState(() {
                                    _obscurePassword = !_obscurePassword;
                                  });
                                },
                              ),
                            ),
                            validator: (value) {
                              if (value == null || value.isEmpty) {
                                return 'Please enter your password';
                              }
                              if (value.length < AppConstants.minPasswordLength) {
                                return 'Password must be at least ${AppConstants.minPasswordLength} characters';
                              }
                              return null;
                            },
                            textInputAction: TextInputAction.done,
                            onFieldSubmitted: (_) => _onLoginPressed(),
                          ),
                          const SizedBox(height: 24),

                          // Login Button
                          BlocBuilder<AuthBloc, AuthState>(
                            builder: (context, state) {
                              if (state is AuthLoading) {
                                return const LoadingWidget();
                              }

                              return SizedBox(
                                width: double.infinity,
                                child: CustomButton(
                                  text: 'Sign In',
                                  onPressed: _onLoginPressed,
                                  icon: Icons.login,
                                ),
                              );
                            },
                          ),
                          const SizedBox(height: 16),

                          // Demo Credentials
                          Container(
                            padding: const EdgeInsets.all(12),
                            decoration: BoxDecoration(
                              color: Colors.blue[50],
                              borderRadius: BorderRadius.circular(8),
                              border: Border.all(color: Colors.blue[200]!),
                            ),
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                Text(
                                  'Demo Credentials:',
                                  style: AppTheme.labelLarge.copyWith(
                                    color: AppTheme.primaryBlue,
                                  ),
                                ),
                                const SizedBox(height: 4),
                                Text(
                                  'Username: admin\nPassword: admin123',
                                  style: AppTheme.bodyMedium.copyWith(
                                    fontFamily: 'monospace',
                                  ),
                                ),
                              ],
                            ),
                          ),
                        ],
                      ),
                    ),
                  ),
                ),
              ),
            ),
          ),
        ),
      ),
    );
  }
}