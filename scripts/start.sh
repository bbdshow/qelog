#!/bin/bash
process='qelog_admin'
nohup ./${process} -f ../configs/config.toml > ./nohup.out 2>&1 &
sleep 1
tail -n 50 ./nohup.out
