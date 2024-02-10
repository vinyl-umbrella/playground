cpdef int tarai(int n):
    if n == 0 or n == 1:
        return 0

    elif n == 2:
        return 1

    else:
        return tarai(n - 1) + tarai(n - 2) + tarai(n - 3)

