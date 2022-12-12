#!/bin/bash
app='qelog_receiver'
nohup ./${app} -f ../configs/config.toml > ./nohup.out 2>&1 &
sleep 1
tail -n 20 ./nohup.out
