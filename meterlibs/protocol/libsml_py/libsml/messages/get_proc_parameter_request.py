"""SML Get Proc Parameter Request message type."""

from typing import Optional
from ..buffer import SmlBuffer
from ..octet_string import OctetString, octet_string_parse, octet_string_write
from ..tree import SmlTreePath, tree_path_parse, tree_path_write
from ..shared import SML_TYPE_LIST


class SmlGetProcParameterRequest:
    """SML Get Proc Parameter Request class."""
    
    def __init__(self):
        """Initialize a new SmlGetProcParameterRequest."""
        self.server_id: Optional[OctetString] = None
        self.username: Optional[OctetString] = None
        self.password: Optional[OctetString] = None
        self.parameter_tree_path: Optional[SmlTreePath] = None
        self.attribute: Optional[OctetString] = None
    
    @classmethod
    def parse(cls, buf: SmlBuffer) -> Optional['SmlGetProcParameterRequest']:
        """Parse get proc parameter request from buffer."""
        msg = cls()
        
        if buf.get_next_type() != SML_TYPE_LIST:
            buf.error = 1
            return None
        
        if buf.get_next_length() != 5:
            buf.error = 1
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
        
        msg.parameter_tree_path = tree_path_parse(buf)
        if buf.has_errors():
            return None
        
        msg.attribute = octet_string_parse(buf)
        if buf.has_errors():
            return None
        
        return msg
    
    def write(self, buf: SmlBuffer) -> None:
        """Write get proc parameter request to buffer."""
        buf.set_type_and_length(SML_TYPE_LIST, 5)
        octet_string_write(self.server_id, buf)
        octet_string_write(self.username, buf)
        octet_string_write(self.password, buf)
        tree_path_write(self.parameter_tree_path, buf)
        octet_string_write(self.attribute, buf)
    
    def free(self) -> None:
        """Free get proc parameter request."""
        pass

