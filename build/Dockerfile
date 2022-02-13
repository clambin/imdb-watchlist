FROM alpine

WORKDIR /app

ARG TARGETOS
ARG TARGETARCH

COPY bld/imdb-watchlist-$TARGETOS-$TARGETARCH /app/imdb-watchlist

RUN /usr/sbin/addgroup app
RUN /usr/sbin/adduser app -G app -D
USER app

EXPOSE 8080
ENTRYPOINT ["/app/imdb-watchlist"]
CMD []
