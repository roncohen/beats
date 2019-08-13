For manual testing and development of this module, start metricbeat with the standard configuration:

```
["source","yaml",subs="attributes"]
------------------------------------------------------------------------------
- module: statsd
  metricsets: ["server"]
  host: "localhost"
  port: "8125"
  #reservoir_size: 1000
  #ttl: "30s"
------------------------------------------------------------------------------
```

then use a statsd client to test the features:

```
$ npm install statsd-client
$ node
> var SDC = require('statsd-client'),
    sdc = new SDC({host: 'localhost', port: 8125});

> sdc.increment('systemname.subsystem.value'); // Increment by one
> sdc.gauge('what.you.gauge', 100);
> sdc.gaugeDelta('what.you.gauge', -70); // Will now count 50
> sdc.gauge('gauge.with.tags', 100, {foo: 'bar'});
> sdc.set('set.with.tags', 100, {foo: 'bar'});
> sdc.set('set.with.tags', 200, {foo: 'bar'});
> sdc.set('set.with.tags', 100, {foo: 'baz'});
....

<CTRL+D>
```
