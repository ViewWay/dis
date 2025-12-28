"""SML Open Request message type."""

from typing import Optional
from ..buffer import SmlBuffer
from ..octet_string import OctetString, octet_string_parse, octet_string_write
from ..number import u8_parse, u8_write
from ..shared import SML_TYPE_LIST


class SmlOpenRequest:
    """SML Open Request class."""
    
    def __init__(self):
        """Initialize a new SmlOpenRequest."""
        self.codepage: Optional[OctetString] = None
        self.client_id: Optional[OctetString] = None
        self.req_file_id: Optional[OctetString] = None
        self.server_id: Optional[OctetString] = None
        self.username: Optional[OctetString] = None
        self.password: Optional[OctetString] = None
        self.sml_version: Optional[bytes] = None  # u8
    
    @classmethod
    def parse(cls, buf: SmlBuffer) -> Optional['SmlOpenRequest']:
        """Parse open request from buffer."""
        msg = cls()
        
        if buf.get_next_type() != SML_TYPE_LIST:
            buf.error = 1
            return None
        
        if buf.get_next_length() != 7:
            buf.error = 1
            return None
        
        msg.codepage = octet_string_parse(buf)
        if buf.has_errors():
            return None
        
        msg.client_id = octet_string_parse(buf)
        if buf.has_errors():
            return None
        
        msg.req_file_id = octet_string_parse(buf)
        if buf.has_errors():
            return None
        
        msg.server_id = octet_string_parse(buf)
        if buf.has_errors():
            return None
        
        msg.username = octet_string_parse(buf)
        if buf.has_errors():
            return None
        
        msg.password = octet_string_parse(buf)
        if buf.has_errors():
            return None
        
        msg.sml_version = u8_parse(buf)
        if buf.has_errors():
            return None
        
        return msg
    
    def write(self, buf: SmlBuffer) -> None:
        """Write open request to buffer."""
        buf.set_type_and_length(SML_TYPE_LIST, 7)
        octet_string_write(self.codepage, buf)
        octet_string_write(self.client_id, buf)
        octet_string_write(self.req_file_id, buf)
        octet_string_write(self.server_id, buf)
        octet_string_write(self.username, buf)
        octet_string_write(self.password, buf)
        u8_write(self.sml_version, buf)
    
    def free(self) -> None:
        """Free open request."""
        pass

