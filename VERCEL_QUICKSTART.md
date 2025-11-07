# Quick Vercel Deployment Guide

## ✅ What's Ready

- ✅ PostgreSQL database support (auto-detects `DATABASE_URL` or `POSTGRES_URL`)
- ✅ SQLite fallback for local development
- ✅ Vercel serverless function handler (`api/index.go`)
- ✅ Vercel configuration (`vercel.json`)
- ✅ Frontend static files ready

## 🚀 Deploy in 5 Minutes

### Step 1: Push to GitHub
```bash
git init
git add .
git commit -m "Ready for Vercel"
git remote add origin <your-github-repo-url>
git push -u origin main
```

### Step 2: Create Vercel Postgres
1. Go to [vercel.com](https://vercel.com)
2. Create new project → Import from GitHub
3. Select your repository
4. Go to **Storage** tab → **Create Database** → **Postgres**
5. Vercel automatically sets `POSTGRES_URL` environment variable

### Step 3: Set Environment Variables
In Vercel dashboard → Settings → Environment Variables:

```
JWT_SECRET=<generate-with-openssl-rand-hex-32>
ENV=production
CORS_ALLOWED_ORIGINS=https://your-app.vercel.app
BASE_URL=https://your-app.vercel.app
```

**Important:** Replace `your-app.vercel.app` with your actual Vercel domain!

Generate JWT secret:
```bash
openssl rand -hex 32
```

### Step 4: Deploy!
- Click **Deploy**
- Wait ~2-3 minutes
- Your app is live! 🎉

## 📝 Local Development

Still works with SQLite locally:
```bash
go run cmd/server/main.go
```

No `DATABASE_URL` needed - automatically uses SQLite!

## 🔍 Verify Deployment

1. Visit your Vercel URL: `https://your-app.vercel.app`
2. Should see the frontend
3. Try creating a short URL
4. Check Vercel Postgres dashboard to see tables created automatically

## 🐛 Troubleshooting

**Database not connecting?**
- Make sure Vercel Postgres is created
- Check `POSTGRES_URL` is set (auto-set by Vercel)
- Tables auto-create on first request

**Build failing?**
- Ensure `go.mod` is in root
- Check build logs in Vercel dashboard
- Verify all files are committed

**Static files not loading?**
- Ensure `web/static` directory is committed
- Check `vercel.json` routes configuration

## 📚 Files Created

- `api/index.go` - Vercel serverless handler
- `vercel.json` - Vercel configuration
- `DEPLOYMENT_VERCEL.md` - Detailed guide
- Updated `internal/database/database.go` - PostgreSQL support

## 🎯 That's It!

Your app is now ready for Vercel! Just push to GitHub and deploy.

