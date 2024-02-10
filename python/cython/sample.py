import time

import tarai


def pytarai(n) -> int:
    if n == 0 or n == 1:
        return 0
    elif n == 2:
        return 1
    else:
        return pytarai(n - 1) + pytarai(n - 2) + pytarai(n - 3)


n = 30
s1 = time.perf_counter()
pytarai(n)
total1 = time.perf_counter() - s1
print("python:\t", total1)

s2 = time.perf_counter()
tarai.tarai(n)
total2 = time.perf_counter() - s2
print("cython:\t", total2)


print(total1 / total2)
