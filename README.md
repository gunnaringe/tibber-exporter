# Prometheus Exporter for Tibber

This implements a Prometheus exporter for the energy data exposed through Tibber's API.

The exporter requires that you use Tibber as your energy company, but is in no way endorsed by Tibber.
None of the authors or Tibber can be responsible for your use of this project.

Visit https://tibber.com for more information

For API documentation and demo access token, visit https://developer.tibber.com/explorer

## Build
This project requires a Go version higher than 1.11

```bash
go build
```

## Configure and run

Scraping will occur synchronous when you pull the metrics endpoint.

##### Help
```bash
Usage:
  tibber_exporter [OPTIONS]

Application Options:
  -t, --token=    Authorization token [$TIBBER_TOKEN]
  -e, --endpoint= Endpoint (default: https://api.tibber.com/v1-beta/gql) [$TIBBER_ENDPOINT]
  -l, --listen=   Address and port to listen (default: :9501) [$TIBBER_LISTEN]
  -v, --verbose   Show verbose debug information

Help Options:
  -h, --help      Show this help message
```

##### Example
```bash
export TIBBER_TOKEN=d1007ead2dc84a2b82f0de19451c5fb22112f7ae11d19bf2bedb224a003ff74a

# Run in background
./tibber-exporter &

# Invoke a scrape
curl localhost:9501
```

