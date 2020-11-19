import 'package:bloc/bloc.dart';
import 'package:frontend/data/models/flow.dart';
import 'package:frontend/data/repository/flow_repository.dart';

import 'flow_event.dart';
import 'flow_state.dart';

class FlowsBloc extends Bloc<FlowEvent, FlowState> {
  FlowRepository repository;

  FlowsBloc(FlowState initialState, FlowRepository repository)
      : super(initialState) {
    this.repository = repository;
  }

  @override
  Stream<FlowState> mapEventToState(FlowEvent event) async* {
    if (event is LoadFlowsEvent) {
      yield FlowLoadingState();
      try {
        List<Flow> flows = await repository.getFlows();
        yield FlowLoadedState(flows: flows);
      } catch (e) {
        yield FlowErrorState(message: e.toString());
      }
    }
  }
}
