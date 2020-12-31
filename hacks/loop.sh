#!/bin/bash

make build

for i in {1..120}
do
  ts=$(date +"%s")
  # echo $ts
  # ginkgo
  make run
  mv results/results.xml ./data/$ts.results.xml
  sleep 60
done
