FROM golang:1.20.4-alpine3.18 as buildStage
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN go build -o main main.go 

RUN apk add --no-cache curl \
    && curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.0/migrate.linux-amd64.tar.gz | tar xvz \
    && apk del curl


FROM alpine:3.18
WORKDIR /app 
COPY --from=buildStage /app/main .
COPY --from=buildStage /app/migrate .
COPY --from=buildStage /app/config.yaml .
COPY --from=buildStage /app/init.sh .
COPY --from=buildStage /app/wait-for.sh .
COPY --from=buildStage /app/db/migrations ./migrations

EXPOSE 8080
CMD [ "init.sh", "/app/main" ]
