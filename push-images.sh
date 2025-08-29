#!/bin/bash

# ===== Konfigurasi =====
DOCKER_USER="adisalaras"   # ganti dengan username Docker Hub kamu
TAG="1.0.0"                # versi/tag image

# ===== List service =====
SERVICES=("product-service" "transaction-service" "api-gateway")

# ===== Build, Tag, Push =====
for SERVICE in "${SERVICES[@]}"; do
  echo "🚀 Building $SERVICE ..."
  docker build -t $SERVICE:latest ./$SERVICE

  echo "🏷️  Tagging $SERVICE ..."
  docker tag $SERVICE:latest $DOCKER_USER/$SERVICE:$TAG

  echo "📤 Pushing $SERVICE ..."
  docker push $DOCKER_USER/$SERVICE:$TAG
done

echo "✅ Semua service berhasil di-push ke Docker Hub!"
