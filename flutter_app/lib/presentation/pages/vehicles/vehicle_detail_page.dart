import 'package:flutter/material.dart';

import '../../core/theme/app_theme.dart';

class VehicleDetailPage extends StatelessWidget {
  const VehicleDetailPage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Vehicle Details'),
      ),
      body: const Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(Icons.info, size: 64, color: AppTheme.primaryBlue),
            SizedBox(height: 16),
            Text('Vehicle Details', style: AppTheme.headlineMedium),
            Text('Coming soon...'),
          ],
        ),
      ),
    );
  }
}