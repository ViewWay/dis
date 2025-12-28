"""
SML Number - Number type handling

This module implements parsing and writing of integer types (signed and unsigned)
with proper byte order handling.
"""

import struct
from typing import Optional
from .buffer import SmlBuffer, SmlBufferError
from .shared import (
    SML_TYPE_UNSIGNED,
    SML_TYPE_INTEGER,
    SML_TYPE_NUMBER_8,
    SML_TYPE_NUMBER_16,
    SML_TYPE_NUMBER_32,
    SML_TYPE_NUMBER_64,
    SML_OPTIONAL_SKIPPED,
    u8, u16, u32, u64, i8, i16, i32, i64,
)


def _is_big_endian() -> bool:
    """Check if system is big-endian."""
    return struct.pack('@I', 1) == struct.pack('>I', 1)


def _number_init(number: int, number_type: int, size: int) -> bytes:
    """
    Initialize a number value as bytes in big-endian format.
    
    Args:
        number: The number value
        number_type: SML_TYPE_UNSIGNED or SML_TYPE_INTEGER
        size: Size in bytes (1, 2, 4, or 8)
        
    Returns:
        Bytes representation of the number
    """
    if size == 1:
        fmt = '>B' if number_type == SML_TYPE_UNSIGNED else '>b'
    elif size == 2:
        fmt = '>H' if number_type == SML_TYPE_UNSIGNED else '>h'
    elif size == 4:
        fmt = '>I' if number_type == SML_TYPE_UNSIGNED else '>i'
    elif size == 8:
        fmt = '>Q' if number_type == SML_TYPE_UNSIGNED else '>q'
    else:
        raise ValueError(f"Invalid size: {size}")
    
    return struct.pack(fmt, number)


def _number_parse(buf: SmlBuffer, number_type: int, max_size: int) -> Optional[bytes]:
    """
    Parse a number from buffer.
    
    Args:
        buf: The buffer to parse from
        number_type: SML_TYPE_UNSIGNED or SML_TYPE_INTEGER
        max_size: Maximum size in bytes
        
    Returns:
        Bytes representation of the number, or None if optional and skipped
    """
    if buf.optional_is_skipped() == 1:
        return None
    
    if (buf.cursor + 1) > buf.buffer_len:
        buf.error = 1
        return None
    
    if buf.get_next_type() != number_type:
        buf.error = 1
        return None
    
    length = buf.get_next_length()
    if length < 0 or length > max_size:
        buf.error = 1
        return None
    
    # Allocate result buffer
    result = bytearray(max_size)
    
    # Check if enough bytes available
    if (buf.cursor + length) > buf.buffer_len:
        buf.error = 1
        return None
    
    if (buf.cursor + 1) > buf.buffer_len:
        buf.error = 1
        return None
    
    current_byte = buf.get_current_byte()
    is_negative = (number_type == SML_TYPE_INTEGER and (current_byte & 0x80) != 0)
    
    missing_bytes = max_size - length
    # Copy the actual bytes
    for i in range(length):
        result[missing_bytes + i] = buf.buffer[buf.cursor + i]
    
    # Sign extension for negative integers
    if is_negative:
        for i in range(missing_bytes):
            result[i] = 0xFF
    
    # Convert to big-endian if system is little-endian
    if not _is_big_endian():
        result.reverse()
    
    buf.update_bytes_read(length)
    return bytes(result)


def _number_write(number_bytes: Optional[bytes], number_type: int, size: int, buf: SmlBuffer) -> None:
    """
    Write a number to buffer.
    
    Args:
        number_bytes: The number as bytes, or None for optional skipped
        number_type: SML_TYPE_UNSIGNED or SML_TYPE_INTEGER
        size: Size in bytes
        buf: The buffer to write to
    """
    if number_bytes is None:
        buf.optional_write()
        return
    
    buf.set_type_and_length(number_type, size)
    
    # Ensure buffer capacity
    buf.ensure_capacity(size)
    
    # Convert from big-endian if system is little-endian
    if not _is_big_endian():
        number_bytes = bytes(reversed(number_bytes))
    
    # Write bytes
    for i in range(size):
        buf.buffer[buf.cursor + i] = number_bytes[i]
    
    buf.update_bytes_read(size)


# Unsigned integer functions
def u8_init(n: int) -> bytes:
    """Initialize u8."""
    return _number_init(n, SML_TYPE_UNSIGNED, SML_TYPE_NUMBER_8)


def u16_init(n: int) -> bytes:
    """Initialize u16."""
    return _number_init(n, SML_TYPE_UNSIGNED, SML_TYPE_NUMBER_16)


def u32_init(n: int) -> bytes:
    """Initialize u32."""
    return _number_init(n, SML_TYPE_UNSIGNED, SML_TYPE_NUMBER_32)


