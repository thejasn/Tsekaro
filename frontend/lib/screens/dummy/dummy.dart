import 'package:flutter/material.dart';
import 'package:frontend/screens/home/home.dart';

class Dummy extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text("Dummy"),
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () {
          Navigator.push(
              context, MaterialPageRoute(builder: (context) => HomeScreen()));
        },
        child: Icon(Icons.add),
      ),
    );
  }
}
