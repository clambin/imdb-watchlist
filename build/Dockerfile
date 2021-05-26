FROM alpine

WORKDIR /app

ARG TARGETOS
ARG TARGETARCH

COPY bld/imdb-watchlist-$TARGETOS-$TARGETARCH /app/imdb-watchlist

EXPOSE 8080
ENTRYPOINT ["/app/imdb-watchlist"]
CMD []
