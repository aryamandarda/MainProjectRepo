#include "matrix.h"
#include <stddef.h>
#include <stdio.h>
#include <stdlib.h>
#include <omp.h>

// Include SSE intrinsics
#if defined(_MSC_VER)
#include <intrin.h>
#elif defined(__GNUC__) && (defined(__x86_64__) || defined(__i386__))
#include <immintrin.h>
#include <x86intrin.h>
#endif

/* Below are some intel intrinsics that might be useful
 * void _mm256_storeu_pd (double * mem_addr, __m256d a)
 * __m256d _mm256_set1_pd (double a)
 * __m256d _mm256_set_pd (double e3, double e2, double e1, double e0)
 * __m256d _mm256_loadu_pd (double const * mem_addr)
 * __m256d _mm256_add_pd (__m256d a, __m256d b)
 * __m256d _mm256_sub_pd (__m256d a, __m256d b)
 * __m256d _mm256_fmadd_pd (__m256d a, __m256d b, __m256d c)
 * __m256d _mm256_mul_pd (__m256d a, __m256d b)
 * __m256d _mm256_cmp_pd (__m256d a, __m256d b, const int imm8)
 * __m256d _mm256_and_pd (__m256d a, __m256d b)
 * __m256d _mm256_max_pd (__m256d a, __m256d b)
*/

/* Generates a random double between low and high */
double rand_double(double low, double high) {
    double range = (high - low);
    double div = RAND_MAX / range;
    return low + (rand() / div);
}

/* Generates a random matrix */
void rand_matrix(matrix *result, unsigned int seed, double low, double high) {
    srand(seed);
    for (int i = 0; i < result->rows; i++) {
        for (int j = 0; j < result->cols; j++) {
            set(result, i, j, rand_double(low, high));
        }
    }
}

/*
 * Returns the double value of the matrix at the given row and column.
 * You may assume `row` and `col` are valid. Note that the matrix is in ROW MAJOR ORDER.
 */
double get(matrix *mat, int row, int col) {
	return mat->data[mat->cols * row + col];
}

/*
 * Sets the value at the given row and column to val. You may assume `row` and
 * `col` are valid. Note that the matrix is in ROW MAJOR ORDER.
 */
void set(matrix *mat, int row, int col, double val) {
	mat->data[mat->cols * row + col] = val;
}

/*
 * Allocates space for a matrix struct pointed to by the double pointer mat with
 * `rows` rows and `cols` columns. You should also allocate memory for the data array
 * and initialize all entries to be zeros. `parent` should be set to NULL to indicate that
 * this matrix is not a slice. You should also set `ref_cnt` to 1.
 * You should return -1 if either `rows` or `cols` or both have invalid values. Return -2 if any
 * call to allocate memory in this function fails.
 * Return 0 upon success.
 */
int allocate_matrix(matrix **mat, int rows, int cols) {
    // HINTS: Follow these steps.
    // 1. Check if the dimensions are valid. Return -1 if either dimension is not positive.
    // 2. Allocate space for the new matrix struct. Return -2 if allocating memory failed.
    // 3. Allocate space for the matrix data, initializing all entries to be 0. Return -2 if allocating memory failed.
    // 4. Set the number of rows and columns in the matrix struct according to the arguments provided.
    // 5. Set the `parent` field to NULL, since this matrix was not created from a slice.
    // 6. Set the `ref_cnt` field to 1.
    // 7. Store the address of the allocated matrix struct at the location `mat` is pointing at.
    // 8. Return 0 upon success.
	if (rows <= 0 || cols <= 0) {
		return -1;
	}
	matrix *matrixx = (matrix *) malloc(sizeof(struct matrix));
	if (matrixx == NULL) {
		return -2;
	}
	matrixx->data = calloc(rows * cols, sizeof(double));
	if (matrixx->data == NULL) {
		return -2;
	}	
	matrixx->rows = rows;
	matrixx->cols = cols;
	matrixx->parent = NULL;
	matrixx->ref_cnt = 1;
	*mat = matrixx;
	
	return 0;
}

/*
 * You need to make sure that you only free `mat->data` if `mat` is not a slice and has no existing slices,
 * or that you free `mat->parent->data` if `mat` is the last existing slice of its parent matrix and its parent
 * matrix has no other references (including itself).
 */
