"""Simple import test."""

import sys
import os

# Add current directory to path
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

try:
    import libsml
    print("✓ Import libsml successful")
    
    # Test basic classes
    print("Testing basic classes...")
    
    # Test SmlBuffer
    buf = libsml.SmlBuffer(512)
    print(f"✓ SmlBuffer created: len={buf.buffer_len}, cursor={buf.cursor}")
    
    # Test OctetString
    octet = libsml.OctetString(b"test", 4)
    print(f"✓ OctetString created: len={octet.len}, data={octet.str}")
    
    # Test UUID generation
    uuid = libsml.octet_string_generate_uuid()
    print(f"✓ UUID generated: len={uuid.len}")
    
    # Test CRC16
    from libsml.crc16 import crc16_calculate
    data = b'\x01\x02\x03\x04'
    crc = crc16_calculate(data, len(data))
    print(f"✓ CRC16 calculated: {crc:04X}")
    
    # Test number init
    from libsml.number import u8_init, u16_init, u32_init
    u8 = u8_init(42)
    u16 = u16_init(1000)
    u32 = u32_init(100000)
    print(f"✓ Numbers initialized: u8={len(u8)} bytes, u16={len(u16)} bytes, u32={len(u32)} bytes")
    
    # Test message init
    msg = libsml.message_init()
    print(f"✓ Message initialized: transaction_id len={msg.transaction_id.len if msg.transaction_id else 0}")
    
    # Test file init
    file = libsml.file_init()
    print(f"✓ File initialized: messages={file.messages_len}")
    
    print("\nAll basic tests passed!")
    
except Exception as e:
    print(f"✗ Error: {e}")
    import traceback
    traceback.print_exc()
    sys.exit(1)

