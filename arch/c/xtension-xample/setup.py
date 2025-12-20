from setuptools import setup, Extension

module = Extension(
    "mymath",              # module name
    sources=["mymath.c"],  # C source file
)

setup(
    name="mymath",
    version="0.1",
    ext_modules=[module],
)