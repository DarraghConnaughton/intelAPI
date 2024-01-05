#!/bin/bash

docker build -t intelagent .
docker run -p 8080:8080 intelagent
