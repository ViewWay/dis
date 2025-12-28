"""
SML Tree - Tree structure handling

This module implements tree structures for SML protocol including
SmlTree, SmlTreePath, SmlProcParValue, SmlTupelEntry, and SmlPeriodEntry.
"""

from typing import Optional, List, Union
from .buffer import SmlBuffer
from .octet_string import OctetString, octet_string_parse, octet_string_write
from .time import SmlTime, time_parse, time_write
from .value import SmlValue, value_parse, value_write
from .number import (
    u8_parse, u8_write, u8_init,
    i8_parse, i8_write,
    i64_parse, i64_write,
    u64_parse, u64_write,
)
from .shared import SML_TYPE_LIST, SML_OPTIONAL_SKIPPED

# Proc par value tags
SML_PROC_PAR_VALUE_TAG_VALUE = 0x01
SML_PROC_PAR_VALUE_TAG_PERIOD_ENTRY = 0x02
SML_PROC_PAR_VALUE_TAG_TUPEL_ENTRY = 0x03
SML_PROC_PAR_VALUE_TAG_TIME = 0x04


class SmlTupelEntry:
    """SML Tupel Entry class."""
    
    def __init__(self):
        """Initialize a new SmlTupelEntry."""
        self.server_id: Optional[OctetString] = None
        self.sec_index: Optional[SmlTime] = None
        self.status: Optional[bytes] = None  # u64
        
        self.unit_pA: Optional[bytes] = None  # u8
        self.scaler_pA: Optional[bytes] = None  # i8
        self.value_pA: Optional[bytes] = None  # i64
        
        self.unit_R1: Optional[bytes] = None  # u8
        self.scaler_R1: Optional[bytes] = None  # i8
        self.value_R1: Optional[bytes] = None  # i64
        
        self.unit_R4: Optional[bytes] = None  # u8
        self.scaler_R4: Optional[bytes] = None  # i8
        self.value_R4: Optional[bytes] = None  # i64
        
        self.signature_pA_R1_R4: Optional[OctetString] = None
        
        self.unit_mA: Optional[bytes] = None  # u8
        self.scaler_mA: Optional[bytes] = None  # i8
        self.value_mA: Optional[bytes] = None  # i64
        
        self.unit_R2: Optional[bytes] = None  # u8
        self.scaler_R2: Optional[bytes] = None  # i8
        self.value_R2: Optional[bytes] = None  # i64
        
        self.unit_R3: Optional[bytes] = None  # u8
        self.scaler_R3: Optional[bytes] = None  # i8
        self.value_R3: Optional[bytes] = None  # i64
        
        self.signature_mA_R2_R3: Optional[OctetString] = None
    
    @classmethod
    def parse(cls, buf: SmlBuffer) -> Optional['SmlTupelEntry']:
        """Parse tupel entry from buffer."""
        if buf.optional_is_skipped() == 1:
            return None
        
        tupel = cls()
        
        if buf.get_next_type() != SML_TYPE_LIST:
            buf.error = 1
            return None
        
        if buf.get_next_length() != 23:
            buf.error = 1
            return None
        
        tupel.server_id = octet_string_parse(buf)
        if buf.has_errors():
            return None
        tupel.sec_index = time_parse(buf)
        if buf.has_errors():
            return None
        tupel.status = u64_parse(buf)
        if buf.has_errors():
            return None
        
        tupel.unit_pA = u8_parse(buf)
        if buf.has_errors():
            return None
        tupel.scaler_pA = i8_parse(buf)
        if buf.has_errors():
            return None
        tupel.value_pA = i64_parse(buf)
        if buf.has_errors():
            return None
        
        tupel.unit_R1 = u8_parse(buf)
        if buf.has_errors():
            return None
        tupel.scaler_R1 = i8_parse(buf)
        if buf.has_errors():
            return None
        tupel.value_R1 = i64_parse(buf)
        if buf.has_errors():
            return None
        
        tupel.unit_R4 = u8_parse(buf)
        if buf.has_errors():
            return None
        tupel.scaler_R4 = i8_parse(buf)
        if buf.has_errors():
            return None
        tupel.value_R4 = i64_parse(buf)
        if buf.has_errors():
            return None
        
        tupel.signature_pA_R1_R4 = octet_string_parse(buf)
        if buf.has_errors():
            return None
        
        tupel.unit_mA = u8_parse(buf)
        if buf.has_errors():
            return None
        tupel.scaler_mA = i8_parse(buf)
        if buf.has_errors():
            return None
        tupel.value_mA = i64_parse(buf)
        if buf.has_errors():
            return None
        
        tupel.unit_R2 = u8_parse(buf)
        if buf.has_errors():
            return None
        tupel.scaler_R2 = i8_parse(buf)
        if buf.has_errors():
            return None
        tupel.value_R2 = i64_parse(buf)
        if buf.has_errors():
            return None
        
        tupel.unit_R3 = u8_parse(buf)
        if buf.has_errors():
            return None
        tupel.scaler_R3 = i8_parse(buf)
        if buf.has_errors():
            return None
        tupel.value_R3 = i64_parse(buf)
        if buf.has_errors():
            return None
        
        tupel.signature_mA_R2_R3 = octet_string_parse(buf)
        if buf.has_errors():
            return None
        
        return tupel
    
    def write(self, buf: SmlBuffer) -> None:
        """Write tupel entry to buffer."""
        if self.server_id is None:
            buf.optional_write()
            return
        
        buf.set_type_and_length(SML_TYPE_LIST, 23)
        
        octet_string_write(self.server_id, buf)
        time_write(self.sec_index, buf)
        u64_write(self.status, buf)
        
        u8_write(self.unit_pA, buf)
        i8_write(self.scaler_pA, buf)
        i64_write(self.value_pA, buf)
        
        u8_write(self.unit_R1, buf)
        i8_write(self.scaler_R1, buf)
        i64_write(self.value_R1, buf)
        
        u8_write(self.unit_R4, buf)
        i8_write(self.scaler_R4, buf)
        i64_write(self.value_R4, buf)
        
        octet_string_write(self.signature_pA_R1_R4, buf)
        
        u8_write(self.unit_mA, buf)
        i8_write(self.scaler_mA, buf)
        i64_write(self.value_mA, buf)
        
        u8_write(self.unit_R2, buf)
        i8_write(self.scaler_R2, buf)
        i64_write(self.value_R2, buf)
        
        u8_write(self.unit_R3, buf)
        i8_write(self.scaler_R3, buf)
        i64_write(self.value_R3, buf)
        
        octet_string_write(self.signature_mA_R2_R3, buf)
    
    def free(self) -> None:
        """Free tupel entry."""
        pass


