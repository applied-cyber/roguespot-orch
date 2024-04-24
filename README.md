# roguespot-orch
The central server for receiving and analyzing information about access
points from the Roguespot wardrivers. It currently flags access points that
are seen for the first time or whose vendor is neither Cisco/Juniper (based on
the BSSID).

```bash
$ make
$ ./roguespot-orch
```