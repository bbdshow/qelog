#!/bin/bash
# 放入到对应目录，修改 program = 对应程序名
program='qelog_admin'
#新程序指定启动配置文件,当当前目录不存在config文件夹时，下面的configfile才有意义
configfile='../configs/config.toml'

if [ ! -f nohup.out ]; then
  touch nohup.out
fi

if [ ! -d log ]; then
  mkdir log
fi

chown -R gouser:gouser  $program  nohup.out  log
chmod 700  $program
workspace=`pwd`
status=`ps -ef | grep -w $program |grep -v "grep" `
if [ -n "$status" ]
	then 
	echo "!!!!!!!!  注意：检测到程序正在运行，故本脚本退出，本次未执行启动命令  !!!!!!!" 
	exit
fi


if [  -d config ]; then #老程序启动脚本 
	chown -R gouser:gouser   config 
	su - gouser -c "cd $workspace; nohup ./$program  >> nohup.out 2>&1  &"
elif [ -d configs ]; then #老程序启动脚本 
	chown -R gouser:gouser   configs 
	su - gouser -c "cd $workspace; nohup ./$program  >> nohup.out 2>&1  &"
else  #新程序启动，指定配置文件
	su - gouser -c "cd $workspace; nohup ./$program -f $configfile >> nohup.out 2>&1  &"
fi


#健康检查
echo '程序启动完成，休息3秒后检查程序是否正在运行......'
sleep 3
status=`ps -ef | grep -w  $program |grep -v "grep" `
if [ -n "$status" ];then 
	echo "健康检查：检测到程序正在运行，程序状态正常" 
else
	echo "健康检查：检测到程序没有运行，程序状态异常，请人工检查" 
fi
