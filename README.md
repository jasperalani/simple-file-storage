## simple-file-storage
File storage and retrieval API written in Go.
Currently configured to allow audio files but can be changed to anything.

### Software Requirements
- [Go](https://go.dev/)
- [Docker](https://docker.com/)
- [Docker compose](https://docs.docker.com/compose/)

### Run
- `docker compose up`

### Restart
- `./docker-restart.sh`

### Endpoints
- `POST {localhost}/api/upload`
  - request:
    - form-data:
      - file: test-file.mp3
  - successful response (200):
    - id: uniquely generated id
    - message: success message
  - error response (400-500):
    - error: system error
    - message: error message
- `GET {localhost}/api/retrieve/{:id}`
  - successful response (200):
    - id: uniquely generated id
    - file: original filename and extension
    - link: link to download requested file
  - error response (400-500):
    - error: system error
    - message: error message

### Todo:
- ~~Upload~~
- ~~Retrieval~~
- Standardise all error responses
- Add message broker system
  - Create upload queue
- Tests
- Add logging to all errors
- Upload multiple files at once
- Check if file already exists in storage
  - By filename
- Rate limit
