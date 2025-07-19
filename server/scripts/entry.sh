#!/bin/bash

migrate -database "${MYSQL_URL}" -path ./db/migrations/ up

./server