#!/bin/bash

for i in {1..120}
do
  ts=$(date +"%s")
  echo $ts
  ginkgo
  mv results.xml ./data/$ts.results.xml
  sleep 60
done