def u64_init(n: int) -> bytes:
    """Initialize u64."""
    return _number_init(n, SML_TYPE_UNSIGNED, SML_TYPE_NUMBER_64)


def u8_parse(buf: SmlBuffer) -> Optional[bytes]:
    """Parse u8."""
    return _number_parse(buf, SML_TYPE_UNSIGNED, SML_TYPE_NUMBER_8)


def u16_parse(buf: SmlBuffer) -> Optional[bytes]:
    """Parse u16."""
    return _number_parse(buf, SML_TYPE_UNSIGNED, SML_TYPE_NUMBER_16)


def u32_parse(buf: SmlBuffer) -> Optional[bytes]:
    """Parse u32."""
    return _number_parse(buf, SML_TYPE_UNSIGNED, SML_TYPE_NUMBER_32)


def u64_parse(buf: SmlBuffer) -> Optional[bytes]:
    """Parse u64."""
    return _number_parse(buf, SML_TYPE_UNSIGNED, SML_TYPE_NUMBER_64)


def u8_write(n: Optional[bytes], buf: SmlBuffer) -> None:
    """Write u8."""
    _number_write(n, SML_TYPE_UNSIGNED, SML_TYPE_NUMBER_8, buf)


def u16_write(n: Optional[bytes], buf: SmlBuffer) -> None:
    """Write u16."""
    _number_write(n, SML_TYPE_UNSIGNED, SML_TYPE_NUMBER_16, buf)


def u32_write(n: Optional[bytes], buf: SmlBuffer) -> None:
    """Write u32."""
    _number_write(n, SML_TYPE_UNSIGNED, SML_TYPE_NUMBER_32, buf)


def u64_write(n: Optional[bytes], buf: SmlBuffer) -> None:
    """Write u64."""
    _number_write(n, SML_TYPE_UNSIGNED, SML_TYPE_NUMBER_64, buf)


# Signed integer functions
def i8_init(n: int) -> bytes:
    """Initialize i8."""
    return _number_init(n, SML_TYPE_INTEGER, SML_TYPE_NUMBER_8)


def i16_init(n: int) -> bytes:
    """Initialize i16."""
    return _number_init(n, SML_TYPE_INTEGER, SML_TYPE_NUMBER_16)


def i32_init(n: int) -> bytes:
    """Initialize i32."""
    return _number_init(n, SML_TYPE_INTEGER, SML_TYPE_NUMBER_32)


def i64_init(n: int) -> bytes:
    """Initialize i64."""
    return _number_init(n, SML_TYPE_INTEGER, SML_TYPE_NUMBER_64)


def i8_parse(buf: SmlBuffer) -> Optional[bytes]:
    """Parse i8."""
    return _number_parse(buf, SML_TYPE_INTEGER, SML_TYPE_NUMBER_8)


def i16_parse(buf: SmlBuffer) -> Optional[bytes]:
    """Parse i16."""
    return _number_parse(buf, SML_TYPE_INTEGER, SML_TYPE_NUMBER_16)


def i32_parse(buf: SmlBuffer) -> Optional[bytes]:
    """Parse i32."""
    return _number_parse(buf, SML_TYPE_INTEGER, SML_TYPE_NUMBER_32)


def i64_parse(buf: SmlBuffer) -> Optional[bytes]:
    """Parse i64."""
    return _number_parse(buf, SML_TYPE_INTEGER, SML_TYPE_NUMBER_64)


def i8_write(n: Optional[bytes], buf: SmlBuffer) -> None:
    """Write i8."""
    _number_write(n, SML_TYPE_INTEGER, SML_TYPE_NUMBER_8, buf)


def i16_write(n: Optional[bytes], buf: SmlBuffer) -> None:
    """Write i16."""
    _number_write(n, SML_TYPE_INTEGER, SML_TYPE_NUMBER_16, buf)


def i32_write(n: Optional[bytes], buf: SmlBuffer) -> None:
    """Write i32."""
    _number_write(n, SML_TYPE_INTEGER, SML_TYPE_NUMBER_32, buf)


def i64_write(n: Optional[bytes], buf: SmlBuffer) -> None:
    """Write i64."""
    _number_write(n, SML_TYPE_INTEGER, SML_TYPE_NUMBER_64, buf)


# Unit type (alias for u8)
sml_unit = u8
sml_unit_init = u8_init
sml_unit_parse = u8_parse
sml_unit_write = u8_write
sml_unit_free = lambda x: None  # No-op in Python


# Free functions (no-op in Python, kept for API compatibility)
def number_free(np: Optional[bytes]) -> None:
    """Free number (no-op in Python)."""
    pass


u8_free = number_free
u16_free = number_free
u32_free = number_free
u64_free = number_free
i8_free = number_free
i16_free = number_free
i32_free = number_free
i64_free = number_free

