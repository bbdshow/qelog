#!/bin/bash
process='qelog_admin'
ps -ef | grep -w ${process} | grep -v grep | awk  '{print "kill " $2}' | sh
