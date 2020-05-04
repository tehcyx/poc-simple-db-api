# Simple DB API

## About

To start run `make`. This will give you two binaries in the `/bin` folder. One called `app` is the backend, that  exposes:
- `/` just a hello world endpoint
- `/create` an endpoint to POST via JSON an order created event in the form of 
    ```
    {
        "baseSiteId": "an-id",
        "orderCode": "4caad296-e0c5-491e-98ac-0ed118f9474e"
    }
    ```
    Successfully creating an event on the API will return a 201 JSON response:
    ```
    {
    "date": "2020-04-30T14:14:23.223079-07:00",
    "data": "ewoJIm9yZGVyQ29kZSI6ICI0Y2FhZDI5Ni1lMGM1LTQ5MWUtOThhYy0wZWQxMThmOTQ3NGUiCn0="
    }
    ```
- `/read` an endpoint that returns a JSON array with all stored events

## Run

### Backend
Run the backend with `KYMA_URL="localhost" COMMERCE_URL="localhost" ./bin/app`. This will expose the backend on `http://localhost:8080`

### Frontend
Run the frontend with `BACKEND_URL="localhost:8080" ./bin/frontend`. This will expose the frontend on `http://localhost:8081` and point the frontend to talk to the backend. Data on the frontend is live served from the backend.