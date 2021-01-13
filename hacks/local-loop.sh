#!/bin/bash

export END=${1:-180}

for i in $(seq 1 $END);
do
  ts=$(date +"%s")
  ginkgo -focus="2.1"
  mv results.xml ./data/$ts.results.xml
  sleep 60
done
