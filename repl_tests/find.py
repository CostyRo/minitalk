import subprocess
import threading
import queue
from tqdm import tqdm

def enqueue_output(out, q):
    for line in iter(out.readline, ''):
        q.put(line)
    out.close()

proc = subprocess.Popen(
    ["go", "run", "."],
    stdin=subprocess.PIPE,
    stdout=subprocess.PIPE,
    stderr=subprocess.PIPE,
    encoding='utf-8',
    bufsize=1,
)

q = queue.Queue()
t = threading.Thread(target=enqueue_output, args=(proc.stdout, q))
t.daemon = True
t.start()

tests_with_found = []

with open("repl_tests/tests.txt") as f:
    for line in tqdm(f):
        line = line.strip()
        if not line or line.startswith("@"):
            continue
        inp, _ = line.split(",", 1)
        proc.stdin.write(inp + "\n")
        proc.stdin.flush()

        output_lines = []
        for _ in range(20):
            try:
                l = q.get(timeout=0.01)
            except queue.Empty:
                break
            output_lines.append(l.strip())
            if len(output_lines) >= 10:
                break

        if any(l.startswith("FOUND!") for l in output_lines):
            tests_with_found.append(line)

proc.stdin.close()
proc.terminate()
proc.wait()

print("\n".join(tests_with_found))
