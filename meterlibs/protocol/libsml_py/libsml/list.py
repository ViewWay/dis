"""
SML List - List structure handling

This module implements SmlSequence and SmlList classes for SML protocol,
including compatibility workarounds for DZG meters.
"""

from typing import Optional, Callable, List, Any
from .buffer import SmlBuffer
from .octet_string import OctetString, octet_string_parse, octet_string_write
from .status import SmlStatus, status_parse, status_write
from .time import SmlTime, time_parse, time_write
from .value import SmlValue, value_parse, value_write
from .number import u8_parse, i8_parse, u8_write, i8_write
from .shared import SML_TYPE_LIST, SML_TYPE_UNSIGNED, SML_OPTIONAL_SKIPPED, SML_TYPE_FIELD, SML_ANOTHER_TL, SML_LENGTH_FIELD


class SmlSequence:
    """
    Generic sequence container for SML protocol.
    
    Attributes:
        elems: List of elements
        elems_len: Number of elements
    """
    
    def __init__(self, elem_free: Optional[Callable] = None):
        """
        Initialize a new SmlSequence.
        
        Args:
            elem_free: Optional function to free elements
        """
        self.elems: List[Any] = []
        self.elems_len: int = 0
        self.elem_free: Optional[Callable] = elem_free
    
    @classmethod
    def parse(cls, buf: SmlBuffer, elem_parse: Callable, elem_free: Optional[Callable] = None) -> Optional['SmlSequence']:
        """
        Parse sequence from buffer.
        
        Args:
            buf: The buffer to parse from
            elem_parse: Function to parse individual elements
            elem_free: Optional function to free elements
            
        Returns:
            SmlSequence instance, or None on error
        """
        if buf.get_next_type() != SML_TYPE_LIST:
            buf.error = 1
            return None
        
        seq = cls(elem_free)
        length = buf.get_next_length()
        
        for i in range(length):
            elem = elem_parse(buf)
            if buf.has_errors():
                seq.free()
                return None
            seq.add(elem)
        
        return seq
    
    def write(self, buf: SmlBuffer, elem_write: Callable) -> None:
        """
        Write sequence to buffer.
        
        Args:
            buf: The buffer to write to
            elem_write: Function to write individual elements
        """
        if self.elems_len == 0:
            buf.optional_write()
            return
        
        buf.set_type_and_length(SML_TYPE_LIST, self.elems_len)
        
        for i in range(self.elems_len):
            elem_write(self.elems[i], buf)
    
    def add(self, new_entry: Any) -> None:
        """
        Add element to sequence.
        
        Args:
            new_entry: The element to add
        """
        self.elems.append(new_entry)
        self.elems_len += 1
    
    def free(self) -> None:
        """Free sequence and all elements."""
        if self.elem_free:
            for elem in self.elems:
                if elem is not None:
                    self.elem_free(elem)
        self.elems = []
        self.elems_len = 0


