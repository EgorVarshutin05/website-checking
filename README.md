# website-checking

Приложение, позволяющее в настоящий момент увидеть, в каком состоянии находятся сайты.


Фронт ходит в **2 ручки**.

## 1. `GET /api/v1/sites` — список сайтов с текущим статусом

Вызывается при открытии страницы и потом каждые 30 секунд.

**200 OK**
```json
{
  "sites": [
    {
      "id": "shop",
      "name": "Интернет-магазин",
      "url": "https://shop.example.com",
      "status": "up",
      "response_time_ms": 412,
      "checked_at": "2026-10-27T15:59:30Z"
    },
    {
      "id": "legacy",
      "name": "Старый портал",
      "url": "https://legacy.example.org",
      "status": "down",
      "response_time_ms": null,
      "checked_at": "2026-10-27T15:59:30Z"
    }
  ]
}
```

| Поле | Тип | Описание |
|---|---|---|
| `id` | string | Уникальный ID сайта (можно число строкой или slug) |
| `name` | string | Отображаемое имя |
| `url` | string | Адрес, который проверяем |
| `status` | `"up"` \| `"down"` \| `"unknown"` | Результат **последней** проверки. `unknown` — ещё ни разу не проверяли |
| `response_time_ms` | int \| null | Время ответа последней проверки. `null`, если сайт лежит |
| `checked_at` | string | Когда была последняя проверка |

---

## 2. `GET /api/v1/sites/{id}/history?period=24h` — история за 24 часа

Вызывается при клике на сайт (для графика в модалке). Сейчас `period` всегда `24h`, параметр на будущее.

**200 OK**
```json
{
  "site_id": "shop",
  "from": "2026-10-26T16:00:00Z",
  "to": "2026-10-27T15:55:00Z",
  "step_seconds": 300,
  "points": [
    { "t": "2026-10-26T16:00:00Z", "status": "up",   "response_time_ms": 380 },
    { "t": "2026-10-26T16:05:00Z", "status": "down", "response_time_ms": null },
    { "t": "2026-10-26T16:10:00Z", "status": "up",   "response_time_ms": 402 }
  ]
}
```