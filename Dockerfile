FROM registry.ci.openshift.org/stolostron/builder:go1.18-linux AS builder

ARG OS=linux
ARG ARCH=amd64
WORKDIR /go/src/open-cluster-management.io/registration
COPY . .
ENV GO_PACKAGE open-cluster-management.io/registration

RUN GOOS=${OS} \
    GOARCH=${ARCH} \
    make build --warn-undefined-variables

FROM registry.access.redhat.com/ubi8/ubi-minimal:latest
ENV USER_UID=10001

COPY --from=builder /go/src/open-cluster-management.io/registration/registration /
COPY --from=builder /go/src/open-cluster-management.io/registration/webhook /

USER ${USER_UID}
