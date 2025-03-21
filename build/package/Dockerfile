## Builder
FROM public.ecr.aws/docker/library/golang:1.22.2-alpine3.19 as build_base

RUN apk update && apk upgrade && \
    apk --no-cache --update add bash git make gcc g++ libc-dev

WORKDIR /go/src/github.com/ocuris/go-template-backend
ENV GO111MODULE=on
ENV GODEBUG="madvdontneed=1"
ARG GITHUB_TOKEN

# Set the time zone data
RUN apk --no-cache add tzdata
ENV TZ=Europe/Madrid

# We want to populate the module cache based on the go.mod and go.sum files.
COPY go.mod .
COPY go.sum .

RUN git config --global url."https://${GITHUB_TOKEN}:x-oauth-basic@github.com/ocuris".insteadOf "https://github.com/ocuris"
RUN go mod download

FROM build_base AS server_builder

WORKDIR /go/src/github.com/ocuris/go-template-backend
COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 make engine

# In this last stage, we start from a fresh alpine image, to reduce the image size and the Go compiler in our production artifacts.
FROM alpine:latest
RUN apk update && apk upgrade && \
    apk --no-cache --update add tzdata

WORKDIR /go-template-backend/app
EXPOSE 11001

COPY --from=server_builder /go/src/github.com/ocuris/go-template-backend/engine .
COPY --from=server_builder /go/src/github.com/ocuris/workflow-engine/configs ./configs/
COPY --from=server_builder /go/src/github.com/ocuris/workflow-engine/files/ ./files/

CMD /workflow-engine/app/engine
