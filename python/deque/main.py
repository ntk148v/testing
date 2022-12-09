import collections

# init
de = collections.deque([1, 2, 3, 3, 4, 2, 4])
print("deque:", de)

# append element at right end
de.append(4)
print("deque:", de)

# append element at left end (start)
de.appendleft(6)
print("deque:", de)

# pop (delete) element from right end
de.pop()
print("deque:", de)

# popleft element from left end
de.popleft()
print("deque:", de)

# using index() to print the first occurrence of 4
# index(ele, beg, end):- This function returns the first index of the
# value mentioned in arguments, starting searching from beg till end index.
print("the number 4 first occurs at a position:", de.index(4, 2, 5))

# insert the value at exact position
de.insert(4, 3)
print("deque:", de)

# count the occurrences of 3
print("the count of 3 in deque:", de.count(3))

# remove() the first occurrence of 3
de.remove(3)
print("deque:", de)

# extend to add numbers to right end
de.extend([4, 5, 6])
print("deque:", de)

# same with extendleft (left end)
de.extendleft([7, 8, 9])
print("deque:", de)

# rotate() to rotate the deque
# rotate():- This function rotates the deque by the number specified in
# arguments. If the number specified is negative, rotation occurs to the
# left. Else rotation is to right.
de.rotate(-3)
print("deque:", de)

# reverse the deque
de.reverse()
print("deque:", de)
