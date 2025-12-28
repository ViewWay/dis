"""
SML Status - Status handling

This module implements status parsing and writing for SML protocol.
"""

from typing import Optional
from .buffer import SmlBuffer
from .number import (
    u8_parse, u16_parse, u32_parse, u64_parse,
    u8_write, u16_write, u32_write, u64_write,
)
from .shared import SML_TYPE_UNSIGNED, SML_TYPE_FIELD, SML_LENGTH_FIELD, SML_OPTIONAL_SKIPPED


class SmlStatus:
    """
    SML Status class.
    
    Attributes:
        type: Status type (includes type and size information)
        data: Status data as bytes
    """
    
    def __init__(self):
        """Initialize a new SmlStatus."""
        self.type: int = SML_TYPE_UNSIGNED
        self.data: Optional[bytes] = None
    
    @classmethod
    def parse(cls, buf: SmlBuffer) -> Optional['SmlStatus']:
        """
        Parse status from buffer.
        
        Args:
            buf: The buffer to parse from
            
        Returns:
            SmlStatus instance, or None if optional and skipped
        """
        if buf.optional_is_skipped() == 1:
            return None
        
        max_size = 1
        status_type = buf.get_next_type()
        byte = buf.get_current_byte()
        
        status = cls()
        status.type = status_type
        
        if status_type == SML_TYPE_UNSIGNED:
            # Get maximal size, if not all bytes are used (example: only 6 bytes for a u64)
            while max_size < ((byte & SML_LENGTH_FIELD) - 1):
                max_size <<= 1
            
            if max_size == 1:
                status.data = u8_parse(buf)
            elif max_size == 2:
                status.data = u16_parse(buf)
            elif max_size == 4:
                status.data = u32_parse(buf)
            elif max_size == 8:
                status.data = u64_parse(buf)
            else:
                buf.error = 1
            
            if status.data is not None:
                status.type |= max_size
        else:
            buf.error = 1
        
        if buf.has_errors():
            return None
        
        return status
    
    def write(self, buf: SmlBuffer) -> None:
        """
        Write status to buffer.
        
        Args:
            buf: The buffer to write to
        """
        if self.data is None:
            buf.optional_write()
            return
        
        status_type = self.type & SML_TYPE_FIELD
        size = self.type & SML_LENGTH_FIELD
        
        if status_type == SML_TYPE_UNSIGNED:
            if size == 1:
                u8_write(self.data, buf)
            elif size == 2:
                u16_write(self.data, buf)
            elif size == 4:
                u32_write(self.data, buf)
            elif size == 8:
                u64_write(self.data, buf)
    
    def free(self) -> None:
        """Free status (no-op in Python, kept for API compatibility)."""
        pass


def status_init() -> SmlStatus:
    """
    Initialize a new SmlStatus.
    
    Returns:
        New SmlStatus instance
    """
    return SmlStatus()


def status_parse(buf: SmlBuffer) -> Optional[SmlStatus]:
    """
    Parse status from buffer.
    
    Args:
        buf: The buffer to parse from
        
    Returns:
        SmlStatus instance, or None if optional and skipped
    """
    return SmlStatus.parse(buf)


def status_write(status: Optional[SmlStatus], buf: SmlBuffer) -> None:
    """
    Write status to buffer.
    
    Args:
        status: The status to write, or None
        buf: The buffer to write to
    """
    if status is None:
        buf.optional_write()
    else:
        status.write(buf)


def status_free(status: Optional[SmlStatus]) -> None:
    """
    Free status (no-op in Python).
    
    Args:
        status: The status to free
    """
    pass

