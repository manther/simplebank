From golang:1.22rc1-alpine3.19
WORKDIR /app
COPY . .
RUN go build -o main main.go

EXPOSE 8080

CMD [ "/app/main" ]