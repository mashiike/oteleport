# oteleport

`oteleport` is an OpenTelemetry Signal Receiver Server that stores received signals in S3 and provides a REST API to retrieve them. The project focuses on "teleporting" OpenTelemetry signals, acting as a buffering server that enables custom signal replays. It is designed to offer flexibility in how signals are managed and retrieved, making it easy to handle and replay telemetry data as needed.

## Features

- **OpenTelemetry Signal Receiver**: Receives OpenTelemetry signals and stores them in S3.
- **REST API**: Provides a REST API to retrieve stored signals.

## Installation

Pre-built release binaries are provided.
Download the latest binary from the [Releases page](https://github.com/mashiike/cflog2otel/releases).

`oteleport-server` is server, `oteleport` is client.

### Homebrew for client 

```
$ brew install mashiike/tap/oteleport
```

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


## Storage Flatten Options

if you followoing config, `oteleport` save OpenTelemetry signals convert to flat structure and json lines.
this option is useful for Amazon Ahena

```jsonnet
local must_env = std.native('must_env');

{
  access_keys: [
    must_env('OTELEPORT_ACCESS_KEY'),
  ],
  storage: {
    cursor_encryption_key: must_env('OTELEPORT_CURSOR_ENCRYPTION_KEY'),
    location: 's3://' + must_env('OTELEPORT_S3_BUCKET') + '/',
    flatten: true, // <- add this option
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

<details>
<summary> traces table schema </summary>

```jsonnet
CREATE EXTERNAL TABLE IF NOT EXISTS oteleport_traces (
    traceId STRING,
    spanId STRING,
    parentSpanId STRING,
    name STRING,
    kind INT,
    startTimeUnixNano BIGINT,
    endTimeUnixNano BIGINT,
    traceState STRING,
    resourceAttributes ARRAY<STRUCT<key: STRING, value: STRUCT<
        stringValue: STRING,
        boolValue: BOOLEAN,
        intValue: BIGINT,
        doubleValue: DOUBLE,
        arrayValue: STRUCT<values: ARRAY<STRUCT<stringValue: STRING>>>,
        kvlistValue: STRUCT<values: ARRAY<STRUCT<key: STRING, value: STRUCT<stringValue: STRING>>>>
    >>>,
    droppedResourceAttributesCount INT,
    resourceSpanSchemaUrl STRING,
    scopeName STRING,
    scopeVersion STRING,
    scopeAttributes ARRAY<STRUCT<key: STRING, value: STRUCT<
        stringValue: STRING,
        boolValue: BOOLEAN,
        intValue: BIGINT,
        doubleValue: DOUBLE,
        arrayValue: STRUCT<values: ARRAY<STRUCT<stringValue: STRING>>>,
        kvlistValue: STRUCT<values: ARRAY<STRUCT<key: STRING, value: STRUCT<stringValue: STRING>>>>
    >>>,
    droppedScopeAttributesCount INT,
    scopeSpanSchemaUrl STRING,
    attributes ARRAY<STRUCT<key: STRING, value: STRUCT<
        stringValue: STRING,
        boolValue: BOOLEAN,
        intValue: BIGINT,
        doubleValue: DOUBLE,
        arrayValue: STRUCT<values: ARRAY<STRUCT<stringValue: STRING>>>,
        kvlistValue: STRUCT<values: ARRAY<STRUCT<key: STRING, value: STRUCT<stringValue: STRING>>>>
    >>>,
    droppedAttributesCount INT,
    events ARRAY<STRUCT<name: STRING, timeUnixNano: BIGINT, attributes: ARRAY<STRUCT<key: STRING, value: STRUCT<
        stringValue: STRING,
        boolValue: BOOLEAN,
        intValue: BIGINT,
        doubleValue: DOUBLE,
        arrayValue: STRUCT<values: ARRAY<STRUCT<stringValue: STRING>>>,
        kvlistValue: STRUCT<values: ARRAY<STRUCT<key: STRING, value: STRUCT<stringValue: STRING>>>>
    >>>>>,
    droppedEventsCount INT,
    links ARRAY<STRUCT<traceId: STRING, spanId: STRING, attributes: ARRAY<STRUCT<key: STRING, value: STRUCT<
        stringValue: STRING,
        boolValue: BOOLEAN,
        intValue: BIGINT,
        doubleValue: DOUBLE,
        arrayValue: STRUCT<values: ARRAY<STRUCT<stringValue: STRING>>>,
        kvlistValue: STRUCT<values: ARRAY<STRUCT<key: STRING, value: STRUCT<stringValue: STRING>>>>
    >>>>>,
    droppedLinksCount INT,
    status STRUCT<code: INT, message: STRING>,
    flags INT
)
PARTITIONED BY (
    partition STRING
)
ROW FORMAT SERDE 'org.openx.data.jsonserde.JsonSerDe'
WITH SERDEPROPERTIES (
    'ignore.malformed.json' = 'true'
)
LOCATION 's3://<your s3 bucket name>/traces/'
TBLPROPERTIES (
    'projection.enabled' = 'true',
    'projection.partition.type' = 'date',
    'projection.partition.format' = 'yyyy/MM/dd/HH',
    'projection.partition.range' = '2023/01/01/00,NOW',
    'projection.partition.interval' = '1',
    'projection.partition.interval.unit' = 'HOURS',
    'storage.location.template' = 's3://<your s3 bucket name>/traces/${partition}/'
);
```

</details>

<details>
<summary> metrics table schema </summary>

```jsonnet
CREATE EXTERNAL TABLE IF NOT EXISTS oteleport_metrics (
    description STRING,
    name STRING,
    unit STRING,
    startTimeUnixNano BIGINT,
    timeUnixNano BIGINT,
    resourceAttributes ARRAY<STRUCT<key: STRING, value: STRUCT<
        stringValue: STRING,
        boolValue: BOOLEAN,
        intValue: BIGINT,
        doubleValue: DOUBLE,
        arrayValue: STRUCT<values: ARRAY<STRUCT<stringValue: STRING>>>,
        kvlistValue: STRUCT<values: ARRAY<STRUCT<key: STRING, value: STRUCT<stringValue: STRING>>>>
    >>>,
    droppedResourceAttributesCount INT,
    resourceMetricSchemaUrl STRING,
    scopeName STRING,
    scopeVersion STRING,
    scopeAttributes ARRAY<STRUCT<key: STRING, value: STRUCT<
        stringValue: STRING,
        boolValue: BOOLEAN,
        intValue: BIGINT,
        doubleValue: DOUBLE,
        arrayValue: STRUCT<values: ARRAY<STRUCT<stringValue: STRING>>>,
        kvlistValue: STRUCT<values: ARRAY<STRUCT<key: STRING, value: STRUCT<stringValue: STRING>>>>
    >>>,
    droppedScopeAttributesCount INT,
    scopeMetricSchemaUrl STRING,
    histogram STRUCT<
        dataPoint: STRUCT<
            attributes: ARRAY<STRUCT<key: STRING, value: STRUCT<
                stringValue: STRING,
                boolValue: BOOLEAN,
                intValue: BIGINT,
                doubleValue: DOUBLE,
                arrayValue: STRUCT<values: ARRAY<STRUCT<stringValue: STRING>>>,
                kvlistValue: STRUCT<values: ARRAY<STRUCT<key: STRING, value: STRUCT<stringValue: STRING>>>>
            >>>,
            startTimeUnixNano: BIGINT,
            timeUnixNano: BIGINT,
            count: BIGINT,
            sum: DOUBLE,
            bucketCounts: ARRAY<BIGINT>,
            explicitBounds: ARRAY<DOUBLE>,
            exemplars: ARRAY<STRUCT<value: DOUBLE, timestampUnixNano: BIGINT, filteredAttributes: ARRAY<STRUCT<key: STRING, value: STRUCT<stringValue: STRING>>>>>,
            flags: INT,
            min: DOUBLE,
            max: DOUBLE
        >,
        aggregationTemporality: INT
    >,
    exponentialHistogram STRUCT<
        dataPoint: STRUCT<
            attributes: ARRAY<STRUCT<key: STRING, value: STRUCT<
                stringValue: STRING,
                boolValue: BOOLEAN,
                intValue: BIGINT,
                doubleValue: DOUBLE,
                arrayValue: STRUCT<values: ARRAY<STRUCT<stringValue: STRING>>>,
                kvlistValue: STRUCT<values: ARRAY<STRUCT<key: STRING, value: STRUCT<stringValue: STRING>>>>
            >>>,
            startTimeUnixNano: BIGINT,
            timeUnixNano: BIGINT,
            count: BIGINT,
            sum: DOUBLE,
            scale: INT,
            zeroCount: BIGINT,
            positive: STRUCT<bucketCounts: ARRAY<BIGINT>, offset: INT>,
            negative: STRUCT<bucketCounts: ARRAY<BIGINT>, offset: INT>,
            flags: INT,
            exemplars: ARRAY<STRUCT<value: DOUBLE, timestampUnixNano: BIGINT, filteredAttributes: ARRAY<STRUCT<key: STRING, value: STRUCT<stringValue: STRING>>>>>,
            min: DOUBLE,
            max: DOUBLE,
            zeroThreshold: DOUBLE
        >,
        aggregationTemporality: INT
    >,
    summary STRUCT<
        dataPoint: STRUCT<
            attributes: ARRAY<STRUCT<key: STRING, value: STRUCT<
                stringValue: STRING,
                boolValue: BOOLEAN,
                intValue: BIGINT,
                doubleValue: DOUBLE,
                arrayValue: STRUCT<values: ARRAY<STRUCT<stringValue: STRING>>>,
                kvlistValue: STRUCT<values: ARRAY<STRUCT<key: STRING, value: STRUCT<stringValue: STRING>>>>
            >>>,
            startTimeUnixNano: BIGINT,
            timeUnixNano: BIGINT,
            count: BIGINT,
            sum: DOUBLE,
            quantileValues: ARRAY<STRUCT<quantile: DOUBLE, value: DOUBLE>>
        >
    >,
    gauge STRUCT<dataPoint: STRUCT<
        attributes: ARRAY<STRUCT<key: STRING, value: STRUCT<
            stringValue: STRING,
            boolValue: BOOLEAN,
            intValue: BIGINT,
            doubleValue: DOUBLE,
            arrayValue: STRUCT<values: ARRAY<STRUCT<stringValue: STRING>>>,
            kvlistValue: STRUCT<values: ARRAY<STRUCT<key: STRING, value: STRUCT<stringValue: STRING>>>>
        >>>,
        startTimeUnixNano: BIGINT,
        timeUnixNano: BIGINT,
        asDouble: DOUBLE, 
        asInt: BIGINT,
        exemplars: ARRAY<STRUCT<value: DOUBLE, timestampUnixNano: BIGINT, filteredAttributes: ARRAY<STRUCT<key: STRING, value: STRUCT<stringValue: STRING>>>>>,
        flags: INT
    >>,
    sum STRUCT<dataPoint: STRUCT<
        attributes: ARRAY<STRUCT<key: STRING, value: STRUCT<
            stringValue: STRING,
            boolValue: BOOLEAN,
            intValue: BIGINT,
            doubleValue: DOUBLE,
            arrayValue: STRUCT<values: ARRAY<STRUCT<stringValue: STRING>>>,
            kvlistValue: STRUCT<values: ARRAY<STRUCT<key: STRING, value: STRUCT<stringValue: STRING>>>>
        >>>,
        startTimeUnixNano: BIGINT,
        timeUnixNano: BIGINT,
        asDouble: DOUBLE, 
        asInt: BIGINT,
        exemplars: ARRAY<STRUCT<value: DOUBLE, timestampUnixNano: BIGINT, filteredAttributes: ARRAY<STRUCT<key: STRING, value: STRUCT<stringValue: STRING>>>>>,
        flags: INT
    >, aggregationTemporality: INT, isMonotonic: BOOLEAN>,
    metadata ARRAY<STRUCT<key: STRING, value: STRUCT<stringValue: STRING>>>
)
PARTITIONED BY (
    partition STRING
)
ROW FORMAT SERDE 'org.openx.data.jsonserde.JsonSerDe'
WITH SERDEPROPERTIES (
    'ignore.malformed.json' = 'true'
)
LOCATION 's3://<your s3 bucket name>/metrics/'
TBLPROPERTIES (
    'projection.enabled' = 'true',
    'projection.partition.type' = 'date',
    'projection.partition.format' = 'yyyy/MM/dd/HH',
    'projection.partition.range' = '2023/01/01/00,NOW',
    'projection.partition.interval' = '1',
    'projection.partition.interval.unit' = 'HOURS',
    'storage.location.template' = 's3://<your s3 bucket name>/metrics/${partition}/'
);
```

</details>

<details>
<summary> logs table schema </summary>

```jsonnet
REATE EXTERNAL TABLE IF NOT EXISTS oteleport_logs (
    resourceAttributes ARRAY<STRUCT<key: STRING, value: STRUCT<
        stringValue: STRING,
        boolValue: BOOLEAN,
        intValue: BIGINT,
        doubleValue: DOUBLE,
        arrayValue: STRUCT<values: ARRAY<STRUCT<stringValue: STRING>>>,
        kvlistValue: STRUCT<values: ARRAY<STRUCT<key: STRING, value: STRUCT<stringValue: STRING>>>>
    >>>,
    droppedResourceAttributesCount INT,
    resourceLogSchemaUrl STRING,
    scopeName STRING,
    scopeVersion STRING,
    scopeAttributes ARRAY<STRUCT<key: STRING, value: STRUCT<
        stringValue: STRING,
        boolValue: BOOLEAN,
        intValue: BIGINT,
        doubleValue: DOUBLE,
        arrayValue: STRUCT<values: ARRAY<STRUCT<stringValue: STRING>>>,
        kvlistValue: STRUCT<values: ARRAY<STRUCT<key: STRING, value: STRUCT<stringValue: STRING>>>>
    >>>,
    droppedScopeAttributesCount INT,
    scopeLogSchemaUrl STRING,
    timeUnixNano BIGINT,
    severityNumber INT,
    severityText STRING,
    body STRUCT<stringValue: STRING>,
    attributes ARRAY<STRUCT<key: STRING, value: STRUCT<
        stringValue: STRING,
        boolValue: BOOLEAN,
        intValue: BIGINT,
        doubleValue: DOUBLE,
        arrayValue: STRUCT<values: ARRAY<STRUCT<stringValue: STRING>>>,
        kvlistValue: STRUCT<values: ARRAY<STRUCT<key: STRING, value: STRUCT<stringValue: STRING>>>>
    >>>,
    droppedAttributesCount INT,
    flags INT,
    traceId STRING,
    spanId STRING,
    observedTimeUnixNano BIGINT
)
PARTITIONED BY (
    partition STRING
)
ROW FORMAT SERDE 'org.openx.data.jsonserde.JsonSerDe'
WITH SERDEPROPERTIES (
    'ignore.malformed.json' = 'true'
)
LOCATION 's3://<your s3 bucket name>/logs/'
TBLPROPERTIES (
    'projection.enabled' = 'true',
    'projection.partition.type' = 'date',
    'projection.partition.format' = 'yyyy/MM/dd/HH',
    'projection.partition.range' = '2023/01/01/00,NOW',
    'projection.partition.interval' = '1',
    'projection.partition.interval.unit' = 'HOURS',
    'storage.location.template' = 's3://<your s3 bucket name>/logs/${partition}/'
);
```

</details>

example query:
```sql
select 
    cast(from_unixtime(traces.startTimeUnixNano/1000000000) as date) as ymd,
    rs.value.stringValue as service_name,
    count(distinct traceId) as trace_count,
    count(distinct spanId) as span_count
from oteleport_traces as traces
cross join unnest(traces.resourceAttributes) with ordinality as t(rs, rs_index)
where traces.partition like '2024/11/%' and rs.key = 'service.name'
group by 1,2 
```

## License

This project is licensed under the MIT License. 
See the [LICENSE](./LICENSE) file for more details.
