# we're using Golang 1.x (1.16.3 at the time of writing)
#             and alpine
FROM golang:alpine

# - this section can be used to use the binary directly
#   instead of building it inside the container.
#   This binary will have to be built with
#   `CGO_ENABLED=0 go build`
# - you might need this if you're developing grofer (making
#   changes) and want to test these in docker. Building
#   grofer from source inside the container can be slow and
#   leads to a wastage of time as well as (internet) bandwith
#   due to `go get` being run everytime.
# - uncomment the following line and comment out the next
#   section to enable this
# ADD grofer /go/bin/

# - this section uses the current directory (grofer repo) and
#   builds grofer inside the container which is then used when
#   the container is run
ADD . /src
WORKDIR /src
RUN go install
WORKDIR /

# host's root (/) is assumed to be mounted at /host
# set env vars for gopsutil to find the right files
# refer: https://github.com/shirou/gopsutil#usage
ENV HOST_PROC /host/proc
ENV HOST_SYS /host/sys
ENV HOST_ETC /host/etc
ENV HOST_VAR /host/var
ENV HOST_RUN /host/run
ENV HOST_DEV /host/dev

ENTRYPOINT ["grofer"]
