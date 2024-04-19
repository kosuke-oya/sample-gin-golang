FROM golang:1.21-alpine as dev

WORKDIR /app

COPY ./ ./
RUN go mod download
RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN go install golang.org/x/tools/gopls@latest
RUN go install github.com/swaggo/swag/cmd/swag@latest


ARG ENV_KEY
ENV ENV_KEY=${ENV_KEY}

# RUN go build -o /httpserver

FROM alpine:latest as production
ARG ENV_KEY
ENV ENV_KEY=${ENV_KEY}
COPY --from=dev /httpserver ./httpserver
WORKDIR /app
ENV PORT=8080
ENTRYPOINT [ "/httpserver" ]

