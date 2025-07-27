FROM golang:1.23-alpine

# Install git
RUN apk add --no-cache git

# Clone the codeforgeai.go repo
RUN git clone https://github.com/codeforge-ide/codeforgeai.go /app

WORKDIR /app

# Download Go dependencies
RUN go mod download

# Build the codeforgeai binary
RUN go build -o /usr/local/bin/codeforgeai .

# Set the default command to run codeforgeai
ENTRYPOINT ["codeforgeai"]
