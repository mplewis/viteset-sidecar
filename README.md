# Viteset Sidecar

The Viteset Sidecar is an easy way to wire your Viteset blobs into your cloud-native app without adding complexity to your app.

# Configuration

The Sidecar is configured by setting the environment variables below:

| key        | type   | default                   | mandatory?         | description                                                                      |
| ---------- | ------ | ------------------------- | ------------------ | -------------------------------------------------------------------------------- |
| `SECRET`   | string | _none_                    | :white_check_mark: | The client secret used to fetch blobs                                            |
| `BLOB`     | string | _none_                    |                    | If set, only fetches the blob with this specific key, regardless of request path |
| `FRESH`    | int    | 15                        |                    | How long to wait, in seconds, before checking if a blob has a new value          |
| `HOST`     | string | 0.0.0.0                   |                    | The host to listen on                                                            |
| `PORT`     | int    | 80                        |                    | The port to listen on                                                            |
| `ENDPOINT` | string | `https://api.viteset.com` |                    | The Viteset API endpoint to use                                                  |
