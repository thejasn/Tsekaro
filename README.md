# Tsekaro

## Summary

Tsekaro is an open system framework for creating/executing/managing API tests and integration tests. It is mainly aimed to help testing in a `Microservices` rich ecosystem.

Tsekaro is built with the idea of `Testing as a Service`, meaning it can easily be deployed as a microservice to help teams during development as well as during regression.

## Index

- [Tsekaro](#tsekaro)
  - [Summary](#summary)
  - [Index](#index)
  - [API Documentation](#api-documentation)
    - [1. Add Flow](#1-add-flow)
    - [2. Add Testcase](#2-add-testcase)
    - [3. Delete Flow](#3-delete-flow)
    - [4. Delete Testcase](#4-delete-testcase)
    - [5. Execute Flow](#5-execute-flow)
    - [6. Execute Testcase](#6-execute-testcase)
    - [7. Get All Flows](#7-get-all-flows)
        - [I. Example Request: Get All Flows](#i-example-request-get-all-flows)
        - [I. Example Response: Get All Flows](#i-example-response-get-all-flows)
    - [8. Get All Testcases](#8-get-all-testcases)
    - [9. Get Flow](#9-get-flow)
    - [10. Get Testcase](#10-get-testcase)
    - [11. Status](#11-status)
    - [12. Update Flows](#12-update-flows)
    - [13. Update Testcase](#13-update-testcase)

---

## API Documentation

### 1. Add Flow

**_Endpoint:_**

```bash
Method: POST
Type: RAW
URL: http://localhost:8080/v1/flows
```

**_Body:_**

```js
{
    "name": "flow2"
}
```

### 2. Add Testcase

**_Endpoint:_**

```bash
Method: POST
Type: RAW
URL: http://localhost:8080/v1/testcases
```

**_Body:_**

```js
{
    "name": "step1",
    "flow_id": 11,
    "api": "GRPC",
    "test_case_id": 1,
    "operation": "EQUAL",
    "expected": {
        "data": "Hello thejas!"
    },
    "actual": "message",
    "scheme": "http",
    "host": "localhost",
    "port": 50051,
    "path": "helloworld.Greeter/SayHello",
    "body": "{\"message\": \"thejas\"}"
}
```

### 3. Delete Flow

**_Endpoint:_**

```bash
Method: DELETE
Type:
URL: http://localhost:8080/v1/flows/1
```

### 4. Delete Testcase

**_Endpoint:_**

```bash
Method: DELETE
Type:
URL: http://localhost:8080/v1/testcases/3
```

### 5. Execute Flow

**_Endpoint:_**

```bash
Method: GET
Type:
URL: http://localhost:8080/v1/flows/execute/11
```

### 6. Execute Testcase

**_Endpoint:_**

```bash
Method: GET
Type:
URL: http://localhost:8080/v1/testcases/execute/5
```

### 7. Get All Flows

**_Endpoint:_**

```bash
Method: GET
Type:
URL: http://localhost:8080/v1/flows
```

**_Query params:_**

| Key      | Value | Description |
| -------- | ----- | ----------- |
| page     | 5     |             |
| pagesize | 1     |             |
| order    | name  |             |

**_More example Requests/Responses:_**

##### I. Example Request: Get All Flows

**_Query:_**

| Key      | Value | Description |
| -------- | ----- | ----------- |
| page     | 5     |             |
| pagesize | 1     |             |
| order    | name  |             |

##### I. Example Response: Get All Flows

```js
{
    "page": 5,
    "page_size": 1,
    "data": [],
    "total_records": 2
}
```

**_Status Code:_** 200

<br>

### 8. Get All Testcases

**_Endpoint:_**

```bash
Method: GET
Type:
URL: http://localhost:8080/v1/testcases
```

### 9. Get Flow

**_Endpoint:_**

```bash
Method: GET
Type:
URL: http://localhost:8080/v1/flows/10
```

### 10. Get Testcase

**_Endpoint:_**

```bash
Method: GET
Type:
URL: http://localhost:8080/v1/testcases/5
```

### 11. Status

**_Endpoint:_**

```bash
Method: GET
Type:
URL: http://localhost:8080/v1/status
```

### 12. Update Flows

**_Endpoint:_**

```bash
Method: PUT
Type: RAW
URL: http://localhost:8080/v1/flows/2
```

**_Body:_**

```js
{
    "name": "name2"
}
```

### 13. Update Testcase

**_Endpoint:_**

```bash
Method: PUT
Type: RAW
URL: http://localhost:8080/v1/testcases/6
```

**_Body:_**

```js
{
    "name": "step1",
    "flow_id": 11,
    "api": "GRPC",
    "test_case_id": 1,
    "operation": "EQUAL",
    "expected": {
        "data": "Hello thejas"
    },
    "actual": "message",
    "scheme": "http",
    "host": "localhost",
    "port": 50051,
    "path": "helloworld.Greeter/SayHello",
    "body": "{\"name\": \"thejas\"}"
}
```

---

[Back to top](#tester)
