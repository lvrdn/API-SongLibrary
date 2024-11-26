# API - библиотека песен

**Структура проекта**
```
│   go.mod
│   go.sum
│   Makefile
│   README.md
│
├───cmd
│       main.go
│
├───config
│       app.env
│       cfg.go
│
├───docs
│       docs.go
│       swagger.json
│       swagger.yaml
│
├───migration
│       000001_init_schema.down.sql
│       000001_init_schema.up.sql
│
└───pkg
    ├───song
    │       handlers.go
    │       models.go
    │
    └───storage
            song_storage.go
```

Примеры CRUD-запросов:
1. **Получение всех песен:**
```
curl -X 'GET' \
  'http://127.0.0.1:8080/api/songs?limit=0&offset=0' \
  -H 'accept: application/json'


{"response":[{"id":1,"song":"demons","group":"imagine dragons","releaseDate":"28.01.2013","text":"line1\nline2\nline3\n\nline4\nline5\nline6 ","link":""},{"id":2,"song":"diamonds","group":"rihanna","releaseDate":"27.09.2012","text":"line1\nline2\n\nline3\nline4\n\nline5\nline6 ","link":""}]}
```

2. **Добавление новой песни:**
```
curl -X 'PUT' \
  'http://127.0.0.1:8080/api/songs' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{"song":"supermassive black hole","group":"muse"}'


{"response":{"id":3}}
```

3. **Получение всех песен с пагниацией и фильтрами:**

```
curl -X 'GET' \
  'http://127.0.0.1:8080/api/songs?limit=0&offset=0' \
  -H 'accept: application/json'


{"response":[{"id":1,"song":"demons","group":"imagine dragons","releaseDate":"28.01.2013","text":"line1\nline2\nline3\n\nline4\nline5\nline6 ","link":""},{"id":2,"song":"diamonds","group":"rihanna","releaseDate":"27.09.2012","text":"line1\nline2\n\nline3\nline4\n\nline5\nline6 ","link":""},{"id":3,"song":"supermassive black hole","group":"muse","releaseDate":"16.07.2006","text":"Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight","link":"https://www.youtube.com/watch?v=Xsp3_a-PMTw"}]}
```

```
curl -X 'GET' \
  'http://127.0.0.1:8080/api/songs?limit=0&offset=0&link=true' \
  -H 'accept: application/json'


{"response":[{"id":3,"song":"supermassive black hole","group":"muse","releaseDate":"16.07.2006","text":"Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight","link":"https://www.youtube.com/watch?v=Xsp3_a-PMTw"}]}
```

```
curl -X 'GET' \
  'http://127.0.0.1:8080/api/songs?limit=0&offset=0&song=mon&group=dragons&releaseDate=2013&text=line1' \
  -H 'accept: application/json'


{"response":[{"id":1,"song":"demons","group":"imagine dragons","releaseDate":"28.01.2013","text":"line1\nline2\nline3\n\nline4\nline5\nline6 ","link":""}]}
```

4. **Получение текста песни с пагинацией по куплетам:**
```
curl -X 'GET' \
  'http://127.0.0.1:8080/api/songs/2?limit=2&offset=0' \
  -H 'accept: application/json'


{"response":{"id":2,"verses":"line1\nline2\n\nline3\nline4","versesInSong":3}}
```

```
curl -X 'GET' \
  'http://127.0.0.1:8080/api/songs/1?limit=1&offset=1' \
  -H 'accept: application/json'


{"response":{"id":1,"verses":"line4\nline5\nline6 ","versesInSong":2}}
```

5. **Изменение данных песни (date, text, link):**
```
curl -X 'POST' \
  'http://127.0.0.1:8080/api/songs' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{"id":1,"releaseDate":"05.10.2016","text":"new text","link":"https://www.somevideohosting.com/123"}'
```

6. **Удаление песни:**
```
curl -X 'DELETE' \
  'http://127.0.0.1:8080/api/songs' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{"id":2}'
```

7. **Получение всех песен после п.5 и п.6 :**
```
curl -X 'GET' \
  'http://127.0.0.1:8080/api/songs?limit=0&offset=0' \
  -H 'accept: application/json'


{"response":[{"id":1,"song":"demons","group":"imagine dragons","releaseDate":"05.10.2016","text":"new text","link":"https://www.somevideohosting.com/123"},{"id":3,"song":"supermassive black hole","group":"muse","releaseDate":"16.07.2006","text":"Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight","link":"https://www.youtube.com/watch?v=Xsp3_a-PMTw"}]}
```
