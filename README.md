# Viteset Sidecar

The Viteset Sidecar is an easy way to wire your Viteset blobs into your cloud-native app without adding any libraries to your app.

# Usage

## Docker

Try the sidecar by running it locally:

```sh
docker run \
  -e SECRET=YOUR_VITESET_CLIENT_SECRET_HERE \
  -e BLOB=YOUR_VITESET_BLOB_NAME_HERE \
  -p 8174:8174
  -it mplewis/viteset-sidecar
```

Then run `curl http://localhost:8174` to retrive the value of your blob from Viteset.

If your app is only using one blob for configuration, you usually want to configure both `SECRET` and `BLOB`. This lets your app make a GET request to the sidecar without having to know the name of the config blob.

If you want to grant your app access to all blobs for its client, you can omit `BLOB`:

```sh
docker run \
  -e SECRET=YOUR_VITESET_CLIENT_SECRET_HERE \
  -p 8174:8174
  -it mplewis/viteset-sidecar
```

Now, try `curl http://localhost:8174/YOUR_VITESET_BLOB_NAME_HERE` to get a blob by name.

## Docker Compose

See the Compose file in [examples/docker-compose.yaml](examples/docker-compose.yaml) which runs an `app` container to poll the value of a blob you define. Make sure to replace the `SECRET` and `BLOB` placeholders with values from your own account.

## Kubernetes

Coming soon.

# Caching

The sidecar assumes that app configuration changes infrequently, but that you want to see changes in production relatively soon after you update your blobs. To reduce load on Viteset servers, the sidecar caches your blobs locally for 15 seconds by default.

You can change this value by setting `FRESH` to your desired cache expiry time, e.g. `FRESH=120` for 2 minutes. Please don't lower `FRESH` below 15 seconds unless you **really** need a shorter caching period – this increases load on Viteset servers dramatically.

# Configuration

The Sidecar is configured by setting the environment variables below:

| key                  | type   | default                   | mandatory?         | description                                                                                      |
| -------------------- | ------ | ------------------------- | ------------------ | ------------------------------------------------------------------------------------------------ |
| `SECRET`             | string | _none_                    | :white_check_mark: | The client secret used to fetch blobs                                                            |
| `BLOB`               | string | _none_                    |                    | If set, only fetches the blob with this specific key, regardless of request path                 |
| `FRESH`              | int    | 15                        |                    | How long to wait, in seconds, before checking if a blob has a new value                          |
| `HOST`               | string | 0.0.0.0                   |                    | The host to listen on                                                                            |
| `PORT`               | int    | 8174                      |                    | The port to listen on                                                                            |
| `ENDPOINT`           | string | `https://api.viteset.com` |                    | The Viteset API endpoint to use                                                                  |
| `DO_NOT_FINGERPRINT` | bool   | _none_                    |                    | If set, omits OS information (e.g `osx El Capitan 11.16`) from the mandatory `User-Agent` header |
