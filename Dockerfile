FROM golang:1.10.3-alpine3.7 AS build

# Install required tools for the project (both to build and to run). This is the requirements section
# Run `trash --file vendor.yaml` to get application dependencies
RUN apk add --no-cache git

# Copy the entire project to the image directory
COPY . /go/src/github.com/tqoliver/empsvc/
WORKDIR /go/src/github.com/tqoliver/empsvc/

RUN CGO_ENABLED=0 GOOS=linux go build -o /empsvc

# scratch is the smallest available container size
FROM scratch
# RUN apk update && apk add curl && apk add mysql-client && apk add sudo \
#        && addgroup -g 1001 -S exampleapp \
#        && adduser -u 1001 -D -S -G exampleapp exampleapp
COPY --from=build /empsvc /
EXPOSE 8000
USER 1001
ENTRYPOINT [ "/empsvc" ]