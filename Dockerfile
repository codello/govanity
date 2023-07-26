FROM scratch

COPY build/govanity /govanity

ENTRYPOINT ["/govanity", "server"]
