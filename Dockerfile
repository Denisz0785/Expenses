FROM golang:1.22
WORKDIR /app1
ADD go.mod .
COPY . .
RUN go build -o main main.go



CMD [ "/app1/main" ]
