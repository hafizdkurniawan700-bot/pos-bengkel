// Placeholder vehicle bloc files
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:equatable/equatable.dart';

import '../../../data/repositories/vehicle_repository.dart';
import '../../../data/models/vehicle_model.dart';

// Events
abstract class VehicleEvent extends Equatable {
  const VehicleEvent();
  @override
  List<Object?> get props => [];
}

class LoadVehicles extends VehicleEvent {
  const LoadVehicles();
}

// States
abstract class VehicleState extends Equatable {
  const VehicleState();
  @override
  List<Object?> get props => [];
}

class VehicleInitial extends VehicleState {
  const VehicleInitial();
}

class VehicleLoading extends VehicleState {
  const VehicleLoading();
}

class VehicleLoaded extends VehicleState {
  final List<Vehicle> vehicles;
  const VehicleLoaded(this.vehicles);
  
  @override
  List<Object?> get props => [vehicles];
}

class VehicleError extends VehicleState {
  final String message;
  const VehicleError(this.message);
  
  @override
  List<Object?> get props => [message];
}

// Bloc
class VehicleBloc extends Bloc<VehicleEvent, VehicleState> {
  final VehicleRepository _repository;

  VehicleBloc(this._repository) : super(const VehicleInitial()) {
    on<LoadVehicles>(_onLoadVehicles);
  }

  Future<void> _onLoadVehicles(
    LoadVehicles event,
    Emitter<VehicleState> emit,
  ) async {
    emit(const VehicleLoading());
    try {
      final vehicles = await _repository.getVehicles();
      emit(VehicleLoaded(vehicles));
    } catch (e) {
      emit(VehicleError(e.toString()));
    }
  }
}