class SmlList:
    """
    SML List entry (linked list structure).
    
    Attributes:
        obj_name: Object name (OctetString)
        status: Status (optional, SmlStatus)
        val_time: Value time (optional, SmlTime)
        unit: Unit (optional, bytes for u8)
        scaler: Scaler (optional, bytes for i8)
        value: Value (SmlValue)
        value_signature: Value signature (optional, OctetString)
        next: Next list entry (optional, SmlList)
    """
    
    def __init__(self):
        """Initialize a new SmlList entry."""
        self.obj_name: Optional[OctetString] = None
        self.status: Optional[SmlStatus] = None
        self.val_time: Optional[SmlTime] = None
        self.unit: Optional[bytes] = None
        self.scaler: Optional[bytes] = None
        self.value: Optional[SmlValue] = None
        self.value_signature: Optional[OctetString] = None
        self.next: Optional['SmlList'] = None
    
    @classmethod
    def parse(cls, buf: SmlBuffer) -> Optional['SmlList']:
        """
        Parse list from buffer.
        
        Args:
            buf: The buffer to parse from
            
        Returns:
            SmlList instance (linked list), or None if optional and skipped
        """
        if buf.optional_is_skipped() == 1:
            return None
        
        if buf.get_next_type() != SML_TYPE_LIST:
            buf.error = 1
            return None
        
        elems = buf.get_next_length()
        
        # Workarounds for DZG meters
        workarounds = {'old_dzg_meter': False}
        
        ret = None
        pos = None
        
        while elems > 0:
            entry = cls._parse_entry(buf, workarounds)
            if buf.has_errors():
                if ret:
                    ret.free()
                return None
            
            if entry:
                if ret is None:
                    ret = entry
                    pos = entry
                else:
                    pos.next = entry
                    pos = entry
            
            elems -= 1
        
        return ret
    
    @classmethod
    def _parse_entry(cls, buf: SmlBuffer, workarounds: dict) -> Optional['SmlList']:
        """
        Parse a single list entry.
        
        Args:
            buf: The buffer to parse from
            workarounds: Workarounds dictionary for compatibility
            
        Returns:
            SmlList entry, or None on error
        """
        # DZG meter workaround constants
        dzg_serial_name = bytes([1, 0, 96, 1, 0, 255])
        dzg_serial_start = bytes([0x0a, 0x01, ord('D'), ord('Z'), ord('G'), 0x00])
        dzg_serial_fixed = bytes([0x0a, 0x01, ord('D'), ord('Z'), ord('G'), 0x00, 0x03, 0x93, 0x87, 0x00])
        dzg_power_name = bytes([1, 0, 16, 7, 0, 255])
        
        if buf.get_next_type() != SML_TYPE_LIST:
            buf.error = 1
            return None
        
        if buf.get_next_length() != 7:
            buf.error = 1
            return None
        
        entry = cls()
        
        entry.obj_name = octet_string_parse(buf)
        if buf.has_errors():
            return None
        
        entry.status = status_parse(buf)
        if buf.has_errors():
            return None
        
        entry.val_time = time_parse(buf)
        if buf.has_errors():
            return None
        
        entry.unit = u8_parse(buf)
        if buf.has_errors():
            return None
        
        entry.scaler = i8_parse(buf)
        if buf.has_errors():
            return None
        
        if buf.cursor >= buf.buffer_len:
            return None
        
        value_tl = buf.get_current_byte()
        value_len_more = value_tl & (SML_ANOTHER_TL | SML_LENGTH_FIELD)
        
        entry.value = value_parse(buf)
        if buf.has_errors():
            return None
        
        entry.value_signature = octet_string_parse(buf)
        if buf.has_errors():
            return None
        
        # DZG meter workaround
        if (entry.obj_name and entry.obj_name.len == len(dzg_serial_name) and
            entry.obj_name.str == dzg_serial_name and
            entry.value and entry.value.type == SML_TYPE_OCTET_STRING and
            isinstance(entry.value.data, OctetString) and
            entry.value.data.len >= len(dzg_serial_start) and
            entry.value.data.str[:len(dzg_serial_start)] == dzg_serial_start and
            entry.value.data.str < dzg_serial_fixed):
            workarounds['old_dzg_meter'] = True
        elif (workarounds.get('old_dzg_meter') and
              entry.obj_name and entry.obj_name.len == len(dzg_power_name) and
              entry.obj_name.str == dzg_power_name and
              entry.value and
              (value_len_more == 1 or value_len_more == 2 or value_len_more == 3)):
            # Fix value type from signed to unsigned
            entry.value.type &= ~SML_TYPE_FIELD
            entry.value.type |= SML_TYPE_UNSIGNED
        
        return entry
    
    def write(self, buf: SmlBuffer) -> None:
        """
        Write list to buffer.
        
        Args:
            buf: The buffer to write to
        """
        if self.value is None:
            buf.optional_write()
            return
        
        # Count entries
        count = 0
        current = self
        while current:
            count += 1
            current = current.next
        
        buf.set_type_and_length(SML_TYPE_LIST, count)
        
        # Write all entries
        current = self
        while current:
            current._write_entry(buf)
            current = current.next
    
    def _write_entry(self, buf: SmlBuffer) -> None:
        """Write a single list entry."""
        buf.set_type_and_length(SML_TYPE_LIST, 7)
        octet_string_write(self.obj_name, buf)
        status_write(self.status, buf)
        time_write(self.val_time, buf)
        u8_write(self.unit, buf)
        i8_write(self.scaler, buf)
        value_write(self.value, buf)
        octet_string_write(self.value_signature, buf)
    
    def add(self, new_entry: 'SmlList') -> None:
        """
        Add entry to list (sets as next).
        
        Args:
            new_entry: The entry to add
        """
        self.next = new_entry
    
    def free(self) -> None:
        """Free list and all entries."""
        current = self
        while current:
            next_entry = current.next
            # Free fields
            if current.obj_name:
                current.obj_name.free()
            if current.status:
                current.status.free()
            if current.val_time:
                current.val_time.free()
            if current.value:
                current.value.free()
            if current.value_signature:
                current.value_signature.free()
            current = next_entry


