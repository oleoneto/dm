ARG DATABASE_URL
ARG MIGRATIONS_DIRECTORY
ARG MIGRATIONS_TABLE

FROM golang:alpine AS build

WORKDIR /app

ADD . /app

RUN \
  apk update; apk add build-base; \
  cd /app; \
  make build-for-docker; \
  make build-api; \
  rm -rf /var/cache/apk/*

# --------------------------

FROM alpine

ARG DATABASE_URL
ARG MIGRATIONS_DIRECTORY
ARG MIGRATIONS_TABLE

ENV DATABASE_URL ${DATABASE_URL}
ENV MIGRATIONS_DIRECTORY ${MIGRATIONS_DIRECTORY}
ENV MIGRATIONS_TABLE ${MIGRATIONS_TABLE}

WORKDIR /app

COPY --from=build /app/dm* /app

EXPOSE 3809

CMD [ "./dm-api" ]
