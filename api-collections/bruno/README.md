## Bruno Collection

Thư mục này chứa collection cho Bruno (file text `.bru`, dễ review và commit).

### Cách dùng (GUI)

1. Mở Bruno
2. Open Collection → chọn folder `api-collections/bruno`
3. Chọn environment:
   - `local`
   - `staging`
4. Chạy requests trong `orders/`

### Cách dùng (CLI)

Cài CLI:

```bash
npm i -g @usebruno/cli
```

Run toàn bộ collection:

```bash
bru run --env local
```

Run riêng folder orders:

```bash
bru run orders --env local
```

### Reuse OpenAPI

Bạn có thể import trực tiếp OpenAPI vào Bruno để generate requests:

- `backend/orders/api/http/openapi.yaml`

Nếu bạn import OpenAPI, bạn vẫn có thể giữ collection curated này để:

- set biến `order_uuid` tự động sau khi tạo order
- viết tests/asserts ở mức request
