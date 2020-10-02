import 'dart:developer';

import 'package:bloc/bloc.dart';
import 'package:frontend/screens/home/bloc/home_event.dart';
import 'package:frontend/screens/home/bloc/home_state.dart';

class HomeBloc extends Bloc<HomeOpenEvent, HomeState> {
  HomeBloc() : super(DefaultHomeState());

  @override
  Stream<HomeState> mapEventToState(HomeOpenEvent event) async* {
    if (event is HomeOpenEvent) {
      yield* _mapHomeOpenToState(event);
    }
  }

  Stream<HomeState> _mapHomeOpenToState(HomeOpenEvent event) async* {
    yield HomeLoading();
    log('data:');
    yield DefaultHomeState();
  }

  @override
  void onTransition(Transition<HomeOpenEvent, HomeState> transition) {
    print(transition);
    super.onTransition(transition);
  }
}
