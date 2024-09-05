import matplotlib.pyplot as plt

def read_benchmark_results(filename):
    queue_sizes = []
    durations = []

    with open(filename, 'r') as file:
        for line in file:
            parts = line.strip().split(',')
            if len(parts) == 2:
                try:
                    queue_size = int(parts[0].split(':')[1].strip())
                    duration = int(parts[1].split(':')[1].strip())
                    queue_sizes.append(queue_size)
                    durations.append(duration)
                except ValueError:
                    print(f"Ошибка чтения строки: {line}")

    return queue_sizes, durations

def plot_benchmark(queue_sizes, durations):
    plt.figure(figsize=(10, 6))

    if max(queue_sizes) > 1000:
        plt.xscale('log')
    if max(durations) > 1000:
        plt.yscale('log')

    plt.plot(queue_sizes, durations, label="Duration per queue size", marker='o')
    plt.xlabel("Queue Size")
    plt.ylabel("Duration (ns)")
    plt.title("Benchmark Results by Queue Size")
    plt.legend()
    plt.grid(True)
    plt.show()

# Основной код
if __name__ == "__main__":
    queue_sizes, durations = read_benchmark_results('benchmark_results.txt')
    if queue_sizes and durations:
        plot_benchmark(queue_sizes, durations)
    else:
        print("Нет данных для отображения графика.")
