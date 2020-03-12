
# Go exporting metrics sample project

## Run
```
go run main.go
```

---

## APIs

### Echo server
```
:8080/echo
:8080/echo/{msg}
```

### Prometheus metrics
```
:9090/metrics
```

---

## Please note

Here we use `promauto` module instead of normal `prometheus` one, so we can avoid to manually register the Prometheus collector with kind of following command

```
prometheus.MustRegister(myCustomMetric)
```

---

## Links
* https://prometheus.io/docs/guides/go-application/
* https://scot.coffee/2018/12/monitoring-go-applications-with-prometheus/
* https://levelup.gitconnected.com/multi-stage-docker-builds-with-go-modules-df23b7f91a67
