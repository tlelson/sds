## Setup & notes for mongodb

- Some modules will need to implement the Entity interface.
- Mongo needs the `_id` tag
- IDs must be incrementing use ksuid.
- Probably needs indexing for queries

```bash
docker run -d -p 27017:27017 --name example-mongo mongo:latest
```