void deallocate_matrix(matrix *mat) {
    // HINTS: Follow these steps.
    // 1. If the matrix pointer `mat` is NULL, return.
    // 2. If `mat` has no parent: decrement its `ref_cnt` field by 1. If the `ref_cnt` field becomes 0, then free `mat` and its `data` field.
    // 3. Otherwise, recursively call `deallocate_matrix` on `mat`'s parent, then free `mat`.
	if (mat == NULL) {
		return;
	}	
	if (mat->parent == NULL) {
		mat->ref_cnt -= 1;
		if (mat->ref_cnt == 0) {
			free(mat->data);
			free(mat);
		}
	} else {
		deallocate_matrix(mat->parent);
		free(mat);
	}
}

/*
 * Allocates space for a matrix struct pointed to by `mat` with `rows` rows and `cols` columns.
 * Its data should point to the `offset`th entry of `from`'s data (you do not need to allocate memory)
 * for the data field. `parent` should be set to `from` to indicate this matrix is a slice of `from`
 * and the reference counter for `from` should be incremented. Lastly, do not forget to set the
 * matrix's row and column values as well.
 * You should return -1 if either `rows` or `cols` or both have invalid values. Return -2 if any
 * call to allocate memory in this function fails.
 * Return 0 upon success.
 * NOTE: Here we're allocating a matrix struct that refers to already allocated data, so
 * there is no need to allocate space for matrix data.
 */
int allocate_matrix_ref(matrix **mat, matrix *from, int offset, int rows, int cols) {
    // HINTS: Follow these steps.
    // 1. Check if the dimensions are valid. Return -1 if either dimension is not positive.
    // 2. Allocate space for the new matrix struct. Return -2 if allocating memory failed.
    // 3. Set the `data` field of the new struct to be the `data` field of the `from` struct plus `offset`.
    // 4. Set the number of rows and columns in the new struct according to the arguments provided.
    // 5. Set the `parent` field of the new struct to the `from` struct pointer.
    // 6. Increment the `ref_cnt` field of the `from` struct by 1.
    // 7. Store the address of the allocated matrix struct at the location `mat` is pointing at.
    // 8. Return 0 upon success.
	if (rows <= 0 || cols <= 0) {
		return -1;
	}
	matrix *new_mat = (matrix *) malloc(sizeof(struct matrix));
	if (new_mat == NULL) {
		return -2;
	}
	new_mat->data = from->data + offset;
	new_mat->rows = rows;
	new_mat->cols = cols;
	new_mat->parent = from;
	from->ref_cnt += 1;
	*mat = new_mat;

	return 0;
}

/*
 * Sets all entries in mat to val. Note that the matrix is in ROW MAJOR ORDER.
 */
void fill_matrix(matrix *mat, double val) {
	__m256d val_vec = {val, val, val, val};
	#pragma omp parallel for
	for (int i = 0; i < (mat->rows * mat->cols) / 4 * 4; i+=4) {
		_mm256_storeu_pd(mat->data + i, val_vec);
	}
	for (int i = (mat->rows * mat->cols) / 4 * 4; i < mat->rows * mat->cols; i++) {
		mat->data[i] = val;
	}
}

/*
 * Store the result of taking the absolute value element-wise to `result`.
 * Return 0 upon success.
 * Note that the matrix is in ROW MAJOR ORDER.
 */
int abs_matrix(matrix *result, matrix *mat) {
	if (mat->rows > 100 || mat->cols > 100) {
		__m256d op = {-0.0, -0.0, -0.0, -0.0};
		#pragma omp parallel for
		for (int i = 0; i < (mat->rows * mat->cols) / 4 * 4 ; i+=4) {
			__m256d temp = _mm256_loadu_pd(mat->data + i);
			__m256d final = _mm256_andnot_pd(op, temp);
			_mm256_storeu_pd(result->data + i, final);
		}
		for (int i = (mat->rows * mat->cols) / 4 * 4; i < mat->rows * mat->cols; i++) {
			if (mat->data[i] < 0) {
				result->data[i] = mat->data[i] * -1.0;
			} else {
				result->data[i] = mat->data[i];
			}
		}
	} else {
		for (int i = 0; i < (mat->rows * mat->cols) / 4 * 4; i+=4) {
			if (mat->data[i] < 0) {
				result->data[i] = mat->data[i] * -1.0;
			} else {
				result->data[i] = mat->data[i];
			} if (mat->data[i + 1] < 0) {
				result->data[i + 1] = mat->data[i + 1] * -1.0;
			} else {
				result->data[i + 1] = mat->data[i + 1];
			} if (mat->data[i + 2] < 0) {
				result->data[i + 2] = mat->data[i + 2] * -1.0;
			} else {
				result->data[i + 2] = mat->data[i + 2];
			} if (mat->data[i + 3] < 0) {
				result->data[i + 3] = mat->data[i + 3] * -1.0;
			} else {
				result->data[i + 3] = mat->data[i + 3];
			}
		}
		for (int i = (mat->rows * mat->cols) / 4 * 4; i < mat->rows * mat->cols; i++) {
			if (mat->data[i] < 0) {
				result->data[i] = mat->data[i] * -1.0;
			} else {
				result->data[i] = mat->data[i];
			}
		}
	}
	return 0;
}

