import cv2
import numpy as np


img = cv2.imread('test.jpg', 0)

# create a CLAHE object
clahe1 = cv2.createCLAHE(clipLimit=1.0, tileGridSize=(8, 8))
cl1 = clahe1.apply(img)

clahe2 = cv2.createCLAHE(clipLimit=2.0, tileGridSize=(8, 8))
cl2 = clahe2.apply(img)

cv2.imwrite('test-1.jpg', cl1)
cv2.imwrite('test-2.jpg', cl2)
