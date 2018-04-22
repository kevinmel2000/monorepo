# SimpleNSQ

This purpose of library is to replace 'real' NSQ with the 'fake' or called `simplensq` when doing **test**

We can't spawn a full-fledged NSQ cluster in a CI platform, and might be too heavy for testing.

## How to use

To use `simplensq` along with NSQ client, you need to provide an interface to satisfy both.