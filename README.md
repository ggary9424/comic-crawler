# Comic Crawler

Mainly, this project is created for me to learn the `Golang`. Besides, I want to use this project to save the time for checkout the comics I followed by browsing each URL everyday.

## Goal

1. Familiar with `Golang`.
2. Try to build a project that is not only a toy but also a project can run at production environment.

## Limitation

Now, only support the [cartoonmad](https://www.cartoonmad.com/) website.

## Environment

1. Golang v1.12
2. Dep v0.5.0
3. Mongo DB v2.6 or higher

## Execution

```bash
dep ensure
go run main.go
```

## TODO

- [ ] An queue-worker architecture can broadcast the notification of following comics updating.
- [ ] An User Interface can edit the following comics. Try to use [GRPC Web](https://www.google.com/search?q=grpc+web&oq=grpc+web&aqs=chrome..69i57j69i60l2j69i59j69i60j0.2801j0j4&sourceid=chrome&ie=UTF-8)
- [ ] Comprehensive log system.
- [ ] Dockerization version.
- [ ] Elegant way to deal with callback hell when using colly. (Maybe use `Go Channel`)
