# concurrency-tester

## start test

1. create config file:

```bash
cp .env.template .env
vim .env
```

2. start test

```bash
make run
```

## examine results

1. get test containers

```bash
docker ps -a | grep concurrency-tester
```

2. get logs

```bash
docker logs <container>
```

## our results

Test details:
- the test app was running in `5` containers
- each test app executed creations on `2` go routines
- each go routine created `15` records

Test results:
- [container 1](results/0527e508fd55.log)
- [container 2](results/22d5f3446a81.log)
- [container 3](results/c775e5e4989b.log)
- [container 4](results/de55d2c9ddf8.log)
- [container 5](results/ff61e1cd2a3c.log)

Maximum number of parallel operations: `5 * 2 = 10`

Total number of records created: `150`

Total number of conflicts: `461`

Total execution time: `~64s`
