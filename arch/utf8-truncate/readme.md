# unicode

a big fucking map between char and special char to integer (hex to be more precise)
u+hex for puny human reading

# utf8 (pretty interesting actually)

encode unicode to byte. 1 byte to 4 byte maxium (i guess if it goes above the maximum we can just add another byte)
each byte have a header. the leading byte header tell howmny byte there is to read. (0x110 - 2 bytes, 0x - 1 byte)
the continuation byte start with 0x10
