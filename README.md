# Description
An application that would fetch data from `https://test.deribit.com/api/v2/`.

### Run the project in Development Mode
```
make run
```


Additional commands:

```bash
run                            Start application
dep-build-up                   Start and build container
dep-up                         Create and start container
dep-down                       Stop and remove container, networks
dep-stop                       Stop services
```
## Available Endpoint

In the project directory, you can call:

### `POST /auth`

To get access_token and refresh token

### `GET /price`

For getting price of a particular currency

### `POST /buy`

To buy a currency

### `POST /sell`

To sell a currency

### `GET /swagger/index.html`

Helps to get the swagger details