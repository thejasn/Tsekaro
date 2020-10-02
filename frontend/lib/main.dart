import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:frontend/routes.dart';
import 'package:frontend/screens/home/bloc/home_bloc.dart';
import 'package:frontend/theme/style.dart';

void main() {
  runApp(TesterApp());
}

class TesterApp extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MultiBlocProvider(
        providers: [
          BlocProvider<HomeBloc>(
            create: (context) => HomeBloc(),
          )
        ],
        child: MaterialApp(
          theme: appTheme(),
          initialRoute: "/",
          routes: routes,
        ));
  }
}