/*
 * (OPTIONAL)
 * Store the result of element-wise negating mat's entries to `result`.
 * Return 0 upon success.
 * Note that the matrix is in row-major order.
 */
int neg_matrix(matrix *result, matrix *mat) {
    // Task 1.5 TODO
}

/*
 * Store the result of adding mat1 and mat2 to `result`.
 * Return 0 upon success.
 * You may assume `mat1` and `mat2` have the same dimensions.
 * Note that the matrix is in ROW MAJOR ORDER.
 */
int add_matrix(matrix *result, matrix *mat1, matrix *mat2) {
	if (mat1->rows > 100 || mat1->cols > 100) {
		#pragma omp parallel for
		for (int i = 0; i < (mat1->rows * mat1->cols) / 4 * 4; i+=4) {
			__m256d temp1 = _mm256_loadu_pd(mat1->data + i);
			__m256d temp2 = _mm256_loadu_pd(mat2->data + i); 
			__m256d sum = _mm256_add_pd(temp1, temp2);
			_mm256_storeu_pd(result->data + i, sum);
		}
		for (int i = (mat1->rows * mat1->cols) / 4 * 4; i < mat1->rows * mat1->cols; i++) {
			result->data[i] = mat1->data[i] + mat2->data[i];
		}
	} else {
		for (int i = 0; i < (mat1->rows * mat1->cols) / 4 * 4; i+=4) {
			result->data[i] = mat1->data[i] + mat2->data[i];
			result->data[i + 1] = mat1->data[i + 1] + mat2->data[i + 1];
			result->data[i + 2] = mat1->data[i + 2] + mat2->data[i + 2];
			result->data[i + 3] = mat1->data[i + 3] + mat2->data[i + 3];
		}
		for (int i = (mat1->rows * mat1->cols) / 4 * 4; i < mat1->rows * mat1->cols; i++) {
			result->data[i] = mat1->data[i] + mat2->data[i];
		}
		//for (int i = 0; i < mat1->rows * mat1->cols; i++) {
		//	result->data[i] = mat1->data[i] + mat2->data[i];
		//}
	}	
	return 0;
}

/*
 * (OPTIONAL)
 * Store the result of subtracting mat2 from mat1 to `result`.
 * Return 0 upon success and a nonzero value upon failure.
 * You may assume `mat1` and `mat2` have the same dimensions.
 * Note that the matrix is in ROW MAJOR ORDER.
 */
int sub_matrix(matrix *result, matrix *mat1, matrix *mat2) {
    // Task 1.5 TODO
}

/*
 * Store the result of multiplying mat1 and mat2 to `result`.
 * Return 0 upon success.
 * Remember that matrix multiplication is not the same as multiplying individual elements.
 * You may assume `mat1`'s number of columns is equal to `mat2`'s number of rows.
 * Note that the matrix is in ROW MAJOR ORDER.
 */
void transpose_matrix(matrix *new, matrix *original) {
	#pragma omp parallel for
	for (int i = 0; i < original->rows; i++) {
		for (int j = 0; j < original->cols; j++) {
			new->data[original->rows * j + i] = original->data[original->cols * i + j];	
		}
	}
}

