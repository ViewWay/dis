# libSML Python Implementation

Python implementation of the Smart Message Language (SML) protocol specified by VDE's Forum Netztechnik/Netzbetrieb (FNN).

This is a port of the original C library [libSML](https://github.com/volkszaehler/libsml) to Python, maintaining binary protocol compatibility.

## Features

- Full SML protocol implementation
- Binary compatibility with C library
- Support for all SML message types
- Pythonic API while maintaining similarity to C API

## Installation

```bash
cd libsml_py
pip install -e .
```

## Usage

### Basic Usage

```python
import libsml

# Parse SML file from bytes
data = b'...'  # SML binary data
sml_file = libsml.file_parse(data, len(data))

# Access messages
for message in sml_file.messages:
    if message.message_body:
        tag_value = int.from_bytes(message.message_body.tag, 'big')
        print(f"Message type: {tag_value:04X}")
        
        # Access message data based on type
        if tag_value == libsml.SML_MESSAGE_GET_LIST_RESPONSE:
            response = message.message_body.data
            if response and response.val_list:
                # Process list entries
                current = response.val_list
                while current:
                    if current.value:
                        value = current.value.to_double()
                        print(f"Value: {value}")
                    current = current.next
```

### Creating Messages

```python
import libsml

# Create a new message
msg = libsml.message_init()

# Set message body
from libsml.messages import get_list_request
request = get_list_request.SmlGetListRequest()
request.client_id = libsml.OctetString(b"client123", 9)
request.server_id = libsml.OctetString(b"server456", 9)

msg.message_body = libsml.message_body_init(
    libsml.SML_MESSAGE_GET_LIST_REQUEST,
    request
)

# Write message to buffer
buf = libsml.SmlBuffer(512)
msg.write(buf)
output = buf.to_bytes()
```

### Using Transport Layer

```python
import libsml

# Read from file-like object
with open('sml_data.bin', 'rb') as f:
    data = libsml.transport_read(f, 8096)
    if data:
        # Remove transport protocol wrapper
        # (data contains escape sequences)
        sml_file = libsml.file_parse(data, len(data))

# Write to file-like object
sml_file = libsml.file_init()
# ... add messages ...
with open('output.bin', 'wb') as f:
    libsml.transport_write(f, sml_file)
```

## API Documentation

### Core Classes

- `SmlFile`: Container for multiple SML messages
- `SmlMessage`: Individual SML message with transaction ID, body, and CRC
- `SmlMessageBody`: Message body containing tag and data
- `SmlBuffer`: Buffer for parsing and writing SML data

### Data Types

- `OctetString`: Byte string type
- `SmlValue`: Value type supporting multiple data types
- `SmlList`: Linked list of value entries
- `SmlTime`: Time representation
- `SmlStatus`: Status information
- `SmlTree`: Tree structure for parameters

### Message Types

All message types are in `libsml.messages`:
- `open_request`, `open_response`
- `close_request`, `close_response`
- `get_list_request`, `get_list_response`
- `get_profile_pack_request`, `get_profile_pack_response`
- `get_profile_list_request`, `get_profile_list_response`
- `get_proc_parameter_request`, `get_proc_parameter_response`
- `set_proc_parameter_request`
- `attention_response`

## License

Copyright 2011 Juri Glass, Mathias Runge, Nadim El Sayed - DAI-Labor, TU-Berlin
Copyright 2014-2018 libSML contributors

libSML is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

See the file LICENSE for the full license text.
