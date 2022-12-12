#!/bin/bash
app='qelog_receiver'
ps -ef | grep -w ${app} | grep -v grep | awk  '{print "kill -9 " $2}' | sh
