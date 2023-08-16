import numpy as np

def einsum_practice():
    # Part 1
    A = np.random.rand(5, 5)
    einTrace = np.einsum("ii->", A)
    trace = np.trace(A)
    result1 = np.abs(trace - einTrace)
    print("Results for part a: " + str(result1))
    
    # Part 2
    b = np.random.rand(5, 1)
    einMul = np.einsum("ij,jk->ik", A, b)
    mul = np.dot(A, b)
    result2 = np.linalg.norm(mul - einMul)
    print("Results for part b: " + str(result2))

    # Part 3
    a = np.random.rand(5, 1)
    einVecProd = np.einsum("ij,kj->ik", a, b)
    vecProd = np.dot(a, b.T)
    result3 = np.linalg.norm(vecProd - einVecProd)
    print("Results for part c: " + str(result3))

einsum_practice()