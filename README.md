# Users API

This is a rest api built with [Go](https://golang.org/).

### Postman

You can interact with the rest api via [postman](https://www.getpostman.com/downloads/).

```
open postman
select import
select users.postman_collection.json from project root
select the collections tab to see a list of some examples on interacting with the api
```

### API

#### Users

| Verb      | Endpoint    | Description               | Query Parameters                |
|:--------- |:----------- |:------------------------- |:------------------------------- |
| GET       | /users      | get a collection of users | ?username=abeatrice             |
|           |             |                           | ?first_name=andrew              |
|           |             |                           | ?last_name=beatrice             |
|           |             |                           | ?email=abeatrice.mail@gmail.com |
| GET       | /users/{id} | get a user                |                                 |
| POST      | /users      | create a user             |                                 |
| PATCH     | /users/{id} | update a user             |                                 |
| DELETE    | /users/{id} | delete a user             |                                 |
