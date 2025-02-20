# Address API cURL Examples

## Get All Addresses

```bash
curl -X GET 'http://localhost:3000/api/v1/address' \
-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDQ1MzAzNDAsImlzX2FkbWluIjp0cnVlLCJ1c2VyX2lkIjozNn0.-bEMCd7zed5LRsXWzbntOPWRq8-q6wkWbHSBuZFxZTY'
```

## Get Address by ID

```bash
curl -X GET 'http://localhost:3000/api/address/1' \
-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDQ1MzAzNDAsImlzX2FkbWluIjp0cnVlLCJ1c2VyX2lkIjozNn0.-bEMCd7zed5LRsXWzbntOPWRq8-q6wkWbHSBuZFxZTY'
```

## Create New Address

```bash
curl -X POST 'http://localhost:3000/api/address' \
-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDQ1MzAzNDAsImlzX2FkbWluIjp0cnVlLCJ1c2VyX2lkIjozNn0.-bEMCd7zed5LRsXWzbntOPWRq8-q6wkWbHSBuZFxZTY' \
-H 'Content-Type: application/json' \
-d '{
    "judul_alamat": "Rumah",
    "nama_penerima": "John Doe",
    "no_telp": "08123456789",
    "detail_alamat": "Jl. Example No. 123",
    "id_provinsi": 1,
    "id_kota": 1
}'
```

## Update Address

```bash
curl -X PUT 'http://localhost:3000/api/address/1' \
-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDQ1MzAzNDAsImlzX2FkbWluIjp0cnVlLCJ1c2VyX2lkIjozNn0.-bEMCd7zed5LRsXWzbntOPWRq8-q6wkWbHSBuZFxZTY' \
-H 'Content-Type: application/json' \
-d '{
    "nama_penerima": "John Doe Updated",
    "no_telp": "08987654321",
    "detail_alamat": "Jl. Example Updated No. 456",
    "id_provinsi": 2,
    "id_kota": 2
}'
```

## Delete Address

```bash
curl -X DELETE 'http://localhost:3000/api/address/1' \
-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDQ1MzAzNDAsImlzX2FkbWluIjp0cnVlLCJ1c2VyX2lkIjozNn0.-bEMCd7zed5LRsXWzbntOPWRq8-q6wkWbHSBuZFxZTY'
```

Note: This token belongs to user_id: 36 with admin privileges and expires in 2025.
