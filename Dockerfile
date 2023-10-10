FROM golang:1.21
MAINTAINER Christoph Kappestein <christoph.kappestein@apioo.de>
LABEL description="SDKgen Code Generator"

ENV TYPE "client-typescript"
ENV CLIENT_ID ""
ENV CLIENT_SECRET ""
ENV BASE_URL ""
ENV NAMESPACE ""

VOLUME /usr/src/sdkgen/output

WORKDIR /usr/src/sdkgen

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/sdkgen ./...

CMD ["sdkgen", "generate", "$TYPE", "output/typeapi.json", "output", "--client-id", "$CLIENT_ID", "--client-secret", "$CLIENT_SECRET", "--base-url", "$BASE_URL", "--namespace", "$NAMESPACE"]
