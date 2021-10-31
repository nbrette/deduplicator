# Deduplicator

This service job is to deduplicate messages based on the attributes given. 
When a message is sent it generates a hash with the concatenated attributes which is used to determine if the message has aldready been sent previously. If the hash already exists in the database a counter is incremented to know how many times the same message has been received. 
Otherwise if the hash does not exist, it is inserted in a redis server using the hash as the key and the attributes as data. The message is then forwarded. Both HTTP request and Pub/Sub method can be used to forward the message.

## Configuration
The following env parameters are needed for this service:
* **REDIS_URL:** The URL of the redis server (in this format: localhost:6379)

## Build
Use this command to create a docker image of this program:
```bash
docker build -t deduplicator .
```

## Local run
This service requires a redis server to run. If you have docker installed you can run one locally with this command:
```bash
docker run -d -p 6379:6379 redis
```

Once your redis server is ready, set the needed variables in .env.local and launch:
```
source ./scripts/.env.local
go run main.go
```
OR
```
go build . -o exe
source ./scripts/.env.local
./exe
```

## Use
Once the deduplicator is running it is ready to receive messages.
Here is the the format that need to be used for theses messages.
```json
{
    "Message": {
        "Attributes": {
            "next_hop": "URL_OR_PUBSUB TOPIC",
            "next_hop_method": "METHODE_TO_FORWARD",
            "hash_1": "HASH_ATTRIBUTE_1",
            "hash_2": "HASH_ATTRIBUTE_2"
        },
        "Data": [
            12,12,12
        ]
    }
}
```

Examples can be found in the folder `internal/services/data`.

Messages must be sent via an HTTP request on the port 

