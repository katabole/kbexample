# This Dockerfile provides the basics needed for deployment or docker-based builds.
# You'll need to modify it to suite your deployment needs.

FROM golang:1.24.7-alpine3.22 AS go-builder

WORKDIR /src

# Install node for webpack bundling. Also psql since for now we use it to create test databases.
RUN apk add --no-cache postgresql-client nodejs npm

# Before we add the whole project in, download just the dependencies, so we can cache that layer despite app changes.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN npm install
RUN go build -o /bin/app


FROM alpine:3.21 AS app

COPY --from=go-builder /bin/app /app

ENV DEPLOY_ENV=production ADDR=0.0.0.0
EXPOSE 3000

# If you want to run DB migrations using this container, you'll need the following changes:
# - Install atlas in the build container: RUN curl -sSf https://atlasgo.sh | sh -s -- -y
# - Copy the migrations folder (or schema.sql file) into the container: COPY migrations /migrations
# - Run the migrations before starting the app: CMD ["sh", "-c", "/atlas migrate apply --dir file://migrations --url $DATABASE_URL && /app"]

CMD ["/app"]
