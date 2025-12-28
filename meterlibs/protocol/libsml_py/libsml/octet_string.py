"""
SML Octet String - Byte string handling

This module implements the OctetString class for handling byte strings in SML protocol.
"""

import uuid
from typing import Optional
from .buffer import SmlBuffer
from .shared import SML_TYPE_OCTET_STRING, SML_OPTIONAL_SKIPPED


class OctetString:
    """
    Octet string class for SML protocol.
    
    Attributes:
        str: The byte data
        len: Length of the byte data
    """
    
    def __init__(self, data: Optional[bytes] = None, length: int = 0):
        """
        Initialize an OctetString.
        
        Args:
            data: The byte data (if None, creates empty string)
            length: Length of data (if data is provided, uses len(data))
        """
        if data is not None:
            self.str: bytes = data
            self.len: int = len(data)
        else:
            self.str: bytes = b''
            self.len: int = length
    
    @classmethod
    def init_from_hex(cls, hex_str: str) -> 'OctetString':
        """
        Initialize from hex string.
        
        Args:
            hex_str: Hex string (e.g., "1A2B3C")
            
        Returns:
            New OctetString instance
            
        Raises:
            ValueError: If hex string length is not even
        """
        if len(hex_str) % 2 != 0:
            raise ValueError("Hex string length must be even")
        
        data = bytes.fromhex(hex_str)
        return cls(data, len(data))
    
    @classmethod
    def parse(cls, buf: SmlBuffer) -> Optional['OctetString']:
        """
        Parse octet string from buffer.
        
        Args:
            buf: The buffer to parse from
            
        Returns:
            OctetString instance, or None if optional and skipped
        """
        if buf.optional_is_skipped() == 1:
            return None
        
        if buf.get_next_type() != SML_TYPE_OCTET_STRING:
            buf.error = 1
            return None
        
        length = buf.get_next_length()
        if length < 0 or length >= (buf.buffer_len - buf.cursor):
            buf.error = 1
            return None
        
        # Extract bytes
        data = bytes(buf.buffer[buf.cursor:buf.cursor + length])
        buf.update_bytes_read(length)
        
        return cls(data, length)
    
    def write(self, buf: SmlBuffer) -> None:
        """
        Write octet string to buffer.
        
        Args:
            buf: The buffer to write to
        """
        if self.str is None or self.len == 0:
            buf.optional_write()
            return
        
        buf.set_type_and_length(SML_TYPE_OCTET_STRING, self.len)
        buf.ensure_capacity(self.len)
        
        for i in range(self.len):
            buf.buffer[buf.cursor + i] = self.str[i]
        
        buf.cursor += self.len
    
    def cmp(self, other: 'OctetString') -> int:
        """
        Compare two octet strings.
        
        Args:
            other: The other octet string to compare
            
        Returns:
            -1 if different lengths, 0 if equal, non-zero if different
        """
        if self.len != other.len:
            return -1
        if self.str == other.str:
            return 0
        return 1
    
    def cmp_with_hex(self, hex_str: str) -> int:
        """
        Compare with hex string.
        
        Args:
            hex_str: Hex string to compare with
            
        Returns:
            -1 if different lengths, 0 if equal, non-zero if different
        """
        try:
            other = self.init_from_hex(hex_str)
            return self.cmp(other)
        except ValueError:
            return -1
    
    def free(self) -> None:
        """Free octet string (no-op in Python, kept for API compatibility)."""
        pass
    
    def __repr__(self) -> str:
        return f"OctetString(len={self.len}, data={self.str.hex()})"


def octet_string_init(data: bytes, length: int) -> OctetString:
    """
    Initialize octet string from data.
    
    Args:
        data: The byte data
        length: Length (usually len(data))
        
    Returns:
        New OctetString instance
    """
    return OctetString(data, length)


def octet_string_init_from_hex(hex_str: str) -> OctetString:
    """
    Initialize octet string from hex string.
    
    Args:
        hex_str: Hex string
        
    Returns:
        New OctetString instance
    """
    return OctetString.init_from_hex(hex_str)


def octet_string_parse(buf: SmlBuffer) -> Optional[OctetString]:
    """
    Parse octet string from buffer.
    
    Args:
        buf: The buffer to parse from
        
    Returns:
        OctetString instance, or None if optional and skipped
    """
    return OctetString.parse(buf)


def octet_string_write(octet_str: Optional[OctetString], buf: SmlBuffer) -> None:
    """
    Write octet string to buffer.
    
    Args:
        octet_str: The octet string to write, or None
        buf: The buffer to write to
    """
    if octet_str is None:
        buf.optional_write()
    else:
        octet_str.write(buf)


def octet_string_free(octet_str: Optional[OctetString]) -> None:
    """
    Free octet string (no-op in Python).
    
    Args:
        octet_str: The octet string to free
    """
    pass


def octet_string_generate_uuid() -> OctetString:
    """
    Generate UUID as octet string.
    
    Returns:
        OctetString containing UUID bytes (16 bytes)
    """
    uuid_obj = uuid.uuid4()
    uuid_bytes = uuid_obj.bytes
    return OctetString(uuid_bytes, 16)


def octet_string_cmp(s1: OctetString, s2: OctetString) -> int:
    """
    Compare two octet strings.
    
    Args:
        s1: First octet string
        s2: Second octet string
        
    Returns:
        -1 if different lengths, 0 if equal, non-zero if different
    """
    return s1.cmp(s2)


def octet_string_cmp_with_hex(octet_str: OctetString, hex_str: str) -> int:
    """
    Compare octet string with hex string.
    
    Args:
        octet_str: The octet string
        hex_str: Hex string to compare with
        
    Returns:
        -1 if different lengths, 0 if equal, non-zero if different
    """
    return octet_str.cmp_with_hex(hex_str)


# Signature type (alias for OctetString)
SmlSignature = OctetString
sml_signature_parse = octet_string_parse
sml_signature_write = octet_string_write
sml_signature_free = octet_string_free

