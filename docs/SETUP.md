# Setup

## Production
```bash
cp .env.example .env
# edit the .env file in your favorite text editor

docker compose up
```

## Development
1. Enure Golang is installed. [More on this can be found on the Go website.](https://go.dev/doc/install)
2. Copy the starter `.env` from `.env.example` to `.env` at the root. Edit the variables as needed.
    1. To generate a secure JWT signing token, you can use OpenSSL: `openssl rand -base64 64`
    2. For a free SMTP server, [consider using Gmail](https://support.google.com/a/answer/176600?hl=en) (this is capped, so be aware of your usage)
3. Enter the frontend folder, and copy the sample `.env.example` file to `.env`, and edit the configuration variables as needed.
4. For development, run `go run cmd\hn\main.go` from the root to start the web server on `localhost` port 3000.