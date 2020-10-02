// Copyright 2018 The Flutter Architecture Sample Authors. All rights reserved.
// Use of this source code is governed by the MIT license that can be found
// in the LICENSE file.

import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:frontend/blocs/simple_bloc_delegate.dart';
import 'package:frontend/localization.dart';
import 'package:frontend/routes.dart';
import 'package:frontend/screens/flow/flows_bloc.dart';
import 'package:frontend/screens/flow/flows_event.dart';
import 'package:frontend/screens/flow/flows_state.dart';
import 'package:frontend/screens/home/home.dart';
import 'package:frontend/theme/style.dart';

void runBlocLibraryApp() {
  // BlocSupervisor oversees Blocs and delegates to BlocDelegate.
  // We can set the BlocSupervisor's delegate to an instance of `SimpleBlocDelegate`.
  // This will allow us to handle all transitions and errors in SimpleBlocDelegate.
  Bloc.observer = SimpleBlocDelegate();

  runApp(
    BlocProvider<FlowsBloc>(
      create: (context) {
        return FlowsBloc(FlowsLoading())..add(LoadFlows());
      },
      child: TesterApp(),
    ),
  );
}

class TesterApp extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    final flowsBloc = BlocProvider.of<FlowsBloc>(context);

    return MaterialApp(
      onGenerateTitle: (context) =>
          FlutterBlocLocalizations.of(context).appTitle,
      theme: ArchSampleTheme.theme,
      routes: {
        Routes.home: (context) {
          return MultiBlocProvider(
            child: HomeScreen(),
          );
        },
        Routes.addTodo: (context) {
          return AddEditScreen(
            key: ArchSampleKeys.addTodoScreen,
            onSave: (task, note) {
              todosBloc.add(AddTodo(Todo(task, note: note)));
            },
            isEditing: false,
          );
        },
      },
    );
  }
}
