import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:frontend/blocs/flow_blocs/flow_bloc.dart';
import 'package:frontend/blocs/flow_blocs/flow_event.dart';
import 'package:frontend/blocs/flow_blocs/flow_state.dart';
import 'package:frontend/data/models/flow.dart' as Model;
import 'package:frontend/screens/flow/flow.dart';

class HomePage extends StatefulWidget {
  @override
  _HomePageState createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  FlowsBloc flowsBloc;

  @override
  void initState() {
    super.initState();
    flowsBloc = BlocProvider.of<FlowsBloc>(context);
    flowsBloc.add(LoadFlowsEvent());
  }

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      home: Builder(
        builder: (context) {
          return Material(
            child: Scaffold(
              appBar: AppBar(
                title: Text("Tsekaro"),
                actions: <Widget>[
                  IconButton(
                    icon: Icon(Icons.refresh),
                    onPressed: () {
                      flowsBloc.add(LoadFlowsEvent());
                    },
                  ),
                  IconButton(
                    icon: Icon(Icons.info),
                    onPressed: () {
                      navigateToAoutPage(context);
                    },
                  )
                ],
              ),
              body: Container(
                child: BlocListener<FlowsBloc, FlowState>(
                  listener: (context, state) {
                    if (state is FlowErrorState) {
                      Scaffold.of(context).showSnackBar(
                        SnackBar(
                          content: Text(state.message),
                        ),
                      );
                    }
                  },
                  child: BlocBuilder<FlowsBloc, FlowState>(
                    builder: (context, state) {
                      if (state is FlowInitialState) {
                        return buildLoading();
                      } else if (state is FlowLoadingState) {
                        return buildLoading();
                      } else if (state is FlowLoadedState) {
                        return buildFlowList(state.flows);
                      } else if (state is FlowErrorState) {
                        return buildErrorUi(state.message);
                      }
                    },
                  ),
                ),
              ),
            ),
          );
        },
      ),
    );
  }
}

Widget buildLoading() {
  return Center(
    child: CircularProgressIndicator(),
  );
}

void navigateToAoutPage(BuildContext context) {
  Navigator.push(context, MaterialPageRoute(builder: (context) {
    return HomePage();
  }));
}

Widget buildErrorUi(String message) {
  return Center(
    child: Padding(
      padding: const EdgeInsets.all(8.0),
      child: Text(
        message,
        style: TextStyle(color: Colors.red),
      ),
    ),
  );
}

Widget buildFlowList(List<Model.Flow> flows) {
  return ListView.builder(
    itemCount: flows.length,
    itemBuilder: (ctx, pos) {
      return Padding(
        padding: const EdgeInsets.all(8.0),
        child: InkWell(
          child: ListTile(
            title: Text(flows[pos].name),
          ),
          onTap: () {
            navigateToArticleDetailPage(ctx, flows[pos]);
          },
        ),
      );
    },
  );
}

void navigateToArticleDetailPage(BuildContext context, Model.Flow flow) {
  Navigator.push(context, MaterialPageRoute(builder: (context) {
    return FlowDetailPage(
      flow: flow,
    );
  }));
}
