#!/bin/bash

for i in {1..10}
do
  ts=$(date +"%s")
  echo $ts
  # ginkgo
  make run
  mv results/results.xml ./data/$ts.results.xml
  sleep 60
done
