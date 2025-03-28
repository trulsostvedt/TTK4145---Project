# Network Module

This module provides essential networking utilities for the elevator control system. It includes functionality for broadcasting messages, establishing connections, discovering peers, and retrieving local IP addresses.

---

## Components

### `bcast`
The `bcast` package is responsible for broadcasting messages over the network using UDP. It allows multiple nodes to communicate in a peer-to-peer fashion without requiring a central server.

### `conn`
The `conn` package provides a wrapper for reliable TCP connections. It ensures that critical messages are delivered reliably between nodes in the system.

### `localip`
The `localip` package retrieves the local IP address of the machine. This is used to identify the node on the network and facilitate communication.

### `peers`
The `peers` package handles peer discovery and management. It allows nodes to detect and maintain a list of active peers in the network.

---

## Source Attribution

The `bcast`, `conn`, `localip`, and `peers` packages are all taken from the [TTK4145 Real-Time Programming course repository](https://github.com/TTK4145/). These packages have been adapted and integrated into this project to meet the specific requirements of the elevator control system.

---

## Usage

To use the network module in your project, import the relevant packages and initialize them as needed. For example:

```go
import (
    "network/bcast"
    "network/peers"
)

// Example: Initialize a broadcaster
go bcast.Transmitter(15657, myChannel)

// Example: Initialize peer discovery
peerUpdateCh := make(chan peers.PeerUpdate)
go peers.Transmitter(15658, "elevator")
go peers.Receiver(15658, peerUpdateCh)