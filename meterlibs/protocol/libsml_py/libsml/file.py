"""
SML File - File parsing and writing

This module implements SML file parsing and writing, which consists of
multiple SML messages.
"""

from typing import List, Optional
from .buffer import SmlBuffer, SmlBufferError
from .message import SmlMessage, message_parse, message_write
from .shared import SML_MESSAGE_END

# EDL meter must provide at least 250 bytes as a receive buffer
SML_FILE_BUFFER_LENGTH = 512


class SmlFile:
    """
    SML File class.
    
    A SML file consists of multiple SML messages.
    
    Attributes:
        messages: List of SML messages
        messages_len: Number of messages
        buf: Buffer for reading/writing
    """
    
    def __init__(self):
        """Initialize a new SmlFile."""
        self.messages: List[SmlMessage] = []
        self.messages_len: int = 0
        self.buf: Optional[SmlBuffer] = None
    
    @classmethod
    def init(cls) -> 'SmlFile':
        """
        Initialize a new SmlFile.
        
        Returns:
            New SmlFile instance
        """
        file = cls()
        file.buf = SmlBuffer(SML_FILE_BUFFER_LENGTH)
        return file
    
    @classmethod
    def parse(cls, buffer: bytes, buffer_len: int) -> 'SmlFile':
        """
        Parse a SML file from bytes.
        
        Args:
            buffer: The bytes to parse
            buffer_len: Length of buffer
            
        Returns:
            SmlFile instance
        """
        file = cls()
        file.buf = SmlBuffer.from_bytes(buffer[:buffer_len])
        
        # Parse all messages
        while file.buf.cursor < file.buf.buffer_len:
            if file.buf.get_current_byte() == SML_MESSAGE_END:
                # Reading trailing zeroed bytes
                file.buf.update_bytes_read(1)
                continue
            
            msg = message_parse(file.buf)
            
            if file.buf.has_errors():
                print("libsml: warning: could not read the whole file")
                break
            
            if msg:
                file.add_message(msg)
        
        return file
    
    def add_message(self, message: SmlMessage) -> None:
        """
        Add a message to the file.
        
        Args:
            message: The message to add
        """
        self.messages.append(message)
        self.messages_len += 1
    
    def write(self) -> None:
        """Write all messages to the buffer."""
        if self.buf is None:
            self.buf = SmlBuffer(SML_FILE_BUFFER_LENGTH)
        
        self.buf.cursor = 0
        
        if self.messages and self.messages_len > 0:
            for message in self.messages:
                message_write(message, self.buf)
    
    def free(self) -> None:
        """Free file and all messages."""
        for message in self.messages:
            message.free()
        self.messages = []
        self.messages_len = 0
        if self.buf:
            self.buf.free()
    
    def print(self) -> None:
        """Print file information."""
        cursor_pos = self.buf.cursor if self.buf else 0
        print(f"SML file ({self.messages_len} SML messages, {cursor_pos} bytes)")
        for i, msg in enumerate(self.messages):
            if msg.message_body and msg.message_body.tag:
                tag_value = int.from_bytes(msg.message_body.tag, 'big')
                print(f"SML message {tag_value:04X}")


# Free functions for API compatibility
def file_init() -> SmlFile:
    """Initialize a new SmlFile."""
    return SmlFile.init()


def file_parse(buffer: bytes, buffer_len: int) -> SmlFile:
    """Parse a SML file from bytes."""
    return SmlFile.parse(buffer, buffer_len)


def file_add_message(file: SmlFile, message: SmlMessage) -> None:
    """Add a message to the file."""
    file.add_message(message)


def file_write(file: SmlFile) -> None:
    """Write all messages to the buffer."""
    file.write()


def file_free(file: Optional[SmlFile]) -> None:
    """Free file."""
    if file:
        file.free()


def file_print(file: SmlFile) -> None:
    """Print file information."""
    file.print()

