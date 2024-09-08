
# Интеграция Prometheus

## Описание изменений:

1. В `docker-compose.yaml` добавил сервис Prometheus.
2. В `Makefile` добавил команду для проверки его доступности: `make check_prometheus`.
3. При инициализации проекта метрики экспортируются через эндпоинт `/metrics`, доступный для Prometheus.
4. Добавил несколько метрик:
    - **Общий счётчик запросов**:
        - Для всех запросов используется метрика `requests_total`.
    - **PingHandler**:
        - Добавлен счётчик запросов `requests_ping` для каждого запроса к `/ping`.
        - Добавлена гистограмма `ping_latency` для измерения времени обработки запросов.
    - **GetHandler**:
        - Добавлен счётчик запросов `requests_get` для каждого запроса к `/get`.
        - Добавлена гистограмма `get_latency` для измерения времени обработки запросов.

## Пример работы **requests_total**

![prom_1.png](prom_1.png)

## Пример работы **requests_ping**

![prom_2.png](prom_2.png)

## Пример работы **ping_latency**

![prom_3.png](prom_3.png)

## Пример работы **requests_get**

![prom_4.png](prom_4.png)

## Пример работы **get_latency**

![prom_5.png](prom_5.png)