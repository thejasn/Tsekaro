import 'package:flutter/material.dart';

class CollectionCart extends StatelessWidget {
  final String title;

  const CollectionCart({Key key, this.title}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Card(
      child: Row(
        children: [ListTile(title: Text(this.title))],
      ),
    );
  }
}
