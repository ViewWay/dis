"""
SML Transport - Transport layer handling

This module implements SML transport protocol for reading and writing
SML files over file-like objects, handling escape sequences.
"""

from typing import Optional, Callable
from .file import SmlFile
from .crc16 import crc16_calculate
from .shared import SML_MESSAGE_END

# Transport protocol constants
ESC_SEQ = bytes([0x1b, 0x1b, 0x1b, 0x1b])
START_SEQ = bytes([0x1b, 0x1b, 0x1b, 0x1b, 0x01, 0x01, 0x01, 0x01])
END_SEQ = bytes([0x1b, 0x1b, 0x1b, 0x1b, 0x1a])

MC_SML_BUFFER_LEN = 8096


def transport_read(file_obj, max_len: int) -> Optional[bytes]:
    """
    Read continuously from file object and scan for SML transport protocol
    escape sequences. If a SML file is detected, it will be copied into
    the buffer.
    
    Args:
        file_obj: File-like object to read from
        max_len: Maximum length of buffer
        
    Returns:
        Bytes containing the SML file with transport protocol, or None on error
    """
    if max_len < 8:
        print("libsml: error: sml_transport_read(): passed buffer too small!")
        return None
    
    buf = bytearray(max_len)
    length = 0
    
    # Wait for start sequence
    while length < 8:
        try:
            byte = file_obj.read(1)
            if not byte:
                return None  # EOF
            
            byte_val = byte[0]
            if (byte_val == 0x1b and length < 4) or (byte_val == 0x01 and length >= 4):
                buf[length] = byte_val
                length += 1
            else:
                length = 0
        except (IOError, OSError) as e:
            print(f"libsml: sml_read(): read error: {e}")
            return None
    
    # Found start sequence, read until end sequence
    while (length + 8) < max_len:
        try:
            chunk = file_obj.read(4)
            if len(chunk) < 4:
                return None
            
            if chunk == ESC_SEQ:
                # Found escape sequence
                for i in range(4):
                    buf[length + i] = chunk[i]
                length += 4
                
                chunk = file_obj.read(4)
                if len(chunk) < 4:
                    return None
                
                if chunk[0] == 0x1a:
                    # Found end sequence
                    for i in range(4):
                        buf[length + i] = chunk[i]
                    length += 4
                    return bytes(buf[:length])
                else:
                    # Don't read other escaped sequences yet
                    print("libsml: error: unrecognized sequence")
                    return None
            else:
                for i in range(4):
                    buf[length + i] = chunk[i]
                length += 4
        except (IOError, OSError) as e:
            print(f"libsml: sml_read(): read error: {e}")
            return None
    
    return None


def transport_listen(file_obj, receiver: Callable[[bytes, int], None]) -> None:
    """
    Endless loop which reads continuously via transport_read and calls
    the receiver function.
    
    Args:
        file_obj: File-like object to read from
        receiver: Callback function that receives (buffer, buffer_len)
    """
    buffer = bytearray(MC_SML_BUFFER_LEN)
    
    while True:
        bytes_read = transport_read(file_obj, MC_SML_BUFFER_LEN)
        if bytes_read:
            receiver(bytes_read, len(bytes_read))
        else:
            break


def transport_write(file_obj, sml_file: SmlFile) -> int:
    """
    Add the SML transport protocol escape sequences and write the given
    file to file object. The file must be in the parsed format.
    
    Args:
        file_obj: File-like object to write to
        sml_file: The SML file to write
        
    Returns:
        Number of bytes written, 0 if there was an error
    """
    if sml_file.buf is None:
        sml_file.buf = SmlBuffer(SML_FILE_BUFFER_LENGTH)
    
    sml_file.buf.cursor = 0
    
    # Add start sequence
    sml_file.buf.ensure_capacity(8)
    for i in range(8):
        sml_file.buf.buffer[sml_file.buf.cursor + i] = START_SEQ[i]
    sml_file.buf.cursor += 8
    
    # Add file
    sml_file.write()
    
    # Add padding bytes
    padding = (sml_file.buf.cursor % 4) if (sml_file.buf.cursor % 4) else 0
    if padding:
        padding = 4 - padding
        sml_file.buf.ensure_capacity(padding)
        for i in range(padding):
            sml_file.buf.buffer[sml_file.buf.cursor + i] = 0
        sml_file.buf.cursor += padding
    
    # Begin end sequence
    sml_file.buf.ensure_capacity(5)
    for i in range(5):
        sml_file.buf.buffer[sml_file.buf.cursor + i] = END_SEQ[i]
    sml_file.buf.cursor += 5
    
    # Add padding info
    sml_file.buf.ensure_capacity(1)
    sml_file.buf.buffer[sml_file.buf.cursor] = padding
    sml_file.buf.cursor += 1
    
    # Add CRC checksum
    crc = crc16_calculate(sml_file.buf.buffer, sml_file.buf.cursor)
    sml_file.buf.ensure_capacity(2)
    sml_file.buf.buffer[sml_file.buf.cursor] = (crc & 0xFF00) >> 8
    sml_file.buf.cursor += 1
    sml_file.buf.buffer[sml_file.buf.cursor] = crc & 0x00FF
    sml_file.buf.cursor += 1
    
    try:
        written = file_obj.write(sml_file.buf.buffer[:sml_file.buf.cursor])
        if written == sml_file.buf.cursor:
            return written
        return 0
    except (IOError, OSError) as e:
        print(f"libsml: sml_transport_write(): write error: {e}")
        return 0

