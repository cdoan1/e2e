default: build

help::
	@echo "help"

build::
	ginkgo build

run::
	docker run --volume $(pwd)/config:/opt/.kube/config \
	--volume $(pwd):/results \
	--volume $(pwd)/resources/options.yaml:/resources/options.yaml \
	-e GINKGO_FOCUS="g0" \
	open-cluster-management-e2e:latest
