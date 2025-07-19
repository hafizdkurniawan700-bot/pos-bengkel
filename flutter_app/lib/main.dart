import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:get_it/get_it.dart';

import 'core/theme/app_theme.dart';
import 'core/network/api_client.dart';
import 'data/repositories/auth_repository.dart';
import 'data/repositories/vehicle_repository.dart';
import 'data/repositories/customer_repository.dart';
import 'data/datasources/remote_datasource.dart';
import 'presentation/bloc/auth/auth_bloc.dart';
import 'presentation/bloc/vehicle/vehicle_bloc.dart';
import 'presentation/bloc/customer/customer_bloc.dart';
import 'presentation/pages/auth/login_page.dart';
import 'presentation/pages/dashboard/admin_dashboard.dart';
import 'presentation/pages/dashboard/sales_dashboard.dart';
import 'presentation/pages/dashboard/customer_dashboard.dart';
import 'presentation/pages/vehicles/vehicle_list_page.dart';
import 'presentation/pages/vehicles/vehicle_detail_page.dart';

final GetIt getIt = GetIt.instance;

void main() {
  setupDependencyInjection();
  runApp(const MyApp());
}

void setupDependencyInjection() {
  // Core
  getIt.registerLazySingleton(() => ApiClient());
  
  // Data Sources
  getIt.registerLazySingleton(() => RemoteDataSource(getIt()));
  
  // Repositories
  getIt.registerLazySingleton<AuthRepository>(() => AuthRepositoryImpl(getIt()));
  getIt.registerLazySingleton<VehicleRepository>(() => VehicleRepositoryImpl(getIt()));
  getIt.registerLazySingleton<CustomerRepository>(() => CustomerRepositoryImpl(getIt()));
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MultiBlocProvider(
      providers: [
        BlocProvider(create: (context) => AuthBloc(getIt())),
        BlocProvider(create: (context) => VehicleBloc(getIt())),
        BlocProvider(create: (context) => CustomerBloc(getIt())),
      ],
      child: MaterialApp(
        title: 'POS Bengkel - Vehicle Sales Management',
        theme: AppTheme.lightTheme,
        darkTheme: AppTheme.darkTheme,
        themeMode: ThemeMode.system,
        initialRoute: '/login',
        routes: {
          '/login': (context) => const LoginPage(),
          '/admin-dashboard': (context) => const AdminDashboard(),
          '/sales-dashboard': (context) => const SalesDashboard(),
          '/customer-dashboard': (context) => const CustomerDashboard(),
          '/vehicles': (context) => const VehicleListPage(),
          '/vehicle-detail': (context) => const VehicleDetailPage(),
        },
        debugShowCheckedModeBanner: false,
      ),
    );
  }
}