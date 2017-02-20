Shortly : a simple url shortening service

TODO: User can claim custom slug

# Design

### Use cases
User can provide a url and the service will return a unique
shorter url which, when visited, will redirect (302) to the provided url. 

User can provide 

This is exposed as an http service with form/post and a simple web ui.

TODO user may optional provide a custom "slug" which is the unique bit of the Url.

### Calculating unique slugs

Considerations:
* short keys preferred
* collisions in keys costly (cost extra DB access at shorten-time)
* human type-ability preferred

Design:

* use hash on the url to compute random value that should be unique per url
* use only alpha-numeric characters in slug to support human input. (36 possible values per char)
* trim the slug to some small number of characters


Calculate the slug as follows:
```
slug(url, size):
    return base_36_encode(md5(url+seed))[:size]
```

The size of the slug (num chars) represents a balance of the two considerations. As we store more and more slugs, the 
chance of collision with existing slugs increases. 
TODO, more  here

After calculating a slug, we must ensure its not already taken in the DB. If it is,
we add some seed the the url (md5 input) and retry.



### Data storage
* Its just a map from slug to url.
* Durability is good (- for redis?)
* Want high throughput, able to scale horizontally

Design to support mongo and redis for storage

# Building

* Install Docker
* run `make`. This will build, start services, run unit/performance test. 
Shortly service is at `http://localhost:8080/` 
* run `make clean` to stop services

# Performance Metrics

Run via docker, with test suit in separate container. 100 concurrent workers making 1000 requests each with 90/10 mix of reads/writes   

Memory:
Read operations: 8051 (1493.706501 op/sec)
Write operations: 1041 (193.137308 op/sec)

Redis:
Read operations: 8185 (839.326457 op/sec)
Write operations: 1023 (104.902989 op/sec)

Mongo:
Read operations: 8185 (877.275351 op/sec)
Write operations: 1062 (113.826075 op/sec)