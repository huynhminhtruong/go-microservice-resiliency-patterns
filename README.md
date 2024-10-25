# 1. Design client-server communication in microservice system:

## Resiliency and fault tolerant:

## 1. Timeout pattern:

## 2. Retry pattern:

- middleware means interceptor in gRPC:

  - https://github.com/grpc-ecosystem/go-grpc-middleware

- In gRPC, there are two types of interceptor usage:
  - WithUnaryInterceptor for unary connection
  - WithStreamingInterceptor for streaming connection respectively

## 3. Circuit breaker pattern:

- Go circuit breaker:

  - https://github.com/sony/gobreaker

- Các đặc trưng khi thiết kế mô hình clien-server theo circuit-breaker pattern:

  - The maximum allowed failure (failed_request_count / total_request_count) ratio is 0.6
  - The MaxRequests allowed during the half-open state is 3
  - The timeout value needed for state transition from open to half-open is 4 seconds
  - Print a log statement whenever the state changes

- Các thuật ngữ trong circuit breaker pattern:

  - MaxRequests:

    - Giới hạn số lượng yêu cầu được phép thông qua khi mạch ở trạng thái half-open (bán mở)
    - Khác biệt giữa mạch mở và bán mở là mạch bán mở cho phép một số yêu cầu được thông qua dựa trên cấu hình của bạn, trong khi mạch mở sẽ chặn toàn bộ các yêu cầu

  - Interval:

    - Xác định khoảng thời gian mà số lượng lỗi được tính trong khi mạch ở trạng thái closed (đóng)
    - Đối với một số trường hợp, bạn có thể không muốn xóa bộ đếm lỗi ngay cả khi đã qua một thời gian dài kể từ lần lỗi cuối
    - Tuy nhiên, trong hầu hết các trường hợp, bộ đếm sẽ được xóa để cho phép các lỗi trong một khoảng thời gian hợp lý

  - Timeout:

    - Quyết định thời gian khi nào mạch sẽ chuyển từ trạng thái open (mở) sang trạng thái half-open
    - Ưu điểm của sự thay đổi này là khi mạch ở trạng thái mở, các yêu cầu sẽ thất bại nhanh chóng
    - Trong khi đó, ở trạng thái bán mở, circuit breaker cho phép thông qua một số lượng yêu cầu giới hạn

  - ReadyToTrip:

    - Kiểm tra ngưỡng lỗi sau lần lỗi cuối cùng và quyết định xem mạch có nên chuyển sang trạng thái open hoàn toàn hay không
    - Đây là điều kiện quyết định khi nào mạch sẽ được mở dựa trên số lần thất bại đã xảy ra

  - OnStateChange:
    - Được sử dụng để theo dõi các thay đổi trạng thái khi xử lý các mô hình kinh doanh trong một hàm bao bọc
    - Điều này hữu ích để log hoặc kích hoạt các sự kiện dựa trên việc thay đổi trạng thái của circuit breaker

# 2. TLS and mTLS handshake:

## 1. TLS:

- The Order service connects the Payment service
- The Payment service shows its certificate to the Order service
- The Order service verifies that certificate
- The Order service sends data to the Payment service in an encrypted communication channel

## 2. mTLS:

- The Order service connects to the Payment service
- The Payment service shows its certificate to the Order service
- The Order service verifies that certificate
- The Order service shows its certificate to the Payment service
- The Payment service verifies that certificate and allows the Order service to send requests
- The Order service sends data to the Payment service in an encrypted communication channel

# 3. Summary:

- In a typical microservices architecture, it is normal for one service to depend on one or more other services. If one of the dependent services is down, it will affect the availability of the entire microservice application. We use resiliency patterns such as retry, timeout, and circuit breaker to prevent these situations

- Once the dependent service is down, we can use the retry strategy on the consumer side to make gRPC communication eventually succeed

- Retry logic can be triggered for certain status codes, such as Unavailable or ResourceExhausted, instead of blindly retrying on each failure. For example, it is not wise to retry a request if you received a validation exception because it shows that you sent an invalid payload; you should fix it to make it succeed

- Using retry logic blocks the actual gRPC call since you have to wait for dependent service. It is hard to detect recovery time for a dependent service, which can create long wait times for retries. To prevent this situation, we use context timeout and deadline to put a time limit on blocking the actual execution

- In the retry mechanism, you can redo an operation for specific time intervals, but this can also put an extra load on the dependent service since it has to retry all the time, even if the dependent service is not ready. To solve this problem, we use a circuit breaker to open a circuit once we reach the failure limit, retry the request after some time, and finally close the circuit once the dependent service is back online

- Error handling is important in interservice communication because the next step is decided from error codes or messages in the gRPC response. We use a status package to return customized errors from a service, and we can convert them on the client side once needed

- Resiliency is important not only in communication patterns, but also in zerotrust environments, in which we use TLS-enabled communications, as the server and client verify their certificates during gRPC communications. This is also called mutual TLS
