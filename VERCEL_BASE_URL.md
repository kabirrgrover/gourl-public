# ✅ Vercel Environment Variables Setup

## Quick Answer: YES!

When hosting on Vercel, you **MUST** set `BASE_URL` in their environment variables.

## Step-by-Step Instructions

### 1. Go to Vercel Dashboard
- Visit [vercel.com](https://vercel.com)
- Select your project

### 2. Navigate to Environment Variables
- Click **Settings** (top menu)
- Click **Environment Variables** (left sidebar)

### 3. Add BASE_URL
Click **Add New** and enter:

**Key:** `BASE_URL`  
**Value:** `https://your-app.vercel.app`

**Important:** 
- Replace `your-app.vercel.app` with your **actual Vercel domain**
- Use `https://` (not `http://`)
- **No trailing slash** (don't add `/` at the end)
- Example: `https://gourl-abc123.vercel.app`

### 4. Set Environment (Optional but Recommended)
Select which environments this applies to:
- ✅ **Production** (required)
- ✅ **Preview** (optional, for preview deployments)
- ❌ **Development** (not needed)

### 5. Save and Redeploy
- Click **Save**
- Vercel will automatically redeploy
- Or manually trigger a new deployment

## Complete Environment Variables List

Set all of these in Vercel:

| Variable | Example Value | Required |
|----------|---------------|----------|
| `BASE_URL` | `https://your-app.vercel.app` | ✅ **YES** |
| `JWT_SECRET` | `a1b2c3d4e5f6...` (32+ chars) | ✅ Yes |
| `ENV` | `production` | ✅ Yes |
| `CORS_ALLOWED_ORIGINS` | `https://your-app.vercel.app` | ✅ Yes |
| `POSTGRES_URL` | Auto-set by Vercel | ✅ Auto |
| `DATABASE_URL` | Auto-set by Vercel | ✅ Auto |

## How to Find Your Vercel Domain

1. After deploying, go to your project dashboard
2. Look at the top - you'll see: `https://your-project-name.vercel.app`
3. Copy that exact URL (including `https://`)
4. Use it for both `BASE_URL` and `CORS_ALLOWED_ORIGINS`

## Example

If your Vercel app is at: `https://gourl-abc123.vercel.app`

Set:
```
BASE_URL=https://gourl-abc123.vercel.app
CORS_ALLOWED_ORIGINS=https://gourl-abc123.vercel.app
```

## Verify It Works

After setting `BASE_URL` and redeploying:

1. Create a short URL via API or frontend
2. Check the response - it should show:
   ```json
   {
     "short_url": "https://gourl-abc123.vercel.app/hello",
     ...
   }
   ```
3. ✅ Should **NOT** show `http://localhost:8080/hello`

## Troubleshooting

**Still seeing localhost?**
- Make sure `BASE_URL` is set correctly
- Check there's no trailing slash
- Verify you redeployed after setting the variable
- Check the variable is set for **Production** environment

**Wrong domain?**
- Double-check your Vercel domain
- Make sure you're using `https://` not `http://`
- Verify no typos in the domain name

## That's It!

Once `BASE_URL` is set in Vercel, all your short URLs will use your production domain! 🎉

