# Render deployment guide

## Deploy to Render

### Steps:

1. **Push code to GitHub** (if not already done)

2. **Go to [render.com](https://render.com)** and sign up/login

3. **Create New Web Service:**
   - Click "New +" → "Web Service"
   - Connect your GitHub repository
   - Select your repository

4. **Configure Build Settings:**
   - **Name:** gourl (or your choice)
   - **Environment:** Go
   - **Build Command:** `go build -o gourl ./cmd/server/main.go`
   - **Start Command:** `./gourl`

5. **Set Environment Variables:**
   ```
   PORT=10000
   DB_PATH=/opt/render/project/src/gourl.db
   JWT_SECRET=<generate-with-openssl-rand-hex-32>
   ENV=production
   CORS_ALLOWED_ORIGINS=https://your-app.onrender.com
   RATE_LIMIT_RPS=10
   RATE_LIMIT_BURST=20
   ```

6. **Add Persistent Disk (for database):**
   - In Render dashboard → Your service → Settings
   - Add Disk: `/opt/render/project/src` (or your preferred path)
   - Update `DB_PATH` to match

7. **Deploy:**
   - Click "Create Web Service"
   - Render will build and deploy automatically
   - Get your URL: `https://your-app.onrender.com`

## Free Tier Notes

- Free tier available
- Spins down after 15 minutes of inactivity
- First request after spin-down takes ~30 seconds
- Perfect for testing/demos

