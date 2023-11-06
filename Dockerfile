FROM golang:1.21
MAINTAINER Christoph Kappestein <christoph.kappestein@apioo.de>
LABEL description="SDKgen Code Generator"

ENV SDKGEN_TYPE "client-java"
ENV SDKGEN_CLIENT_ID ""
ENV SDKGEN_CLIENT_SECRET ""
ENV SDKGEN_BASE_URL ""
ENV SDKGEN_NAMESPACE "app.sdkgen"

VOLUME /usr/src/sdkgen/output

WORKDIR /usr/src/sdkgen

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/sdkgen

CMD ["sh", "-c", "/usr/local/bin/sdkgen generate $SDKGEN_TYPE output/typeapi.json output --client-id=\"$SDKGEN_CLIENT_ID\" --client-secret=\"$SDKGEN_CLIENT_SECRET\" --base-url=\"$SDKGEN_BASE_URL\" --namespace=\"$SDKGEN_NAMESPACE\""]
