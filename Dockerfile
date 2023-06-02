FROM golang:1.20.4-alpine3.18 as buildStage
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN go build -o main main.go 

FROM alpine:3.18
WORKDIR /app 
COPY --from=buildStage /app/main .
COPY --from=buildStage /app/config.yaml .
COPY --from=buildStage /app/wait-for.sh .
COPY --from=buildStage /app/db/migrations ./db/migrations

EXPOSE 8080
CMD [ "/app/main" ]
