# The version parameter in the url. (i.e v1)
ARG API_VERSION
# The namespace of the migrations resources.(i.e migrations)
ARG API_NAMESPACE
# The connection string to the database. (i.e postgres://<username>:<password>@host:5432/database)
ARG DATABASE_URL
# The directory where the migrator should look for migration fiels. (i.e migrations or ./)
ARG MIGRATIONS_DIRECTORY
# The table in the database where migration state information will be recorded.
ARG MIGRATIONS_TABLE

FROM golang:alpine AS build

WORKDIR /app

ENV PATH /app:$PATH

ADD . /app

RUN \
  apk update; apk add ca-certificates build-base git; \
  make build-for-docker; \
  rm -rf /var/cache/apk/*

# --------------------------

FROM alpine

ARG API_VERSION
ARG API_NAMESPACE
ARG DATABASE_URL
ARG MIGRATIONS_DIRECTORY
ARG MIGRATIONS_TABLE

ENV API_VERSION=${API_VERSION}
ENV API_NAMESPACE=${API_NAMESPACE}
ENV DATABASE_URL=${DATABASE_URL}
ENV MIGRATIONS_DIRECTORY=${MIGRATIONS_DIRECTORY}
ENV MIGRATIONS_TABLE=${MIGRATIONS_TABLE}

ENV PATH /app:$PATH

WORKDIR /app

COPY --from=build /app/dm /app

EXPOSE 3809

CMD [ "dm", "api", "-o", "json", "-p", "3809"]
