"""
SML Boolean - Boolean value handling

This module implements boolean value parsing and writing for SML protocol.
"""

from typing import Optional
from .buffer import SmlBuffer
from .shared import SML_TYPE_BOOLEAN, SML_OPTIONAL_SKIPPED

SML_BOOLEAN_TRUE = 0xFF
SML_BOOLEAN_FALSE = 0x00


def boolean_init(value: bool) -> int:
    """
    Initialize boolean value.
    
    Args:
        value: Boolean value
        
    Returns:
        SML_BOOLEAN_TRUE (0xFF) or SML_BOOLEAN_FALSE (0x00)
    """
    return SML_BOOLEAN_TRUE if value else SML_BOOLEAN_FALSE


def boolean_parse(buf: SmlBuffer) -> Optional[int]:
    """
    Parse boolean from buffer.
    
    Args:
        buf: The buffer to parse from
        
    Returns:
        SML_BOOLEAN_TRUE or SML_BOOLEAN_FALSE, or None if optional and skipped
    """
    if buf.optional_is_skipped() == 1:
        return None
    
    if buf.get_next_type() != SML_TYPE_BOOLEAN:
        buf.error = 1
        return None
    
    length = buf.get_next_length()
    if length != 1:
        buf.error = 1
        return None
    
    if buf.cursor >= buf.buffer_len:
        buf.error = 1
        return None
    
    value = buf.get_current_byte()
    buf.update_bytes_read(1)
    
    if value:
        return SML_BOOLEAN_TRUE
    else:
        return SML_BOOLEAN_FALSE


def boolean_write(boolean: Optional[int], buf: SmlBuffer) -> None:
    """
    Write boolean to buffer.
    
    Args:
        boolean: Boolean value (SML_BOOLEAN_TRUE or SML_BOOLEAN_FALSE), or None
        buf: The buffer to write to
    """
    if boolean is None:
        buf.optional_write()
        return
    
    buf.set_type_and_length(SML_TYPE_BOOLEAN, 1)
    buf.ensure_capacity(1)
    
    if boolean == SML_BOOLEAN_FALSE:
        buf.buffer[buf.cursor] = SML_BOOLEAN_FALSE
    else:
        buf.buffer[buf.cursor] = SML_BOOLEAN_TRUE
    
    buf.cursor += 1


def boolean_free(boolean: Optional[int]) -> None:
    """
    Free boolean (no-op in Python, kept for API compatibility).
    
    Args:
        boolean: The boolean value to free
    """
    pass

