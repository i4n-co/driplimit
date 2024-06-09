# How to use Driplimit ?

Start the driplimit API.

``` bash
$ export DATA_DIR="./"
$ driplimit

Jun  9 18:11:57.050 INF starting driplimit... component=api addr=127.0.0.1:7131
```

 In order to create your first keys, you'll need an `Admin Service Key`. Open another terminal. 

```bash
$ export DATA_DIR="./"
$ driplimit -init-admin

Admin service key created successfully.

Please store the key in a safe place. It will not be shown again.
Service Key ID:         sk_syndftqsmmjzwsciuidnqy
Service Key Token:      ************ redacted ****************
Description:            cli generated admin service key at 2024-06-09T15:44:54+02:00
Creation time:          2024-06-09T15:44:54+02:00
```
With this service key, you can then create a `Keyspace` to hold all your keys.

```bash
$ curl \
-H "Content-Type: application/json" \
-H "Authorization: Bearer *********" \
localhost:7131/v1/keyspaces.create -d \
'{
    "name": "api.yourcompany.io - dev",
    "keys_prefix": "dev_"
}'
```
```json
{
  "ksid": "ks_oorvwflallocxhruynieim",
  "name": "api.yourcompany.io - dev",
  "keys_prefix": "dev_"
}
```

Once your keyspace is created, you can start providing api keys for your users.

```bash
# headers are omitted for easy reading
$ curl [...] localhost:7131/v1/keys.create -d \
'{
    "ksid": "ks_oorvwflallocxhruynieim",
    "ratelimit": {
        "limit": 5, 
        "refill_rate": 1, 
        "refill_interval": 1000
    }
}'
```

```json
{
  "kid": "k_bgvzmoqgdpjbwhmxsrinnl",
  "ksid": "ks_oorvwflallocxhruynieim",
  "token": "dev_5D80xsGN+YBBv...%s9qUwIFZ6duHefdtvO3N9s",
  "created_at": "2024-06-09T16:36:24.957229514+02:00",
  "ratelimit": {
    "state": {
      "remaining": 5,
      "last_refilled": "2024-06-09T16:36:24.957233822+02:00"
    },
    "limit": 5,
    "refill_rate": 1,
    "refill_interval": 1000
  }
}
```
On creation, key's `token` is shown only once. **Don't store it in your systems, your user is the only proprietary of this secret.**

When your user starts consuming your api, you can check the provided token against Driplimit if the key is still valid and respects the configured rate limit.

```bash
# headers are omitted for easy reading
$ curl [...] localhost:7131/v1/keys.check -d '{
    "ksid": "ks_oorvwflallocxhruynieim",
    "token": "dev_5D80xsGN+YBBv...%s9qUwIFZ6duHefdtvO3N9s" 
}'
```
```json
{
  "kid": "k_bgvzmoqgdpjbwhmxsrinnl",
  "ksid": "ks_oorvwflallocxhruynieim",
  "created_at": "2024-06-09T16:36:24.957229514+02:00",
  "ratelimit": {
    "state": {
      "remaining": 4,
      "last_refilled": "2024-06-09T16:45:20.61031098+02:00"
    },
    "limit": 5,
    "refill_rate": 1,
    "refill_interval": 1000
  },
  "last_used": "2024-06-09T16:45:20.610501898+02:00"
}
```