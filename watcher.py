import time
import os
import subprocess
import sys

# Command to execute when the file changes
cs = [
    ["echo", "\n---\n"],
    ["python3", "db/dbms-py/db/encode.py"],
]

def get_mtime(path):
    try:
        return os.path.getmtime(path)
    except FileNotFoundError:
        return None

def main():
    if len(sys.argv) < 2:
        print("Usage: python3 watcher.py <file_to_watch>")
        sys.exit(1)

    file_to_watch = sys.argv[1]
    last_mtime = get_mtime(file_to_watch)
    if last_mtime is None:
        print(f"File {file_to_watch} not found.")
        sys.exit(1)

    print(f"Watching {file_to_watch} for changes...")
    while True:
        time.sleep(1)
        current_mtime = get_mtime(file_to_watch)
        if current_mtime is not None and current_mtime != last_mtime:
            print(f"{file_to_watch} changed!")
            for c in cs:
                subprocess.run(c)
            last_mtime = current_mtime

if __name__ == "__main__":
    main()
