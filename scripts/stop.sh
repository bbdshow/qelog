#!/bin/bash
app='qelog'
ps -ef | grep -w ${app} | grep -v grep | awk  '{print "kill -9 " $2}' | sh
