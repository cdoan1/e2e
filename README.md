# e2e

ginkgo based end to end testing for Red Hat Open Cluster Management

## Setup

The `resources/options.yaml` specifies the target cluster to run the test against.
If this file does not exist, the cluster settings will be taken from environment variables.

### From Source

```bash
git clone 
vi resources/options.yaml
ginkgo
```

### From Container

Run the ginkgo test from docker container. 

* Only `headless` mode is supported when running from docker container.
* Screen shot is taken of reach view.

```bash
make
```
