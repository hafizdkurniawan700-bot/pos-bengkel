import 'package:flutter/material.dart';

import '../../core/theme/app_theme.dart';
import '../../core/constants/app_constants.dart';

class AdminDashboard extends StatefulWidget {
  const AdminDashboard({super.key});

  @override
  State<AdminDashboard> createState() => _AdminDashboardState();
}

class _AdminDashboardState extends State<AdminDashboard> {
  int _selectedIndex = 0;

  final List<Widget> _pages = [
    const DashboardOverview(),
    const VehicleManagement(),
    const CustomerManagement(),
    const TransactionManagement(),
    const UserManagement(),
  ];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Admin Dashboard'),
        actions: [
          IconButton(
            icon: const Icon(Icons.notifications),
            onPressed: () {
              // Show notifications
            },
          ),
          PopupMenuButton<String>(
            onSelected: (value) {
              switch (value) {
                case 'profile':
                  // Navigate to profile
                  break;
                case 'settings':
                  // Navigate to settings
                  break;
                case 'logout':
                  Navigator.of(context).pushReplacementNamed('/login');
                  break;
              }
            },
            itemBuilder: (context) => [
              const PopupMenuItem(
                value: 'profile',
                child: ListTile(
                  leading: Icon(Icons.person),
                  title: Text('Profile'),
                ),
              ),
              const PopupMenuItem(
                value: 'settings',
                child: ListTile(
                  leading: Icon(Icons.settings),
                  title: Text('Settings'),
                ),
              ),
              const PopupMenuDivider(),
              const PopupMenuItem(
                value: 'logout',
                child: ListTile(
                  leading: Icon(Icons.logout),
                  title: Text('Logout'),
                ),
              ),
            ],
          ),
        ],
      ),
      body: _pages[_selectedIndex],
      bottomNavigationBar: BottomNavigationBar(
        type: BottomNavigationBarType.fixed,
        currentIndex: _selectedIndex,
        onTap: (index) {
          setState(() {
            _selectedIndex = index;
          });
        },
        items: const [
          BottomNavigationBarItem(
            icon: Icon(Icons.dashboard),
            label: 'Overview',
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.directions_car),
            label: 'Vehicles',
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.people),
            label: 'Customers',
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.receipt),
            label: 'Transactions',
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.admin_panel_settings),
            label: 'Users',
          ),
        ],
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () {
          // Quick action based on current tab
          switch (_selectedIndex) {
            case 1:
              Navigator.of(context).pushNamed('/add-vehicle');
              break;
            case 2:
              Navigator.of(context).pushNamed('/add-customer');
              break;
            case 3:
              Navigator.of(context).pushNamed('/add-transaction');
              break;
          }
        },
        child: const Icon(Icons.add),
      ),
    );
  }
}

class DashboardOverview extends StatelessWidget {
  const DashboardOverview({super.key});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(16.0),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            'Welcome back, Admin!',
            style: AppTheme.headlineMedium,
          ),
          const SizedBox(height: 20),
          
          // Statistics Cards
          Row(
            children: [
              Expanded(
                child: _StatCard(
                  title: 'Total Vehicles',
                  value: '45',
                  icon: Icons.directions_car,
                  color: AppTheme.primaryBlue,
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: _StatCard(
                  title: 'Available',
                  value: '32',
                  icon: Icons.check_circle,
                  color: AppTheme.successColor,
                ),
              ),
            ],
          ),
          const SizedBox(height: 12),
          Row(
            children: [
              Expanded(
                child: _StatCard(
                  title: 'Customers',
                  value: '128',
                  icon: Icons.people,
                  color: AppTheme.primaryOrange,
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: _StatCard(
                  title: 'Transactions',
                  value: '23',
                  icon: Icons.receipt,
                  color: AppTheme.warningColor,
                ),
              ),
            ],
          ),
          const SizedBox(height: 24),
          
          // Recent Activity
          Text(
            'Recent Activity',
            style: AppTheme.titleLarge,
          ),
          const SizedBox(height: 12),
          Expanded(
            child: ListView(
              children: const [
                _ActivityTile(
                  title: 'New vehicle added',
                  subtitle: 'Toyota Camry 2024',
                  time: '2 hours ago',
                  icon: Icons.add_circle,
                ),
                _ActivityTile(
                  title: 'Transaction completed',
                  subtitle: 'Honda Civic sold to John Doe',
                  time: '4 hours ago',
                  icon: Icons.monetization_on,
                ),
                _ActivityTile(
                  title: 'New customer registered',
                  subtitle: 'Jane Smith',
                  time: '6 hours ago',
                  icon: Icons.person_add,
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

class _StatCard extends StatelessWidget {
  final String title;
  final String value;
  final IconData icon;
  final Color color;

  const _StatCard({
    required this.title,
    required this.value,
    required this.icon,
    required this.color,
  });

  @override
  Widget build(BuildContext context) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Icon(icon, color: color, size: 24),
                Text(
                  value,
                  style: AppTheme.headlineMedium.copyWith(color: color),
                ),
              ],
            ),
            const SizedBox(height: 8),
            Text(
              title,
              style: AppTheme.bodyMedium.copyWith(
                color: Colors.grey[600],
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _ActivityTile extends StatelessWidget {
  final String title;
  final String subtitle;
  final String time;
  final IconData icon;

  const _ActivityTile({
    required this.title,
    required this.subtitle,
    required this.time,
    required this.icon,
  });

  @override
  Widget build(BuildContext context) {
    return ListTile(
      leading: CircleAvatar(
        backgroundColor: AppTheme.primaryBlue.withOpacity(0.1),
        child: Icon(icon, color: AppTheme.primaryBlue, size: 20),
      ),
      title: Text(title, style: AppTheme.bodyLarge),
      subtitle: Text(subtitle),
      trailing: Text(
        time,
        style: AppTheme.bodyMedium.copyWith(
          color: Colors.grey[600],
        ),
      ),
    );
  }
}

// Placeholder pages for other tabs
class VehicleManagement extends StatelessWidget {
  const VehicleManagement({super.key});

  @override
  Widget build(BuildContext context) {
    return const Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(Icons.directions_car, size: 64, color: Colors.grey),
          SizedBox(height: 16),
          Text('Vehicle Management'),
          Text('Coming soon...'),
        ],
      ),
    );
  }
}

class CustomerManagement extends StatelessWidget {
  const CustomerManagement({super.key});

  @override
  Widget build(BuildContext context) {
    return const Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(Icons.people, size: 64, color: Colors.grey),
          SizedBox(height: 16),
          Text('Customer Management'),
          Text('Coming soon...'),
        ],
      ),
    );
  }
}

class TransactionManagement extends StatelessWidget {
  const TransactionManagement({super.key});

  @override
  Widget build(BuildContext context) {
    return const Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(Icons.receipt, size: 64, color: Colors.grey),
          SizedBox(height: 16),
          Text('Transaction Management'),
          Text('Coming soon...'),
        ],
      ),
    );
  }
}

class UserManagement extends StatelessWidget {
  const UserManagement({super.key});

  @override
  Widget build(BuildContext context) {
    return const Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(Icons.admin_panel_settings, size: 64, color: Colors.grey),
          SizedBox(height: 16),
          Text('User Management'),
          Text('Coming soon...'),
        ],
      ),
    );
  }
}