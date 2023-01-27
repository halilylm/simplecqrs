
# Simple CQRS(Mimic)

However, write and query services are coupled, only queries are decoupled on the same service. Postgresql and Elasticsearch are being synced via Logstash.


## Usage

```bash
docker-compose up
```

Insert a product
```bash
curl --location --request POST 'http://localhost:8080/v1/products' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "Alarm paketi",
    "price": 5000000
}'
```

Delete a product

```bash
curl --location --request DELETE 'http://localhost:8080/v1/products/37198ff6-bde2-47a7-8f8f-8c16a5cd0c11'
```

Search products
```bash
curl --location --request GET 'http://localhost:8080/v1/products?name=alarm'
```
