#include <Python.h>

/* add(a, b): return a + b */
static PyObject* mymath_add(PyObject *self, PyObject *args) {
    long a, b;

    // parse two long ints from Python arguments
    if (!PyArg_ParseTuple(args, "ll", &a, &b)) {
        return NULL; // error already set
    }

    long result = a + b;
    return PyLong_FromLong(result);
}

// Method table
static PyMethodDef MyMathMethods[] = {
    {"add", mymath_add, METH_VARARGS, "Add two integers."},
    {NULL, NULL, 0, NULL}  // sentinel
};

// Module definition
static struct PyModuleDef mymathmodule = {
    PyModuleDef_HEAD_INIT,
    "mymath",              // module name in Python
    "Minimal C extension", // docstring
    -1,                    // per-interpreter state size (or -1)
    MyMathMethods
};

// Initialization function: must be PyInit_<modulename>
PyMODINIT_FUNC PyInit_mymath(void) {
    return PyModule_Create(&mymathmodule);
}
