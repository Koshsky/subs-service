#!/bin/bash

# Скрипт для остановки всех сервисов
set -e

echo "🛑 Stopping all services..."

# Останавливаем все сервисы
docker-compose down

# Опционально удаляем volumes (раскомментируйте если нужно)
# docker-compose down -v

echo "✅ All services stopped!"

# Показываем использование диска Docker
echo "💾 Docker disk usage:"
docker system df