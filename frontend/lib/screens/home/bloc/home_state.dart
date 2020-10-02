import 'package:equatable/equatable.dart';

abstract class HomeState extends Equatable {
  @override
  List<Object> get props => [];
}

class DefaultHomeState extends HomeState {
  @override
  bool get stringify => true;
}

class HomeLoading extends HomeState {
  @override
  String toString() => 'HomeScreenLoading';
}
