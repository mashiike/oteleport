{
  storage: {
    cursor_encryption_key: 'r0JwTGIzoOpTi+gH9t+6i/kIwxDi7kR23uwKAeSxxEE=',
    location: 's3://oteleport-local/',
    aws: {
      endpoint: 'http://localhost:9000',
      use_s3_path_style: true,
      credentials: {
        access_key_id: 'oteleport0000',
        secret_access_key: 'oteleport0000',
      },
    },
  },
  otlp: {
    grpc: {
      enable: true,
      address: '0.0.0.0:4317',
    },
    http: {
      enable: true,
      address: '0.0.0.0:4318',
    },
  },
  api: {
    http: {
      enable: true,
      address: '0.0.0.0:8080',
    },
  },
}
