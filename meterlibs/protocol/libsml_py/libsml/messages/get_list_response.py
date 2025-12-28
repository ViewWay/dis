"""SML Get List Response message type."""

from typing import Optional
from ..buffer import SmlBuffer
from ..octet_string import OctetString, octet_string_parse, octet_string_write
from ..time import SmlTime, time_parse, time_write
from ..list import SmlList, list_parse, list_write
from ..shared import SML_TYPE_LIST


class SmlGetListResponse:
    """SML Get List Response class."""
    
    def __init__(self):
        """Initialize a new SmlGetListResponse."""
        self.client_id: Optional[OctetString] = None
        self.server_id: Optional[OctetString] = None
        self.list_name: Optional[OctetString] = None
        self.act_sensor_time: Optional[SmlTime] = None
        self.val_list: Optional[SmlList] = None
        self.list_signature: Optional[OctetString] = None
        self.act_gateway_time: Optional[SmlTime] = None
    
    @classmethod
    def parse(cls, buf: SmlBuffer) -> Optional['SmlGetListResponse']:
        """Parse get list response from buffer."""
        msg = cls()
        
        if buf.get_next_type() != SML_TYPE_LIST:
            buf.error = 1
            return None
        
        if buf.get_next_length() != 7:
            buf.error = 1
            return None
        
        msg.client_id = octet_string_parse(buf)
        if buf.has_errors():
            return None
        
        msg.server_id = octet_string_parse(buf)
        if buf.has_errors():
            return None
        
        msg.list_name = octet_string_parse(buf)
        if buf.has_errors():
            return None
        
        msg.act_sensor_time = time_parse(buf)
        if buf.has_errors():
            return None
        
        msg.val_list = list_parse(buf)
        if buf.has_errors():
            return None
        
        msg.list_signature = octet_string_parse(buf)
        if buf.has_errors():
            return None
        
        msg.act_gateway_time = time_parse(buf)
        if buf.has_errors():
            return None
        
        return msg
    
    def write(self, buf: SmlBuffer) -> None:
        """Write get list response to buffer."""
        buf.set_type_and_length(SML_TYPE_LIST, 7)
        octet_string_write(self.client_id, buf)
        octet_string_write(self.server_id, buf)
        octet_string_write(self.list_name, buf)
        time_write(self.act_sensor_time, buf)
        list_write(self.val_list, buf)
        octet_string_write(self.list_signature, buf)
        time_write(self.act_gateway_time, buf)
    
    def free(self) -> None:
        """Free get list response."""
        pass