int mul_matrix(matrix *result, matrix *mat1, matrix *mat2) {
	if (mat1->rows < 100 || mat1->cols < 100 || mat2->cols < 100) {
		for (int i = 0; i < result->rows; i++) {
			for (int j = 0; j < result->cols; j++) {
				result->data[result->cols * i + j] = 0;
			}
		}
		for (int j = 0; j < mat2->cols; j++) {
			for (int k = 0; k < mat2->rows; k++) {
				for (int i = 0; i < mat1->rows; i++) {
					result->data[result->cols * i + j] += mat1->data[mat1->cols * i + k] * mat2->data[mat2->cols * k + j];
				}
			}
		}
	} else {
		const int UNROLL = 4;
		matrix *mat2_t = NULL;
		allocate_matrix(&mat2_t, mat2->cols, mat2->rows);
		transpose_matrix(mat2_t, mat2);
		#pragma omp parallel for
		for (int i = 0; i < mat1->rows; i++) {
			for (int j = 0; j < mat2_t->rows; j++) {
				__m256d sum = {0.0, 0.0, 0.0, 0.0};
				for (int k = 0; k < mat2_t->cols / 16 * 16; k+=16) {
					__m256d temp1 = _mm256_loadu_pd(mat1->data + mat1->cols * i + k);
					__m256d temp2 = _mm256_loadu_pd(mat2_t->data + mat2_t->cols * j + k);
					sum = _mm256_fmadd_pd(temp1, temp2, sum);
					temp1 = _mm256_loadu_pd(mat1->data + mat1->cols * i + (k+4));
					temp2 = _mm256_loadu_pd(mat2_t->data + mat2_t->cols * j + (k+4));
					sum = _mm256_fmadd_pd(temp1, temp2, sum);
					temp1 = _mm256_loadu_pd(mat1->data + mat1->cols * i + (k+8));
					temp2 = _mm256_loadu_pd(mat2_t->data + mat2_t->cols * j + (k+8));
					sum = _mm256_fmadd_pd(temp1, temp2, sum);
					temp1 = _mm256_loadu_pd(mat1->data + mat1->cols * i + (k+12));
					temp2 = _mm256_loadu_pd(mat2_t->data + mat2_t->cols * j + (k+12));
					sum = _mm256_fmadd_pd(temp1, temp2, sum);
				}
				double tmp_arr[4];
				_mm256_storeu_pd((__m256d *) tmp_arr, sum);
				double total = tmp_arr[0] + tmp_arr[1] + tmp_arr[2] + tmp_arr[3];
				for (int k = mat2_t->cols / 16 * 16; k < mat2_t->cols; k++) {
					total += mat1->data[mat1->cols * i + k] * mat2_t->data[mat2_t->cols * j + k];
				}
				result->data[result->cols * i + j] = total;
			}
		}
	}	
	return 0;	
}

/*
 * Store the result of raising mat to the (pow)th power to `result`.
 * Return 0 upon success.
 * Remember that pow is defined with matrix multiplication, not element-wise multiplication.
 * You may assume `mat` is a square matrix and `pow` is a non-negative integer.
 * Note that the matrix is in ROW MAJOR ORDER.
 */
int pow_matrix(matrix *result, matrix *mat, int pow) {
	if (pow == 0) {
		for (int i = 0; i < mat->rows; i++) {
			for (int j = 0; j < mat->cols; j++) {
				if (i == j) {
					result->data[result->cols * i + j] = 1;
				} else {
					result->data[result->cols * i + j] = 0;
				}
			}
		}	
		return 0;
	}
	
	matrix *temp_mat = NULL;
	allocate_matrix(&temp_mat, mat->rows, mat->rows);
	matrix *temp_mat2 = NULL;
	allocate_matrix(&temp_mat2, mat->rows, mat->rows);
	for (int i = 0; i < mat->rows; i++) {
		for (int j = 0; j < mat->cols; j++) {
			double temp = mat->data[mat->cols * i + j];
			temp_mat2->data[temp_mat2->cols * i + j] = temp;
			if (i == j) {
				temp_mat->data[temp_mat->cols * i + j] = 1;
			} else {
				temp_mat->data[temp_mat->cols * i + j] = 0;
			}
		}
	}
	while (pow > 1) {
		if (pow % 2 == 1) {
			mul_matrix(result, temp_mat2, temp_mat);
			for (int r = 0; r < result->rows; r++) {
				for (int c = 0; c < result->cols; c++) {
					double temp = result->data[result->cols * r + c];
					temp_mat->data[temp_mat->cols * r + c] = temp;
				}	
			}		
			mul_matrix(result, temp_mat2, temp_mat2);
			pow = (pow - 1) / 2;
		} else {
			mul_matrix(result, temp_mat2, temp_mat2);
			pow = pow / 2;
		}
		for (int r = 0; r < temp_mat2->rows; r++) {
			for (int c = 0; c < temp_mat2->cols; c++) {
				double temp = result->data[result->cols * r + c];
				temp_mat2->data[temp_mat2->cols * r + c] = temp;
			}	
		}
	}
	mul_matrix(result, temp_mat2, temp_mat);
	return 0;
}

