## Готовый `docker-compose.yml` для запуска стека observability (PostgreSQL, Jaeger, Prometheus, Grafana)

### Как запустить:
```bash
docker-compose up -d
```

### Доступ к сервисам:
- **Grafana**: http://localhost:3000 (admin/admin)
- **Prometheus**: http://localhost:9090
- **Jaeger UI**: http://localhost:16686
- **PostgreSQL**: `postgres://postgres:password@localhost:5432/db`
- **Ваше приложение**: http://localhost:8080

### Особенности:
1. Jaeger настроен в режиме All-in-One (удобно для разработки)
2. Prometheus автоматически собирает метрики с вашего приложения
3. Grafana предварительно настроена для работы с Prometheus
4. PostgreSQL с сохранением данных между перезапусками