import 'package:equatable/equatable.dart';
import 'package:frontend/data/models/flow.dart';
import 'package:meta/meta.dart';

abstract class FlowState extends Equatable {}

class FlowInitialState extends FlowState {
  @override
  // TODO: implement props
  List<Object> get props => [];
}

class FlowLoadingState extends FlowState {
  @override
  // TODO: implement props
  List<Object> get props => [];
}

class FlowLoadedState extends FlowState {
  List<Flow> flows;

  FlowLoadedState({@required this.flows});

  @override
  // TODO: implement props
  List<Object> get props => [flows];
}

class FlowErrorState extends FlowState {
  String message;

  FlowErrorState({@required this.message});

  @override
  // TODO: implement props
  List<Object> get props => [message];
}
