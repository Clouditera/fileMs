##!/bin/bash
#
#### 构建要求
#### 1. go >= 1.16
#
#VERSION=`echo 'v1.4'`
#INSTALL_DIR=`echo '/home/cloud/IntegrationTesting/bin'`
#BIN=`echo 'fileMS'`
#
#funcBuildfileMS() {
#  echo 'fileMS module building...'
#  go build -v -ldflags "-s -X main.Version=${VERSION}" -o ./ ./
#  echo 'fileMS module building...finished'
#}
#
#funcInstallfileMS() {
#  echo 'fileMS module installing...'
#  cp $BIN $INSTALL_DIR -f
#  echo 'fileMS module installing...finished'
#}
#
#funcRunfileMS() {
#  echo 'fileMS module starting...'
#  sudo $INSTALL_DIR/$BIN
#  echo 'fileMS module starting...finished'
#}
#
#funcBuildfileMS
#funcInstallfileMS
#funcRunfileMS
#
#echo 'done!!!'