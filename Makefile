default: run

help::
	@echo "help"

build::
	docker build -t e2e .

run::
	docker run -it --rm \
	--volume "$$(pwd)"/results:/go/src/open-cluster-management-e2e/results \
	--volume "$$(pwd)"/resources/options.yaml:/go/src/open-cluster-management-e2e/resources/options.yaml \
	--name e2e e2e

clean::
	rm -rf ./results/*.png

test::
	./hacks/loop.sh

aggregate::
	python etl/junit2csv.py

