import 'package:equatable/equatable.dart';

class HomeOpenEvent extends Equatable {
  final String user;

  const HomeOpenEvent({this.user});

  @override
  List<Object> get props => [user];

  @override
  bool get stringify => true;
}
