ARG RUN_FROM=gcr.io/distroless/static-debian11@sha256:8ad6f3ec70dad966479b9fb48da991138c72ba969859098ec689d1450c2e6c97
FROM ${RUN_FROM} AS image

COPY output/bin/ansible-automation-platform-connector /bin

USER nonroot:nonroot

ENTRYPOINT [ "/bin/ansible-automation-platform-connector" ]
