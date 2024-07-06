FROM ubuntu:latest
RUN apt-get update && apt-get install -y build-essential gcc libc6-dev sqlite3 curl git
RUN curl -sL https://dl.google.com/go/go1.18.linux-amd64.tar.gz | tar -zxv -C /usr/local
ENV PATH=$PATH:/usr/local/go/bin
ENV CGO_ENABLED=1
WORKDIR /app
COPY . .
RUN go mod download && go mod verify
RUN make build
RUN make initdb
WORKDIR /app
EXPOSE 8080
CMD ["make", "server"]