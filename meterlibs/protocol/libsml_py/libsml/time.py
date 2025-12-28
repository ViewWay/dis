"""
SML Time - Time handling

This module implements time parsing and writing for SML protocol,
including compatibility workarounds for specific meter types.
"""

from typing import Optional
from .buffer import SmlBuffer
from .number import u8_parse, u32_parse, u8_write, u32_write, u8_init
from .shared import SML_TYPE_LIST, SML_TYPE_UNSIGNED, SML_OPTIONAL_SKIPPED

SML_TIME_SEC_INDEX = 0x01
SML_TIME_TIMESTAMP = 0x02


class SmlTime:
    """
    SML Time class.
    
    Attributes:
        tag: Time type (SML_TIME_SEC_INDEX or SML_TIME_TIMESTAMP)
        data: Time data (sec_index or timestamp as bytes)
    """
    
    def __init__(self):
        """Initialize a new SmlTime."""
        self.tag: Optional[bytes] = None
        self.data: Optional[bytes] = None  # sec_index or timestamp
    
    @classmethod
    def parse(cls, buf: SmlBuffer) -> Optional['SmlTime']:
        """
        Parse time from buffer.
        
        Args:
            buf: The buffer to parse from
            
        Returns:
            SmlTime instance, or None if optional and skipped
        """
        if buf.optional_is_skipped() == 1:
            return None
        
        time_obj = cls()
        
        # Workaround for Holley DTZ541
        # If SML_ListEntry valTime (SML_Time) is given there are missing bytes:
        # 0x72: indicate a list for SML_Time with 2 entries
        # 0x62 0x01: indicate secIndex
        # Instead, the DTZ541 starts with 0x65 + 4 bytes secIndex
        # The workaround will add this information during parsing
        if buf.cursor < buf.buffer_len and buf.get_current_byte() == (SML_TYPE_UNSIGNED | 5):
            time_obj.tag = u8_init(SML_TIME_SEC_INDEX)
        else:
            if buf.get_next_type() != SML_TYPE_LIST:
                buf.error = 1
                return None
            
            if buf.get_next_length() != 2:
                buf.error = 1
                return None
            
            time_obj.tag = u8_parse(buf)
            if buf.has_errors() or time_obj.tag is None:
                return None
        
        value_type = buf.get_next_type()
        if value_type == SML_TYPE_UNSIGNED:
            time_obj.data = u32_parse(buf)
            if buf.has_errors():
                return None
        elif value_type == SML_TYPE_LIST:
            # Some meters (e.g. FROETEC Multiflex ZG22) giving not one uint32
            # as timestamp, but a list of 3 values.
            # Ignoring these values, so that parsing does not fail.
            buf.get_next_length()  # Should we check the length here?
            t1 = u32_parse(buf)
            if buf.has_errors():
                return None
            from .number import i16_parse
            t2 = i16_parse(buf)
            if buf.has_errors():
                return None
            t3 = i16_parse(buf)
            if buf.has_errors():
                return None
            # Log warning (in Python we can use print or logging)
            print(f"libsml: error: sml_time as list[3]: ignoring value[0]={t1} value[1]={t2} value[2]={t3}")
        else:
            return None
        
        return time_obj
    
    def write(self, buf: SmlBuffer) -> None:
        """
        Write time to buffer.
        
        Args:
            buf: The buffer to write to
        """
        if self.tag is None or self.data is None:
            buf.optional_write()
            return
        
        buf.set_type_and_length(SML_TYPE_LIST, 2)
        u8_write(self.tag, buf)
        u32_write(self.data, buf)
    
    def free(self) -> None:
        """Free time (no-op in Python, kept for API compatibility)."""
        pass


def time_init() -> SmlTime:
    """
    Initialize a new SmlTime.
    
    Returns:
        New SmlTime instance
    """
    return SmlTime()


def time_parse(buf: SmlBuffer) -> Optional[SmlTime]:
    """
    Parse time from buffer.
    
    Args:
        buf: The buffer to parse from
        
    Returns:
        SmlTime instance, or None if optional and skipped
    """
    return SmlTime.parse(buf)


def time_write(time_obj: Optional[SmlTime], buf: SmlBuffer) -> None:
    """
    Write time to buffer.
    
    Args:
        time_obj: The time to write, or None
        buf: The buffer to write to
    """
    if time_obj is None:
        buf.optional_write()
    else:
        time_obj.write(buf)


def time_free(time_obj: Optional[SmlTime]) -> None:
    """
    Free time (no-op in Python).
    
    Args:
        time_obj: The time to free
    """
    pass

