# Fly.io deployment guide

## Deploy to Fly.io

### Steps:

1. **Install Fly CLI:**
   ```bash
   curl -L https://fly.io/install.sh | sh
   ```

2. **Login:**
   ```bash
   fly auth login
   ```

3. **Initialize Fly app:**
   ```bash
   fly launch
   ```
   - Follow prompts
   - Don't deploy yet (we'll configure first)

4. **Create fly.toml** (already created, but verify):
   ```toml
   app = "your-app-name"
   primary_region = "iad"  # Choose closest region

   [build]
     builder = "paketobuildpacks/builder:base"

   [http_service]
     internal_port = 8080
     force_https = true
     auto_stop_machines = true
     auto_start_machines = true
     min_machines_running = 0

   [[vm]]
     memory_mb = 256
   ```

5. **Add persistent volume for database:**
   ```bash
   fly volumes create data --size 1 --region iad
   ```

6. **Set secrets:**
   ```bash
   fly secrets set JWT_SECRET=$(openssl rand -hex 32)
   fly secrets set ENV=production
   fly secrets set DB_PATH=/data/gourl.db
   fly secrets set CORS_ALLOWED_ORIGINS=https://your-app.fly.dev
   ```

7. **Deploy:**
   ```bash
   fly deploy
   ```

8. **Get your URL:**
   - `https://your-app.fly.dev`

## Free Tier

- 3 shared-cpu VMs
- 3GB persistent volumes
- 160GB outbound data transfer
- Great for production use!

