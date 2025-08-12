import subprocess

def read_tests(filename):
    with open(filename) as f:
        lines = [l.strip() for l in f if l.strip() and not l.startswith("@")]
    inputs, outputs = [], []
    for line in lines:
        if ',' not in line:
            continue
        inp, out = line.split(",", 1)
        inp = inp.encode().decode("unicode_escape")
        out_decoded = out.encode().decode("unicode_escape")
        if out_decoded == "":
            outputs.append([])
        else:
            outputs.append([o.strip() for o in out_decoded.split("\n")])
        inputs.append(inp)
    return inputs, outputs

def run_repl(inputs):
    proc = subprocess.Popen(
        ["go", "run", "."],
        stdin=subprocess.PIPE,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        text=True
    )
    stdout, stderr = proc.communicate(input="\n".join(inputs) + "\nexit\n")
    out_lines = [l.strip() for l in stdout.splitlines()]
    err_lines = [l.strip() for l in stderr.splitlines()]
    return out_lines, err_lines

def check_results(inputs, expected_outputs, actual_outputs, stream_name):
    passed, failed = 0, 0
    idx = 0
    for i, expected in enumerate(expected_outputs):
        got = actual_outputs[idx:idx + len(expected)]
        if got != expected:
            print(f"[{stream_name} FAIL]")
            print(f"  Input:    {inputs[i]}")
            print(f"  Expected: {expected}")
            print(f"  Got:      {got}\n")
            failed += 1
        else:
            passed += 1
        idx += len(expected)
    return passed, failed

def main():
    SETUP = 1
    err_inputs, err_expected = read_tests("repl_tests/syntax_errors.txt")
    good_inputs, good_expected = read_tests("repl_tests/tests.txt")

    _, err_output = run_repl(err_inputs)
    err_passed, err_failed = check_results(err_inputs, err_expected, err_output, "stderr")
    err_passed -= SETUP

    good_output, _ = run_repl(good_inputs)
    good_passed, good_failed = check_results(good_inputs, good_expected, good_output, "stdout")

    print("\n--- Test Summary ---")
    print(f"Syntax Errors: {err_passed} passed, {err_failed} failed")
    print(f"Valid Tests:   {good_passed} passed, {good_failed} failed")
    print(f"All Tests:     {err_passed + good_passed} passed, {err_failed + good_failed} failed")

    if err_failed == 0 and good_failed == 0:
        print("\n✅ All tests passed!")
    else:
        print("\n❌ Some tests failed.")

if __name__ == "__main__":
    main()