class SmlPeriodEntry:
    """SML Period Entry class."""
    
    def __init__(self):
        """Initialize a new SmlPeriodEntry."""
        self.obj_name: Optional[OctetString] = None
        self.unit: Optional[bytes] = None  # u8
        self.scaler: Optional[bytes] = None  # i8
        self.value: Optional[SmlValue] = None
        self.value_signature: Optional[OctetString] = None
    
    @classmethod
    def parse(cls, buf: SmlBuffer) -> Optional['SmlPeriodEntry']:
        """Parse period entry from buffer."""
        if buf.optional_is_skipped() == 1:
            return None
        
        period = cls()
        
        if buf.get_next_type() != SML_TYPE_LIST:
            buf.error = 1
            return None
        
        if buf.get_next_length() != 5:
            buf.error = 1
            return None
        
        period.obj_name = octet_string_parse(buf)
        if buf.has_errors():
            return None
        
        period.unit = u8_parse(buf)
        if buf.has_errors():
            return None
        
        period.scaler = i8_parse(buf)
        if buf.has_errors():
            return None
        
        period.value = value_parse(buf)
        if buf.has_errors():
            return None
        
        period.value_signature = octet_string_parse(buf)
        if buf.has_errors():
            return None
        
        return period
    
    def write(self, buf: SmlBuffer) -> None:
        """Write period entry to buffer."""
        if self.obj_name is None:
            buf.optional_write()
            return
        
        buf.set_type_and_length(SML_TYPE_LIST, 5)
        
        octet_string_write(self.obj_name, buf)
        u8_write(self.unit, buf)
        i8_write(self.scaler, buf)
        value_write(self.value, buf)
        octet_string_write(self.value_signature, buf)
    
    def free(self) -> None:
        """Free period entry."""
        pass


