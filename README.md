# Ozon Comments GraphQL

## Возможности

- Создание, чтение и отключение комментариев для постов.
- Поддержка пагинации для комментариев.
- Конфигурируемое хранилище данных (InMemory или PostgreSQL).
- Unit-тесты
- Docker-образ для легкого развертывания.

## Начало работы

### Установка

1. Клонируйте репозиторий:
    ```sh
    git clone https://github.com/yourusername/ozon-comments-graphql.git
    cd ozon-comments-graphql
    ```

2. Установите зависимости:
    ```sh
    go mod download
    ```

### Конфигурация

Приложение использует Viper для управления конфигурацией. Вы можете настроить приложение с помощью файла `config.yaml`. Пример конфигурации:

```yaml
app:
  port: 8080

db:
  host: localhost
  port: 5432
  user: postgres
  password: 7070
  name: ozon-comments
```

### Запуск приложения

#### С использованием Go:
```
go run .\cmd\server\server.go
```

```
go run .\cmd\server\server.go --storage postgres
```

#### С использованием Docker:
1. Постройте Docker-образ
```
docker build -t ozon-comments-graphql .
```
2. Запустите Docker-контейнер
```
docker run -p 8080:8080 ozon-comments-graphql
```
#### С использованием docker-compose:
```
docker-compose up
```

## Тестирование
Для запуска тестов нужно выполнить программу:
```
go test ./...
```

## Использование GraphQL Playground

После запуска приложения, вы можете получить доступ к GraphQL Playground по адресу `http://localhost:8080/` для выполнения запросов и тестирования API.

### Примеры запросов:
#### Получить все посты:
```graphql
query {
  posts {
    id
    title
    content
    comments {
      id
      content
    }
    commentsDisabled
  }
}
```

#### Получить пост по ID:
```graphql
query {
  post(id: "1") {
    id
    title
    content
    comments {
      id
      content
    }
    commentsDisabled
  }
}
```

#### Получить комментарии с пагинацией для определенного поста:
```graphql
query {
  comments(postId: "1", limit: 10, offset: 0) {
    comments {
      id
      postId
      parentId
      content
      children {
        id
        content
      }
      createdAt
    }
    totalCount
  }
}
```

### Примеры мутаций:
#### Создать пост:
```graphql
mutation {
  createPost(title: "Заголовок поста", content: "Содержание поста", commentsDisabled: false) {
    id
    title
    content
    commentsDisabled
  }
}
```

#### Создать комментарий:
```graphql
mutation {
  createComment(postId: "1", parentId: null, content: "Текст комментария") {
    id
    postId
    parentId
    content
    createdAt
  }
}
```

#### Отключить комментарии для поста:
```graphql
mutation {
  disableComments(postId: "1") {
    id
    title
    content
    commentsDisabled
  }
}
```

### Подписки:
#### Подписаться на добавление комментария к определенному посту:
```graphql
subscription {
  commentAdded(postId: "1") {
    id
    postId
    parentId
    content
    createdAt
  }
}
```
