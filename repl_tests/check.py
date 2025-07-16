import subprocess
import sys

def main():
    debug = "--debug" in sys.argv

    with open("repl_tests/tests.txt", "r") as f:
        lines = [line.strip() for line in f if line.strip() and not line.strip().startswith("#")]

    inputs = []
    expected_outputs = []
    for line in lines:
        if ',' not in line:
            print(f"Ignored invalid line (no comma): {line}")
            continue
        inp, outp = line.split(",", 1)
        inp = inp.strip().encode().decode("unicode_escape")
        outp = outp.strip().encode().decode("unicode_escape")
        inputs.append(inp.strip())
        expected_outputs.append([o.strip() for o in outp.split("\n")])

    repl_input = "\n".join(inputs) + "\nexit\n"

    proc = subprocess.Popen(
        ["go", "run", "."],
        stdin=subprocess.PIPE,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        text=True
    )

    stdout, stderr = proc.communicate(input=repl_input)

    output_lines = [
        line.strip() for line in stderr.strip().splitlines()+stdout.strip().splitlines() if line.strip()
    ]

    min_len = min(len(output_lines), len(expected_outputs))
    all_pass = True
    output_index = 0
    all_pass = True

    for i in range(len(inputs)):
        inp = inputs[i]
        exp_lines = expected_outputs[i]
        got_lines = output_lines[output_index : output_index + len(exp_lines)]

        if got_lines != exp_lines:
            all_pass = False
            print(f"Test input:    {inp}")
            print(f"Expected out:  {exp_lines}")
            print(f"Got output:    {got_lines}")
            print("")
        output_index += len(exp_lines)

    if all_pass:
        print(f"All {min_len} tests passed!")
    else:
        if len(output_lines) != len(expected_outputs):
            print(f"Warning: output lines count ({len(output_lines)}) != expected count ({len(expected_outputs)})")

    if debug:
        print("\nDebug output (expected vs actual):")
        for i in range(len(inputs)):
            inp = inputs[i]
            exp = expected_outputs[i] if i < len(expected_outputs) else ""
            got = output_lines[i] if i < len(output_lines) else ""
            print(f"{inp} | {exp} | {got}")

if __name__ == "__main__":
    main()
