#define PY_SSIZE_T_CLEAN
#include <Python.h>
#include <stdint.h>

static PyObject* cvarint_encode(PyObject* self, PyObject* args) {
  uint64_t n;
  if (!PyArg_ParseTuple(args, "k", &n)) {
    return NULL;
  }

  uint64_t temp_n = n;
  size_t cnt = 0;
  do {
    temp_n >>= 7;
    cnt += 1;
  } while (temp_n);

  uint8_t* out = (uint8_t*)calloc(cnt, sizeof *out);
  if (!out) {
    PyErr_SetString(PyExc_MemoryError, "Failed to allocate memory");
    return NULL;
  }

  for (size_t i = 0; i < cnt; i++) {
    uint8_t b = n & 0x7F;
    n >>= 7;
    out[i] = n ? (b | 0x80) : b;
  }

  PyObject* result = PyBytes_FromStringAndSize((char*)out, cnt);
  free(out);
  return result;
}

static PyObject* cvarint_decode(PyObject* self, PyObject* args) {
  char* varint;
  if (!PyArg_ParseTuple(args, "s#", &varint)) {
    return NULL;
  }

  uint64_t out = 0, shift = 0;
  for (size_t i = 0;; i++) {
    uint8_t byte = varint[i];
    uint8_t payload = byte & 0x7F;  // get the last 7 bits
    out |= (uint64_t)payload << shift; // concat bits
    if (!(byte & 0x80)) {  // break if continuation bit == 0
      break;
    }
    shift += 7;
  }

  return PyLong_FromUnsignedLongLong(out);
}

static PyMethodDef CVarintMethods[] = {
    {"encode", cvarint_encode, METH_VARARGS, "Encode an integer as varint."},
    {"decode", cvarint_decode, METH_VARARGS,
     "Decode varint bytes to an integer."},
    {NULL, NULL, 0, NULL}};

static struct PyModuleDef cvarintmodule = {
    PyModuleDef_HEAD_INIT, "cvarint",
    "A C implementation of protobuf varint encoding", -1, CVarintMethods};

PyMODINIT_FUNC PyInit_cvarint(void) { return PyModule_Create(&cvarintmodule); }