def list_init() -> SmlList:
    """
    Initialize a new SmlList.
    
    Returns:
        New SmlList instance
    """
    return SmlList()


def list_parse(buf: SmlBuffer) -> Optional[SmlList]:
    """
    Parse list from buffer.
    
    Args:
        buf: The buffer to parse from
        
    Returns:
        SmlList instance (linked list), or None if optional and skipped
    """
    return SmlList.parse(buf)


def list_write(list_obj: Optional[SmlList], buf: SmlBuffer) -> None:
    """
    Write list to buffer.
    
    Args:
        list_obj: The list to write, or None
        buf: The buffer to write to
    """
    if list_obj is None:
        buf.optional_write()
    else:
        list_obj.write(buf)


def list_add(list_obj: SmlList, new_entry: SmlList) -> None:
    """
    Add entry to list.
    
    Args:
        list_obj: The list to add to
        new_entry: The entry to add
    """
    list_obj.add(new_entry)


def list_free(list_obj: Optional[SmlList]) -> None:
    """
    Free list.
    
    Args:
        list_obj: The list to free
    """
    if list_obj:
        list_obj.free()


# Sequence functions
def sequence_init(elem_free: Optional[Callable] = None) -> SmlSequence:
    """
    Initialize a new SmlSequence.
    
    Args:
        elem_free: Optional function to free elements
        
    Returns:
        New SmlSequence instance
    """
    return SmlSequence(elem_free)


def sequence_parse(buf: SmlBuffer, elem_parse: Callable, elem_free: Optional[Callable] = None) -> Optional[SmlSequence]:
    """
    Parse sequence from buffer.
    
    Args:
        buf: The buffer to parse from
        elem_parse: Function to parse individual elements
        elem_free: Optional function to free elements
        
    Returns:
        SmlSequence instance, or None on error
    """
    return SmlSequence.parse(buf, elem_parse, elem_free)


def sequence_write(seq: Optional[SmlSequence], buf: SmlBuffer, elem_write: Callable) -> None:
    """
    Write sequence to buffer.
    
    Args:
        seq: The sequence to write, or None
        buf: The buffer to write to
        elem_write: Function to write individual elements
    """
    if seq is None:
        buf.optional_write()
    else:
        seq.write(buf, elem_write)


def sequence_add(seq: SmlSequence, new_entry: Any) -> None:
    """
    Add element to sequence.
    
    Args:
        seq: The sequence
        new_entry: The element to add
    """
    seq.add(new_entry)


def sequence_free(seq: Optional[SmlSequence]) -> None:
    """
    Free sequence.
    
    Args:
        seq: The sequence to free
    """
    if seq:
        seq.free()

