#!/bin/bash
# Build script for Vercel
cd /vercel/path0 || cd /vercel/path1 || cd .
go build -o bootstrap api/index.go

