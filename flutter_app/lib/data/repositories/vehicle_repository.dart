// Placeholder repositories
import '../models/vehicle_model.dart';

abstract class VehicleRepository {
  Future<List<Vehicle>> getVehicles({VehicleSearchParams? params});
  Future<Vehicle> getVehicleById(int id);
  Future<Vehicle> createVehicle(VehicleRequest request);
  Future<Vehicle> updateVehicle(int id, VehicleRequest request);
  Future<void> deleteVehicle(int id);
  Future<List<Vehicle>> searchVehicles(String query);
}

class VehicleRepositoryImpl implements VehicleRepository {
  final dynamic _dataSource;
  const VehicleRepositoryImpl(this._dataSource);

  @override
  Future<List<Vehicle>> getVehicles({VehicleSearchParams? params}) async {
    throw UnimplementedError();
  }

  @override
  Future<Vehicle> getVehicleById(int id) async {
    throw UnimplementedError();
  }

  @override
  Future<Vehicle> createVehicle(VehicleRequest request) async {
    throw UnimplementedError();
  }

  @override
  Future<Vehicle> updateVehicle(int id, VehicleRequest request) async {
    throw UnimplementedError();
  }

  @override
  Future<void> deleteVehicle(int id) async {
    throw UnimplementedError();
  }

  @override
  Future<List<Vehicle>> searchVehicles(String query) async {
    throw UnimplementedError();
  }
}

abstract class CustomerRepository {
  // Placeholder methods
}

class CustomerRepositoryImpl implements CustomerRepository {
  final dynamic _dataSource;
  const CustomerRepositoryImpl(this._dataSource);
}