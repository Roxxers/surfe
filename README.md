# Surfe Tech Test

## Installation

To get up and running with the API

```sh
git clone https://github.com/roxxers/surfe
cd surfe
go mod tidy # grab dependencies without needing to build
go run ./cmd/main.go # or go build ...
```

## Usage

All of the API is mounted to the path `/api/v1`. If you see an endpoint listed below, it will be subseded by `/api/v1`. Examples will show the full path to request when running the server unmodified.

### Task 1

#### GET: Fetch user

`/user/:id` replacing :id with the user you want to fetch.

##### Example
```sh
> curl 127.0.0.1:8080/api/v1/user/1
{"id":0,"name":"Allyson","created_at":"0001-01-01T00:00:00Z"}
```

### Task 2

#### GET: Get count of all actions performed by user

`/user/:id/actioncount` replacing :id with the user you want to query

##### Example
```sh
> curl 127.0.0.1:8080/api/v1/user/1/actioncount
{"count":49}
```


### Task 3

#### POST: Get probablity of next 

`/actions/probablity` using a JSON request body like:

```json
{
    "action": "EXAMPLE_ACTION",
}
```

##### Example
```sh
> curl -X POST 127.0.0.1:8080/api/v1/actions/probablity -d '{"action": "CONNECT_CRM"}'
{"ADD_CONTACT":0.32,"EDIT_CONTACT":0.33,"REFER_USER":0.03,"VIEW_CONTACTS":0.31}
```

### Task 4

#### GET: Show user referal indexes for all users

`/users/referalindex`

##### Example
```sh
> curl 127.0.0.1:8080/api/v1/users/referalindex
{"0":0,"1":1,"10":1,"100":0,"101":0,"102":0,"103":0,"104":3,"105":0,"106":0,"107":0,"108":0,....
```