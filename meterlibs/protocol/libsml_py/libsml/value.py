"""
SML Value - Value type handling

This module implements the SmlValue class which can hold various data types
(boolean, bytes, signed/unsigned integers of various sizes).
"""

import struct
from typing import Optional, Union
from .buffer import SmlBuffer
from .octet_string import OctetString
from .boolean import boolean_parse, boolean_write, SML_BOOLEAN_TRUE, SML_BOOLEAN_FALSE
from .number import (
    u8_parse, u16_parse, u32_parse, u64_parse,
    i8_parse, i16_parse, i32_parse, i64_parse,
    u8_write, u16_write, u32_write, u64_write,
    i8_write, i16_write, i32_write, i64_write,
)
from .shared import (
    SML_TYPE_OCTET_STRING,
    SML_TYPE_BOOLEAN,
    SML_TYPE_INTEGER,
    SML_TYPE_UNSIGNED,
    SML_TYPE_FIELD,
    SML_LENGTH_FIELD,
    SML_OPTIONAL_SKIPPED,
)


class SmlValue:
    """
    SML Value class that can hold various data types.
    
    Attributes:
        type: The type of the value (includes type and size information)
        data: The actual data (can be boolean, OctetString, or bytes for numbers)
    """
    
    def __init__(self):
        """Initialize a new SmlValue."""
        self.type: int = SML_TYPE_OCTET_STRING
        self.data: Union[Optional[int], Optional[OctetString], Optional[bytes]] = None
    
    @classmethod
    def parse(cls, buf: SmlBuffer) -> Optional['SmlValue']:
        """
        Parse value from buffer.
        
        Args:
            buf: The buffer to parse from
            
        Returns:
            SmlValue instance, or None if optional and skipped
        """
        if buf.optional_is_skipped() == 1:
            return None
        
        max_size = 1
        value_type = buf.get_next_type()
        byte = buf.get_current_byte()
        
        value = cls()
        value.type = value_type
        
        if value_type == SML_TYPE_OCTET_STRING:
            value.data = OctetString.parse(buf)
        elif value_type == SML_TYPE_BOOLEAN:
            value.data = boolean_parse(buf)
        elif value_type == SML_TYPE_UNSIGNED or value_type == SML_TYPE_INTEGER:
            # Get maximal size, if not all bytes are used (example: only 6 bytes for a u64)
            while max_size < ((byte & SML_LENGTH_FIELD) - 1):
                max_size <<= 1
            
            if value_type == SML_TYPE_UNSIGNED:
                if max_size == 1:
                    value.data = u8_parse(buf)
                elif max_size == 2:
                    value.data = u16_parse(buf)
                elif max_size == 4:
                    value.data = u32_parse(buf)
                elif max_size == 8:
                    value.data = u64_parse(buf)
                else:
                    buf.error = 1
            else:  # SML_TYPE_INTEGER
                if max_size == 1:
                    value.data = i8_parse(buf)
                elif max_size == 2:
                    value.data = i16_parse(buf)
                elif max_size == 4:
                    value.data = i32_parse(buf)
                elif max_size == 8:
                    value.data = i64_parse(buf)
                else:
                    buf.error = 1
            
            if value.data is not None:
                value.type |= max_size
        else:
            buf.error = 1
        
        if buf.has_errors():
            return None
        
        return value
    
    def write(self, buf: SmlBuffer) -> None:
        """
        Write value to buffer.
        
        Args:
            buf: The buffer to write to
        """
        if self.data is None:
            buf.optional_write()
            return
        
        value_type = self.type & SML_TYPE_FIELD
        
        if value_type == SML_TYPE_OCTET_STRING:
            if isinstance(self.data, OctetString):
                self.data.write(buf)
        elif value_type == SML_TYPE_BOOLEAN:
            boolean_write(self.data, buf)
        elif value_type == SML_TYPE_UNSIGNED or value_type == SML_TYPE_INTEGER:
            size = self.type & SML_LENGTH_FIELD
            if value_type == SML_TYPE_UNSIGNED:
                if size == 1:
                    u8_write(self.data, buf)
                elif size == 2:
                    u16_write(self.data, buf)
                elif size == 4:
                    u32_write(self.data, buf)
                elif size == 8:
                    u64_write(self.data, buf)
            else:  # SML_TYPE_INTEGER
                if size == 1:
                    i8_write(self.data, buf)
                elif size == 2:
                    i16_write(self.data, buf)
                elif size == 4:
                    i32_write(self.data, buf)
                elif size == 8:
                    i64_write(self.data, buf)
    
    def to_double(self) -> float:
        """
        Convert value to double.
        
        Returns:
            The value as a float
        """
        if self.data is None:
            return 0.0
        
        # Type values: 0x51=i8, 0x52=i16, 0x54=i32, 0x58=i64,
        #               0x61=u8, 0x62=u16, 0x64=u32, 0x68=u64
        if self.type == 0x51:  # i8
            return struct.unpack('>b', self.data)[0]
        elif self.type == 0x52:  # i16
            return struct.unpack('>h', self.data)[0]
        elif self.type == 0x54:  # i32
            return struct.unpack('>i', self.data)[0]
        elif self.type == 0x58:  # i64
            return struct.unpack('>q', self.data)[0]
        elif self.type == 0x61:  # u8
            return struct.unpack('>B', self.data)[0]
        elif self.type == 0x62:  # u16
            return struct.unpack('>H', self.data)[0]
        elif self.type == 0x64:  # u32
            return struct.unpack('>I', self.data)[0]
        elif self.type == 0x68:  # u64
            return struct.unpack('>Q', self.data)[0]
        else:
            return 0.0
    
    def to_strhex(self, mixed: bool = False) -> Optional[str]:
        """
        Convert SML octet string to a printable hex string.
        
        Args:
            mixed: If True, print printable characters as-is
            
        Returns:
            Hex string representation, or None if not an octet string
        """
        if self.type != SML_TYPE_OCTET_STRING:
            return None
        
        if not isinstance(self.data, OctetString) or self.data.str is None:
            return None
        
        hex_str = "0123456789abcdef"
        length = self.data.len
        data_str = self.data.str
        
        result = []
        for i in range(length):
            if mixed and (0x20 < data_str[i] < 0x7B):
                result.append(chr(data_str[i]))
            else:
                mixed = False
                result.append(hex_str[(data_str[i] >> 4) & 0x0F])
                result.append(hex_str[data_str[i] & 0x0F])
                result.append(' ')
        
        return ''.join(result).rstrip()
    
    def free(self) -> None:
        """Free value (no-op in Python, kept for API compatibility)."""
        pass


