"""SML Get Profile List Response message type."""

from typing import Optional
from ..buffer import SmlBuffer
from ..octet_string import OctetString, octet_string_parse, octet_string_write
from ..time import SmlTime, time_parse, time_write
from ..tree import SmlTreePath, tree_path_parse, tree_path_write
from ..number import u32_parse, u32_write, u64_parse, u64_write
from ..list import SmlSequence, sequence_parse, sequence_write
from ..tree import SmlPeriodEntry, period_entry_parse, period_entry_write
from ..shared import SML_TYPE_LIST


class SmlGetProfileListResponse:
    """SML Get Profile List Response class."""
    
    def __init__(self):
        """Initialize a new SmlGetProfileListResponse."""
        self.server_id: Optional[OctetString] = None
        self.act_time: Optional[SmlTime] = None
        self.reg_period: Optional[bytes] = None  # u32
        self.parameter_tree_path: Optional[SmlTreePath] = None
        self.val_time: Optional[SmlTime] = None
        self.status: Optional[bytes] = None  # u64
        self.period_list: Optional[SmlSequence] = None
        self.rawdata: Optional[OctetString] = None
        self.period_signature: Optional[OctetString] = None
    
    @classmethod
    def parse(cls, buf: SmlBuffer) -> Optional['SmlGetProfileListResponse']:
        """Parse get profile list response from buffer."""
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
        
        msg.act_time = time_parse(buf)
        if buf.has_errors():
            return None
        
        msg.reg_period = u32_parse(buf)
        if buf.has_errors():
            return None
        
        msg.parameter_tree_path = tree_path_parse(buf)
        if buf.has_errors():
            return None
        
        msg.val_time = time_parse(buf)
        if buf.has_errors():
            return None
        
        msg.status = u64_parse(buf)
        if buf.has_errors():
            return None
        
        def period_entry_parse_func(b):
            return period_entry_parse(b)
        
        def period_entry_free_func(e):
            if e:
                e.free()
        
        msg.period_list = sequence_parse(buf, period_entry_parse_func, period_entry_free_func)
        if buf.has_errors():
            return None
        
        msg.rawdata = octet_string_parse(buf)
        if buf.has_errors():
            return None
        
        msg.period_signature = octet_string_parse(buf)
        if buf.has_errors():
            return None
        
        return msg
    
    def write(self, buf: SmlBuffer) -> None:
        """Write get profile list response to buffer."""
        buf.set_type_and_length(SML_TYPE_LIST, 9)
        octet_string_write(self.server_id, buf)
        time_write(self.act_time, buf)
        u32_write(self.reg_period, buf)
        tree_path_write(self.parameter_tree_path, buf)
        time_write(self.val_time, buf)
        u64_write(self.status, buf)
        
        def period_entry_write_func(e, b):
            if e:
                period_entry_write(e, b)
        
        sequence_write(self.period_list, buf, period_entry_write_func)
        octet_string_write(self.rawdata, buf)
        octet_string_write(self.period_signature, buf)
    
    def free(self) -> None:
        """Free get profile list response."""
        pass

