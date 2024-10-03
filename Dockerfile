# # Gunakan image Go versi 1.21.1
# FROM golang:1.21.1 AS builder

# WORKDIR /app

# # Salin file go.mod dan plugin
# COPY go.mod .
# COPY response_logger.go .

# # Tampilkan isi direktori
# RUN ls -la

# # Tampilkan isi file go.mod
# RUN cat go.mod

# # Tampilkan isi file response_logger.go
# RUN cat response_logger.go

# # Download dependensi dan build plugin
# RUN go mod download && \
#     go build -buildmode=plugin -v -o response_logger.so response_logger.go

# # Periksa apakah plugin berhasil dibuat
# RUN ls -l response_logger.so

# Gunakan image KrakenD official
FROM devopsfaith/krakend:latest AS runner

# Buat direktori untuk plugin jika belum ada
RUN mkdir -p /etc/krakend/plugins/

# Salin plugin yang sudah di-build
# COPY --from=builder /app/response_logger.so /etc/krakend/plugins/
COPY krakend.json /etc/krakend/krakend.json
COPY response_logger.so /etc/krakend/plugins/response_logger.so

# Periksa apakah plugin berhasil disalin
RUN ls -l /etc/krakend/plugins/

# Set environment variable untuk mengaktifkan plugin
ENV KRAKEND_PLUGIN_LOAD=1

# Jalankan KrakenD
CMD [ "run", "-c", "/etc/krakend/krakend.json" ]