"""SML Get Profile Pack Response message type."""

from typing import Optional, List
from ..buffer import SmlBuffer
from ..octet_string import OctetString, octet_string_parse, octet_string_write
from ..time import SmlTime, time_parse, time_write
from ..tree import SmlTreePath, tree_path_parse, tree_path_write
from ..number import u32_parse, u32_write, u64_parse, u64_write, i8_parse, i8_write, u8_parse, u8_write
from ..value import SmlValue, value_parse, value_write
from ..list import SmlSequence, sequence_parse, sequence_write
from ..shared import SML_TYPE_LIST


class SmlProfObjHeaderEntry:
    """SML Profile Object Header Entry."""
    
    def __init__(self):
        """Initialize a new SmlProfObjHeaderEntry."""
        self.obj_name: Optional[OctetString] = None
        self.unit: Optional[bytes] = None  # u8
        self.scaler: Optional[bytes] = None  # i8
    
    @classmethod
    def parse(cls, buf: SmlBuffer) -> Optional['SmlProfObjHeaderEntry']:
        """Parse profile object header entry from buffer."""
        entry = cls()
        
        if buf.get_next_type() != SML_TYPE_LIST:
            buf.error = 1
            return None
        
        if buf.get_next_length() != 3:
            buf.error = 1
            return None
        
        entry.obj_name = octet_string_parse(buf)
        if buf.has_errors():
            return None
        
        entry.unit = u8_parse(buf)
        if buf.has_errors():
            return None
        
        entry.scaler = i8_parse(buf)
        if buf.has_errors():
            return None
        
        return entry
    
    def write(self, buf: SmlBuffer) -> None:
        """Write profile object header entry to buffer."""
        buf.set_type_and_length(SML_TYPE_LIST, 3)
        octet_string_write(self.obj_name, buf)
        u8_write(self.unit, buf)
        i8_write(self.scaler, buf)
    
    def free(self) -> None:
        """Free profile object header entry."""
        pass


class SmlValueEntry:
    """SML Value Entry."""
    
    def __init__(self):
        """Initialize a new SmlValueEntry."""
        self.value: Optional[SmlValue] = None
        self.value_signature: Optional[OctetString] = None
    
    @classmethod
    def parse(cls, buf: SmlBuffer) -> Optional['SmlValueEntry']:
        """Parse value entry from buffer."""
        entry = cls()
        
        if buf.get_next_type() != SML_TYPE_LIST:
            buf.error = 1
            return None
        
        if buf.get_next_length() != 2:
            buf.error = 1
            return None
        
        entry.value = value_parse(buf)
        if buf.has_errors():
            return None
        
        entry.value_signature = octet_string_parse(buf)
        if buf.has_errors():
            return None
        
        return entry
    
    def write(self, buf: SmlBuffer) -> None:
        """Write value entry to buffer."""
        buf.set_type_and_length(SML_TYPE_LIST, 2)
        value_write(self.value, buf)
        octet_string_write(self.value_signature, buf)
    
    def free(self) -> None:
        """Free value entry."""
        pass


class SmlProfObjPeriodEntry:
    """SML Profile Object Period Entry."""
    
    def __init__(self):
        """Initialize a new SmlProfObjPeriodEntry."""
        self.val_time: Optional[SmlTime] = None
        self.status: Optional[bytes] = None  # u64
        self.value_list: Optional[SmlSequence] = None
        self.period_signature: Optional[OctetString] = None
    
    @classmethod
    def parse(cls, buf: SmlBuffer) -> Optional['SmlProfObjPeriodEntry']:
        """Parse profile object period entry from buffer."""
        entry = cls()
        
        if buf.get_next_type() != SML_TYPE_LIST:
            buf.error = 1
            return None
        
        if buf.get_next_length() != 4:
            buf.error = 1
            return None
        
        entry.val_time = time_parse(buf)
        if buf.has_errors():
            return None
        
        entry.status = u64_parse(buf)
        if buf.has_errors():
            return None
        
        def value_entry_parse(b):
            return SmlValueEntry.parse(b)
        
        def value_entry_free(e):
            if e:
                e.free()
        
        entry.value_list = sequence_parse(buf, value_entry_parse, value_entry_free)
        if buf.has_errors():
            return None
        
        entry.period_signature = octet_string_parse(buf)
        if buf.has_errors():
            return None
        
        return entry
    
    def write(self, buf: SmlBuffer) -> None:
        """Write profile object period entry to buffer."""
        buf.set_type_and_length(SML_TYPE_LIST, 4)
        time_write(self.val_time, buf)
        u64_write(self.status, buf)
        
        def value_entry_write(e, b):
            if e:
                e.write(b)
        
        sequence_write(self.value_list, buf, value_entry_write)
        octet_string_write(self.period_signature, buf)
    
    def free(self) -> None:
        """Free profile object period entry."""
        pass


class SmlGetProfilePackResponse:
    """SML Get Profile Pack Response class."""
    
    def __init__(self):
        """Initialize a new SmlGetProfilePackResponse."""
        self.server_id: Optional[OctetString] = None
        self.act_time: Optional[SmlTime] = None
        self.reg_period: Optional[bytes] = None  # u32
        self.parameter_tree_path: Optional[SmlTreePath] = None
        self.header_list: Optional[SmlSequence] = None
        self.period_list: Optional[SmlSequence] = None
        self.rawdata: Optional[OctetString] = None
        self.profile_signature: Optional[OctetString] = None
    
    @classmethod
    def parse(cls, buf: SmlBuffer) -> Optional['SmlGetProfilePackResponse']:
        """Parse get profile pack response from buffer."""
        msg = cls()
        
        if buf.get_next_type() != SML_TYPE_LIST:
            buf.error = 1
            return None
        
        if buf.get_next_length() != 8:
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
        
        def header_entry_parse(b):
            return SmlProfObjHeaderEntry.parse(b)
        
        def header_entry_free(e):
            if e:
                e.free()
        
        msg.header_list = sequence_parse(buf, header_entry_parse, header_entry_free)
        if buf.has_errors():
            return None
        
        def period_entry_parse(b):
            return SmlProfObjPeriodEntry.parse(b)
        
        def period_entry_free(e):
            if e:
                e.free()
        
        msg.period_list = sequence_parse(buf, period_entry_parse, period_entry_free)
        if buf.has_errors():
            return None
        
        msg.rawdata = octet_string_parse(buf)
        if buf.has_errors():
            return None
        
        msg.profile_signature = octet_string_parse(buf)
        if buf.has_errors():
            return None
        
        return msg
    
    def write(self, buf: SmlBuffer) -> None:
        """Write get profile pack response to buffer."""
        buf.set_type_and_length(SML_TYPE_LIST, 8)
        octet_string_write(self.server_id, buf)
        time_write(self.act_time, buf)
        u32_write(self.reg_period, buf)
        tree_path_write(self.parameter_tree_path, buf)
        
        def header_entry_write(e, b):
            if e:
                e.write(b)
        
        sequence_write(self.header_list, buf, header_entry_write)
        
        def period_entry_write(e, b):
            if e:
                e.write(b)
        
        sequence_write(self.period_list, buf, period_entry_write)
        octet_string_write(self.rawdata, buf)
        octet_string_write(self.profile_signature, buf)
    
    def free(self) -> None:
        """Free get profile pack response."""
        pass

