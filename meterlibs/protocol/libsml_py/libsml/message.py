"""
SML Message - Message encapsulation

This module implements SML message parsing and writing, including
message body handling for all message types.
"""

from typing import Optional, Any
from .buffer import SmlBuffer
from .octet_string import OctetString, octet_string_parse, octet_string_write, octet_string_generate_uuid
from .number import u8_parse, u16_parse, u32_parse, u8_write, u16_write, u32_write, u32_init
from .crc16 import crc16_calculate, crc16kermit_calculate
from .shared import SML_TYPE_LIST, SML_MESSAGE_END

# Message type constants
SML_MESSAGE_OPEN_REQUEST = 0x00000100
SML_MESSAGE_OPEN_RESPONSE = 0x00000101
SML_MESSAGE_CLOSE_REQUEST = 0x00000200
SML_MESSAGE_CLOSE_RESPONSE = 0x00000201
SML_MESSAGE_GET_PROFILE_PACK_REQUEST = 0x00000300
SML_MESSAGE_GET_PROFILE_PACK_RESPONSE = 0x00000301
SML_MESSAGE_GET_PROFILE_LIST_REQUEST = 0x00000400
SML_MESSAGE_GET_PROFILE_LIST_RESPONSE = 0x00000401
SML_MESSAGE_GET_PROC_PARAMETER_REQUEST = 0x00000500
SML_MESSAGE_GET_PROC_PARAMETER_RESPONSE = 0x00000501
SML_MESSAGE_SET_PROC_PARAMETER_REQUEST = 0x00000600
SML_MESSAGE_SET_PROC_PARAMETER_RESPONSE = 0x00000601  # This doesn't exist in the spec
SML_MESSAGE_GET_LIST_REQUEST = 0x00000700
SML_MESSAGE_GET_LIST_RESPONSE = 0x00000701
SML_MESSAGE_ATTENTION_RESPONSE = 0x0000FF01


