import 'package:equatable/equatable.dart';
import 'package:frontend/models/flow.dart';

abstract class FlowsEvent extends Equatable {
  const FlowsEvent();

  @override
  List<Object> get props => [];
}

class LoadFlows extends FlowsEvent {}

class AddFlow extends FlowsEvent {
  final Flow flow;

  const AddFlow(this.flow);

  @override
  List<Object> get props => [flow];

  @override
  String toString() => 'AddFlow { flow: $flow }';
}

class UpdateFlow extends FlowsEvent {
  final Flow updatedFlow;

  const UpdateFlow(this.updatedFlow);

  @override
  List<Object> get props => [updatedFlow];

  @override
  String toString() => 'UpdateFlow { updatedFlow: $updatedFlow }';
}

class DeleteFlow extends FlowsEvent {
  final Flow flow;

  const DeleteFlow(this.flow);

  @override
  List<Object> get props => [flow];

  @override
  String toString() => 'DeleteFlow { flow: $flow }';
}

class ClearCompletedFlows extends FlowsEvent {}

class ToggleAllFlows extends FlowsEvent {}
