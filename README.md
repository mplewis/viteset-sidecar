# Viteset Sidecar

[Viteset](https://viteset.com) lets you configure your app at lightspeed.

The Viteset Sidecar keeps your cloud-native app updated with your latest Viteset configs.

You don't need to add any libraries to your app. To get your latest config, your app makes a GET request to the Viteset Sidecar. The Sidecar handles polling for updates and caching the last known blob value.

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

## Docker Compose

See the Compose file in [examples/docker-compose.yaml](examples/docker-compose.yaml) which runs an `app` container to poll the value of a blob you define.

Make sure to replace the `SECRET` and `BLOB` placeholders with values from your own account.

## Kubernetes

See the Deployment in [examples/k8s-deployment.yaml](examples/k8s-deployment.yaml) which runs curl in a loop, to simulate your application requesting the latest blob value from the Sidecar.

Make sure to replace the `YOUR_VITESET_CLIENT_SECRET_GOES_HERE` and `YOUR_VITESET_BLOB_NAME_GOES_HERE` placeholders with values from your own account.

# Caching

The sidecar assumes that app configuration changes infrequently, but that you want to see changes in production relatively soon after you update your blobs. To reduce load on Viteset servers, the sidecar caches your blobs locally for 15 seconds by default.

You can change this value by setting `INTERVAL` to your desired cache expiry time, e.g. `INTERVAL=120` for 2 minutes. Please don't lower `INTERVAL` below 15 seconds unless you **really** need a shorter caching period – this increases load on Viteset servers dramatically.

# Configuration

The Sidecar is configured by setting the environment variables below:

| key        | type   | default                   | mandatory?         | description                                              |
| ---------- | ------ | ------------------------- | ------------------ | -------------------------------------------------------- |
| `SECRET`   | string | _none_                    | :white_check_mark: | The client secret with access to the given blob          |
| `BLOB`     | string | _none_                    | :white_check_mark: | The name of the blob to fetch                            |
| `INTERVAL` | int    | 15                        |                    | How often the sidecar polls for blob updates, in seconds |
| `HOST`     | string | 0.0.0.0                   |                    | The host to listen on                                    |
| `PORT`     | int    | 8174                      |                    | The port to listen on                                    |
| `ENDPOINT` | string | `https://api.viteset.com` |                    | The Viteset API endpoint to use                          |