class SmlMessageBody:
    """
    SML Message Body class.
    
    Attributes:
        tag: Message type tag (u32)
        data: Message data (type depends on tag)
    """
    
    def __init__(self, tag: int = 0, data: Any = None):
        """
        Initialize a new SmlMessageBody.
        
        Args:
            tag: Message type tag
            data: Message data
        """
        self.tag: Optional[bytes] = None
        self.data: Any = data
        if tag != 0:
            self.tag = u32_init(tag)
    
    @classmethod
    def parse(cls, buf: SmlBuffer) -> Optional['SmlMessageBody']:
        """
        Parse message body from buffer.
        
        Args:
            buf: The buffer to parse from
            
        Returns:
            SmlMessageBody instance, or None on error
        """
        if (buf.cursor + 1) > buf.buffer_len:
            buf.error = 1
            return None
        
        if buf.get_next_type() != SML_TYPE_LIST:
            buf.error = 1
            return None
        
        if buf.get_next_length() != 2:
            buf.error = 1
            return None
        
        tag = u32_parse(buf)
        if buf.has_errors() or tag is None:
            return None
        
        tag_value = int.from_bytes(tag, 'big')
        
        # Import message parsers dynamically to avoid circular imports
        from .messages import (
            open_request, open_response,
            close_request, close_response,
            get_list_request, get_list_response,
            get_profile_pack_request, get_profile_pack_response,
            get_profile_list_response,
            get_proc_parameter_request, get_proc_parameter_response,
            set_proc_parameter_request,
            attention_response,
        )
        
        msg_body = cls()
        msg_body.tag = tag
        
        if tag_value == SML_MESSAGE_OPEN_REQUEST:
            msg_body.data = open_request.SmlOpenRequest.parse(buf)
        elif tag_value == SML_MESSAGE_OPEN_RESPONSE:
            msg_body.data = open_response.SmlOpenResponse.parse(buf)
        elif tag_value == SML_MESSAGE_CLOSE_REQUEST:
            msg_body.data = close_request.SmlCloseRequest.parse(buf)
        elif tag_value == SML_MESSAGE_CLOSE_RESPONSE:
            msg_body.data = close_response.SmlCloseResponse.parse(buf)
        elif tag_value == SML_MESSAGE_GET_PROFILE_PACK_REQUEST:
            msg_body.data = get_profile_pack_request.SmlGetProfilePackRequest.parse(buf)
        elif tag_value == SML_MESSAGE_GET_PROFILE_PACK_RESPONSE:
            msg_body.data = get_profile_pack_response.SmlGetProfilePackResponse.parse(buf)
        elif tag_value == SML_MESSAGE_GET_PROFILE_LIST_REQUEST:
            # GetProfileListRequest is the same as GetProfilePackRequest
            msg_body.data = get_profile_pack_request.SmlGetProfilePackRequest.parse(buf)
        elif tag_value == SML_MESSAGE_GET_PROFILE_LIST_RESPONSE:
            msg_body.data = get_profile_list_response.SmlGetProfileListResponse.parse(buf)
        elif tag_value == SML_MESSAGE_GET_PROC_PARAMETER_REQUEST:
            msg_body.data = get_proc_parameter_request.SmlGetProcParameterRequest.parse(buf)
        elif tag_value == SML_MESSAGE_GET_PROC_PARAMETER_RESPONSE:
            msg_body.data = get_proc_parameter_response.SmlGetProcParameterResponse.parse(buf)
        elif tag_value == SML_MESSAGE_SET_PROC_PARAMETER_REQUEST:
            msg_body.data = set_proc_parameter_request.SmlSetProcParameterRequest.parse(buf)
        elif tag_value == SML_MESSAGE_GET_LIST_REQUEST:
            msg_body.data = get_list_request.SmlGetListRequest.parse(buf)
        elif tag_value == SML_MESSAGE_GET_LIST_RESPONSE:
            msg_body.data = get_list_response.SmlGetListResponse.parse(buf)
        elif tag_value == SML_MESSAGE_ATTENTION_RESPONSE:
            msg_body.data = attention_response.SmlAttentionResponse.parse(buf)
        else:
            print(f"libsml: error: message type {tag_value:04X} not yet implemented")
            return None
        
        return msg_body
    
    def write(self, buf: SmlBuffer) -> None:
        """
        Write message body to buffer.
        
        Args:
            buf: The buffer to write to
        """
        if self.tag is None:
            return
        
        buf.set_type_and_length(SML_TYPE_LIST, 2)
        u32_write(self.tag, buf)
        
        tag_value = int.from_bytes(self.tag, 'big')
        
        # Import message writers dynamically
        from .messages import (
            open_request, open_response,
            close_request, close_response,
            get_list_request, get_list_response,
            get_profile_pack_request, get_profile_pack_response,
            get_profile_list_response,
            get_proc_parameter_request, get_proc_parameter_response,
            set_proc_parameter_request,
            attention_response,
        )
        
        if tag_value == SML_MESSAGE_OPEN_REQUEST:
            if isinstance(self.data, open_request.SmlOpenRequest):
                self.data.write(buf)
        elif tag_value == SML_MESSAGE_OPEN_RESPONSE:
            if isinstance(self.data, open_response.SmlOpenResponse):
                self.data.write(buf)
        elif tag_value == SML_MESSAGE_CLOSE_REQUEST:
            if isinstance(self.data, close_request.SmlCloseRequest):
                self.data.write(buf)
        elif tag_value == SML_MESSAGE_CLOSE_RESPONSE:
            if isinstance(self.data, close_response.SmlCloseResponse):
                self.data.write(buf)
        elif tag_value == SML_MESSAGE_GET_PROFILE_PACK_REQUEST:
            if isinstance(self.data, get_profile_pack_request.SmlGetProfilePackRequest):
                self.data.write(buf)
        elif tag_value == SML_MESSAGE_GET_PROFILE_PACK_RESPONSE:
            if isinstance(self.data, get_profile_pack_response.SmlGetProfilePackResponse):
                self.data.write(buf)
        elif tag_value == SML_MESSAGE_GET_PROFILE_LIST_REQUEST:
            # GetProfileListRequest is the same as GetProfilePackRequest
            if isinstance(self.data, get_profile_pack_request.SmlGetProfilePackRequest):
                self.data.write(buf)
        elif tag_value == SML_MESSAGE_GET_PROFILE_LIST_RESPONSE:
            if isinstance(self.data, get_profile_list_response.SmlGetProfileListResponse):
                self.data.write(buf)
        elif tag_value == SML_MESSAGE_GET_PROC_PARAMETER_REQUEST:
            if isinstance(self.data, get_proc_parameter_request.SmlGetProcParameterRequest):
                self.data.write(buf)
        elif tag_value == SML_MESSAGE_GET_PROC_PARAMETER_RESPONSE:
            if isinstance(self.data, get_proc_parameter_response.SmlGetProcParameterResponse):
                self.data.write(buf)
        elif tag_value == SML_MESSAGE_SET_PROC_PARAMETER_REQUEST:
            if isinstance(self.data, set_proc_parameter_request.SmlSetProcParameterRequest):
                self.data.write(buf)
        elif tag_value == SML_MESSAGE_GET_LIST_REQUEST:
            if isinstance(self.data, get_list_request.SmlGetListRequest):
                self.data.write(buf)
        elif tag_value == SML_MESSAGE_GET_LIST_RESPONSE:
            if isinstance(self.data, get_list_response.SmlGetListResponse):
                self.data.write(buf)
        elif tag_value == SML_MESSAGE_ATTENTION_RESPONSE:
            if isinstance(self.data, attention_response.SmlAttentionResponse):
                self.data.write(buf)
        else:
            print(f"libsml: error: message type {tag_value:04X} not yet implemented")
    
    def free(self) -> None:
        """Free message body."""
        pass


