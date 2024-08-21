POST http://localhost:8080/tagsl/v1
{
    "port": 1,
    "payload": "8002cdcd1300744f5e166018040b14341a",
    "devEui": ""
}

HTTP 200
Content-Type: application/json

[Asserts]
bytes count == 138
jsonpath "$.moving" == false
jsonpath "$.latitude" == 47.041811
jsonpath "$.longitude" == 7.622494
jsonpath "$.altitude" == 572.8
jsonpath "$.year" == 24
jsonpath "$.month" == 4
jsonpath "$.day" == 11
jsonpath "$.hour" == 20
jsonpath "$.minute" == 52
jsonpath "$.second" == 26