class SmlProcParValue:
    """SML Proc Par Value class."""
    
    def __init__(self):
        """Initialize a new SmlProcParValue."""
        self.tag: Optional[bytes] = None  # u8
        self.data: Union[Optional[SmlValue], Optional[SmlPeriodEntry], Optional[SmlTupelEntry], Optional[SmlTime]] = None
    
    @classmethod
    def parse(cls, buf: SmlBuffer) -> Optional['SmlProcParValue']:
        """Parse proc par value from buffer."""
        if buf.optional_is_skipped() == 1:
            return None
        
        ppv = cls()
        
        if buf.get_next_type() != SML_TYPE_LIST:
            buf.error = 1
            return None
        
        if buf.get_next_length() != 2:
            buf.error = 1
            return None
        
        ppv.tag = u8_parse(buf)
        if buf.has_errors() or ppv.tag is None:
            return None
        
        tag_value = ppv.tag[0] if isinstance(ppv.tag, bytes) else ppv.tag
        
        if tag_value == SML_PROC_PAR_VALUE_TAG_VALUE:
            ppv.data = value_parse(buf)
        elif tag_value == SML_PROC_PAR_VALUE_TAG_PERIOD_ENTRY:
            ppv.data = SmlPeriodEntry.parse(buf)
        elif tag_value == SML_PROC_PAR_VALUE_TAG_TUPEL_ENTRY:
            ppv.data = SmlTupelEntry.parse(buf)
        elif tag_value == SML_PROC_PAR_VALUE_TAG_TIME:
            ppv.data = time_parse(buf)
        else:
            buf.error = 1
            return None
        
        return ppv
    
    def write(self, buf: SmlBuffer) -> None:
        """Write proc par value to buffer."""
        if self.tag is None:
            buf.optional_write()
            return
        
        buf.set_type_and_length(SML_TYPE_LIST, 2)
        u8_write(self.tag, buf)
        
        tag_value = self.tag[0] if isinstance(self.tag, bytes) else self.tag
        
        if tag_value == SML_PROC_PAR_VALUE_TAG_VALUE:
            value_write(self.data, buf)
        elif tag_value == SML_PROC_PAR_VALUE_TAG_PERIOD_ENTRY:
            if isinstance(self.data, SmlPeriodEntry):
                self.data.write(buf)
        elif tag_value == SML_PROC_PAR_VALUE_TAG_TUPEL_ENTRY:
            if isinstance(self.data, SmlTupelEntry):
                self.data.write(buf)
        elif tag_value == SML_PROC_PAR_VALUE_TAG_TIME:
            time_write(self.data, buf)
    
    def free(self) -> None:
        """Free proc par value."""
        pass


class SmlTree:
    """SML Tree class."""
    
    def __init__(self):
        """Initialize a new SmlTree."""
        self.parameter_name: Optional[OctetString] = None
        self.parameter_value: Optional[SmlProcParValue] = None
        self.child_list: List['SmlTree'] = []
        self.child_list_len: int = 0
    
    @classmethod
    def parse(cls, buf: SmlBuffer) -> Optional['SmlTree']:
        """Parse tree from buffer."""
        if buf.optional_is_skipped() == 1:
            return None
        
        tree = cls()
        
        if buf.get_next_type() != SML_TYPE_LIST:
            buf.error = 1
            return None
        
        if buf.get_next_length() != 3:
            buf.error = 1
            return None
        
        tree.parameter_name = octet_string_parse(buf)
        if buf.has_errors():
            return None
        
        tree.parameter_value = SmlProcParValue.parse(buf)
        if buf.has_errors():
            return None
        
        if buf.optional_is_skipped() != 1:
            if buf.get_next_type() != SML_TYPE_LIST:
                buf.error = 1
                return None
            
            elems = buf.get_next_length()
            for _ in range(elems):
                child = SmlTree.parse(buf)
                if buf.has_errors():
                    return None
                if child:
                    tree.add_tree(child)
        
        return tree
    
    def write(self, buf: SmlBuffer) -> None:
        """Write tree to buffer."""
        if self.parameter_name is None:
            buf.optional_write()
            return
        
        buf.set_type_and_length(SML_TYPE_LIST, 3)
        
        octet_string_write(self.parameter_name, buf)
        if self.parameter_value:
            self.parameter_value.write(buf)
        else:
            buf.optional_write()
        
        if self.child_list and self.child_list_len > 0:
            buf.set_type_and_length(SML_TYPE_LIST, self.child_list_len)
            for child in self.child_list:
                child.write(buf)
        else:
            buf.optional_write()
    
    def add_tree(self, tree: 'SmlTree') -> None:
        """Add child tree."""
        self.child_list.append(tree)
        self.child_list_len += 1
    
    def free(self) -> None:
        """Free tree and all children."""
        for child in self.child_list:
            child.free()
        self.child_list = []
        self.child_list_len = 0


