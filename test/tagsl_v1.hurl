POST http://localhost:8080/tagsl/v1
{
    "port": 1,
    "payload": "8002cdcd1300744f5e166018040b14341a",
    "devEui": ""
}

HTTP 200
Content-Type: application/json

[Asserts]
bytes count == 239
jsonpath "$.data.moving" == false
jsonpath "$.data.latitude" == 47.041811
jsonpath "$.data.longitude" == 7.622494
jsonpath "$.data.altitude" == 572.8
jsonpath "$.data.year" == 24
jsonpath "$.data.month" == 4
jsonpath "$.data.day" == 11
jsonpath "$.data.hour" == 20
jsonpath "$.data.minute" == 52
jsonpath "$.data.second" == 26
jsonpath "$.metadata.moving" == false
jsonpath "$.metadata.configChangeId" == 0
jsonpath "$.metadata.configChangeSuccess" == false
jsonpath "$.metadata.dutyCycle" == true