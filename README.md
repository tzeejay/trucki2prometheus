## trucki2prometheus
Pure Go(-lang) prometheus exporter for the Lumentree inverter feat. Trucki T2SG (Trucki2Sun Gateway) stick enabling easy state & energy production monitoring in Grafana of your grid-tied solar systems' zero-export-control (nulleinspeisung) operations.  
The only non-go-stdlib dependencies are the Go Prometheus client SDKs, no other import should be required to operate this tool making it very easy to maintain, compile and transpile to various platforms

### What problem am I trying to solve
The Trucki stick comes with MQTT tooling built in, though as I don't have a MQTT deployment and don't see myself using MQTT ever, this capability did not solve my problem, desire really, to monitor everything happening on the inverter side of the system. One possible solution to my desire to record & monitor what the Trucki stick is doing would have been to run some kind of MQTT broker in my network as well as a MQTT to Prometheus exporter/bridge.  
Following down this path would have prevented this project from being created in the first place but as I already said I see no value in setting up MQTT-anything in my life, this would have been a lot of overhead and additional complexity that I did not want to introduce which is why I wrote this little Go tool instead. I hope it helps you as well :)


## ¡Unknown values! Please help
There are a handful of values that I am not yet sure how to export/translate/document properly. If you have any of these set due to leveraging a different configuration of devices please open a new issue and include what the values mean in your PV-system's context or how they are used so that the Prometheus help text can be phrased accordingly 


### Currently unknown JSON values
```
- SUN2ROUNDTRIP <- this should be an int to me. Why is this a string?
- SUN3ROUNDTRIP <- this should be an int to me. Why is this a string?
- DAYENERGY <- day as in date borders or over the last 24 hours?
- TOTALENERGY <- total as in persistent since start of use (last factory reset) or total as in since last reboot?
- METERDAYENERGY <- day as in date borders or over the last 24 hours?
- ZEPCPOWER <- how is this value calculated exactly?
- WIFI <- are there more states than 'Connected' and possibly 'DISCONNECTED'? Are those two correct at all?
- RSSI <- do string values here follow the common RSSI string descriptions? Are there other possible values?
```

## Build
A functional Go(-lang) toolchain needs to be a setup on your computer in order to compile this project. A release build (eg. for an ARM based host like a Raspberry Pi) should be created with

`GOOS=linux GOARCH=arm go build -a -v -o trucki2prometheus -ldflags "-X main.t2PromVersion=$(git rev-parse --short HEAD) -X main.buildDate=$(date +"%Y-%m-%dT%H:%M:%S")"`

## Automated operations
I currently run the exporter as a systemd unit on a Linux based ARM host. It possible to set this up by placing the following configuration into `/etc/systemd/system/trucki2prometheus.service`
```
[Unit]
Description=Trucki stick to Prometheus exporter
After=network.target

[Service]
Type=exec
User=root
ExecStart=/opt/trucki2prometheus/trucki2prometheus -t 
Restart=on-failure

[Install]
WantedBy=default.target
```
