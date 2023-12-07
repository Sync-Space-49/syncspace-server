FROM golang:1.20-alpine 
WORKDIR /app
COPY . ./
# This is where one could build the application code as well.
RUN go build -o ./main

# Run on container startup.
CMD ["./main"]