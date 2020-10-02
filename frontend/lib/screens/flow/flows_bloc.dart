import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:frontend/screens/flow/flows_event.dart';
import 'package:frontend/screens/flow/flows_state.dart';

class FlowsBloc extends Bloc<FlowsEvent, FlowsState> {
  FlowsBloc(FlowsState initialState) : super(initialState);

  @override
  Stream<FlowsState> mapEventToState(FlowsEvent event) async* {
    if (event is LoadFlows) {
      yield* _mapLoadFlowsToState();
    }
    yield* _mapLoadFlowsToState();
  }

  Stream<FlowsState> _mapLoadFlowsToState() async* {
    try {
      yield FlowsLoaded();
    } catch (_) {
      yield FlowsNotLoaded();
    }
  }
}
