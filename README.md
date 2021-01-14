# Phppoolstat (a Telegraf exec)

Phppoolstat tallies the PHP-FPM pool workers per pool. The output of
phppoolstat can be fed into [Telegraf](https://www.influxdata.com/time-series-platform/telegraf/)
as influx format.

This tool is build for speed and a low impact on the host machine. So data is
gathered from `/proc/*/cmdline` only. In the `cmdline` of a process is stated
for which PHP-FPM pool the process is. More detailed explanation of the output
can be found below at [measurements](#measurements).

## Example usage in Telegraf

```
[[inputs.exec]]
  commands = ["/usr/bin/phppoolstat --tag env=production"]
  data_format = "influx"
```

The above configuration would result in output like:
```
> phppoolstat,host=tengu,env=production,pool=prd_pool count=12i 1562471764000000000
> phppoolstat,host=tengu,env=production,pool=acc_pool count=6i 1562471764000000000
> phppoolstat,host=tengu,env=production,pool=_all_ count=18i 1562471764000000000
```

## Measurements

For each pool the number of running workers is shown as an integer:

phppoolstat.count value=12

## Package building

To build a package, the tool [fpm](https://github.com/jordansissel/fpm) is
used, please follow the [installation guide](https://fpm.readthedocs.io/en/latest/installing.html)
to be able to build a package.

The following OSes are supported:

- Debian, Ubuntu: `make deb`
- CentOS, RHEL: `make rpm`
