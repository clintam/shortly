Shortly : a simple url shortening service designed to scale out

TODO: User can claim custom slug

# Design

### Use cases
User can provide a url and the service will return a unique
shorter url which, when visited, will redirect (302) to the provided url. 

TODO user may optional provide a custom "slug" which is the unique bit of the Url.

### High level design

* Exposed as an http service with form/post and a simple web ui.
** form to provide url and get the "redirecting" url
** http GET handler on the redirecting url to  
* The urls that are returned will be of the form `http://${host}/{$slug}.`
* 


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

After calculating a slug, we must ensure its not already taken in the DB. If it is,
we add some seed the the url (md5 input) and retry.
 
The size of the slug (num chars) represents a balance of the two considerations. As we store more and more slugs, the 
chance of collision with existing slugs increases. 

We could allow the slug size to increase as we scale out, or fix it with some expectation on the maximum number of 
urls we expect to handle. At full capacity, we want to avoid excessive number of conflicts (and thus re-hashes). Thus, calculate the size as follows:
```
slug_size(expected_total):
    log_base_36(expected_total * 2))
```

Here we make the assumptions that we expect to run for 5 years at 100 shortens/second = 15B rows

Thus we get a slug_size recommendation of ~6.5 so will use 7 charecters in the slug. 

### Data storage
* Its just a map from slug to url.
* Durability is good (- for redis?)
* Want high throughput (quick lookup by slug), and able to scale horizontally to support large number of entries.

Design to support mongo and redis for storage.

# Implementation

* Use go
* Layers for http server, link shortener, and link storage
* Interface to support pluggable storage layers

# Scaling 

* horizontal scale the http server by placing behind a LB
* horizontal scale the data layer via more mongo/redis nodes (redis looks trickier-but-possible to do here)

# Building

* Install Docker
* run `make`. This will build, start services, run unit tests.
* run `make perftest` to run performance tests against the various storages. 
Shortly service is at `http://localhost:8080/` 
* run `make clean` to stop services

# Performance Metrics

Run via docker, with test suit in separate container. 

First do a warmup with 100K url-shorten requests (writes). 
Then, test with 100 concurrent workers making 100 requests each with 90/10 mix of reads/writes   

Memory:
Read operations: 8935 (1104.636671 op/sec)
Write operations: 1065 (131.666262 op/sec)

Redis:
Read operations: 8906 (1029.608662 op/sec)
Write operations: 1094 (126.475621 op/sec)

Mongo:
Read operations: 8882 (913.073504 op/sec)
Write operations: 1118 (114.930891 op/sec)