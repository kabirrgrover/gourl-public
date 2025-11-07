# Railway deployment guide

## Quick Deploy to Railway

### Option 1: Railway CLI (Recommended)

1. **Install Railway CLI:**
   ```bash
   npm i -g @railway/cli
   ```

2. **Login:**
   ```bash
   railway login
   ```

3. **Initialize Railway project:**
   ```bash
   railway init
   ```

4. **Set environment variables in Railway dashboard:**
   - `JWT_SECRET` - Generate with: `openssl rand -hex 32`
   - `ENV=production`
   - `CORS_ALLOWED_ORIGINS` - Your Railway domain (e.g., `https://your-app.railway.app`)
   - `DB_PATH=/data/gourl.db` (Railway will create persistent volume)

5. **Add persistent volume for database:**
   - In Railway dashboard → Your service → Volumes
   - Add volume: `/data` (mount path)

6. **Deploy:**
   ```bash
   railway up
   ```

### Option 2: GitHub Integration (Easier)

1. **Push code to GitHub:**
   ```bash
   git init
   git add .
   git commit -m "Initial commit"
   git remote add origin <your-github-repo-url>
   git push -u origin main
   ```

2. **Connect Railway to GitHub:**
   - Go to [railway.app](https://railway.app)
   - Click "New Project"
   - Select "Deploy from GitHub repo"
   - Choose your repository

3. **Configure in Railway dashboard:**
   - Add environment variables (see above)
   - Add persistent volume: `/data`
   - Railway will auto-detect Go and deploy

4. **Get your URL:**
   - Railway provides a URL like: `https://your-app.railway.app`
   - Your frontend will be accessible at this URL!

## Environment Variables

Set these in Railway dashboard:

```
JWT_SECRET=<generate-random-secret>
ENV=production
DB_PATH=/data/gourl.db
CORS_ALLOWED_ORIGINS=https://your-app.railway.app
RATE_LIMIT_RPS=10
RATE_LIMIT_BURST=20
```

## Generate JWT Secret

```bash
openssl rand -hex 32
```

## Notes

- Railway automatically sets `PORT` environment variable
- Database persists in `/data` volume
- Free tier includes $5/month credit
- Auto-deploys on git push

