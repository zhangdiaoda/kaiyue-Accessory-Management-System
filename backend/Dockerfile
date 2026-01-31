FROM alpine:latest
WORKDIR /app
RUN apk --no-cache add ca-certificates tzdata curl
ENV TZ=Asia/Shanghai
ARG TARGETARCH
COPY warehouse-server-${TARGETARCH} ./warehouse-server
COPY config/ ./config/
RUN chmod +x ./warehouse-server
EXPOSE 8080
CMD ["./warehouse-server"]
