# Create file
with open("large.txt", "wb") as f:
    f.truncate(1024 * 1024 * 1024 * 1024)


def regular_io(filename):
    with open(filename, mode="r", encoding="utf8") as file_obj:
        text = file_obj.read()
        print(text.find("abc"))


regular_io("large.txt")
