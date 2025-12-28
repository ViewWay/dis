"""
libSML - Python implementation of Smart Message Language (SML) protocol

This library implements the SML protocol specified by VDE's Forum Netztechnik/Netzbetrieb (FNN).
It can be utilized to communicate to FNN specified Smart Meters or Smart Meter components (EDL/MUC).
"""

__version__ = "1.0.0"

# Core classes
from .buffer import SmlBuffer
from .file import SmlFile, file_parse, file_init, file_add_message, file_write, file_free, file_print
from .message import (
    SmlMessage, SmlMessageBody,
    message_init, message_parse, message_write, message_free,
    message_body_init, message_body_parse, message_body_write, message_body_free,
    SML_MESSAGE_OPEN_REQUEST, SML_MESSAGE_OPEN_RESPONSE,
    SML_MESSAGE_CLOSE_REQUEST, SML_MESSAGE_CLOSE_RESPONSE,
    SML_MESSAGE_GET_PROFILE_PACK_REQUEST, SML_MESSAGE_GET_PROFILE_PACK_RESPONSE,
    SML_MESSAGE_GET_PROFILE_LIST_REQUEST, SML_MESSAGE_GET_PROFILE_LIST_RESPONSE,
    SML_MESSAGE_GET_PROC_PARAMETER_REQUEST, SML_MESSAGE_GET_PROC_PARAMETER_RESPONSE,
    SML_MESSAGE_SET_PROC_PARAMETER_REQUEST,
    SML_MESSAGE_GET_LIST_REQUEST, SML_MESSAGE_GET_LIST_RESPONSE,
    SML_MESSAGE_ATTENTION_RESPONSE,
)

# Data types
from .octet_string import OctetString, octet_string_init, octet_string_parse, octet_string_write, octet_string_generate_uuid
from .value import SmlValue, value_init, value_parse, value_write, value_to_double, value_to_strhex
from .list import SmlList, SmlSequence, list_init, list_parse, list_write, list_free
from .time import SmlTime, time_init, time_parse, time_write, time_free
from .status import SmlStatus, status_init, status_parse, status_write, status_free
from .tree import (
    SmlTree, SmlTreePath, SmlProcParValue,
    tree_init, tree_parse, tree_write, tree_free,
    tree_path_init, tree_path_parse, tree_path_write, tree_path_free,
)

# Transport
from .transport import transport_read, transport_listen, transport_write

# CRC
from .crc16 import crc16_calculate, crc16kermit_calculate

# Constants
from .shared import (
    SML_TYPE_OCTET_STRING, SML_TYPE_BOOLEAN, SML_TYPE_INTEGER,
    SML_TYPE_UNSIGNED, SML_TYPE_LIST,
    SML_OPTIONAL_SKIPPED, SML_MESSAGE_END,
)

__all__ = [
    # Core
    'SmlBuffer', 'SmlFile', 'SmlMessage', 'SmlMessageBody',
    'file_parse', 'file_init', 'file_add_message', 'file_write', 'file_free', 'file_print',
    'message_init', 'message_parse', 'message_write', 'message_free',
    # Data types
    'OctetString', 'SmlValue', 'SmlList', 'SmlSequence', 'SmlTime', 'SmlStatus',
    'SmlTree', 'SmlTreePath', 'SmlProcParValue',
    # Transport
    'transport_read', 'transport_listen', 'transport_write',
    # CRC
    'crc16_calculate', 'crc16kermit_calculate',
    # Constants
    'SML_MESSAGE_OPEN_REQUEST', 'SML_MESSAGE_OPEN_RESPONSE',
    'SML_MESSAGE_CLOSE_REQUEST', 'SML_MESSAGE_CLOSE_RESPONSE',
    'SML_MESSAGE_GET_LIST_REQUEST', 'SML_MESSAGE_GET_LIST_RESPONSE',
    'SML_MESSAGE_GET_PROFILE_PACK_REQUEST', 'SML_MESSAGE_GET_PROFILE_PACK_RESPONSE',
    'SML_MESSAGE_GET_PROFILE_LIST_REQUEST', 'SML_MESSAGE_GET_PROFILE_LIST_RESPONSE',
    'SML_MESSAGE_GET_PROC_PARAMETER_REQUEST', 'SML_MESSAGE_GET_PROC_PARAMETER_RESPONSE',
    'SML_MESSAGE_SET_PROC_PARAMETER_REQUEST',
    'SML_MESSAGE_ATTENTION_RESPONSE',
]
