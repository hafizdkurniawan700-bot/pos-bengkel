// Placeholder customer bloc
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:equatable/equatable.dart';

import '../../../data/repositories/customer_repository.dart';

// Events
abstract class CustomerEvent extends Equatable {
  const CustomerEvent();
  @override
  List<Object?> get props => [];
}

// States  
abstract class CustomerState extends Equatable {
  const CustomerState();
  @override
  List<Object?> get props => [];
}

class CustomerInitial extends CustomerState {
  const CustomerInitial();
}

// Bloc
class CustomerBloc extends Bloc<CustomerEvent, CustomerState> {
  final CustomerRepository _repository;

  CustomerBloc(this._repository) : super(const CustomerInitial());
}