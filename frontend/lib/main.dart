import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:frontend/blocs/flow_blocs/flow_bloc.dart';
import 'package:frontend/blocs/flow_blocs/flow_state.dart';
import 'package:frontend/data/repository/flow_repository.dart';
import 'package:frontend/routes.dart';
import 'package:frontend/theme/style.dart';

void main() {
  runApp(TesterApp());
}

class TesterApp extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MultiBlocProvider(
        providers: [
          BlocProvider<FlowsBloc>(
              create: (context) =>
                  FlowsBloc(FlowInitialState(), FlowRepositoryImpl()))
        ],
        child: MaterialApp(
          theme: appTheme(),
          initialRoute: "/",
          routes: routes,
        ));
  }
}
