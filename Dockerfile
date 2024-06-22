# Build stage
FROM golang:1.22 as build

# Create working directory
WORKDIR /app

# Download and extract the piper binary
COPY ./download-piper.sh ./
RUN chmod +x ./download-piper.sh
RUN ./download-piper.sh

# Copy language models
COPY models /app/models

# Copy source code
COPY go.mod go.sum main.go /app/

# Build the application
RUN go build .

CMD ["./piper-tts"]
