# Pub/Sub Direct Push Bridge

This repository provides a lightweight gRPC server that acts as a bridge between the Google Cloud Pub/Sub Emulator and
HTTP-based receivers. It is designed to forward publish events directly as HTTP POST requests, bypassing the need for
topic and subscription configuration.

## Purpose

The bridge enables push-style event delivery in local development environments. It is particularly useful when
simulating Google Cloud Storage notifications in combination
with [`fake-gcs-server`](https://github.com/fsouza/fake-gcs-server).

## Installation

To install the application using Go:

```bash
go install github.com/devnsi/pubsub-direct-push/cmd/bridge@latest
````

## Usage

Start the gRPC server and configure the Pub/Sub Emulator to publish messages to it. The server will forward incoming
messages to the specified HTTP endpoint.
