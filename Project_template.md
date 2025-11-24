# Задание 1

**Контейнерная диаграмма в нотации С4.**  
[Диаграмма png](https://github.com/belyaev-a-g/architecture-cinemaabyss./blob/cinema/diagrams/container/Cinema_C4_Container.png)  
[Puml файл](https://github.com/belyaev-a-g/architecture-cinemaabyss./blob/cinema/diagrams/container/Cinema_C4_Container.puml)  

# Задание 2

### 1. Proxy
[API Gateway](https://github.com/belyaev-a-g/architecture-cinemaabyss./tree/cinema/src/microservices/proxy)  

### 2. Kafka
[Тесты 1](https://github.com/belyaev-a-g/architecture-cinemaabyss./blob/cinema/tests/screenshots/test_1.png)  
[Тесты 1](https://github.com/belyaev-a-g/architecture-cinemaabyss./blob/cinema/tests/screenshots/test_2.png)  
[Kafka topics 1](https://github.com/belyaev-a-g/architecture-cinemaabyss./blob/cinema/tests/screenshots/kafka_topics_main.png)  
[Kafka topics 2](https://github.com/belyaev-a-g/architecture-cinemaabyss./blob/cinema/tests/screenshots/kafka_topics.png)  
[Kafka topic payment-events](https://github.com/belyaev-a-g/architecture-cinemaabyss./blob/cinema/tests/screenshots/kafka_topic_payments-events.png)  
[Kafka topic user-events](https://github.com/belyaev-a-g/architecture-cinemaabyss./blob/cinema/tests/screenshots/kafka_topic_user-events.png)  

# Задание 3
### CI/CD
[Actions GitHUb](https://github.com/belyaev-a-g/architecture-cinemaabyss/actions)

### Proxy в Kubernetes
[Log events service](https://github.com/belyaev-a-g/architecture-cinemaabyss/blob/cinema/tests/screenshots/kuber_events_log.png)  
[Tests_1](https://github.com/belyaev-a-g/architecture-cinemaabyss/blob/cinema/tests/screenshots/kuber_tests_1.png)  
[Tests_2](https://github.com/belyaev-a-g/architecture-cinemaabyss/blob/cinema/tests/screenshots/kuber_tests_2.png)  
[Movies_1](https://github.com/belyaev-a-g/architecture-cinemaabyss/blob/cinema/tests/screenshots/kuber_movies_1.png)  
[Movies_2](https://github.com/belyaev-a-g/architecture-cinemaabyss/blob/cinema/tests/screenshots/kuber_movies_2.png)  

# Задание 4
Для простоты дальнейшего обновления и развертывания вам как архитектуру необходимо так же реализовать helm-чарты для прокси-сервиса и проверить работу 

Для этого:
1. Перейдите в директорию helm и отредактируйте файл values.yaml

```yaml
# Proxy service configuration
proxyService:
  enabled: true
  image:
    repository: ghcr.io/db-exp/cinemaabysstest/proxy-service
    tag: latest
    pullPolicy: Always
  replicas: 1
  resources:
    limits:
      cpu: 300m
      memory: 256Mi
    requests:
      cpu: 100m
      memory: 128Mi
  service:
    port: 80
    targetPort: 8000
    type: ClusterIP
```

- Вместо ghcr.io/db-exp/cinemaabysstest/proxy-service напишите свой путь до образа для всех сервисов
- для imagePullSecret проставьте свое значение (скопируйте из конфигурации kubernetes)
  ```yaml
  imagePullSecrets:
      dockerconfigjson: ewoJImF1dGhzIjogewoJCSJnaGNyLmlvIjogewoJCQkiYXV0aCI6ICJaR0l0Wlhod09tZG9jRjl2UTJocVZIa3dhMWhKVDIxWmFVZHJOV2hRUW10aFVXbFZSbTVaTjJRMFNYUjRZMWM9IgoJCX0KCX0sCgkiY3JlZHNTdG9yZSI6ICJkZXNrdG9wIiwKCSJjdXJyZW50Q29udGV4dCI6ICJkZXNrdG9wLWxpbnV4IiwKCSJwbHVnaW5zIjogewoJCSIteC1jbGktaGludHMiOiB7CgkJCSJlbmFibGVkIjogInRydWUiCgkJfQoJfSwKCSJmZWF0dXJlcyI6IHsKCQkiaG9va3MiOiAidHJ1ZSIKCX0KfQ==
  ```

2. В папке ./templates/services заполните шаблоны для proxy-service.yaml и events-service.yaml (опирайтесь на свою kubernetes конфигурацию - смысл helm'а сделать шаблоны для быстрого обновления и установки)

```yaml
template:
    metadata:
      labels:
        app: proxy-service
    spec:
      containers:
       Тут ваша конфигурация
```

3. Проверьте установку
Сначала удалим установку руками

```bash
kubectl delete all --all -n cinemaabyss
kubectl delete  namespace cinemaabyss
```
Запустите 
```bash
helm install cinemaabyss .\src\kubernetes\helm --namespace cinemaabyss --create-namespace
```
Если в процессе будет ошибка
```code
[2025-04-08 21:43:38,780] ERROR Fatal error during KafkaServer startup. Prepare to shutdown (kafka.server.KafkaServer)
kafka.common.InconsistentClusterIdException: The Cluster ID OkOjGPrdRimp8nkFohYkCw doesn't match stored clusterId Some(sbkcoiSiQV2h_mQpwy05zQ) in meta.properties. The broker is trying to join the wrong cluster. Configured zookeeper.connect may be wrong.
```

Проверьте развертывание:
```bash
kubectl get pods -n cinemaabyss
minikube tunnel
```

Потом вызовите 
https://cinemaabyss.example.com/api/movies
и приложите скриншот развертывания helm и вывода https://cinemaabyss.example.com/api/movies

## Удаляем все

```bash
kubectl delete all --all -n cinemaabyss
kubectl delete namespace cinemaabyss
```