class SmlTreePath:
    """SML Tree Path class."""
    
    def __init__(self):
        """Initialize a new SmlTreePath."""
        self.path_entries: List[OctetString] = []
        self.path_entries_len: int = 0
    
    @classmethod
    def parse(cls, buf: SmlBuffer) -> Optional['SmlTreePath']:
        """Parse tree path from buffer."""
        if buf.optional_is_skipped() == 1:
            return None
        
        tree_path = cls()
        
        if buf.get_next_type() != SML_TYPE_LIST:
            buf.error = 1
            return None
        
        elems = buf.get_next_length()
        for _ in range(elems):
            entry = octet_string_parse(buf)
            if buf.has_errors():
                return None
            if entry:
                tree_path.add_path_entry(entry)
        
        return tree_path
    
    def write(self, buf: SmlBuffer) -> None:
        """Write tree path to buffer."""
        if self.path_entries_len == 0:
            buf.optional_write()
            return
        
        buf.set_type_and_length(SML_TYPE_LIST, self.path_entries_len)
        for entry in self.path_entries:
            octet_string_write(entry, buf)
    
    def add_path_entry(self, entry: OctetString) -> None:
        """Add path entry."""
        self.path_entries.append(entry)
        self.path_entries_len += 1
    
    def free(self) -> None:
        """Free tree path."""
        for entry in self.path_entries:
            entry.free()
        self.path_entries = []
        self.path_entries_len = 0


# Free functions for API compatibility
def tree_init() -> SmlTree:
    """Initialize a new SmlTree."""
    return SmlTree()


def tree_parse(buf: SmlBuffer) -> Optional[SmlTree]:
    """Parse tree from buffer."""
    return SmlTree.parse(buf)


def tree_write(tree: Optional[SmlTree], buf: SmlBuffer) -> None:
    """Write tree to buffer."""
    if tree is None:
        buf.optional_write()
    else:
        tree.write(buf)


def tree_add_tree(base_tree: SmlTree, tree: SmlTree) -> None:
    """Add child tree."""
    base_tree.add_tree(tree)


def tree_free(tree: Optional[SmlTree]) -> None:
    """Free tree."""
    if tree:
        tree.free()


def tree_path_init() -> SmlTreePath:
    """Initialize a new SmlTreePath."""
    return SmlTreePath()


def tree_path_parse(buf: SmlBuffer) -> Optional[SmlTreePath]:
    """Parse tree path from buffer."""
    return SmlTreePath.parse(buf)


def tree_path_write(tree_path: Optional[SmlTreePath], buf: SmlBuffer) -> None:
    """Write tree path to buffer."""
    if tree_path is None:
        buf.optional_write()
    else:
        tree_path.write(buf)


def tree_path_add_path_entry(tree_path: SmlTreePath, entry: OctetString) -> None:
    """Add path entry."""
    tree_path.add_path_entry(entry)


def tree_path_free(tree_path: Optional[SmlTreePath]) -> None:
    """Free tree path."""
    if tree_path:
        tree_path.free()


def proc_par_value_init() -> SmlProcParValue:
    """Initialize a new SmlProcParValue."""
    return SmlProcParValue()


def proc_par_value_parse(buf: SmlBuffer) -> Optional[SmlProcParValue]:
    """Parse proc par value from buffer."""
    return SmlProcParValue.parse(buf)


def proc_par_value_write(value: Optional[SmlProcParValue], buf: SmlBuffer) -> None:
    """Write proc par value to buffer."""
    if value is None:
        buf.optional_write()
    else:
        value.write(buf)


def proc_par_value_free(value: Optional[SmlProcParValue]) -> None:
    """Free proc par value."""
    pass


def tupel_entry_init() -> SmlTupelEntry:
    """Initialize a new SmlTupelEntry."""
    return SmlTupelEntry()


def tupel_entry_parse(buf: SmlBuffer) -> Optional[SmlTupelEntry]:
    """Parse tupel entry from buffer."""
    return SmlTupelEntry.parse(buf)


def tupel_entry_write(tupel: Optional[SmlTupelEntry], buf: SmlBuffer) -> None:
    """Write tupel entry to buffer."""
    if tupel is None:
        buf.optional_write()
    else:
        tupel.write(buf)


def tupel_entry_free(tupel: Optional[SmlTupelEntry]) -> None:
    """Free tupel entry."""
    pass


def period_entry_init() -> SmlPeriodEntry:
    """Initialize a new SmlPeriodEntry."""
    return SmlPeriodEntry()


def period_entry_parse(buf: SmlBuffer) -> Optional[SmlPeriodEntry]:
    """Parse period entry from buffer."""
    return SmlPeriodEntry.parse(buf)


def period_entry_write(period: Optional[SmlPeriodEntry], buf: SmlBuffer) -> None:
    """Write period entry to buffer."""
    if period is None:
        buf.optional_write()
    else:
        period.write(buf)


def period_entry_free(period: Optional[SmlPeriodEntry]) -> None:
    """Free period entry."""
    pass

