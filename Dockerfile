FROM node:22.21.0-bookworm

RUN apt-get update && apt-get install -y wget ca-certificates \
    && wget https://go.dev/dl/go1.25.6.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go1.25.6.linux-amd64.tar.gz \
    && rm go1.25.6.linux-amd64.tar.gz

ENV PATH="/usr/local/go/bin:${PATH}"

WORKDIR /app

COPY server ./server
WORKDIR /app/server
RUN go mod download
RUN go build -o server

WORKDIR /app
COPY web ./web
WORKDIR /app/web
RUN npm install

WORKDIR /app
COPY run.sh .
RUN chmod +x run.sh

EXPOSE 8080
EXPOSE 5173

CMD ["./run.sh"]

