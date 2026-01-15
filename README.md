# gosumsub

**gosumsub** is a lightweight and idiomatic Go SDK for integrating with the **Sumsub API**.

It is designed for simplicity, correctness, and ease of integration in production systems.

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.22-blue)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

## Features

- **Simple & Idiomatic API**
  Clean, readable interfaces that follow Go best practices.

- **Built-in Request Signing**
  Secure HMAC-SHA256 request signing handled internally.

- **Context Support**
  Full support for `context.Context` for cancellation and timeouts.

- **Zero External Dependencies**
  Uses only the Go standard library.

## Installation

```bash
go get github.com/andyle182810/gosumsub
```

## Quick Start

```go
client, err := gosumsub.NewClient(
    "https://api.sumsub.com",
    appToken,
    secretKey,
    gosumsub.WithDebug(true),
)
if err != nil {
    log.Fatal(err)
}

// Check API health
err = client.GetAPIHealthStatus(context.Background())
if err != nil {
    log.Fatal(err)
}

// Generate Web SDK link for KYC verification
resp, err := client.GenerateExternalWebSDKLink(
    context.Background(),
    &gosumsub.GenerateExternalWebSDKLinkRequest{
        LevelName: "basic-kyc-level",
        UserID:    "user-123",
        TTLInSecs: 3600,
        ApplicantIdentifiers: &gosumsub.ApplicantIdentifiers{
            Email: "user@example.com",
        },
    },
)
if err != nil {
    log.Fatal(err)
}

fmt.Println("Web SDK URL:", resp.URL)
```

## Testing

Integration tests automatically skip when required credentials are missing.

Required environment variables:

```bash
SUMSUB_APP_TOKEN=your_app_token
SUMSUB_SECRET_KEY=your_secret_key
```

Run tests:

```bash
go test ./...
```

## Contributing

Contributions are welcome and appreciated.

1. Fork the repository
2. Create a feature branch

   ```bash
   git checkout -b feature/my-feature
   ```

3. Commit your changes

   ```bash
   git commit -m "Add my feature"
   ```

4. Push to your fork

   ```bash
   git push origin feature/my-feature
   ```

5. Open a Pull Request

## Support

For bugs, questions, or feature requests:

- Open an issue on GitHub
  [https://github.com/andyle182810/gosumsub/issues](https://github.com/andyle182810/gosumsub/issues)

## License

**gosumsub** is licensed under the **MIT License**.
See the [LICENSE](LICENSE) file for details.
