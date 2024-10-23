local must_env = std.native('must_env');

{
  storage: {
    cursor_encryption_key: must_env('OTELEPORT_CURSOR_ENCRYPTION_KEY'),
    location: 's3://' + must_env('OTELEPORT_S3_BUCKET') + '/',
  },
  otlp: {
    grpc: {
      enable: true,
      address: '0.0.0.0:4317',
    },
  },
  api: {
    http: {
      enable: true,
      address: '0.0.0.0:8080',
    },
  },
}
