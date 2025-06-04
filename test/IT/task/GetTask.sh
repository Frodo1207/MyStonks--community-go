#!/bin/bash

# 设置接口地址
URL="http://localhost:8000/api/v1/tasks?category=daily"

echo "Sending request to $URL"
echo "------------------------------------"
curl -s -X GET "$URL" | jq .
echo "------------------------------------"

# 设置接口地址
URL="http://localhost:8000/api/v1/tasks?category=newbie"

echo "Sending request to $URL"
echo "------------------------------------"
curl -s -X GET "$URL" | jq .
echo "------------------------------------"


# 设置接口地址
URL="http://localhost:8000/api/v1/tasks?category=other"

echo "Sending request to $URL"
echo "------------------------------------"
curl -s -X GET "$URL" | jq .
echo "------------------------------------"




