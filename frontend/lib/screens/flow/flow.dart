import 'package:flutter/material.dart';
import 'package:frontend/data/models/flow.dart' as Model;

class FlowDetailPage extends StatelessWidget {
  Model.Flow flow;

  FlowDetailPage({this.flow});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text("Tsekaro"),
      ),
      body: ListView(
        padding: EdgeInsets.all(5.0),
        children: <Widget>[
          Container(
            alignment: Alignment.center,
            padding: EdgeInsets.all(10.0),
            child: Text(
              flow.name,
              style: TextStyle(fontSize: 20.0),
            ),
          ),
          Container(
            alignment: Alignment.topRight,
            margin: EdgeInsets.all(5.0),
            child: Text(flow.createdAt),
          ),
        ],
      ),
    );
  }
}
