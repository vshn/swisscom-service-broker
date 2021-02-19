FROM gcr.io/distroless/static:nonroot

COPY swisscom-service-broker /

ENTRYPOINT [ "/swisscom-service-broker" ]
