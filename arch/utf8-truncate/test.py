data = open("in", "rb")

for line in data:
    for  l in line:
        print(l, end=' ')
        if l == b'\n': print()