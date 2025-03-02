FROM docker.io/library/e2e-extension as ext
FROM ghcr.io/linuxsuren/api-testing:master

WORKDIR /workspace
COPY . .
COPY --from=ext /usr/local/bin/atest-store-redis /usr/local/bin/atest-store-redis

CMD [ "/workspace/entrypoint.sh" ]
