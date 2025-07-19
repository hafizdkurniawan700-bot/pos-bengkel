class AppConstants {
  static const String appName = 'POS Bengkel';
  static const String appVersion = '1.0.0';
  
  // Storage keys
  static const String tokenKey = 'auth_token';
  static const String userKey = 'user_data';
  static const String themeKey = 'theme_mode';
  
  // User roles
  static const String adminRole = 'admin';
  static const String salesRole = 'sales';
  static const String cashierRole = 'cashier';
  static const String customerRole = 'customer';
  
  // Vehicle statuses
  static const String availableStatus = 'available';
  static const String soldStatus = 'sold';
  static const String reservedStatus = 'reserved';
  static const String maintenanceStatus = 'maintenance';
  
  // Transaction statuses
  static const String pendingStatus = 'pending';
  static const String completedStatus = 'completed';
  static const String cancelledStatus = 'cancelled';
  static const String refundedStatus = 'refunded';
  
  // Test drive statuses
  static const String scheduledStatus = 'scheduled';
  static const String testDriveCompletedStatus = 'completed';
  static const String testDriveCancelledStatus = 'cancelled';
  static const String noShowStatus = 'no_show';
  
  // Pagination
  static const int defaultPageSize = 10;
  static const int maxPageSize = 100;
  
  // Image settings
  static const int maxImageSize = 5 * 1024 * 1024; // 5MB
  static const List<String> allowedImageTypes = ['jpg', 'jpeg', 'png', 'webp'];
  
  // Validation
  static const int minPasswordLength = 6;
  static const int maxNameLength = 100;
  static const int maxDescriptionLength = 500;
  
  // Formats
  static const String dateFormat = 'yyyy-MM-dd';
  static const String timeFormat = 'HH:mm';
  static const String dateTimeFormat = 'yyyy-MM-dd HH:mm';
  static const String currencyFormat = '\$#,##0.00';
}