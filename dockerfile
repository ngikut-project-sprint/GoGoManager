# # Gunakan Ubuntu 20.04 sebagai base image untuk platform linux/amd64
# FROM --platform=linux/amd64 ubuntu:20.04
FROM golang:1.23.4

# # Perbarui repositori, instal dependensi dasar, unduh, dan pasang Go
# RUN apt-get update && apt-get install -y \
#     wget \
#     tar \
#     git \
#     build-essential \
#     && apt-get clean \
#     && rm -rf /var/lib/apt/lists/*

# Atur direktori kerja
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o main ./main.go

RUN chmod +x main

EXPOSE 8080

CMD ["./main"]

#docker build -t my-app .
#docker run -p 8080:8080 -tid 83f3195e5e3b
#docker cp f5b125250c7c:/app/main .
#sudo lsof -i -P -n
#GOARCH=amd64 GOOS=linux go build -o main ./main.go    