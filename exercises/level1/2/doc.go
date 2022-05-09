//Use next curl string to test encoding:
//curl localhost:8080/encode
//Use next curl string to test decoding:
//curl -s -XGET -H "Content-type: application/json" -d '{ "FirstName": "James", "LastName": "Bond", "Age": 33 }' 'localhost:8080/decode'
package main
