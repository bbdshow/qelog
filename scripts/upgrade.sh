#!/bin/bash
app="qelog_receiver"
ps -ef | grep -w ${app} | grep -v grep | awk  '{print "kill -9 " $2}' | sh

sleep 1

mv ${app} ${app}`date +%m%d%H%M%S`

chmod +x ${app}_new

mv ${app}_new ${app}

nohup ./${app} -f ../configs/config.toml > ./nohup.out 2>&1 &

sleep 2

tail -n 20 ./nohup.out

