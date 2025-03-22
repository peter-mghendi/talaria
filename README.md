# Talaria

Fast, Lightweight Email Rendering Over gRPC, gRPC-Web, HTTP, and Connect

Talaria is a high-performance email rendering service built around the excellent [Hermes](https://github.com/go-hermes/hermes) library. It provides a structured way to generate HTML and plain-text emails using a flexible API that supports **gRPC, gRPC-Web, HTTP, and Connect**.

## Features

- **Multi-Protocol Support** – Works with **gRPC, gRPC-Web, HTTP, and Connect**.
- **Efficient Email Rendering** – Convert structured email data into **HTML and plain-text** outputs.
- **Streaming & Batch Processing** – Supports both **unary and streaming RPCs** for high-throughput rendering.
- **Lightweight & Scalable** – Optimized for speed and minimal resource usage.
- **Fully OpenAPI & gRPC Compatible** – Can be easily integrated into different backend systems.
- **Hermes-Based Email Templating** – Uses the [Hermes](https://github.com/go-hermes/hermes) templating engine to generate beautiful emails.

## Installation

### Running with Docker

Talaria is available on both **DockerHub** and **GitHub Container Registry (GHCR)**.

#### **From DockerHub**
```sh
docker run -p 9999:9999 petermghendi/talaria:latest
```

#### **From GHCR**
```sh
docker run -p 9999:9999 ghcr.io/peter-mghendi/talaria:latest
```

### Pre-Built Binaries

Pre-compiled binaries for major platforms are available on the [Releases Page](https://github.com/peter-mghendi/talaria/releases).

### Running Locally

Talaria is built in **Go** and can be run directly:

```sh
git clone https://github.com/peter-mghendi/talaria.git
cd talaria
go run cmd/server/main.go
```

## Usage

### API Overview

> See [render.proto](https://github.com/peter-mghendi/talaria/blob/main/render/v1/render.proto) for use with gRPC clients.

Talaria exposes the following APIs:

- `Render` – Converts structured email data into HTML and plain text (Unary RPC).
- `RenderStream` – Handles batch or streaming requests for rendering multiple emails (Bidirectional Streaming RPC).

### Example Request (Unary)

The unary RPC accepts a `request` consisting a `hermes.Hermes` object and a `hermes.Email` object.

> [!NOTE]
> Themes are passed in as lowercase strings, i.e. to apply `hermes.Default` pass in `theme: default` and to apply `hermes.Flat` pass in `theme: flat`.
> `hermes.Default` will be used if no theme is selected.

#### **gRPC**
Using `grpcurl`:

```sh
grpcurl -plaintext -d @ -proto render/v1/render.proto \
    localhost:9999 render.v1.RenderService/Render < render_request.json
```

Example JSON Payload (`render_request.json`):

```json
{
  "hermes": {
    "theme": "default",
    "text_direction": "ltr",
    "disable_css_inlining": false,
    "product": {
      "name": "Hermes",
      "link": "https://example-hermes.com/",
      "logo": "http://www.duchess-france.org/wp-content/uploads/2016/01/gopher.png"
    }
  },
  "email": {
    "body": {
      "name": "Jon Snow",
      "intros": ["Welcome to Hermes!"],
      "outros": ["Thanks for using Hermes."]
    }
  }
}
```

#### **HTTP (Connect/REST)**
Using `curl`:

```sh
curl -X POST http://localhost:9999/render.v1.RenderService/Render \
    -H "Content-Type: application/json" \
    -d @render_request.json
```

### Example Response

```json
{
  "html": "<html><body>...</body></html>",
  "text": "Welcome to Hermes!"
}
```

## Streaming Requests

For batch processing or live email previews, Talaria supports bidirectional streaming via `RenderStream`.

The `identifier` can be any string and is used to corrrelate streaming responses with the original requests.
The `request` is the same payload as above, consisting a `hermes.Hermes` object and a `hermes.Email` object.

#### **gRPC Streaming Request**
```sh
grpcurl -plaintext -d '{"identifier": "stream-1", "request": YOUR_JSON_PAYLOAD}' \
    -proto render/v1/render.proto \
    localhost:9999 render.v1.RenderService/RenderStream
```

## Load Testing

To benchmark Talaria's performance:

### **Unary Load Test**
```sh
ghz --insecure \
    --proto render/v1/render.proto \
    --call render.v1.RenderService/Render \
    -d @render_request.json -n 10000 -c 50 \
    localhost:9999
```

### **Streaming Load Test**
```sh
ghz --insecure \
    --proto render/v1/render.proto \
    --call render.v1.RenderService/RenderStream \
    -d @render_request.json -n 5000 -c 20 --stream-call-duration 10s \
    localhost:9999
```

## Deployment

### **Using Docker Compose**
For easy deployment, you can use `docker-compose`:

```yaml
version: "3"
services:
  talaria:
    image: petermghendi/talaria:latest # or ghcr.io/peter-mghendi/talaria:latest
    ports:
      - "9999:9999"
    environment:
      - LOG_LEVEL=info
```

Start the service:

```sh
docker-compose up -d
```

## Development

### **Building from Source**
```sh
go build -o talaria cmd/server/main.go
```

### **Running Tests**
```sh
go test ./...
```

### **Generating Protobuf Files**
```sh
buf lint
buf generate
```

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository.
2. Create a new branch (`git switch --create feature-name`).
3. Commit changes (`git commit --message "Add feature"`).
4. Push the branch (`git push origin feature-name`).
5. Open a Pull Request.

## License

This project is licensed under the **MIT License**.
