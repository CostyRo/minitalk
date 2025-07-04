import subprocess

def main():
    with open("repl_tests/tests.txt", "r") as f:
        lines = [line.strip() for line in f if line.strip() and not line.strip().startswith("#")]

    inputs = []
    expected_outputs = []
    for line in lines:
        if ',' not in line:
            print(f"Ignored invalid line (no comma): {line}")
            continue
        inp, outp = line.split(",", 1)
        inputs.append(inp.strip())
        expected_outputs.append(outp.strip())

    repl_input = "\n".join(inputs) + "\nexit\n"

    proc = subprocess.Popen(
        ["cargo", "run", "--quiet",],
        stdin=subprocess.PIPE,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        text=True
    )

    stdout, _ = proc.communicate(input=repl_input)

    output_lines = [line.strip() for line in stdout.strip().splitlines() if line.strip()]

    min_len = min(len(output_lines), len(expected_outputs))
    all_pass = True
    for i in range(min_len):
        inp = inputs[i]
        got = output_lines[i]
        exp = expected_outputs[i]
        if got != exp:
            all_pass = False
            print(f"Test input:    {inp}")
            print(f"Expected out:  {exp}")
            print(f"Got output:    {got}")
            print("")

    if all_pass and len(output_lines) == len(expected_outputs):
        print("All tests passed!")
    else:
        if len(output_lines) != len(expected_outputs):
            print(f"Warning: output lines count ({len(output_lines)}) != expected count ({len(expected_outputs)})")

if __name__ == "__main__":
    main()
