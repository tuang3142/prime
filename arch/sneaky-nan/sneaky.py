# encode the message to be nan
# but we can deode it to see the msg
# 64 bit
# 1 bit: sign bit
# 11 bits: exponent - all 1s
# 1 bit: always on - so that it can't be INF, 3 bits: length 8?
# 48 bits: content, utf8-encoded, 6 bites

def conceal(msg):
    bs = msg.encode('utf-8')
    n = len(bs)
    if n > 6:
        raise ValueError()
    first = b'x7f'
    second = (0xf8 ^ n).tobytes(1, 'big') 
    padding = b'0x00' * (6 - n) 
    payload = bs
    return struct.unpack('>d, first + second + padding + payload)[0]

def extract(x):
    return

msg = "hello"
# do some assertion here

if __name__ == '__main__':
    msg = 'secret'
    x = conceal(msg)
    assert isinstance(x, float)
    assert repr(x) == 'nan'
    assert extract(x) == msg
