# Address Book

This is a practice API written in Go simulating a simple address book. The intent of this repo 
is to showcase a simple API which has test coverage, CRUD endpoints and a csv export capability.

## API Endpoints

Baseurl: `http://localhost:9090/api/`

| Path | Operation | Comment |
|---|---|---|
| `/address` | `GET` | retrieve all entries |
| `/address` | `POST` | create a new entry|
| `/address/<pk>` | `PUT` | update one or more fields for a single entry |
| `/address/<pk>` | `GET` | retrieve a single entries |
| `/address/<pk>` | `DELETE` | delete a entries |
| `/address/export` | `GET` | export the address book to csv format |

## Data

The address book must contain the following fields

| Name | Type |
|---|---|
| fullName | string |
| lastName | string |
| email | string |
| phoneNumber | int |

## How to test


## How to run

To run the application,
