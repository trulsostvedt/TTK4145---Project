# Network Module for Go (UDP Broadcast and Peer-to-Peer Communication)

This module provides networking utilities for the elevator control system, enabling UDP broadcasting, peer discovery, and local IP retrieval. The module has been updated and optimized for the final version of the project.

---

## Usage

To include this module in your project, add the following lines to your `go.mod` file:

```go
require Network-go v0.0.0
replace Network-go => ./Network-go
```

Where `./Network-go` is the relative path to this folder.

---

## Features

### Broadcasting and Receiving Messages
- **Channel-based Communication**: The module supports channel-in/channel-out pairs for custom or built-in data types. 
- **Serialization and Deserialization**: Data sent to the transmitter is serialized and broadcast on the specified port. Received messages are deserialized and sent to the corresponding channel.
- **Usage**: See [bcast.Transmitter and bcast.Receiver](network/bcast/bcast.go).

### Peer Discovery
- **Dynamic Peer Management**: Detect peers on the local network by supplying your own ID to a transmitter. Receive updates about new, current, and lost peers.
- **Usage**: See [peers.Transmitter and peers.Receiver](network/peers/peers.go).

### Local IP Retrieval
- **Convenience Function**: Retrieve the local IP address of the machine using the [LocalIP](network/localip/localip.go) function. Note: This requires an active internet connection.

---

## Recent Updates

- **Improved Serialization**: Enhanced support for custom data types, ensuring compatibility with the elevator control system's requirements.
- **Optimized Peer Discovery**: Reduced latency in detecting new peers and handling lost connections.
- **Error Handling**: Added robust error handling for network interruptions and invalid data.
- **Documentation**: Updated inline comments and examples for better clarity.

---

## Limitations

- The `LocalIP` function may not work in offline environments.
- Peer discovery is limited to devices on the same local network.

---

## License

This module is distributed under the MIT License. See the `LICENSE` file for more details.

---
