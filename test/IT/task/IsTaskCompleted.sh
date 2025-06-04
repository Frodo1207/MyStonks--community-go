#!/bin/bash

# 请求参数
USER_ID=1
TASK_ID=101

# 接口地址
URL="http://localhost:8000/api/v1/task/complete?user_id=${USER_ID}&task_id=${TASK_ID}"

# 发送 POST 请求
echo "Sending POST request to $URL"
echo "------------------------------------"
curl -s -X GET "$URL" | jq .
echo "------------------------------------"