# Simple DB API

## About

To start run `make`. This will give you two binaries in the `/bin` folder. One called `app` is the backend, that  exposes:
- `/` just a hello world endpoint
- `/create` an endpoint to POST via JSON an order created event in the form of 
    ```
    {
        "baseSiteUid": "an-id",
        "orderCode": "4caad296-e0c5-491e-98ac-0ed118f9474e"
    }
    ```
    Successfully creating an event on the API will return a 201 JSON response:
    ```
    {
        "id": 0,
        "firstName": "",
        "lastName": "",
        "orderCode": "4caad296-e0c5-491e-98ac-0ed118f9474e",
        "baseSiteUid": "an-id"
    }
    ```
- `/read` an endpoint that returns a JSON array with all stored events

## Run

Run the everything with `docker-compose up`. This will expose the backend on `http://localhost:8080` and the frontend on `http://localhost:8081`. In this docker-compose setup the backend will be supported by a postgres database on port 5432.

If you wish to change the demo setup to run without a database you can do so by using the `InMemoryStore`. Changing the implementation is easily done in the [`cmd\simple-db-api\main.go`](./cmd/simple-db-api/main.go) file, by replacing the service initialisation with this line:
```
svc := service.NewSimpleDBAPI().WithStorage(store.NewInMemoryStore())
```