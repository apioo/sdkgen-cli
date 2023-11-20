FROM golang:1.21
MAINTAINER Christoph Kappestein <christoph.kappestein@apioo.de>
LABEL description="SDKgen Code Generator"

ENV SDKGEN_CLIENT_ID ""
ENV SDKGEN_CLIENT_SECRET ""

VOLUME /usr/src/sdkgen/output
RUN mkdir -p /usr/src/sdkgen/output/sdk

WORKDIR /usr/src/sdkgen

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/sdkgen

WORKDIR /usr/src/sdkgen/output
CMD ["sh", "-c", "/usr/local/bin/sdkgen install --remove"]
