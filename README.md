# oteleport

`oteleport` is an OpenTelemetry Signal Receiver Server that stores received signals in S3 and provides a REST API to retrieve them. The project focuses on "teleporting" OpenTelemetry signals, acting as a buffering server that enables custom signal replays. It is designed to offer flexibility in how signals are managed and retrieved, making it easy to handle and replay telemetry data as needed.

## Features

- **OpenTelemetry Signal Receiver**: Receives OpenTelemetry signals and stores them in S3.
- **REST API**: Provides a REST API to retrieve stored signals.

## Installation

Pre-built release binaries are provided.
Download the latest binary from the [Releases page](https://github.com/mashiike/cflog2otel/releases).

## Usage

simple config file example. `oteleport.jsonnet`
```jsonnet
local must_env = std.native('must_env');

{
  access_keys: [
    must_env('OTELEPORT_ACCESS_KEY'),
  ],
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
```

set environment variables for your S3 bucket and encryption key.

```shell
$ oteleport --config oteleport.jsonnet
```

and send signals to the oteleport. for example using [`otel-cli`](https://github.com/equinix-labs/otel-cli).

```shell
$ otel-cli exec --protocol grpc --endpoint http://localhost:4317/ --service my-service --otlp-headers Oteleport-Access-Key=$OTELEPORT_ACCESS_KEY --name "curl google" curl https://www.google.com
```

after get traces from the oteleport.

```
$ curl -X POST -H 'Content-Type: application/json' -H "Oteleport-Access-Key: $OTELEPORT_ACCESS_KEY" -d "{\"startTimeUnixNano\":$(date -v -5M +%s)000000000, \"limit\": 100}" http://localhost:8080/api/traces/fetch | jq
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100  1582  100  1527  100    55  87970   3168 --:--:-- --:--:-- --:--:-- 93058
{
  "resourceSpans": [
    {
      "resource": {
        "attributes": [
          {
            "key": "service.name",
            "value": {
              "stringValue": "my-service"
            }
          }
        ]
      },
      "schemaUrl": "https://opentelemetry.io/schemas/1.17.0",
      "scopeSpans": [
        {
          "schemaUrl": "https://opentelemetry.io/schemas/1.17.0",
          "scope": {
            "name": "github.com/equinix-labs/otel-cli",
            "version": "0.4.5 0d4b8a9c49f60a6fc25ed22863259ff573332060 2024-04-01T20:56:07Z"
          },
          "spans": [
            {
              "attributes": [
                {
                  "key": "process.command",
                  "value": {
                    "stringValue": "curl"
                  }
                },
                {
                  "key": "process.command_args",
                  "value": {
                    "arrayValue": {
                      "values": [
                        {
                          "stringValue": "curl"
                        },
                        {
                          "stringValue": "https://www.google.com"
                        }
                      ]
                    }
                  }
                },
                {
                  "key": "process.owner",
                  "value": {
                    "stringValue": "ikeda-masashi"
                  }
                },
                {
                  "key": "process.pid",
                  "value": {
                    "intValue": "23539"
                  }
                },
                {
                  "key": "process.parent_pid",
                  "value": {
                    "intValue": "23540"
                  }
                }
              ],
              "endTimeUnixNano": "1729674137524531000",
              "kind": 3,
              "name": "curl google",
              "spanId": "EB1F5DC940D26FD8",
              "startTimeUnixNano": "1729674137341884000",
              "status": {},
              "traceId": "D69015A574B485137570DB01CC5B8D7D"
            },
            {
              "attributes": [
                {
                  "key": "process.command",
                  "value": {
                    "stringValue": "curl"
                  }
                },
                {
                  "key": "process.command_args",
                  "value": {
                    "arrayValue": {
                      "values": [
                        {
                          "stringValue": "curl"
                        },
                        {
                          "stringValue": "https://www.google.com"
                        }
                      ]
                    }
                  }
                },
                {
                  "key": "process.owner",
                  "value": {
                    "stringValue": "mashiike"
                  }
                },
                {
                  "key": "process.pid",
                  "value": {
                    "intValue": "00000"
                  }
                },
                {
                  "key": "process.parent_pid",
                  "value": {
                    "intValue": "00001"
                  }
                }
              ],
              "endTimeUnixNano": "1729677243239248000",
              "kind": 3,
              "name": "curl google",
              "spanId": "9F53FAE65642617D",
              "startTimeUnixNano": "1729677242838971000",
              "status": {},
              "traceId": "44E4B5F99D51AEACBF0AB243C01849A5"
            }
          ]
        }
      ]
    }
  ]
}
```

## Usage as AWS Lambda function

`oteleport` can be used as an AWS Lambda function bootstrap.
this Lambda function is triggered by Lambda Function URL. work as a `http/otlp` and rest api endpoint. 

See [_examples](./_examples) dir for more details.
Include terraform code and [lambroll](https://github.com/fujiwara/lambroll) configuration.


## License

This project is licensed under the MIT License. 
See the [LICENSE](./LICENSE) file for more details.
