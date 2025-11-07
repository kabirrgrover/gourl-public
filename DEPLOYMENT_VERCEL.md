# Vercel Deployment Guide

## Prerequisites

1. **Vercel Account**: Sign up at [vercel.com](https://vercel.com)
2. **Vercel CLI** (optional): `npm i -g vercel`
3. **PostgreSQL Database**: Use Vercel Postgres (recommended) or external PostgreSQL

## Step 1: Set Up Vercel Postgres

1. Go to your Vercel dashboard
2. Create a new project or select existing
3. Go to **Storage** → **Create Database** → **Postgres**
4. Copy the connection string (it will be set as `POSTGRES_URL` automatically)

## Step 2: Deploy to Vercel

### Option A: Via Vercel Dashboard (Easiest)

1. **Push code to GitHub:**
   ```bash
   git init
   git add .
   git commit -m "Initial commit"
   git remote add origin <your-github-repo>
   git push -u origin main
   ```

2. **Import to Vercel:**
   - Go to [vercel.com/new](https://vercel.com/new)
   - Import your GitHub repository
   - Vercel will auto-detect Go

3. **Configure Environment Variables:**
   In Vercel dashboard → Settings → Environment Variables, add:
   ```
   JWT_SECRET=<generate-with-openssl-rand-hex-32>
   ENV=production
   CORS_ALLOWED_ORIGINS=https://your-app.vercel.app
   BASE_URL=https://your-app.vercel.app
   RATE_LIMIT_RPS=10
   RATE_LIMIT_BURST=20
   ```
   
   **Important:** 
   - Replace `your-app.vercel.app` with your actual Vercel domain
   - `BASE_URL` ensures short URLs use your domain (not localhost)
   - `POSTGRES_URL` or `DATABASE_URL` is automatically set by Vercel Postgres

4. **Deploy:**
   - Click "Deploy"
   - Wait for build to complete
   - Your app will be live!

### Option B: Via Vercel CLI

```bash
# Install Vercel CLI
npm i -g vercel

# Login
vercel login

# Deploy
vercel

# Set environment variables
vercel env add JWT_SECRET
vercel env add ENV production
vercel env add CORS_ALLOWED_ORIGINS

# Deploy to production
vercel --prod
```

## Step 3: Generate JWT Secret

```bash
openssl rand -hex 32
```

## Environment Variables

Set these in Vercel dashboard:

| Variable | Value | Notes |
|----------|-------|-------|
| `JWT_SECRET` | Random hex string | Generate with `openssl rand -hex 32` |
| `ENV` | `production` | |
| `CORS_ALLOWED_ORIGINS` | Your Vercel URL | e.g., `https://your-app.vercel.app` |
| `BASE_URL` | Your Vercel URL | **Important:** e.g., `https://your-app.vercel.app` (no trailing slash) |
| `RATE_LIMIT_RPS` | `10` | Optional |
| `RATE_LIMIT_BURST` | `20` | Optional |
| `POSTGRES_URL` | Auto-set | Set automatically by Vercel Postgres |
| `DATABASE_URL` | Auto-set | Alternative name Vercel uses |

## Database Setup

Vercel Postgres automatically:
- Creates the database
- Sets `POSTGRES_URL` environment variable
- Handles connection pooling
- Provides web UI for data

The app will automatically create tables on first run!

## Troubleshooting

### Database Connection Issues
- Make sure Vercel Postgres is created and linked to your project
- Check that `POSTGRES_URL` or `DATABASE_URL` is set
- Verify tables are created (check Vercel Postgres dashboard)

### Build Failures
- Ensure `go.mod` is in root directory
- Check build logs in Vercel dashboard
- Verify all dependencies are in `go.mod`

### Static Files Not Loading
- Ensure `web/static` directory exists
- Check file paths in `vercel.json`
- Verify routes configuration

## Free Tier Limits

- **Hobby Plan**: Free
- **Serverless Functions**: 100GB-hours/month
- **Bandwidth**: 100GB/month
- **Vercel Postgres**: 256MB storage (free tier)

## Production Checklist

- [ ] Set strong `JWT_SECRET`
- [ ] Set `ENV=production`
- [ ] Configure `CORS_ALLOWED_ORIGINS`
- [ ] **Set `BASE_URL` to your Vercel domain** (prevents localhost in short URLs)
- [ ] Create Vercel Postgres database
- [ ] Test all endpoints
- [ ] Verify static files load
- [ ] Check database tables created
- [ ] Test short URL creation (should show your domain, not localhost)

## Your App URL

After deployment, your app will be available at:
```
https://your-project-name.vercel.app
```

Frontend and API will both be accessible at this URL!