class SmlMessage:
    """
    SML Message class.
    
    Attributes:
        transaction_id: Transaction ID (OctetString)
        group_id: Group ID (optional, u8)
        abort_on_error: Abort on error flag (optional, u8)
        message_body: Message body (SmlMessageBody)
        crc: CRC16 checksum (u16)
    """
    
    def __init__(self):
        """Initialize a new SmlMessage."""
        self.transaction_id: Optional[OctetString] = None
        self.group_id: Optional[bytes] = None  # u8
        self.abort_on_error: Optional[bytes] = None  # u8
        self.message_body: Optional[SmlMessageBody] = None
        self.crc: Optional[bytes] = None  # u16
    
    @classmethod
    def init(cls) -> 'SmlMessage':
        """
        Initialize a new SmlMessage with transaction ID.
        
        Returns:
            New SmlMessage instance
        """
        msg = cls()
        msg.transaction_id = octet_string_generate_uuid()
        return msg
    
    @classmethod
    def parse(cls, buf: SmlBuffer) -> Optional['SmlMessage']:
        """
        Parse message from buffer.
        
        Args:
            buf: The buffer to parse from
            
        Returns:
            SmlMessage instance, or None on error
        """
        msg = cls()
        msg_start = buf.cursor
        
        if buf.get_next_type() != SML_TYPE_LIST:
            buf.error = 1
            return None
        
        if buf.get_next_length() != 6:
            buf.error = 1
            return None
        
        msg.transaction_id = octet_string_parse(buf)
        if buf.has_errors():
            return None
        
        msg.group_id = u8_parse(buf)
        if buf.has_errors():
            return None
        
        msg.abort_on_error = u8_parse(buf)
        if buf.has_errors():
            return None
        
        msg.message_body = SmlMessageBody.parse(buf)
        if buf.has_errors():
            return None
        
        length = buf.cursor - msg_start
        if (buf.buffer_len - buf.cursor) < 3:
            buf.error = 1
            return None
        
        msg.crc = u16_parse(buf)
        if buf.has_errors() or msg.crc is None:
            buf.error = 1
            return None
        
        # Verify CRC
        crc_value = int.from_bytes(msg.crc, 'big')
        calculated_crc = crc16_calculate(buf.buffer[msg_start:buf.cursor], length)
        
        if crc_value != calculated_crc:
            # Try CRC-16/Kermit as workaround for Holley DTZ541
            calculated_crc_kermit = crc16kermit_calculate(buf.buffer[msg_start:buf.cursor], length)
            if crc_value != calculated_crc_kermit:
                return None
        
        if buf.cursor >= buf.buffer_len:
            buf.error = 1
            return None
        
        if buf.get_current_byte() == SML_MESSAGE_END:
            buf.update_bytes_read(1)
        
        return msg
    
    def write(self, buf: SmlBuffer) -> None:
        """
        Write message to buffer.
        
        Args:
            buf: The buffer to write to
        """
        msg_start = buf.cursor
        
        buf.set_type_and_length(SML_TYPE_LIST, 6)
        octet_string_write(self.transaction_id, buf)
        u8_write(self.group_id, buf)
        u8_write(self.abort_on_error, buf)
        
        if self.message_body:
            self.message_body.write(buf)
        
        # Calculate and write CRC
        length = buf.cursor - msg_start
        crc = crc16_calculate(buf.buffer[msg_start:buf.cursor], length)
        self.crc = u16_init(crc)
        u16_write(self.crc, buf)
        
        # End of message
        buf.ensure_capacity(1)
        buf.buffer[buf.cursor] = SML_MESSAGE_END
        buf.cursor += 1
    
    def free(self) -> None:
        """Free message."""
        pass


# u32_init is already imported from .number


# Free functions for API compatibility
def message_init() -> SmlMessage:
    """Initialize a new SmlMessage."""
    return SmlMessage.init()


def message_parse(buf: SmlBuffer) -> Optional[SmlMessage]:
    """Parse message from buffer."""
    return SmlMessage.parse(buf)


def message_write(msg: SmlMessage, buf: SmlBuffer) -> None:
    """Write message to buffer."""
    msg.write(buf)


def message_free(msg: Optional[SmlMessage]) -> None:
    """Free message."""
    pass


def message_body_init(tag: int, data: Any) -> SmlMessageBody:
    """Initialize a new SmlMessageBody."""
    return SmlMessageBody(tag, data)


def message_body_parse(buf: SmlBuffer) -> Optional[SmlMessageBody]:
    """Parse message body from buffer."""
    return SmlMessageBody.parse(buf)


def message_body_write(message_body: SmlMessageBody, buf: SmlBuffer) -> None:
    """Write message body to buffer."""
    message_body.write(buf)


def message_body_free(message_body: Optional[SmlMessageBody]) -> None:
    """Free message body."""
    pass

