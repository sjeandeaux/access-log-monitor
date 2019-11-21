FROM golang:1.13.4-alpine as base

RUN apk --no-cache add ca-certificates make git g++
ENV GO111MODULE on

#get source
WORKDIR /go/src/github.com/sjeandeaux/access-log-parsor

#copy the source
COPY . .

RUN make tools
RUN make dependencies

## test
FROM base AS test
RUN make test

#Build the application
FROM base AS build
RUN make build


FROM scratch AS release

ARG BUILD_VERSION=undefined
ARG BUILD_DATE=undefined
ARG VCS_REF=undefined

#http://label-schema.org/rc1/
LABEL "maintainer"="stephane.jeandeaux@gmail.com" \
      "org.label-schema.vendor"="sjeandeaux" \
      "org.label-schema.schema-version"="1.0.0-rc.1" \
      "org.label-schema.applications.access-log-parsor.version"=${BUILD_VERSION} \
      "org.label-schema.vcs-ref"=$VCS_REF \
      "org.label-schema.build-date"=${BUILD_DATE}

COPY --from=build /go/src/github.com/sjeandeaux/access-log-parsor/target/access-log-parsor /access-log-parsor

VOLUME [ "/tmp" ]
CMD ["watch"]
ENTRYPOINT ["/access-log-parsor"]