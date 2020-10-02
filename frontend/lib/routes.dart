import 'package:flutter/material.dart';
import 'package:frontend/screens/home/home.dart';

final Map<String, WidgetBuilder> routes = <String, WidgetBuilder>{
  "/": (BuildContext context) => HomeScreen()
};

class Routes {
  static final home = '/';
  static final addTodo = '/addTodo';
}
