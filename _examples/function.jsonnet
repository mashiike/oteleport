{
  Architectures: [
    'arm64',
  ],
  EphemeralStorage: {
    Size: 512,
  },
  FunctionName: 'oteleport',
  Handler: 'bootstrap',
  LoggingConfig: {
    LogFormat: 'Text',
    LogGroup: '/aws/lambda/oteleport',
  },
  MemorySize: 128,
  Role: 'arn:aws:iam::314472643515:role/oteleport-lambda',
  Runtime: 'provided.al2023',
  SnapStart: {
    ApplyOn: 'None',
  },
  Timeout: 3,
  TracingConfig: {
    Mode: 'PassThrough',
  },
}
