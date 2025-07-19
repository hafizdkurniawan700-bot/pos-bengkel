import 'package:flutter/material.dart';

import '../../core/theme/app_theme.dart';

class SalesDashboard extends StatelessWidget {
  const SalesDashboard({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Sales Dashboard'),
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
            Icon(Icons.trending_up, size: 64, color: AppTheme.primaryOrange),
            SizedBox(height: 16),
            Text('Sales Dashboard', style: AppTheme.headlineMedium),
            Text('Coming soon...'),
          ],
        ),
      ),
    );
  }
}