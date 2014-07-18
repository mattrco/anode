Anode
=====

**New**: there's now an IRC channel for anode! Direct your client to #anode on Freenode, or use [webchat from your browser](https://webchat.freenode.net/?url=irc%3A%2F%2Firc.freenode.net%2Fanode).

Anode is a tool for experimenting with different analysis algorithms on metrics and other time series.

Data analysis starts with getting to know your data. The overarching goal of anode is to advance the open-source offerings for metrics analysis by creating a proving ground for different approaches.

Anode ships with a Graphite input plugin, capable of fetching a named metric from graphite and streaming updates to registered analysis plugins. Each analysis plugin then streams its result to output plugins. 

Just learning about metrics analysis? Have a look at this [excellent talk on anomaly detection for metrics](http://www.metaforsoftware.com/monitorama-pdx-simple-math-for-anomaly-detection/).

Example
-------

Try the [three sigma](http://www.encyclopediaofmath.org/index.php/Three-sigma_rule) analyzer that's included. This works best for data with a normal (gaussian) distribution.

```
go get github.com/mattrco/anode.exp
anode.exp -metric=app.latency
```

(Other flags exist; [see here](https://github.com/mattrco/anode.exp/blob/master/main.go#L21)).

This will fetch and process the last 24 hours of data. New metrics will appear in graphite under `anode.threesig` which you can then plot alongside your existing metric. The screenshot below shows the original metric (lilac) with anomalous values highlighted in orange. 

![Screenshot of Three Sigma analyzer](threesigma.png)

Coming soon
-----------

* More analyzers, outputs, docs, and better-defined APIs.

* Write-ups on statistical (and other) approaches to metrics analysis. You can [follow along on twitter](https://twitter.com/mattrco) to stay on top of these, but I'll link to them from this doc too.

Contributions
-------------

Contributions welcome! Please open an issue to discuss what you'd like to work on.

Related projects and acknowledgement
------------------------------------

[mozilla-services/heka](https://github.com/mozilla-services/heka)

Heka is a more general data collection and processing system, processing much more than time series. Its data pipeline inspired anode's architecture.

[etsy/skyline](https://github.com/etsy/skyline)

Skyline is an anomaly detection system for time series data.

Key differences:

* Skyline is production-ready and more work to get set up
* Fairly rigid decision logic for determining whether latest updates to a series are anomalous. Algorithms each return a boolean and majority vote wins. There's no way to apply different analyses to different metrics, which is a key design point of anode.

