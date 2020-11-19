import 'dart:convert';

import 'package:frontend/data/models/flow.dart';
import 'package:http/http.dart' as http;

abstract class FlowRepository {
  Future<List<Flow>> getFlows();
}

class FlowRepositoryImpl implements FlowRepository {
  @override
  Future<List<Flow>> getFlows() async {
    var response =
        await http.get(Uri.encodeFull("http://localhost:8080/v1/flows"));
    if (response.statusCode == 200) {
      var data = json.decode(response.body);
      List<Flow> flows = FlowList.fromJson(data).flows;
      return flows;
    } else {
      throw Exception();
    }
  }
}
