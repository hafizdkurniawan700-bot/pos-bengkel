import 'package:flutter/material.dart';

import '../../core/theme/app_theme.dart';

class CustomerDashboard extends StatelessWidget {
  const CustomerDashboard({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Customer Dashboard'),
        actions: [
          IconButton(
            icon: const Icon(Icons.logout),
            onPressed: () {
              Navigator.of(context).pushReplacementNamed('/login');
            },
          ),
        ],
      ),
      body: const Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(Icons.person, size: 64, color: AppTheme.primaryBlue),
            SizedBox(height: 16),
            Text('Customer Dashboard', style: AppTheme.headlineMedium),
            Text('Coming soon...'),
          ],
        ),
      ),
    );
  }
}