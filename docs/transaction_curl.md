# Transaction API cURL Examples

## Get All Transactions

```bash
# Basic request (uses default limit=10, page=1)
curl -X GET 'http://localhost:3000/api/v1/trx' \
-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDQ1MzAzNDAsImlzX2FkbWluIjp0cnVlLCJ1c2VyX2lkIjozNn0.-bEMCd7zed5LRsXWzbntOPWRq8-q6wkWbHSBuZFxZTY'

# With pagination parameters
curl -X GET 'http://localhost:3000/api/v1/trx?limit=10&page=1' \
-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDQ1MzAzNDAsImlzX2FkbWluIjp0cnVlLCJ1c2VyX2lkIjozNn0.-bEMCd7zed5LRsXWzbntOPWRq8-q6wkWbHSBuZFxZTY'
```

## Get Transaction by ID

```bash
curl -X GET 'http://localhost:3000/api/v1/transactions/1' \
-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDQ1MzAzNDAsImlzX2FkbWluIjp0cnVlLCJ1c2VyX2lkIjozNn0.-bEMCd7zed5LRsXWzbntOPWRq8-q6wkWbHSBuZFxZTY'
```

## Create New Transaction

```bash
curl -X POST 'http://localhost:3000/api/v1/trx' \
-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDQ1MzAzNDAsImlzX2FkbWluIjp0cnVlLCJ1c2VyX2lkIjozNn0.-bEMCd7zed5LRsXWzbntOPWRq8-q6wkWbHSBuZFxZTY' \
-H 'Content-Type: application/json' \
-d '{
    "method_bayar": "BANK_TRANSFER",
    "alamat_kirim": 3,
    "detail_trx": [
        {
            "id_produk": 74,
            "quantity": 2
        }
    ]
}'
```

## Update Transaction Status

```bash
curl -X PUT 'http://localhost:3000/api/v1/transactions/1/status' \
-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDQ1MzAzNDAsImlzX2FkbWluIjp0cnVlLCJ1c2VyX2lkIjozNn0.-bEMCd7zed5LRsXWzbntOPWRq8-q6wkWbHSBuZFxZTY' \
-H 'Content-Type: application/json' \
-d '{
    "status": "PAID"
}'
```

## Cancel Transaction

```bash
curl -X PUT 'http://localhost:3000/api/v1/transactions/1/cancel' \
-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDQ1MzAzNDAsImlzX2FkbWluIjp0cnVlLCJ1c2VyX2lkIjozNn0.-bEMCd7zed5LRsXWzbntOPWRq8-q6wkWbHSBuZFxZTY'
```

Note:

- GET request now works without requiring limit/page parameters
- For POST request, make sure:
  1. The address ID (alamat_kirim) exists in the database
  2. The product ID exists in the database
  3. The JSON field names match exactly: "id_produk" and "quantity" (not "product_id" and "kuantitas")
- Address ID 3 belongs to "Reyhan Update" with phone "082269283309"
- Product ID 74 should exist in your database
- The user token belongs to user_id: 36
  Status options: "WAITING_PAYMENT", "PAID", "PROCESSED", "DELIVERED", "CANCELLED"
- The field names must match exactly:
  - Use "id_produk" (not "product_id")
  - Use "quantity" (not "kuantitas")
- Address ID 3 belongs to "Reyhan Update"
- Product ID 74 is "Kemeja Batik" with stock 10
- The JSON must be valid without any comments
- Make sure to use correct field names:
  - "id_produk" for product ID
  - "quantity" for quantity
- This example uses:
  - Address ID 3 (belongs to "Reyhan Update")
  - Product ID 74 (Kemeja Batik with stock 10)
- The field `alamat_kirim` in JSON maps to `alamat_pengiriman` in database
- Make sure address ID exists in the `alamat` table
- Address ID 3 belongs to "Reyhan Update"
- Request body must match TransactionRequest model structure exactly:
  - "method_bayar" - payment method
  - "alamat_kirim" - maps to AlamatPengiriman internally
  - "detail_trx" - array of product details with "id_produk" and "quantity"
- Don't include fields like:
  - kode_invoice (generated automatically)
  - id_user (taken from JWT token)
  - harga_total (calculated from products)
