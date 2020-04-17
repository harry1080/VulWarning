#!/usr/bin/env bash

# Author : Virink <virink@outlook.com>
# Date   : 2019/11/05, 16:00

ssh secscanner 'sudo supervisorctl stop vulwarning'
ssh secscanner 'rm /opt/sss/vulwarning/vulwarning'
# upx -9 build/vulwarning-linux-amd64/vulwarning
scp build/vulwarning-linux-amd64/vulwarning secscanner:/opt/sss/vulwarning/vulwarning
# scp ./config.conf secscanner:/opt/metis/config.conf
ssh secscanner 'sudo supervisorctl start vulwarning'
