#!/bin/bash

/usr/bin/redis-server /etc/redis.conf &
sleep 5s
/app/go-short