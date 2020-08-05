from collections import namedtuple

d = {"name": "kien", "age": "26"}
d_named = namedtuple("Tui", d.keys())(*d.values())
print(d_named)
print(d_named.name)
