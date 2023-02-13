import sys
import fontforge

f = fontforge.open(sys.argv[1])
k = f['k']
i = f ['i']

if i.width == k.width:
    print("Monospaced")
else:
    print("Non-monospaced")