## Sample output (using demo access token)
```
# HELP go_gc_duration_seconds A summary of the GC invocation durations.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 1.0598e-05
go_gc_duration_seconds{quantile="0.25"} 1.508e-05
go_gc_duration_seconds{quantile="0.5"} 1.8913e-05
go_gc_duration_seconds{quantile="0.75"} 2.8612e-05
go_gc_duration_seconds{quantile="1"} 2.8612e-05
go_gc_duration_seconds_sum 7.3203e-05
go_gc_duration_seconds_count 4
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
go_goroutines 8
# HELP go_info Information about the Go environment.
# TYPE go_info gauge
go_info{version="go1.12.1"} 1
# HELP go_memstats_alloc_bytes Number of bytes allocated and still in use.
# TYPE go_memstats_alloc_bytes gauge
go_memstats_alloc_bytes 2.990144e+06
# HELP go_memstats_alloc_bytes_total Total number of bytes allocated, even if freed.
# TYPE go_memstats_alloc_bytes_total counter
go_memstats_alloc_bytes_total 1.1911624e+07
# HELP go_memstats_buck_hash_sys_bytes Number of bytes used by the profiling bucket hash table.
# TYPE go_memstats_buck_hash_sys_bytes gauge
go_memstats_buck_hash_sys_bytes 1.445271e+06
# HELP go_memstats_frees_total Total number of frees.
# TYPE go_memstats_frees_total counter
go_memstats_frees_total 201532
# HELP go_memstats_gc_cpu_fraction The fraction of this program's available CPU time used by the GC since the program started.
# TYPE go_memstats_gc_cpu_fraction gauge
go_memstats_gc_cpu_fraction 7.379787276297374e-07
# HELP go_memstats_gc_sys_bytes Number of bytes used for garbage collection system metadata.
# TYPE go_memstats_gc_sys_bytes gauge
go_memstats_gc_sys_bytes 2.381824e+06
# HELP go_memstats_heap_alloc_bytes Number of heap bytes allocated and still in use.
# TYPE go_memstats_heap_alloc_bytes gauge
go_memstats_heap_alloc_bytes 2.990144e+06
# HELP go_memstats_heap_idle_bytes Number of heap bytes waiting to be used.
# TYPE go_memstats_heap_idle_bytes gauge
go_memstats_heap_idle_bytes 6.144e+07
# HELP go_memstats_heap_inuse_bytes Number of heap bytes that are in use.
# TYPE go_memstats_heap_inuse_bytes gauge
go_memstats_heap_inuse_bytes 4.947968e+06
# HELP go_memstats_heap_objects Number of allocated objects.
# TYPE go_memstats_heap_objects gauge
go_memstats_heap_objects 16698
# HELP go_memstats_heap_released_bytes Number of heap bytes released to OS.
# TYPE go_memstats_heap_released_bytes gauge
go_memstats_heap_released_bytes 6.1382656e+07
# HELP go_memstats_heap_sys_bytes Number of heap bytes obtained from system.
# TYPE go_memstats_heap_sys_bytes gauge
go_memstats_heap_sys_bytes 6.6387968e+07
# HELP go_memstats_last_gc_time_seconds Number of seconds since 1970 of last garbage collection.
# TYPE go_memstats_last_gc_time_seconds gauge
go_memstats_last_gc_time_seconds 1.5539008929008768e+09
# HELP go_memstats_lookups_total Total number of pointer lookups.
# TYPE go_memstats_lookups_total counter
go_memstats_lookups_total 0
# HELP go_memstats_mallocs_total Total number of mallocs.
# TYPE go_memstats_mallocs_total counter
go_memstats_mallocs_total 218230
# HELP go_memstats_mcache_inuse_bytes Number of bytes in use by mcache structures.
# TYPE go_memstats_mcache_inuse_bytes gauge
go_memstats_mcache_inuse_bytes 13888
# HELP go_memstats_mcache_sys_bytes Number of bytes used for mcache structures obtained from system.
# TYPE go_memstats_mcache_sys_bytes gauge
go_memstats_mcache_sys_bytes 16384
# HELP go_memstats_mspan_inuse_bytes Number of bytes in use by mspan structures.
# TYPE go_memstats_mspan_inuse_bytes gauge
go_memstats_mspan_inuse_bytes 63504
# HELP go_memstats_mspan_sys_bytes Number of bytes used for mspan structures obtained from system.
# TYPE go_memstats_mspan_sys_bytes gauge
go_memstats_mspan_sys_bytes 81920
# HELP go_memstats_next_gc_bytes Number of heap bytes when next garbage collection will take place.
# TYPE go_memstats_next_gc_bytes gauge
go_memstats_next_gc_bytes 4.194304e+06
# HELP go_memstats_other_sys_bytes Number of bytes used for other system allocations.
# TYPE go_memstats_other_sys_bytes gauge
go_memstats_other_sys_bytes 1.514337e+06
# HELP go_memstats_stack_inuse_bytes Number of bytes in use by the stack allocator.
# TYPE go_memstats_stack_inuse_bytes gauge
go_memstats_stack_inuse_bytes 720896
# HELP go_memstats_stack_sys_bytes Number of bytes obtained from system for stack allocator.
# TYPE go_memstats_stack_sys_bytes gauge
go_memstats_stack_sys_bytes 720896
# HELP go_memstats_sys_bytes Number of bytes obtained from system.
# TYPE go_memstats_sys_bytes gauge
go_memstats_sys_bytes 7.25486e+07
# HELP go_threads Number of OS threads created.
# TYPE go_threads gauge
go_threads 15
# HELP process_cpu_seconds_total Total user and system CPU time spent in seconds.
# TYPE process_cpu_seconds_total counter
process_cpu_seconds_total 0.09
# HELP process_max_fds Maximum number of open file descriptors.
# TYPE process_max_fds gauge
process_max_fds 1024
# HELP process_open_fds Number of open file descriptors.
# TYPE process_open_fds gauge
process_open_fds 8
# HELP process_resident_memory_bytes Resident memory size in bytes.
# TYPE process_resident_memory_bytes gauge
process_resident_memory_bytes 1.4270464e+07
# HELP process_start_time_seconds Start time of the process since unix epoch in seconds.
# TYPE process_start_time_seconds gauge
process_start_time_seconds 1.55390021588e+09
# HELP process_virtual_memory_bytes Virtual memory size in bytes.
# TYPE process_virtual_memory_bytes gauge
process_virtual_memory_bytes 1.11200256e+09
# HELP process_virtual_memory_max_bytes Maximum amount of virtual memory available in bytes.
# TYPE process_virtual_memory_max_bytes gauge
process_virtual_memory_max_bytes -1
# HELP promhttp_metric_handler_requests_in_flight Current number of scrapes being served.
# TYPE promhttp_metric_handler_requests_in_flight gauge
promhttp_metric_handler_requests_in_flight 1
# HELP promhttp_metric_handler_requests_total Total number of scrapes by HTTP status code.
# TYPE promhttp_metric_handler_requests_total counter
promhttp_metric_handler_requests_total{code="200"} 4
promhttp_metric_handler_requests_total{code="500"} 0
promhttp_metric_handler_requests_total{code="503"} 0
# HELP tibber_current_energy_price Current energy price
# TYPE tibber_current_energy_price gauge
tibber_current_energy_price{currency="",home_id="68e6938b-91a6-4199-a0d4-f24c22be87bb",type="energy"} 0
tibber_current_energy_price{currency="",home_id="68e6938b-91a6-4199-a0d4-f24c22be87bb",type="tax"} 0
tibber_current_energy_price{currency="SEK",home_id="c70dcbe5-4485-4821-933d-a8a86452737b",type="energy"} 0.3407
tibber_current_energy_price{currency="SEK",home_id="c70dcbe5-4485-4821-933d-a8a86452737b",type="tax"} 0.1589
# HELP tibber_home_info Home information
# TYPE tibber_home_info gauge
tibber_home_info{address1="Kungsgatan 8",address2="",address3="",city="Stockholm",country="SE",home_id="c70dcbe5-4485-4821-933d-a8a86452737b",latitude="59.3362066",longitude="18.0675126",postal_code="11759",time_zone="Europe/Stockholm"} 1
tibber_home_info{address1="Winterfell",address2="",address3="",city="Førde",country="NO",home_id="68e6938b-91a6-4199-a0d4-f24c22be87bb",latitude="61.457539",longitude="5.829413",postal_code="6812",time_zone="Europe/Oslo"} 1
# HELP tibber_previous_hour_consumption_watt_hour Total Watt hours used last hourly period
# TYPE tibber_previous_hour_consumption_watt_hour gauge
tibber_previous_hour_consumption_watt_hour{home_id="c70dcbe5-4485-4821-933d-a8a86452737b"} 0
# HELP tibber_previous_hour_energy_price Energy price last hourly period
# TYPE tibber_previous_hour_energy_price gauge
tibber_previous_hour_energy_price{currency="SEK",home_id="c70dcbe5-4485-4821-933d-a8a86452737b",type="energy"} 0.4164875
tibber_previous_hour_energy_price{currency="SEK",home_id="c70dcbe5-4485-4821-933d-a8a86452737b",type="vat"} 0.0832975
# HELP tibber_previous_hour_total_cost Total cost last hourly period
# TYPE tibber_previous_hour_total_cost gauge
tibber_previous_hour_total_cost{currency="SEK",home_id="c70dcbe5-4485-4821-933d-a8a86452737b"} 0
# HELP tibber_requests_total Number of requests done by exporter
# TYPE tibber_requests_total counter
tibber_requests_total 5
# HELP tibber_scrape_duration_seconds Scrape duration of last scrape
# TYPE tibber_scrape_duration_seconds gauge
tibber_scrape_duration_seconds 0.123476135
```
