"""SML Get Proc Parameter Response message type."""

from typing import Optional
from ..buffer import SmlBuffer
from ..octet_string import OctetString, octet_string_parse, octet_string_write
from ..tree import SmlTreePath, tree_path_parse, tree_path_write, SmlTree, tree_parse, tree_write
from ..shared import SML_TYPE_LIST


class SmlGetProcParameterResponse:
    """SML Get Proc Parameter Response class."""
    
    def __init__(self):
        """Initialize a new SmlGetProcParameterResponse."""
        self.server_id: Optional[OctetString] = None
        self.parameter_tree_path: Optional[SmlTreePath] = None
        self.parameter_tree: Optional[SmlTree] = None
    
    @classmethod
    def parse(cls, buf: SmlBuffer) -> Optional['SmlGetProcParameterResponse']:
        """Parse get proc parameter response from buffer."""
        msg = cls()
        
        if buf.get_next_type() != SML_TYPE_LIST:
            buf.error = 1
            return None
        
        if buf.get_next_length() != 3:
            buf.error = 1
            return None
        
        msg.server_id = octet_string_parse(buf)
        if buf.has_errors():
            return None
        
        msg.parameter_tree_path = tree_path_parse(buf)
        if buf.has_errors():
            return None
        
        msg.parameter_tree = tree_parse(buf)
        if buf.has_errors():
            return None
        
        return msg
    
    def write(self, buf: SmlBuffer) -> None:
        """Write get proc parameter response to buffer."""
        buf.set_type_and_length(SML_TYPE_LIST, 3)
        octet_string_write(self.server_id, buf)
        tree_path_write(self.parameter_tree_path, buf)
        tree_write(self.parameter_tree, buf)
    
    def free(self) -> None:
        """Free get proc parameter response."""
        pass

