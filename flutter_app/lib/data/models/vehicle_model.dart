import 'package:json_annotation/json_annotation.dart';
import 'package:equatable/equatable.dart';

part 'vehicle_model.g.dart';

@JsonSerializable()
class Vehicle extends Equatable {
  final int id;
  final String brand;
  final String model;
  final int year;
  final double price;
  final String status;
  final String? images;
  final String? description;
  @JsonKey(name: 'engine_type')
  final String? engineType;
  final String? transmission;
  @JsonKey(name: 'fuel_type')
  final String? fuelType;
  final int mileage;
  final String? color;
  @JsonKey(name: 'created_at')
  final DateTime createdAt;
  @JsonKey(name: 'updated_at')
  final DateTime updatedAt;

  const Vehicle({
    required this.id,
    required this.brand,
    required this.model,
    required this.year,
    required this.price,
    required this.status,
    this.images,
    this.description,
    this.engineType,
    this.transmission,
    this.fuelType,
    required this.mileage,
    this.color,
    required this.createdAt,
    required this.updatedAt,
  });

  factory Vehicle.fromJson(Map<String, dynamic> json) => _$VehicleFromJson(json);
  Map<String, dynamic> toJson() => _$VehicleToJson(this);

  @override
  List<Object?> get props => [
        id,
        brand,
        model,
        year,
        price,
        status,
        images,
        description,
        engineType,
        transmission,
        fuelType,
        mileage,
        color,
        createdAt,
        updatedAt,
      ];

  // Helper methods
  bool get isAvailable => status == 'available';
  bool get isSold => status == 'sold';
  bool get isReserved => status == 'reserved';
  bool get inMaintenance => status == 'maintenance';

  String get displayName => '$brand $model ($year)';
  String get formattedPrice => '\$${price.toStringAsFixed(2)}';
  String get formattedMileage => '${mileage.toString().replaceAllMapped(
        RegExp(r'(\d{1,3})(?=(\d{3})+(?!\d))'),
        (Match m) => '${m[1]},',
      )} miles';

  List<String> get imageUrls {
    if (images == null || images!.isEmpty) return [];
    // Assuming images is a JSON array string
    try {
      // Simple parsing for now - in production, use proper JSON parsing
      return images!.split(',').map((url) => url.trim()).toList();
    } catch (e) {
      return [];
    }
  }

  Vehicle copyWith({
    int? id,
    String? brand,
    String? model,
    int? year,
    double? price,
    String? status,
    String? images,
    String? description,
    String? engineType,
    String? transmission,
    String? fuelType,
    int? mileage,
    String? color,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) {
    return Vehicle(
      id: id ?? this.id,
      brand: brand ?? this.brand,
      model: model ?? this.model,
      year: year ?? this.year,
      price: price ?? this.price,
      status: status ?? this.status,
      images: images ?? this.images,
      description: description ?? this.description,
      engineType: engineType ?? this.engineType,
      transmission: transmission ?? this.transmission,
      fuelType: fuelType ?? this.fuelType,
      mileage: mileage ?? this.mileage,
      color: color ?? this.color,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
    );
  }
}

@JsonSerializable()
class VehicleRequest extends Equatable {
  final String brand;
  final String model;
  final int year;
  final double price;
  final String status;
  final String? images;
  final String? description;
  @JsonKey(name: 'engine_type')
  final String? engineType;
  final String? transmission;
  @JsonKey(name: 'fuel_type')
  final String? fuelType;
  final int mileage;
  final String? color;

  const VehicleRequest({
    required this.brand,
    required this.model,
    required this.year,
    required this.price,
    required this.status,
    this.images,
    this.description,
    this.engineType,
    this.transmission,
    this.fuelType,
    required this.mileage,
    this.color,
  });

  factory VehicleRequest.fromJson(Map<String, dynamic> json) => _$VehicleRequestFromJson(json);
  Map<String, dynamic> toJson() => _$VehicleRequestToJson(this);

  @override
  List<Object?> get props => [
        brand,
        model,
        year,
        price,
        status,
        images,
        description,
        engineType,
        transmission,
        fuelType,
        mileage,
        color,
      ];
}

@JsonSerializable()
class VehicleSearchParams extends Equatable {
  final String? query;
  final String? brand;
  final String? model;
  @JsonKey(name: 'year_min')
  final int? yearMin;
  @JsonKey(name: 'year_max')
  final int? yearMax;
  @JsonKey(name: 'price_min')
  final double? priceMin;
  @JsonKey(name: 'price_max')
  final double? priceMax;
  final String? status;
  @JsonKey(name: 'fuel_type')
  final String? fuelType;
  final String? transmission;
  final int limit;
  final int offset;

  const VehicleSearchParams({
    this.query,
    this.brand,
    this.model,
    this.yearMin,
    this.yearMax,
    this.priceMin,
    this.priceMax,
    this.status,
    this.fuelType,
    this.transmission,
    this.limit = 10,
    this.offset = 0,
  });

  factory VehicleSearchParams.fromJson(Map<String, dynamic> json) => _$VehicleSearchParamsFromJson(json);
  Map<String, dynamic> toJson() => _$VehicleSearchParamsToJson(this);

  @override
  List<Object?> get props => [
        query,
        brand,
        model,
        yearMin,
        yearMax,
        priceMin,
        priceMax,
        status,
        fuelType,
        transmission,
        limit,
        offset,
      ];
}