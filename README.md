# 📡 PanCORS

We all know how it feels when you want to fetch something and CORS screws it all up.

**Meet PanCORS:** super light-weight CORS-anywhere proxy implemented in only 2<sup>6</sup> lines of Go code.

## Installation

1. Clone the repo
2. `docker build -t pancors .`
3. `docker run --rm -d -p 8080:8080/tcp pancors:latest`

_Stop CORS oppresion!_

![How CORS works](assets/cors_explanation.jpg)

> CORS keeps us all safe but sometimes it just gets in the way. _Peace_
