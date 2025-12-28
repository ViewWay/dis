"""SML Attention Response message type."""

from typing import Optional
from ..buffer import SmlBuffer
from ..octet_string import OctetString, octet_string_parse, octet_string_write
from ..tree import SmlTree, tree_parse, tree_write
from ..shared import SML_TYPE_LIST


class SmlAttentionResponse:
    """SML Attention Response class."""
    
    def __init__(self):
        """Initialize a new SmlAttentionResponse."""
        self.server_id: Optional[OctetString] = None
        self.attention_number: Optional[OctetString] = None
        self.attention_message: Optional[OctetString] = None
        self.attention_details: Optional[SmlTree] = None
    
    @classmethod
    def parse(cls, buf: SmlBuffer) -> Optional['SmlAttentionResponse']:
        """Parse attention response from buffer."""
        msg = cls()
        
        if buf.get_next_type() != SML_TYPE_LIST:
            buf.error = 1
            return None
        
        if buf.get_next_length() != 4:
            buf.error = 1
            return None
        
        msg.server_id = octet_string_parse(buf)
        if buf.has_errors():
            return None
        
        msg.attention_number = octet_string_parse(buf)
        if buf.has_errors():
            return None
        
        msg.attention_message = octet_string_parse(buf)
        if buf.has_errors():
            return None
        
        msg.attention_details = tree_parse(buf)
        if buf.has_errors():
            return None
        
        return msg
    
    def write(self, buf: SmlBuffer) -> None:
        """Write attention response to buffer."""
        buf.set_type_and_length(SML_TYPE_LIST, 4)
        octet_string_write(self.server_id, buf)
        octet_string_write(self.attention_number, buf)
        octet_string_write(self.attention_message, buf)
        tree_write(self.attention_details, buf)
    
    def free(self) -> None:
        """Free attention response."""
        pass

