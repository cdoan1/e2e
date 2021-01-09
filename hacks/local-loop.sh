#!/bin/bash

for i in {1..180}
do
  ts=$(date +"%s")
  # echo $ts
  # ginkgo
  ginkgo -focus="2.1"
  mv results.xml ./data/$ts.results.xml
  sleep 60
done
