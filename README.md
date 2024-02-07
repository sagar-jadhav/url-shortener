[![ci](https://github.com/sagar-jadhav/url-shortener/actions/workflows/docker-image.yml/badge.svg?branch=main)](https://github.com/sagar-jadhav/url-shortener/actions/workflows/docker-image.yml)

# URL Shortener written in Golang

## Prerequisites

- Golang
- Docker
- Docker Hub Account

## Steps to run the unit tests 

```bash
make test
```

## Steps to run the server locally

```bash
make run
```

## Steps to run the server in Docker container

- Run the following command to build the Docker image:

  ```bash
  make docker-build
  ```

- Run the following command to start the Docker container:

  ```bash
  make docker-run
  ```

## Steps to call the API's

### URL Shortening API

#### Request

```bash
curl --location 'localhost:3000' \
--header 'Content-Type: application/json' \
--data '{
    "longURL": "https://github.com/sagar-jadhav"
}'
```

#### Response

You should get the response in the below format. Data might vary as the short URL is generated randomly. But once you generated the Short URL you will get the same response till the application is running.

```json
{
    "longURL": "https://github.com/sagar-jadhav",
    "shortURL": "http://localhost:3000/qIFxf"
}
```

### Metrics API

#### Request

```bash
curl --location 'http://localhost:3000/metrics'
```

#### Response

```bash
[
    {
        "domain": "github.com",
        "count": 1
    }
]
```

**Note:**

I have already pushed the image to DockerHub So If you want to run the application without building it then run the following command: 

```bash
docker run -d -p 3000:3000  --name url-shortener developersthought/url-shortener:1.0
```
