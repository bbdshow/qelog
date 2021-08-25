#!/bin/bash
program='qelog_admin'

ps xua| grep -w  $program | grep -v grep|awk '{print $2}'|xargs kill 

#健康检查
result="$?"

if [ $result -ne 0 ];then
	echo -e "\033[41;33m 检测到程序就没有运行，不在关闭程序 \033[0m"
else
	echo '程序关闭完成，休息1秒后检查程序是否结束成功.....'
	sleep 1
	n=1 #循环检测
	while [ $n -le 3 ] 
	do
	    status=`ps -ef | grep -w  $program |grep -v "grep" `
		if [ -n "$status" ];then 
			echo -e "\033[41;33m 检测到程序正在运行，休眠5秒后继续检测...  \033[0m"
			sleep 5
			let n++    #或者写作n=$(( $n + 1 ))
			if [ $n -eq 3 ];then
				echo -e "\033[41;33m 经过三次检测，程序依然没有退出，故退出自动部署脚本，请人工检查...  \033[0m"
				exit
			fi
		else
			echo "程序已经正常关闭了" 
			break 
		fi
	done
fi


