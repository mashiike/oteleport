{
  storage: {
    cursor_encryption_key: 'r0JwTGIzoOpTi+gH9t+6i/kIwxDi7kR23uwKAeSxxEE=',
    location: 's3://oteleport-test/',
    flatten: true,
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
