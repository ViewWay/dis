"""Tests for SmlBuffer."""

import unittest
import sys
import os

sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

import libsml
from libsml.shared import SML_TYPE_LIST, SML_TYPE_OCTET_STRING, SML_OPTIONAL_SKIPPED


class TestBuffer(unittest.TestCase):
    """Tests for SmlBuffer class."""
    
    def test_get_next_type(self):
        """Test getting next type from buffer."""
        # Create buffer with type byte
        data = bytes([SML_TYPE_OCTET_STRING | 0x05])  # Type + length
        buf = libsml.SmlBuffer.from_bytes(data)
        type_val = buf.get_next_type()
        self.assertEqual(type_val, SML_TYPE_OCTET_STRING)
    
    def test_get_current_byte(self):
        """Test getting current byte."""
        data = b'\x42'
        buf = libsml.SmlBuffer.from_bytes(data)
        byte = buf.get_current_byte()
        self.assertEqual(byte, 0x42)
    
    def test_update_bytes_read(self):
        """Test updating bytes read."""
        buf = libsml.SmlBuffer.from_bytes(b'\x01\x02\x03')
        self.assertEqual(buf.cursor, 0)
        buf.update_bytes_read(2)
        self.assertEqual(buf.cursor, 2)
    
    def test_optional_is_skipped(self):
        """Test optional skipped check."""
        data = bytes([SML_OPTIONAL_SKIPPED])
        buf = libsml.SmlBuffer.from_bytes(data)
        result = buf.optional_is_skipped()
        self.assertEqual(result, 1)
        self.assertEqual(buf.cursor, 1)
    
    def test_optional_write(self):
        """Test writing optional skipped marker."""
        buf = libsml.SmlBuffer(10)
        buf.optional_write()
        self.assertEqual(buf.cursor, 1)
        self.assertEqual(buf.buffer[0], SML_OPTIONAL_SKIPPED)
    
    def test_to_bytes(self):
        """Test converting buffer to bytes."""
        buf = libsml.SmlBuffer(10)
        buf.buffer[0] = 0x01
        buf.buffer[1] = 0x02
        buf.cursor = 2
        result = buf.to_bytes()
        self.assertEqual(result, b'\x01\x02')


if __name__ == '__main__':
    unittest.main()

