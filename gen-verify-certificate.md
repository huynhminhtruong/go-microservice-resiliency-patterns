# 1. Tạo CA (Certificate Authority)

- Bước 1: Tạo khóa riêng cho CA

  - openssl genpkey -algorithm RSA -out ca.key -aes256

  - ca.key: Tệp chứa private key của CA
  - aes256: Mã hóa private key bằng AES-256 để bảo mật

- Bước 2: Tạo certificate cho CA

  - openssl req -new -x509 -key ca.key -sha256 -days 365 -out ca.crt

  - ca.crt: Tệp certificate của CA
  - -days 365: Thời hạn certificate là 1 năm

# 2. Tạo certificate cho Server

- Certificate được ký bởi CA

- Bước 1: Tạo khóa riêng cho server
  openssl genpkey -algorithm RSA -out server.key -aes256

- Bước 2: Tạo CSR (Certificate Signing Request) cho server
  openssl req -new -key server.key -out server.csr

- Bước 3: Ký certificate cho server bằng CA
  openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 365 -sha256

- -CAcreateserial: Tạo file serial number cho CA (thường là ca.srl)

# 3. Tạo certificate cho Client

- Certificate được ký bởi CA

- Bước 1: Tạo khóa riêng cho client
  openssl genpkey -algorithm RSA -out client.key -aes256

- Bước 2: Tạo CSR cho client
  openssl req -new -key client.key -out client.csr

- Bước 3: Ký certificate cho client bằng CA
  - openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client.crt -days 365 -sha256

# 4. Thiết lập TLS trên server và client

- Server
  Sử dụng server.crt và server.key để thiết lập TLS trên server
  Cung cấp ca.crt cho client để client có thể xác minh certificate của server

- Client
  Sử dụng client.crt và client.key để thiết lập mTLS (Mutual TLS) trên client
  Cung cấp ca.crt cho server để server có thể xác minh certificate của client

# 5. Xác minh TLS và mTLS

- Kiểm tra certificate của server (TLS)
  openssl s_client -connect server_domain:443 -CAfile ca.crt
  "Lệnh này kết nối tới server qua TLS, xác minh certificate của server bằng CA"

- Kiểm tra mTLS giữa client và server
  openssl s_client -connect server_domain:443 -cert client.crt -key client.key -CAfile ca.crt
  "Lệnh này kết nối tới server qua mTLS, yêu cầu cả server và client xác minh certificate của nhau"

# 6. Xác minh certificate

- Xác minh certificate của server
  openssl verify -CAfile ca.crt server.crt

- Xác minh certificate của client
  openssl verify -CAfile ca.crt client.crt

- Lệnh trên kiểm tra xem certificate của server/client có được ký bởi CA hợp lệ hay không

- TLS: Server chứng thực với client thông qua CA
- mTLS: Cả client và server đều chứng thực lẫn nhau thông qua CA
