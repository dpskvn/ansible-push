ARG RUN_FROM=gcr.io/distroless/static-debian11@sha256:8ad6f3ec70dad966479b9fb48da991138c72ba969859098ec689d1450c2e6c97
FROM ${RUN_FROM} as image

COPY output/bin/push-connector-template /bin

USER nonroot:nonroot

ENTRYPOINT [ "/bin/push-connector-template" ]
