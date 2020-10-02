import 'package:equatable/equatable.dart';
import 'package:frontend/models/flow.dart';

abstract class FlowsState extends Equatable {
  const FlowsState();

  @override
  List<Object> get props => [];
}

class FlowsLoading extends FlowsState {}

class FlowsLoaded extends FlowsState {
  final List<Flow> flows;

  const FlowsLoaded([this.flows = const []]);

  @override
  List<Object> get props => [flows];

  @override
  String toString() => 'FlowsLoaded { flows: $flows }';
}

class FlowsNotLoaded extends FlowsState {}
