import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:frontend/localization.dart';
import 'package:frontend/routes.dart';

import 'bloc/home_bloc.dart';
import 'bloc/home_state.dart';


class HomeScreen extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    // final tabBloc = BlocProvider.of<TabBloc>(context);

        return Scaffold(
          appBar: AppBar(
            title: Text(FlutterBlocLocalizations.of(context).appTitle),
          ),
          floatingActionButton: FloatingActionButton(
            key: ArchSampleKeys.addTodoFab,
            onPressed: () {
              Navigator.pushNamed(context, Routes.addTodo);
            },
            child: Icon(Icons.add),
            tooltip: ArchSampleLocalizations.of(context).addTodo,
          ),
          bottomNavigationBar: TabSelector(
            activeTab: activeTab,
            onTabSelected: (tab) => tabBloc.add(UpdateTab(tab)),
          ),
        );

  }
}

class HomeDefauleState extends State<StatefulWidget> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
        appBar: AppBar(
          title: Text("Tester"),
        ),
        body: BlocBuilder<HomeBloc, HomeState>(builder: (context, homeState) {
          return Scaffold(
            body: homeState == "Default" ? DefaultHomeState() ,
          );
        }));
  }
}
