"""Basic tests for libsml Python implementation."""

import unittest
import sys
import os

# Add parent directory to path
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

import libsml


class TestBasic(unittest.TestCase):
    """Basic functionality tests."""
    
    def test_import(self):
        """Test that the library can be imported."""
        self.assertTrue(hasattr(libsml, 'SmlBuffer'))
        self.assertTrue(hasattr(libsml, 'SmlFile'))
        self.assertTrue(hasattr(libsml, 'SmlMessage'))
    
    def test_buffer_init(self):
        """Test buffer initialization."""
        buf = libsml.SmlBuffer(512)
        self.assertIsNotNone(buf)
        self.assertEqual(buf.buffer_len, 512)
        self.assertEqual(buf.cursor, 0)
        self.assertEqual(buf.error, 0)
    
    def test_buffer_from_bytes(self):
        """Test creating buffer from bytes."""
        data = b'\x01\x02\x03\x04'
        buf = libsml.SmlBuffer.from_bytes(data)
        self.assertEqual(buf.buffer_len, 4)
        self.assertEqual(buf.cursor, 0)
        self.assertEqual(bytes(buf.buffer[:4]), data)
    
    def test_octet_string_init(self):
        """Test octet string initialization."""
        data = b'hello'
        octet = libsml.OctetString(data, len(data))
        self.assertEqual(octet.str, data)
        self.assertEqual(octet.len, 5)
    
    def test_octet_string_from_hex(self):
        """Test octet string from hex."""
        hex_str = "48656c6c6f"  # "Hello"
        octet = libsml.OctetString.init_from_hex(hex_str)
        self.assertEqual(octet.str, b'Hello')
        self.assertEqual(octet.len, 5)
    
    def test_uuid_generation(self):
        """Test UUID generation."""
        uuid1 = libsml.octet_string_generate_uuid()
        uuid2 = libsml.octet_string_generate_uuid()
        self.assertEqual(uuid1.len, 16)
        self.assertEqual(uuid2.len, 16)
        self.assertNotEqual(uuid1.str, uuid2.str)  # Should be different
    
    def test_crc16(self):
        """Test CRC16 calculation."""
        from libsml.crc16 import crc16_calculate
        data = b'\x01\x02\x03\x04'
        crc = crc16_calculate(data, len(data))
        self.assertIsInstance(crc, int)
        self.assertGreaterEqual(crc, 0)
        self.assertLessEqual(crc, 0xFFFF)
    
    def test_number_init(self):
        """Test number initialization."""
        from libsml.number import u8_init, u16_init, u32_init, u64_init
        from libsml.number import i8_init, i16_init, i32_init, i64_init
        
        u8_val = u8_init(42)
        self.assertEqual(len(u8_val), 1)
        
        u16_val = u16_init(1000)
        self.assertEqual(len(u16_val), 2)
        
        u32_val = u32_init(100000)
        self.assertEqual(len(u32_val), 4)
        
        u64_val = u64_init(1000000000)
        self.assertEqual(len(u64_val), 8)
        
        i8_val = i8_init(-42)
        self.assertEqual(len(i8_val), 1)
        
        i16_val = i16_init(-1000)
        self.assertEqual(len(i16_val), 2)
        
        i32_val = i32_init(-100000)
        self.assertEqual(len(i32_val), 4)
        
        i64_val = i64_init(-1000000000)
        self.assertEqual(len(i64_val), 8)
    
    def test_message_init(self):
        """Test message initialization."""
        msg = libsml.message_init()
        self.assertIsNotNone(msg)
        self.assertIsNotNone(msg.transaction_id)
        self.assertEqual(msg.transaction_id.len, 16)  # UUID is 16 bytes
    
    def test_file_init(self):
        """Test file initialization."""
        file = libsml.file_init()
        self.assertIsNotNone(file)
        self.assertEqual(file.messages_len, 0)
        self.assertIsNotNone(file.buf)


if __name__ == '__main__':
    unittest.main()

