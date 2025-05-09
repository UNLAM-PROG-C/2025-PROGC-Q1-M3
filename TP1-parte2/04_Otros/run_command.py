# /// script
# requires-python = ">=3.9"
# dependencies = [
#     "matplotlib",
#     "tqdm",
# ]
# ///

import subprocess
import sys

MIN_CLONES = 2
MAX_CLONES = 12
NUM_CALLS = 6
NUM_THREADS = 4
PLOT_FILE = "results.png"

def parse_output(output: str) -> tuple[int,float]:
    """
    Parses the output of the Rust program to extract the level reached and the time taken.
    """
    lines = list(output.splitlines())
    level = int(lines[0].split(": ")[1])
    seconds = float(lines[1].split(": ")[1])
    return level, seconds

def call_rust_program(num_clones: int) -> tuple[int, float]:
    """
    Calls the Rust program with the specified number of clones and returns the level reached and time taken.
    """
    rust_env = {
        "NUM_CLONES": str(num_clones),
        "NUM_THREADS": str(NUM_THREADS),
    }
    result = subprocess.run(
        ["./target/debug/naruto_training"],
        env=rust_env,
        input=None,
        capture_output=True,
        text=True
    )
    return parse_output(result.stdout.strip())

def generate_results() -> list[tuple[int, int, float]]:
    """
    Generates the results for the Rust program by calling it with different numbers of clones.
    """
    from tqdm import tqdm
    results = []
    progress_bar = tqdm(total=(MAX_CLONES - MIN_CLONES + 1) * NUM_CALLS, desc="Generating results", unit="call")

    for num_clones in range(MIN_CLONES, MAX_CLONES + 1):
        for _ in range(NUM_CALLS):
            level, seconds = call_rust_program(num_clones)
            results.append((num_clones, level, seconds))
            progress_bar.update(1)
    return results

def plot_results(results: list[tuple[int, int, float]]):
    """
    Plots the results of the Rust program.
    """
    import matplotlib.pyplot as plt

    num_clones, levels, times = zip(*results)


    plt.scatter(num_clones, levels, label='Level Reached', color='blue')
    plt.scatter(num_clones, times, label='Time Taken (s)', color='red')


    plt.xlabel('Number of Clones')
    plt.ylabel('Level Reached / Time Taken (s)')
    plt.title('Rust Program Performance')
    plt.legend()

    plt.savefig(PLOT_FILE)
    print(f"Plot saved to {PLOT_FILE}")

def main():
    results = generate_results()
    plot_results(results)
    print("Results generated and plotted.")

if __name__ == "__main__":
    main()
    
