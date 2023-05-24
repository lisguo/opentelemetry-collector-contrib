# DataSet Exporter

<!-- status autogenerated section -->
| Status        |           |
| ------------- |-----------|
| Stability     | [alpha]: logs, traces   |
| Distributions | [] |

[alpha]: https://github.com/open-telemetry/opentelemetry-collector#alpha
<!-- end autogenerated section -->

This exporter sends logs to [DataSet](https://www.dataset.com/).

See the [Getting Started](https://app.scalyr.com/help/getting-started) guide.

## Configuration

### Required Settings

- `dataset_url` (no default): The URL of the DataSet API that ingests the data. Most likely https://app.scalyr.com.
- `api_key` (no default): The "Log Write" API Key required to use API. Instructions how to get [API key](https://app.scalyr.com/help/api-keys).

If you do not want to specify `api_key` in the file, you can use the [builtin functionality](https://opentelemetry.io/docs/collector/configuration/#configuration-environment-variables) and use `api_key: ${env:DATASET_API_KEY}`.

### Optional Settings

- `buffer`:
  - `max_lifetime` (default = 5s): The maximum delay between sending batches from the same source.
  - `group_by` (default = []): The list of attributes based on which events should be grouped.
  - `retry_initial_interval` (default = 5s): Time to wait after the first failure before retrying.
  - `retry_max_interval` (default = 30s): Is the upper bound on backoff.
  - `retry_max_elapsed_time` (default = 300s): Is the maximum amount of time spent trying to send a buffer.
- `traces`:
  - `aggregate` (default = false): Count the number of spans and errors belonging to a trace.
  - `max_wait` (default = 5s): The maximum waiting for all spans from single trace to arrive; ignored if `aggregate` is false.
- `retry_on_failure`: See [retry_on_failure](https://github.com/open-telemetry/opentelemetry-collector/blob/main/exporter/exporterhelper/README.md)
- `sending_queue`: See [sending_queue](https://github.com/open-telemetry/opentelemetry-collector/blob/main/exporter/exporterhelper/README.md)
- `timeout`: See [timeout](https://github.com/open-telemetry/opentelemetry-collector/blob/main/exporter/exporterhelper/README.md)


### Example

```yaml

exporters:
  dataset:
    # DataSet API URL
    dataset_url: https://app.scalyr.com
    # API Key
    api_key: your_api_key
    buffer:
      # Send buffer to the API at least every 10s
      max_lifetime: 10s
      # Group data based on these attributes
      group_by:
        - attributes.container_id

service:
  pipelines:
    logs:
      receivers: [otlp]
      processors: [batch]
      # add dataset among your exporters
      exporters: [dataset]
    traces:
      receivers: [otlp]
      processors: [batch]
      # add dataset among your exporters
      exporters: [dataset]
```