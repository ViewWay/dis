"""Run all tests."""

import sys
import os

# Add current directory to path
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

print("=" * 60)
print("Running libsml Python Implementation Tests")
print("=" * 60)

# Test 1: Import
print("\n[Test 1] Testing imports...")
try:
    import libsml
    print("✓ libsml imported successfully")
except Exception as e:
    print(f"✗ Import failed: {e}")
    sys.exit(1)

# Test 2: Basic classes
print("\n[Test 2] Testing basic classes...")
try:
    buf = libsml.SmlBuffer(512)
    assert buf.buffer_len == 512
    assert buf.cursor == 0
    print("✓ SmlBuffer created")
    
    octet = libsml.OctetString(b"test", 4)
    assert octet.len == 4
    print("✓ OctetString created")
    
    uuid = libsml.octet_string_generate_uuid()
    assert uuid.len == 16
    print("✓ UUID generated")
except Exception as e:
    print(f"✗ Basic classes test failed: {e}")
    import traceback
    traceback.print_exc()
    sys.exit(1)

# Test 3: CRC16
print("\n[Test 3] Testing CRC16...")
try:
    from libsml.crc16 import crc16_calculate, crc16kermit_calculate
    data = b'\x01\x02\x03\x04'
    crc = crc16_calculate(data, len(data))
    assert 0 <= crc <= 0xFFFF
    print(f"✓ CRC16 calculated: {crc:04X}")
    
    crc_kermit = crc16kermit_calculate(data, len(data))
    assert 0 <= crc_kermit <= 0xFFFF
    print(f"✓ CRC-16/Kermit calculated: {crc_kermit:04X}")
except Exception as e:
    print(f"✗ CRC16 test failed: {e}")
    import traceback
    traceback.print_exc()
    sys.exit(1)

# Test 4: Numbers
print("\n[Test 4] Testing number types...")
try:
    from libsml.number import u8_init, u16_init, u32_init, u64_init
    from libsml.number import i8_init, i16_init, i32_init, i64_init
    
    u8 = u8_init(42)
    assert len(u8) == 1
    print("✓ u8 initialized")
    
    u16 = u16_init(1000)
    assert len(u16) == 2
    print("✓ u16 initialized")
    
    u32 = u32_init(100000)
    assert len(u32) == 4
    print("✓ u32 initialized")
    
    u64 = u64_init(1000000000)
    assert len(u64) == 8
    print("✓ u64 initialized")
    
    i8 = i8_init(-42)
    assert len(i8) == 1
    print("✓ i8 initialized")
    
    i16 = i16_init(-1000)
    assert len(i16) == 2
    print("✓ i16 initialized")
    
    i32 = i32_init(-100000)
    assert len(i32) == 4
    print("✓ i32 initialized")
    
    i64 = i64_init(-1000000000)
    assert len(i64) == 8
    print("✓ i64 initialized")
except Exception as e:
    print(f"✗ Number types test failed: {e}")
    import traceback
    traceback.print_exc()
    sys.exit(1)

# Test 5: Message
print("\n[Test 5] Testing message initialization...")
try:
    msg = libsml.message_init()
    assert msg.transaction_id is not None
    assert msg.transaction_id.len == 16
    print("✓ Message initialized with transaction ID")
except Exception as e:
    print(f"✗ Message test failed: {e}")
    import traceback
    traceback.print_exc()
    sys.exit(1)

# Test 6: File
print("\n[Test 6] Testing file initialization...")
try:
    file = libsml.file_init()
    assert file.messages_len == 0
    assert file.buf is not None
    print("✓ File initialized")
except Exception as e:
    print(f"✗ File test failed: {e}")
    import traceback
    traceback.print_exc()
    sys.exit(1)

# Test 7: Buffer operations
print("\n[Test 7] Testing buffer operations...")
try:
    from libsml.shared import SML_OPTIONAL_SKIPPED
    
    buf = libsml.SmlBuffer.from_bytes(bytes([SML_OPTIONAL_SKIPPED]))
    result = buf.optional_is_skipped()
    assert result == 1
    assert buf.cursor == 1
    print("✓ Optional skipped check works")
    
    buf2 = libsml.SmlBuffer(10)
    buf2.optional_write()
    assert buf2.cursor == 1
    assert buf2.buffer[0] == SML_OPTIONAL_SKIPPED
    print("✓ Optional write works")
except Exception as e:
    print(f"✗ Buffer operations test failed: {e}")
    import traceback
    traceback.print_exc()
    sys.exit(1)

print("\n" + "=" * 60)
print("All tests passed! ✓")
print("=" * 60)

