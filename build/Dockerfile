FROM golang:1.14 as BUILD
WORKDIR /tests
COPY . .
ENTRYPOINT ["go", "test", "-v", "./...", "-tags", "integration", "-coverprofile", "cover.out"]
