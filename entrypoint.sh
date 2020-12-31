#!/bin/bash

echo ""
echo "🌅 Version Info"
echo "------------------"

go version
ginkgo version
chromedriver --version
echo "------------------"

echo "😊 Starting ginkgo test ..."
ginkgo open-cluster-management-e2e.test
cp results.xml results
