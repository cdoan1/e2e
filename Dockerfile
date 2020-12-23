FROM golang:1.15-alpine as build
MAINTAINER cdoan <cdoan@redhat.com>

RUN apk update && apk add --no-cache git gcc bash musl-dev xvfb chromium chromium-chromedriver

RUN go get -u github.com/onsi/ginkgo/ginkgo && go get -u github.com/onsi/gomega/...

WORKDIR /go/src/open-cluster-management-e2e
COPY . .

RUN ginkgo build

FROM golang:1.15-alpine
MAINTAINER cdoan <cdoan@redhat.com>

RUN apk update && apk add --no-cache git gcc bash musl-dev xvfb chromium chromium-chromedriver

RUN go get -u github.com/onsi/ginkgo/ginkgo && go get -u github.com/onsi/gomega/...

WORKDIR /go/src/open-cluster-management-e2e
COPY . .

COPY --from=build /go/src/open-cluster-management-e2e/open-cluster-management-e2e.test /go/src/open-cluster-management-e2e

# CMD ["bash", "-c", "./entrypoint.sh"]
# CMD ["bash", "-c", "ginkgo open-cluster-management-e2e.test"]
ENTRYPOINT [ "./entrypoint.sh" ]
