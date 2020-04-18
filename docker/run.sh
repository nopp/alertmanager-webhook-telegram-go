#!/bin/bash

if [ "x" == "x$bottoken" ]; then
  echo "FAIL: bottoken is not set"
  exit 1
else
  echo "bottoken is set (not visible)"
fi

if [ "x" == "x$chatid" ]; then
  echo "FAIL: chatid is not set"
  exit 2
else
  echo "chatid is $chatid"
fi

sed -i s/xxbotTokenxx/"$bottoken"/ alert/alert.go
sed -i s/666777666/"$chatid"/ alert/alert.go

./main