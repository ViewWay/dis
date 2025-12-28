"""SML Close Request message type."""

from typing import Optional
from ..buffer import SmlBuffer
from ..octet_string import OctetString, octet_string_parse, octet_string_write
from ..shared import SML_TYPE_LIST


class SmlCloseRequest:
    """SML Close Request class."""
    
    def __init__(self):
        """Initialize a new SmlCloseRequest."""
        self.global_signature: Optional[OctetString] = None
    
    @classmethod
    def parse(cls, buf: SmlBuffer) -> Optional['SmlCloseRequest']:
        """Parse close request from buffer."""
        msg = cls()
        
        if buf.get_next_type() != SML_TYPE_LIST:
            buf.error = 1
            return None
        
        if buf.get_next_length() != 1:
            buf.error = 1
            return None
        
        msg.global_signature = octet_string_parse(buf)
        if buf.has_errors():
            return None
        
        return msg
    
    def write(self, buf: SmlBuffer) -> None:
        """Write close request to buffer."""
        buf.set_type_and_length(SML_TYPE_LIST, 1)
        octet_string_write(self.global_signature, buf)
    
    def free(self) -> None:
        """Free close request."""
        pass

