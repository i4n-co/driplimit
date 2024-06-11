Driplimit works with a JSON RPC API. All calls are made via HTTP using the POST method. 
The API returns 200 OK or 204 CREATED on success.

See [Errors](/#errors) section in case of failure.


## Keys

### `POST /v1/keys.create`
Create a key

**Headers**

* `Content-Type: application/json` - tells the service you wish to communicate with json
* `Authorization: Bearer <token>`  - the service key token

**Parameters**


* `ksid` <span style="border: 1px #AAA solid; padding: 2px; border-radius: 5px;">string</span> - The id of the keyspace to which the key belongs to (required)


* `expires_in` <span style="border: 1px #AAA solid; padding: 2px; border-radius: 5px;">integer</span> - The duration in milliseconds after which the key expires


* `expires_at` <span style="border: 1px #AAA solid; padding: 2px; border-radius: 5px;">timestamp</span> - The time at which the key expires (expires_at takes precedence over expires_in)


* `ratelimit` <span style="border: 1px #AAA solid; padding: 2px; border-radius: 5px;">object</span> - The rate limit configuration for the key (required)

  * `limit` <span style="border: 1px #AAA solid; padding: 2px; border-radius: 5px;">integer</span> - The rate limit


  * `refill_rate` <span style="border: 1px #AAA solid; padding: 2px; border-radius: 5px;">integer</span> - The rate at which the rate limit refills


  * `refill_interval` <span style="border: 1px #AAA solid; padding: 2px; border-radius: 5px;">integer</span> - The interval at which the rate limit refills




<details>
<summary> <b>cURL example</b> </summary>

```bash
$ curl -X POST
       -H "Content-Type: application/json" \
       -H "Authorization: Bearer <token>" \
       --data '
        {
          "ksid": "ks_abc",
          "expires_in": 300000,
          "expires_at": "0001-01-01T00:00:00Z",
          "ratelimit": {
            "limit": 5,
            "refill_rate": 1,
            "refill_interval": 1000
          }
        }' https://demo.driplim.it/v1/keys.create
```

```json
{
  "kid": "k_xyz",
  "ksid": "ks_abc",
  "token": "demo_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
  "created_at": "2024-06-11T17:10:49.74621998+02:00",
  "ratelimit": {
    "state": {
      "remaining": 5,
      "last_refilled": "2024-06-11T17:10:49.746220248+02:00"
    },
    "limit": 5,
    "refill_rate": 1,
    "refill_interval": 1000
  },
  "expires_at": "2024-06-11T17:15:49.746220023+02:00"
}
```
</details>

### `POST /v1/keys.check`
Check a key

**Headers**

* `Content-Type: application/json` - tells the service you wish to communicate with json
* `Authorization: Bearer <token>`  - the service key token

**Parameters**


* `ksid` <span style="border: 1px #AAA solid; padding: 2px; border-radius: 5px;">string</span> - The id of the keyspace to which the key belongs to (required)


* `token` <span style="border: 1px #AAA solid; padding: 2px; border-radius: 5px;">string</span> - The token to check (required)



<details>
<summary> <b>cURL example</b> </summary>

```bash
$ curl -X POST
       -H "Content-Type: application/json" \
       -H "Authorization: Bearer <token>" \
       --data '
        {
          "ksid": "ks_abc",
          "token": "demo_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
        }' https://demo.driplim.it/v1/keys.check
```

```json
{
  "kid": "k_xyz",
  "ksid": "ks_abc",
  "created_at": "2024-06-11T17:10:49.746222385+02:00",
  "ratelimit": {
    "state": {
      "remaining": 4,
      "last_refilled": "2024-06-11T17:10:49.746222501+02:00"
    },
    "limit": 5,
    "refill_rate": 1,
    "refill_interval": 1000
  },
  "expires_at": "2024-06-11T17:15:49.746222427+02:00"
}
```
</details>

### `POST /v1/keys.list`
List keys

**Headers**

* `Content-Type: application/json` - tells the service you wish to communicate with json
* `Authorization: Bearer <token>`  - the service key token

**Parameters**


* `list` <span style="border: 1px #AAA solid; padding: 2px; border-radius: 5px;">object</span> - The list options

  * `page` <span style="border: 1px #AAA solid; padding: 2px; border-radius: 5px;">integer</span> - The page number


  * `limit` <span style="border: 1px #AAA solid; padding: 2px; border-radius: 5px;">integer</span> - The number of items per page



* `ksid` <span style="border: 1px #AAA solid; padding: 2px; border-radius: 5px;">string</span> - The id of the keyspace to which the keys belong to (required)



<details>
<summary> <b>cURL example</b> </summary>

```bash
$ curl -X POST
       -H "Content-Type: application/json" \
       -H "Authorization: Bearer <token>" \
       --data '
        {
          "list": {
            "page": 1,
            "limit": 10
          },
          "ksid": "ks_abc"
        }' https://demo.driplim.it/v1/keys.list
```

```json
{
  "list": {
    "page": 1,
    "limit": 10,
    "last_page": 1
  },
  "keys": [
    {
      "kid": "k_xyz",
      "ksid": "ks_abc",
      "created_at": "2024-06-11T17:10:49.746223976+02:00",
      "ratelimit": {
        "state": {
          "remaining": 4,
          "last_refilled": "2024-06-11T17:10:49.746224058+02:00"
        },
        "limit": 5,
        "refill_rate": 1,
        "refill_interval": 1000
      },
      "expires_at": "2024-06-11T17:15:49.746224013+02:00"
    }
  ]
}
```
</details>

### `POST /v1/keys.get`
Get a key

**Headers**

* `Content-Type: application/json` - tells the service you wish to communicate with json
* `Authorization: Bearer <token>`  - the service key token

**Parameters**


* `ksid` <span style="border: 1px #AAA solid; padding: 2px; border-radius: 5px;">string</span> - The id of the keyspace to which the key belongs to (required)


* `kid` <span style="border: 1px #AAA solid; padding: 2px; border-radius: 5px;">string</span> - The id of the key to get (kid takes precedence over token if both are provided)


* `token` <span style="border: 1px #AAA solid; padding: 2px; border-radius: 5px;">string</span> - The token of the key to get



<details>
<summary> <b>cURL example</b> </summary>

```bash
$ curl -X POST
       -H "Content-Type: application/json" \
       -H "Authorization: Bearer <token>" \
       --data '
        {
          "ksid": "ks_abc",
          "kid": "k_xyz",
          "token": ""
        }' https://demo.driplim.it/v1/keys.get
```

```json
{
  "kid": "k_xyz",
  "ksid": "ks_abc",
  "created_at": "2024-06-11T17:10:49.746225317+02:00",
  "ratelimit": {
    "state": {
      "remaining": 4,
      "last_refilled": "2024-06-11T17:10:49.746225401+02:00"
    },
    "limit": 5,
    "refill_rate": 1,
    "refill_interval": 1000
  },
  "expires_at": "2024-06-11T17:15:49.746225355+02:00"
}
```
</details>

### `POST /v1/keys.delete`
Delete a key from a keyspace

**Headers**

* `Content-Type: application/json` - tells the service you wish to communicate with json
* `Authorization: Bearer <token>`  - the service key token

**Parameters**


* `ksid` <span style="border: 1px #AAA solid; padding: 2px; border-radius: 5px;">string</span> - The id of the keyspace to which the key belongs to (required)


* `kid` <span style="border: 1px #AAA solid; padding: 2px; border-radius: 5px;">string</span> - The id of the key to delete (required)



<details>
<summary> <b>cURL example</b> </summary>

```bash
$ curl -X POST
       -H "Content-Type: application/json" \
       -H "Authorization: Bearer <token>" \
       --data '
        {
          "ksid": "ks_abc",
          "kid": "k_xyz"
        }' https://demo.driplim.it/v1/keys.delete
```

```json
null
```
</details>


## Keyspaces

### `POST /v1/keyspaces.get`
Get keyspace by ID

**Headers**

* `Content-Type: application/json` - tells the service you wish to communicate with json
* `Authorization: Bearer <token>`  - the service key token

**Parameters**


* `ksid` <span style="border: 1px #AAA solid; padding: 2px; border-radius: 5px;">string</span> - The id of the keyspace to get (required)



<details>
<summary> <b>cURL example</b> </summary>

```bash
$ curl -X POST
       -H "Content-Type: application/json" \
       -H "Authorization: Bearer <token>" \
       --data '
        {
          "ksid": "ks_abc"
        }' https://demo.driplim.it/v1/keyspaces.get
```

```json
{
  "ksid": "ks_abc",
  "name": "demo.yourapi.com (env: production)",
  "keys_prefix": "demo_",
  "ratelimit": {
    "limit": 100,
    "refill_rate": 1,
    "refill_interval": 1000
  }
}
```
</details>

### `POST /v1/keyspaces.list`
Get keyspace by ID

**Headers**

* `Content-Type: application/json` - tells the service you wish to communicate with json
* `Authorization: Bearer <token>`  - the service key token

**Parameters**


* `list` <span style="border: 1px #AAA solid; padding: 2px; border-radius: 5px;">object</span> - The list options

  * `page` <span style="border: 1px #AAA solid; padding: 2px; border-radius: 5px;">integer</span> - The page number


  * `limit` <span style="border: 1px #AAA solid; padding: 2px; border-radius: 5px;">integer</span> - The number of items per page




<details>
<summary> <b>cURL example</b> </summary>

```bash
$ curl -X POST
       -H "Content-Type: application/json" \
       -H "Authorization: Bearer <token>" \
       --data '
        {
          "list": {
            "page": 1,
            "limit": 10
          }
        }' https://demo.driplim.it/v1/keyspaces.list
```

```json
{
  "list": {
    "page": 1,
    "limit": 10,
    "last_page": 1
  },
  "keyspaces": [
    {
      "ksid": "ks_abc",
      "name": "demo.yourapi.com (env: production)",
      "keys_prefix": "demo_",
      "ratelimit": {
        "limit": 100,
        "refill_rate": 1,
        "refill_interval": 1000
      }
    }
  ]
}
```
</details>

### `POST /v1/keyspaces.create`
Create a new keyspace

**Headers**

* `Content-Type: application/json` - tells the service you wish to communicate with json
* `Authorization: Bearer <token>`  - the service key token

**Parameters**


* `name` <span style="border: 1px #AAA solid; padding: 2px; border-radius: 5px;">string</span> - The name of the keyspace (required)


* `keys_prefix` <span style="border: 1px #AAA solid; padding: 2px; border-radius: 5px;">string</span> - The prefix for the keys in the keyspace (required)


* `ratelimit` <span style="border: 1px #AAA solid; padding: 2px; border-radius: 5px;">object</span> - The default rate limit configuration for keys in the keyspace

  * `limit` <span style="border: 1px #AAA solid; padding: 2px; border-radius: 5px;">integer</span> - The rate limit


  * `refill_rate` <span style="border: 1px #AAA solid; padding: 2px; border-radius: 5px;">integer</span> - The rate at which the rate limit refills


  * `refill_interval` <span style="border: 1px #AAA solid; padding: 2px; border-radius: 5px;">integer</span> - The interval at which the rate limit refills




<details>
<summary> <b>cURL example</b> </summary>

```bash
$ curl -X POST
       -H "Content-Type: application/json" \
       -H "Authorization: Bearer <token>" \
       --data '
        {
          "name": "demo.yourapi.com (env: production)",
          "keys_prefix": "demo_",
          "ratelimit": {
            "limit": 100,
            "refill_rate": 1,
            "refill_interval": 1000
          }
        }' https://demo.driplim.it/v1/keyspaces.create
```

```json
{
  "ksid": "ks_abc",
  "name": "demo.yourapi.com (env: production)",
  "keys_prefix": "demo_",
  "ratelimit": {
    "limit": 100,
    "refill_rate": 1,
    "refill_interval": 1000
  }
}
```
</details>

### `POST /v1/keyspaces.delete`
Delete a keyspace

**Headers**

* `Content-Type: application/json` - tells the service you wish to communicate with json
* `Authorization: Bearer <token>`  - the service key token

**Parameters**


* `ksid` <span style="border: 1px #AAA solid; padding: 2px; border-radius: 5px;">string</span> - The id of the keyspace to delete (required)



<details>
<summary> <b>cURL example</b> </summary>

```bash
$ curl -X POST
       -H "Content-Type: application/json" \
       -H "Authorization: Bearer <token>" \
       --data '
        {
          "ksid": "ks_abc"
        }' https://demo.driplim.it/v1/keyspaces.delete
```

```json
null
```
</details>


## ServiceKeys

### `POST /v1/serviceKeys.current`
Get the current authenticated service key

**Headers**

* `Content-Type: application/json` - tells the service you wish to communicate with json
* `Authorization: Bearer <token>`  - the service key token

**Parameters**



<details>
<summary> <b>cURL example</b> </summary>

```bash
$ curl -X POST
       -H "Content-Type: application/json" \
       -H "Authorization: Bearer <token>" \ https://demo.driplim.it/v1/serviceKeys.current
```

```json
{
  "skid": "sk_uvw",
  "description": "cli generated admin service key at 2024-06-11T17:10:49+02:00",
  "admin": true,
  "keyspaces_policies": {
    "ks_abc": {
      "read": true,
      "write": true
    }
  },
  "created_at": "2024-06-11T17:10:49.746257293+02:00"
}
```
</details>

### `POST /v1/serviceKeys.get`
Get the service key by ID or by token

**Headers**

* `Content-Type: application/json` - tells the service you wish to communicate with json
* `Authorization: Bearer <token>`  - the service key token

**Parameters**


* `skid` <span style="border: 1px #AAA solid; padding: 2px; border-radius: 5px;">string</span> - The id of the service key to get (skid takes precedence over token)


* `token` <span style="border: 1px #AAA solid; padding: 2px; border-radius: 5px;">string</span> - The token of the service key to get



<details>
<summary> <b>cURL example</b> </summary>

```bash
$ curl -X POST
       -H "Content-Type: application/json" \
       -H "Authorization: Bearer <token>" \
       --data '
        {
          "skid": "sk_uvw",
          "token": ""
        }' https://demo.driplim.it/v1/serviceKeys.get
```

```json
{
  "skid": "sk_uvw",
  "description": "cli generated admin service key at 2024-06-11T17:10:49+02:00",
  "admin": true,
  "keyspaces_policies": {
    "ks_abc": {
      "read": true,
      "write": true
    }
  },
  "created_at": "2024-06-11T17:10:49.746258936+02:00"
}
```
</details>

### `POST /v1/serviceKeys.list`
List all service keys

**Headers**

* `Content-Type: application/json` - tells the service you wish to communicate with json
* `Authorization: Bearer <token>`  - the service key token

**Parameters**


* `list` <span style="border: 1px #AAA solid; padding: 2px; border-radius: 5px;">object</span> - The list options

  * `page` <span style="border: 1px #AAA solid; padding: 2px; border-radius: 5px;">integer</span> - The page number


  * `limit` <span style="border: 1px #AAA solid; padding: 2px; border-radius: 5px;">integer</span> - The number of items per page




<details>
<summary> <b>cURL example</b> </summary>

```bash
$ curl -X POST
       -H "Content-Type: application/json" \
       -H "Authorization: Bearer <token>" \
       --data '
        {
          "list": {
            "page": 1,
            "limit": 10
          }
        }' https://demo.driplim.it/v1/serviceKeys.list
```

```json
{
  "list": {
    "page": 1,
    "limit": 10,
    "last_page": 1
  },
  "service_keys": [
    {
      "skid": "sk_uvw",
      "description": "cli generated admin service key at 2024-06-11T17:10:49+02:00",
      "admin": true,
      "keyspaces_policies": {
        "ks_abc": {
          "read": true,
          "write": true
        }
      },
      "created_at": "2024-06-11T17:10:49.746260221+02:00"
    }
  ]
}
```
</details>

### `POST /v1/serviceKeys.delete`
Delete a service key

**Headers**

* `Content-Type: application/json` - tells the service you wish to communicate with json
* `Authorization: Bearer <token>`  - the service key token

**Parameters**


* `skid` <span style="border: 1px #AAA solid; padding: 2px; border-radius: 5px;">string</span> - The id of the service key to delete (required)



<details>
<summary> <b>cURL example</b> </summary>

```bash
$ curl -X POST
       -H "Content-Type: application/json" \
       -H "Authorization: Bearer <token>" \
       --data '
        {
          "skid": "sk_uvw"
        }' https://demo.driplim.it/v1/serviceKeys.delete
```

```json
null
```
</details>

### `POST /v1/serviceKeys.create`
Get the service key by ID or by token

**Headers**

* `Content-Type: application/json` - tells the service you wish to communicate with json
* `Authorization: Bearer <token>`  - the service key token

**Parameters**


* `description` <span style="border: 1px #AAA solid; padding: 2px; border-radius: 5px;">string</span> - The description of the service key


* `admin` <span style="border: 1px #AAA solid; padding: 2px; border-radius: 5px;">bool</span> - The admin flag of the service key


* `keyspaces_policies` <span style="border: 1px #AAA solid; padding: 2px; border-radius: 5px;">map</span> - The keyspaces policies of the service key. Map keys are the keyspace ids and the values are the policies for the keyspace

  * <span style="border: 1px #AAA solid; padding: 2px; border-radius: 5px;">string</span> - keys of the map (required)

    * `read` <span style="border: 1px #AAA solid; padding: 2px; border-radius: 5px;">bool</span> - Read permission

    * `write` <span style="border: 1px #AAA solid; padding: 2px; border-radius: 5px;">bool</span> - Write permission




<details>
<summary> <b>cURL example</b> </summary>

```bash
$ curl -X POST
       -H "Content-Type: application/json" \
       -H "Authorization: Bearer <token>" \
       --data '
        {
          "description": "api generated non admin service key",
          "admin": false,
          "keyspaces_policies": {
            "ks_abc": {
              "read": true,
              "write": false
            }
          }
        }' https://demo.driplim.it/v1/serviceKeys.create
```

```json
{
  "skid": "sk_uvw",
  "description": "api generated non admin service key",
  "admin": true,
  "keyspaces_policies": {
    "ks_abc": {
      "read": true,
      "write": false
    }
  },
  "created_at": "2024-06-11T17:10:49.746262098+02:00"
}
```
</details>



## Errors

If HTTP response code is greater than or equal to 400, the api returns a json object indicating the reason of the failure:

```json
{
  "error": "the reason of the failure",
  "invalid_fields": [
    "field1",
    "field2"
  ]
}
```

`invalid_fields` can also be integrated in the error response if one or more input parameters are invalids.

### HTTP response code

* `200` ok
* `204` created
* `400` invalid payload
* `401` unauthorized
* `403` cannot delete itself
* `404` not found
* `409` already exists
* `419` key expired
* `429` rate limit exceeded
