
# Simple CQRS(Mimic)

However, command and query services are coupled, only read and write queries are decoupled on the same service. Postgresql and Elasticsearch are being synced via Logstash.


## Usage

```bash
docker-compose up
```

Migrate database(https://github.com/golang-migrate/migrate you will need to install this to be able to run `migrate`)
```bash 
export PGURL="postgres://product:secret@localhost:5432/product?sslmode=disable"
```
```bash 
migrate -database $PGURL -path business/data/migrations/ up
```

Insert a product
```bash
curl --location --request POST 'http://localhost:8000/v1/products' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "Alarm paketi",
    "price": 5000000
}'
```

Delete a product(don't forget to change id)

```bash
curl --location --request DELETE 'http://localhost:8000/v1/products/37198ff6-bde2-47a7-8f8f-8c16a5cd0c11'
```

Search products
```bash
curl --location --request GET 'http://localhost:8000/v1/products?name=alarm'
```

Jaeger
http://localhost:16686/search
