# TTK-Organazier (backend)

## Launch methods

env example to connect with supabase
```.env
DATABASE_URL="postgresql://postgres.nacuemduatmefatbmnoq:[YOUR-PASSWORD]@aws-0-eu-central-1.pooler.supabase.com:6543/postgres?pgbouncer=true"

DIRECT_URL="postgresql://postgres.nacuemduatmefatbmnoq:[YOUR-PASSWORD]@aws-0-eu-central-1.pooler.supabase.com:5432/postgres"
```

config .yaml
```yaml
env: "local"
token_ttl: 1h
app_secret: "YOUR_JWT_SECRET"
```

Start project
```bash
go run cmd/main.go --config=./config/local.yaml
```

# Run with docker
Pull the image from Docker Hub and run the docker run command.
```bash
docker pull c0dys/ttk-back:latest
sudo docker run -e CONFIG_PATH=/app/config/local.yaml -p 8080:8080  ttk-back
```