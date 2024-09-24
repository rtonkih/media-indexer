# Media Indexer
## Description
Media Indexer is a web application for searching media files by associated tags. It provides an API for creating and searching media items.

## Table of Contents
- [Usage](#usage)
- [API Documentation](#api-documentation)
- [How to run tests](#how-to-run-tests)
- [Design and Technology Choices](#design-and-technology-choices)
- [Thoughts and improvements](#thoughts-and-improvements)

### Usage
1. Start the application:
    ```sh
    make run-dev
    ```
   It will start database and application inside docker containers

2. Start the application in production mode:
    ```sh
    make run
    ```
3. Migration is running automatically on application start, but can be run manually in app docker container:
    ```sh
    make run-migrate
    ```
4. To populate database with fake tags and media data run:
    ```sh
    make populate
    ```
5. By default app is available on 8080 port

### API Documentation
To start using API read [documentation](http://localhost:8080/swagger/index.html)

### How to run tests
```sh
    make test
```

### App details
**General**
- app design is quite simple, folder structure is divided into controllers, models, repositories, services, and storage. Could be switched to DDD, to separate UI, domain, and infrastructure logic, and operate with domain objects instead of models. But for this project, it would be an overkill
- GORM as ORM, which is a good choice for small projects. It provides features out of the box, like migrations, associations, and hooks. Could be switched to plain SQL queries for better performance, if needed
- gin as a web framework. Features out of the box, like routing, and request validation

**Tags**
- app does normalization for tags, so that "Ronaldo" and "ronaldo" are the same tags, also "World Cup" and "WorldCup" becomes "worldcup"
- app stores, returns, and searches tag in normalized form
- in case tag already exist in database, app return the existing tag, otherwise app creates a new tag and saves it to database

**Media**
- when user creates media, user required to provide tag in request. If tag already exist, app uses the existing tag, otherwise app creates a new tag and assign it to media
- app applies the same normalization logic to tags in media
- depending on needs, the app is able to save media to s3 or local storage and is ready to be extended to other storage types by implementing the `StorageProvider` interface
- when the app is starting, media storage can be chosen. This can be configured in `docker-compose.yml` by setting `STORAGE_TYPE` to "s3" or "local"


### Thoughts and improvements
- I suppose the system will have a lot of search media operations by tag name, so I added search by indexed **tag name** field from the **media_tag** table, without JOINing to the tag table. This will improve the performance of search by tag name
- I played with media query by tag on 1 million tags and 1 million media items. It takes under 50ms to search media by tag name, which looks good
- In case of performance issues, a caching system like Redis can be used. For example, store 10% of the most searched tags in Redis and search in Redis first; if not found, then search in the database. Also, need to come up with an invalidation strategy for the cache

### Done
- [x] Creating and searching tags
- [x] Creating and searching media by tags
- [x] Dockerized application
- [x] Swagger documentation
- [x] Migrations
- [x] Fake data population
- [x] Test cases for media and tag controllers

### TODO
- [ ] Add validation for photos in the CreateMedia endpoint
- [ ] Extract meta tags from media files and assign them to media entities
- [ ] If there are performance issues, add a caching system for searching media
- [ ] Add fuzzy search for tags so that users can find media by similar tags
- [ ] Add integration tests
- [ ] Add test cases for storage providers, repositories, and services
- [ ] Implement live reloading for development