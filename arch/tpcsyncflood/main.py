import struct

with open("synflood.pcap", "rb") as f:
    magic_number, vmajor, vminor, _, _, _, llh_type = struct.unpack(
        "<IHHIIII", f.read(24)
    )
    assert magic_number == 0xA1B2C3D4
    print(f"pcap version: {vmajor}.{vminor}")
    assert llh_type == 0  # link layer header type, = 0 means localhost

    package_count = 0
    while True:  # reading per package header
        per_package_header = f.read(16)
        if len(per_package_header) == 0:
            break
        package_count += 1
        _, _, length, untruncate_length = struct.unpack("<IIII", per_package_header)
        assert length == untruncate_length
        packet = f.read(length) # read the package
        assert struct.unpack('<I', packet[:4])[0] == 2
        ihl = (packet[4] & 0x0F) << 2
        assert ihl == 20 # assertion is useful for parsing byte
        src, dst, _, _, flags = struct.unpack('!HHIIH', packet[24:38])
        sync = (flags & 0x0002) > 0
        ack  = (flags & 0x0010) > 0
        print(f'{src} -> {dst}')
    print(package_count)
