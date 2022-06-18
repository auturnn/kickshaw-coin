#!/bin/bash
echo "> server start"
#80번 포트 사용을 위해 root로 진행
nohup sudo /home/ec2-user/build/kickshaw-coin > /dev/null 2> /dev/null < /dev/null &