#!/bin/sh

nohup /server/hello > /data/log/nginx/hello.log 2>&1 &
