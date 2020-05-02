# Sentryceptor

Sentryceptor intercepts the HTTP requests, filters body request and redirects to Sentry instance.

The best method is to [add filters in your code](https://docs.sentry.io/error-reporting/configuration/filtering/), but when you can't modify the code source, 
Sentryceptor can help you to filter the information from each body request to Sentry. You keep the control of you data.


## How to run?

To run with default flags :

```bash
sentryceptor > sentryceptor.log
```

To run with specific config file :

```bash
sentryceptor -config=PATH/YOUR/FILE.json > sentryceptor.log
```


## Configuration

The configuration file by default is `config.json`, but you can use another file (see above).


### Comparator available


Comparator  | Description
------------|------------
  eq        | equals condition
  neq       | not equals condition
  lt        | less than condition
  lte       | less than or equal condition
  gt        | greater than condition
  gte       | greater than or equal condition
  matches   | matches (regex) condition
  contains  | contains condition
  ncontains | not contains condition
  exists    | exists condition (value field must be equal to `""`)
  nexists   | not exists condition (value field must be equal to `""`)

