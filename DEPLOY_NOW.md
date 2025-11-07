# 🚀 Quick Vercel Deployment Guide

## Step 1: Push to GitHub

```bash
# Create a new repository on GitHub first, then:
git remote add origin https://github.com/YOUR_USERNAME/GOURL.git
git branch -M main
git push -u origin main
```

## Step 2: Deploy to Vercel

1. **Go to [vercel.com](https://vercel.com)** and sign in
2. **Click "New Project"**
3. **Import your GitHub repository** (GOURL)
4. **Vercel will auto-detect Go** - no build settings needed!

## Step 3: Create PostgreSQL Database

1. In Vercel dashboard → **Storage** tab
2. Click **Create Database** → **Postgres**
3. Vercel automatically sets `POSTGRES_URL` environment variable ✅

## Step 4: Set Environment Variables

In Vercel dashboard → **Settings** → **Environment Variables**, add:

```
JWT_SECRET=<generate-with-openssl-rand-hex-32>
ENV=production
BASE_URL=https://your-app.vercel.app
CORS_ALLOWED_ORIGINS=https://your-app.vercel.app
```

**Generate JWT secret:**
```bash
openssl rand -hex 32
```

**Important:** Replace `your-app.vercel.app` with your actual Vercel domain!

## Step 5: Deploy!

Click **Deploy** and wait ~2-3 minutes. Your app will be live! 🎉

## ✅ Verify Deployment

1. Visit your Vercel URL
2. Create a short URL
3. Check analytics - should work perfectly!

## 📝 Notes

- Database tables auto-create on first request
- Static files are served automatically
- All routes work (shortening, analytics, QR codes, etc.)
- Geographic analytics will show real countries (not "Local")

## 🐛 Troubleshooting

**Database not connecting?**
- Make sure Vercel Postgres is created
- Check `POSTGRES_URL` is set (auto-set by Vercel)

**Build failing?**
- Check build logs in Vercel dashboard
- Ensure `go.mod` is in root directory

**Static files not loading?**
- Ensure `web/static` directory is committed
- Check `vercel.json` routes configuration

---

**That's it! Your URL shortener is now live on Vercel! 🚀**

