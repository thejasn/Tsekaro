class FlowList {
  int page;
  int pageSize;
  List<Flow> flows;
  int totalRecords;

  FlowList({this.page, this.pageSize, this.flows, this.totalRecords});

  FlowList.fromJson(Map<String, dynamic> json) {
    page = json['page'];
    pageSize = json['page_size'];
    if (json['data'] != null) {
      flows = new List<Flow>();
      json['data'].forEach((v) {
        flows.add(new Flow.fromJson(v));
      });
    }
    totalRecords = json['total_records'];
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = new Map<String, dynamic>();
    data['page'] = this.page;
    data['page_size'] = this.pageSize;
    if (this.flows != null) {
      data['data'] = this.flows.map((v) => v.toJson()).toList();
    }
    data['total_records'] = this.totalRecords;
    return data;
  }
}

class Flow {
  int id;
  String name;
  String createdAt;
  String updatedAt;

  Flow({this.id, this.name, this.createdAt, this.updatedAt});

  Flow.fromJson(Map<String, dynamic> json) {
    id = json['id'];
    name = json['name'];
    createdAt = json['created_at'];
    updatedAt = json['updated_at'];
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = new Map<String, dynamic>();
    data['id'] = this.id;
    data['name'] = this.name;
    data['created_at'] = this.createdAt;
    data['updated_at'] = this.updatedAt;
    return data;
  }
}
