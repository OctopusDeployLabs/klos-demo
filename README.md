# klos-demo

This repository contains a demo for the Kubernetes Live Object Status (LOS) feature, which allows you to use Octopus 
to view the status of Kubernetes resources in real-time. 

## Overview
The demo includes 3 simple applications, intended to be deployed as independent projects in Octopus Deploy:

### KLOSWorker
A simple web application that has three endpoints:
- `/`: A simple web page that returns 200 all the time
- `/ceo`: A simple web page that returns 200 ONLY when `LetTheCeoDoThings` is set to `true`
- `/healthz`: A simple web page that returns 200 when the pod is running
This is to simulate a webserver that can be misconfigured

### KLOSCache
A simple application that loads up 300MiB of data into memory, intended to simulate a long-running process running out
of memory. It has a `/healthz` endpoint that returns 200 when the pod is running.

### LoadGenerator
A simple application that generates load on the KLOSWorker application. 
It has a `/healthz` endpoint that returns 200 when the pod is running.
It frequently attempts to access the `/` and `/ceo` endpoints of the KLOSWorker application, logging the results.