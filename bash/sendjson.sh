#!/bin/bash

# JSON 数据
json_data='{
    "email": "9328123@qq.com",
    "password": "Cc@002300",
    "confirmPassword": "Cc@002300"
}'

# 发送 POST 请求
curl -X POST -H "Content-Type: application/json" -d "$json_data" http://localhost:8080/users/signup

