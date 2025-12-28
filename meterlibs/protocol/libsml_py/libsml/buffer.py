"""
SML Buffer - Buffer operations for parsing and writing SML data

This module implements the SmlBuffer class which handles reading from and writing
to byte buffers for SML protocol parsing and serialization.
"""

from typing import Optional
from .shared import (
    SML_TYPE_FIELD,
    SML_LENGTH_FIELD,
    SML_ANOTHER_TL,
    SML_TYPE_LIST,
    SML_OPTIONAL_SKIPPED,
    SML_MESSAGE_END,
)


class SmlBufferError(Exception):
    """Exception raised for buffer-related errors"""
    pass


class SmlBuffer:
    """
    Buffer class for SML parsing and writing.
    
    Used in two different use-cases:
    - Parsing: raw data is in the buffer, buffer_len is the number of raw bytes,
      cursor points to the current position during parsing
    - Writing: buffer is allocated with a default length (buffer_len),
      cursor points to the position where one can write
    """
    
    def __init__(self, length: int = 512):
        """
        Initialize a new SmlBuffer.
        
        Args:
            length: Initial buffer length (for writing) or will be set from data (for parsing)
        """
        self.buffer: bytearray = bytearray(length)
        self.buffer_len: int = length
        self.cursor: int = 0
        self.error: int = 0
        self.error_msg: Optional[str] = None
    
    @classmethod
    def from_bytes(cls, data: bytes) -> 'SmlBuffer':
        """
        Create a buffer from existing bytes (for parsing).
        
        Args:
            data: The bytes to parse
            
        Returns:
            A new SmlBuffer instance
        """
        buf = cls(len(data))
        buf.buffer = bytearray(data)
        buf.buffer_len = len(data)
        buf.cursor = 0
        return buf
    
    def has_errors(self) -> bool:
        """Check if an error has occurred."""
        return self.error != 0
    
    def get_next_type(self) -> int:
        """
        Returns the type field of the current byte.
        
        Returns:
            The type field value, or 0x100 if invalid
        """
        if self.cursor >= self.buffer_len:
            self.error = 1
            return 0x100  # invalid type
        return self.buffer[self.cursor] & SML_TYPE_FIELD
    
    def get_current_byte(self) -> int:
        """Returns the current byte."""
        if self.cursor >= self.buffer_len:
            raise SmlBufferError("Buffer cursor out of bounds")
        return self.buffer[self.cursor]
    
    def get_current_buf(self) -> bytearray:
        """Returns a reference to the current buffer position."""
        return self.buffer[self.cursor:]
    
    def update_bytes_read(self, bytes_count: int) -> None:
        """Sets the number of bytes read (moves the cursor forward)."""
        self.cursor += bytes_count
    
    def get_next_length(self) -> int:
        """
        Returns the length of the following data structure.
        Sets the cursor position to the value field.
        
        Returns:
            The length value, or -1 on error
        """
        length = 0
        
        # Check if current byte is available
        if (self.cursor + 1) > self.buffer_len:
            self.error = 1
            return -1
        
        byte = self.get_current_byte()
        is_list = ((byte & SML_TYPE_FIELD) == SML_TYPE_LIST)
        list_offset = 0 if is_list else -1
        
        # Read length field (may span multiple bytes)
        while self.cursor < self.buffer_len:
            byte = self.get_current_byte()
            length <<= 4
            length |= (byte & SML_LENGTH_FIELD)
            
            if (byte & SML_ANOTHER_TL) != SML_ANOTHER_TL:
                break
            self.update_bytes_read(1)
            if not is_list:
                list_offset -= 1
        
        if self.cursor < self.buffer_len:
            self.update_bytes_read(1)
        else:
            self.error = 1
            return -1
        
        return length + list_offset
    
    def set_type_and_length(self, type_val: int, length: int) -> None:
        """
        Sets the type and length field at the current cursor position.
        
        Args:
            type_val: The type value
            length: The length value
        """
        # Set the type
        self.buffer[self.cursor] = type_val
        
        if type_val != SML_TYPE_LIST:
            length += 1
        
        if length > SML_LENGTH_FIELD:
            # Calculate how many TL bytes are necessary
            mask_pos = (4 * 2) - 1  # sizeof(unsigned int) * 2 - 1
            
            # The 4 most significant bits of length
            mask = 0xF0 << (8 * (4 - 1))  # sizeof(unsigned int) - 1
            
            # Select the next 4 most significant bits with a bit set until there is something
            while not (mask & length):
                mask >>= 4
                mask_pos -= 1
            
            length += mask_pos  # for every TL-field
            
            if (0x0F << (4 * (mask_pos + 1))) & length:
                # For the rare case that the addition of the number of TL-fields
                # results in another TL-field
                mask_pos += 1
                length += 1
            
            # Copy 4 bits of the number to the buffer
            while mask > SML_LENGTH_FIELD:
                self.buffer[self.cursor] |= SML_ANOTHER_TL
                self.buffer[self.cursor] |= ((mask & length) >> (4 * mask_pos))
                mask >>= 4
                mask_pos -= 1
                self.cursor += 1
        
        self.buffer[self.cursor] |= (length & SML_LENGTH_FIELD)
        self.cursor += 1
    
    def optional_write(self) -> None:
        """Writes an optional skipped marker."""
        self.buffer[self.cursor] = SML_OPTIONAL_SKIPPED
        self.cursor += 1
    
    def optional_is_skipped(self) -> int:
        """
        Checks if the next field is a skipped optional field,
        updates the buffer accordingly.
        
        Returns:
            1 if skipped, 0 if not skipped, -1 on error
        """
        if self.cursor >= self.buffer_len:
            self.error = 1
            return -1
        
        if self.get_current_byte() == SML_OPTIONAL_SKIPPED:
            if (self.cursor + 1) > self.buffer_len:
                self.error = 1
                return -1
            self.update_bytes_read(1)
            return 1
        
        return 0
    
    def ensure_capacity(self, additional_bytes: int) -> None:
        """Ensure buffer has enough capacity for writing."""
        if self.cursor + additional_bytes > self.buffer_len:
            # Expand buffer
            new_size = max(self.buffer_len * 2, self.cursor + additional_bytes)
            self.buffer.extend(bytearray(new_size - self.buffer_len))
            self.buffer_len = new_size
    
    def to_bytes(self) -> bytes:
        """Returns the buffer content as bytes (up to cursor position)."""
        return bytes(self.buffer[:self.cursor])
    
    def free(self) -> None:
        """Free the buffer (Python handles this automatically, but kept for API compatibility)."""
        pass


def hexdump(buffer: bytes, buffer_len: Optional[int] = None) -> None:
    """
    Prints byte string to stdout in hex format.
    
    Args:
        buffer: The bytes to print
        buffer_len: Length to print (defaults to len(buffer))
    """
    if buffer_len is None:
        buffer_len = len(buffer)
    
    for i in range(buffer_len):
        print(f"{buffer[i]:02X} ", end="")
        if (i + 1) % 8 == 0:
            print()
    print()

