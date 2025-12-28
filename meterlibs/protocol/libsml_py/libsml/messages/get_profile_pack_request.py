"""SML Get Profile Pack Request message type."""

from typing import Optional
from ..buffer import SmlBuffer
from ..octet_string import OctetString, octet_string_parse, octet_string_write
from ..time import SmlTime, time_parse, time_write
from ..tree import SmlTreePath, tree_path_parse, tree_path_write, SmlTree, tree_parse, tree_write
from ..boolean import boolean_parse, boolean_write
from ..shared import SML_TYPE_LIST, SML_OPTIONAL_SKIPPED


class SmlObjReqEntryList:
    """SML Object Request Entry List (linked list)."""
    
    def __init__(self):
        """Initialize a new SmlObjReqEntryList."""
        self.object_list_entry: Optional[OctetString] = None
        self.next: Optional['SmlObjReqEntryList'] = None


class SmlGetProfilePackRequest:
    """SML Get Profile Pack Request class."""
    
    def __init__(self):
        """Initialize a new SmlGetProfilePackRequest."""
        self.server_id: Optional[OctetString] = None
        self.username: Optional[OctetString] = None
        self.password: Optional[OctetString] = None
        self.with_rawdata: Optional[int] = None  # boolean
        self.begin_time: Optional[SmlTime] = None
        self.end_time: Optional[SmlTime] = None
        self.parameter_tree_path: Optional[SmlTreePath] = None
        self.object_list: Optional[SmlObjReqEntryList] = None
        self.das_details: Optional[SmlTree] = None
    
    @classmethod
    def parse(cls, buf: SmlBuffer) -> Optional['SmlGetProfilePackRequest']:
        """Parse get profile pack request from buffer."""
        msg = cls()
        
        if buf.get_next_type() != SML_TYPE_LIST:
            buf.error = 1
            return None
        
        if buf.get_next_length() != 9:
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
        
        msg.with_rawdata = boolean_parse(buf)
        if buf.has_errors():
            return None
        
        msg.begin_time = time_parse(buf)
        if buf.has_errors():
            return None
        
        msg.end_time = time_parse(buf)
        if buf.has_errors():
            return None
        
        msg.parameter_tree_path = tree_path_parse(buf)
        if buf.has_errors():
            return None
        
        if buf.optional_is_skipped() != 1:
            if buf.get_next_type() != SML_TYPE_LIST:
                buf.error = 1
                return None
            
            length = buf.get_next_length()
            last = None
            for _ in range(length):
                entry = SmlObjReqEntryList()
                entry.object_list_entry = octet_string_parse(buf)
                if buf.has_errors():
                    return None
                
                if msg.object_list is None:
                    msg.object_list = entry
                    last = entry
                else:
                    last.next = entry
                    last = entry
        
        msg.das_details = tree_parse(buf)
        if buf.has_errors():
            return None
        
        return msg
    
    def write(self, buf: SmlBuffer) -> None:
        """Write get profile pack request to buffer."""
        buf.set_type_and_length(SML_TYPE_LIST, 9)
        octet_string_write(self.server_id, buf)
        octet_string_write(self.username, buf)
        octet_string_write(self.password, buf)
        boolean_write(self.with_rawdata, buf)
        time_write(self.begin_time, buf)
        time_write(self.end_time, buf)
        tree_path_write(self.parameter_tree_path, buf)
        
        if self.object_list:
            count = 0
            current = self.object_list
            while current:
                count += 1
                current = current.next
            
            buf.set_type_and_length(SML_TYPE_LIST, count)
            current = self.object_list
            while current:
                octet_string_write(current.object_list_entry, buf)
                current = current.next
        else:
            buf.optional_write()
        
        tree_write(self.das_details, buf)
    
    def free(self) -> None:
        """Free get profile pack request."""
        pass

