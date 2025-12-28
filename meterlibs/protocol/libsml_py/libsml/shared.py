"""
SML Shared - Base types and constants

This module defines the base types and constants used throughout the SML library.
"""

from typing import TypeAlias

# Type aliases matching C typedefs
u8: TypeAlias = int  # uint8_t
u16: TypeAlias = int  # uint16_t
u32: TypeAlias = int  # uint32_t
u64: TypeAlias = int  # uint64_t

i8: TypeAlias = int  # int8_t
i16: TypeAlias = int  # int16_t
i32: TypeAlias = int  # int32_t
i64: TypeAlias = int  # int64_t

# Message end marker
SML_MESSAGE_END = 0x0

# Type and length field masks
SML_TYPE_FIELD = 0x70
SML_LENGTH_FIELD = 0xF
SML_ANOTHER_TL = 0x80

# SML type constants
SML_TYPE_OCTET_STRING = 0x0
SML_TYPE_BOOLEAN = 0x40
SML_TYPE_INTEGER = 0x50
SML_TYPE_UNSIGNED = 0x60
SML_TYPE_LIST = 0x70

# Optional field marker
SML_OPTIONAL_SKIPPED = 0x1

# Number type sizes (in bytes)
SML_TYPE_NUMBER_8 = 1  # sizeof(u8)
SML_TYPE_NUMBER_16 = 2  # sizeof(u16)
SML_TYPE_NUMBER_32 = 4  # sizeof(u32)
SML_TYPE_NUMBER_64 = 8  # sizeof(u64)