def value_init() -> SmlValue:
    """
    Initialize a new SmlValue.
    
    Returns:
        New SmlValue instance
    """
    return SmlValue()


def value_parse(buf: SmlBuffer) -> Optional[SmlValue]:
    """
    Parse value from buffer.
    
    Args:
        buf: The buffer to parse from
        
    Returns:
        SmlValue instance, or None if optional and skipped
    """
    return SmlValue.parse(buf)


def value_write(value: Optional[SmlValue], buf: SmlBuffer) -> None:
    """
    Write value to buffer.
    
    Args:
        value: The value to write, or None
        buf: The buffer to write to
    """
    if value is None:
        buf.optional_write()
    else:
        value.write(buf)


def value_free(value: Optional[SmlValue]) -> None:
    """
    Free value (no-op in Python).
    
    Args:
        value: The value to free
    """
    pass


def value_to_double(value: SmlValue) -> float:
    """
    Convert value to double.
    
    Args:
        value: The value to convert
        
    Returns:
        The value as a float
    """
    return value.to_double()


def value_to_strhex(value: SmlValue, mixed: bool = False) -> Optional[str]:
    """
    Convert SML octet string to a printable hex string.
    
    Args:
        value: The value to convert
        mixed: If True, print printable characters as-is
        
    Returns:
        Hex string representation, or None if not an octet string
    """
    return value.to_strhex(mixed)

