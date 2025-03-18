# Settle Backend

[![Go 1.23x](https://img.shields.io/badge/Go-1.23.x-blue.svg)](https://go.dev/) [![Go Build](https://github.com/kwa0x2/Settle-Backend/actions/workflows/ci.yml/badge.svg)](https://github.com/kwa0x2//Settle-Backend/actions/workflows/ci.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/kwa0x2/settle-backend?style=flat-square)](https://goreportcard.com/report/github.com/kwa0x2/settle-backend) [![Website](https://img.shields.io/badge/Website-chat.nettasec.com-red.svg)](https://chat.nettasec.com/)

A backend for a gamer-focused chat application where only users with over 500 hours on Steam can join. Built with MongoDB, it integrates OpenID for Steam-based authentication, real-time messaging with Socket.IO, and file handling with Amazon S3. This project follows Uncle Bob's Clean Architecture principles to ensure a robust, maintainable structure. Monitoring is implemented with Grafana and Prometheus, along with Node Exporter for comprehensive dashboard insights.

## Diagram:

![diagram](https://github.com/kwa0x2/Settle-Backend/blob/a05fd4bdbd1a41ee45e24c37d5348ad0a83a1075/diagram.png)

## Arch:

This project is structured following the Clean Architecture principles, inspired by [go-clean-arch](https://github.com/bxcodec/go-clean-arch) by bxcodec. The architecture ensures separation of concerns, making the codebase scalable, maintainable, and testable. Each layer has clearly defined responsibilities, promoting a strong dependency rule from outer to inner layers.

![architecture](https://raw.githubusercontent.com/bxcodec/go-clean-arch/master/clean-arch.png)
