"""SML Get List Request message type."""

from typing import Optional
from ..buffer import SmlBuffer
from ..octet_string import OctetString, octet_string_parse, octet_string_write
from ..shared import SML_TYPE_LIST


class SmlGetListRequest:
    """SML Get List Request class."""
    
    def __init__(self):
        """Initialize a new SmlGetListRequest."""
        self.client_id: Optional[OctetString] = None
        self.server_id: Optional[OctetString] = None
        self.username: Optional[OctetString] = None
        self.password: Optional[OctetString] = None
        self.list_name: Optional[OctetString] = None
    
    @classmethod
    def parse(cls, buf: SmlBuffer) -> Optional['SmlGetListRequest']:
        """Parse get list request from buffer."""
        msg = cls()
        
        if buf.get_next_type() != SML_TYPE_LIST:
            buf.error = 1
            return None
        
        if buf.get_next_length() != 5:
            buf.error = 1
            return None
        
        msg.client_id = octet_string_parse(buf)
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
        
        msg.list_name = octet_string_parse(buf)
        if buf.has_errors():
            return None
        
        return msg
    
    def write(self, buf: SmlBuffer) -> None:
        """Write get list request to buffer."""
        buf.set_type_and_length(SML_TYPE_LIST, 5)
        octet_string_write(self.client_id, buf)
        octet_string_write(self.server_id, buf)
        octet_string_write(self.username, buf)
        octet_string_write(self.password, buf)
        octet_string_write(self.list_name, buf)
    
    def free(self) -> None:
        """Free get list request."""
        pass

