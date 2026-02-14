import pandas as pd
import os

def analyze():
    client_file = 'client_results.csv'
    resource_file = 'raw_resources.log'

    if not os.path.exists(client_file) or not os.path.exists(resource_file):
        print(f"Error: Could not find files.")
        return

    print("Reading and Cleaning files...")
    
    # 1. Load Go Client Data
    go_df = pd.read_csv(client_file)

    # 2. Load and Clean Docker Stats
    docker_rows = []
    with open(resource_file, 'r') as f:
        for line in f:
            parts = line.strip().split(',')
            if len(parts) >= 3:
                ts = float(parts[0])
                cpu = float(parts[1].replace('%', ''))
                mem = float(parts[2].split('MiB')[0])
                docker_rows.append([ts, cpu, mem])

    docker_df = pd.DataFrame(docker_rows, columns=['timestamp', 'cpu_percent', 'mem_mib'])

    # 3. Fuzzy Merge
    df = pd.merge_asof(go_df.sort_values('timestamp'), 
                       docker_df.sort_values('timestamp'), 
                       on='timestamp', 
                       direction='nearest')

    # 4. THROUGHPUT & JITTER LOGIC
    window = 60
    # Throughput (Hz)
    df['throughput_hz'] = df['seq'].diff(window) / df['timestamp'].diff(window)
    
    # Inter-Packet Delay (ms) - This shows the gap between individual inputs
    # A perfect system stays at 16.67ms
    df['inter_packet_ms'] = df['timestamp'].diff() * 1000
    
    # Identify Stalls: Any gap larger than 2 frame durations (33.3ms)
    # This is what the user perceives as a "hitch" or "stutter"
    df['is_stall'] = df['inter_packet_ms'] > 33.33

    # 5. Save and Clean
    clean_df = df.dropna().copy()
    clean_df.to_csv('final_analysis.csv', index=False)

    # 6. CALCULATE FINAL METRICS
    avg_tput = clean_df['throughput_hz'].mean()
    tput_jitter = clean_df['throughput_hz'].std() # Variation in frequency
    avg_ipd = clean_df['inter_packet_ms'].mean()
    ipd_jitter = clean_df['inter_packet_ms'].std() # The "Smoothness" metric
    stall_count = clean_df['is_stall'].sum()
    stall_percent = (stall_count / len(clean_df)) * 100

    print(f"--- RESOURCE METRICS ---")
    print(f"Peak CPU:            {clean_df['cpu_percent'].max():.2f}%")
    print(f"Avg CPU:             {clean_df['cpu_percent'].mean():.2f}%")
    
    print(f"\n--- LATENCY METRICS ---")
    print(f"Avg Injection Latency: {clean_df['latency_ms'].mean():.2f} ms")
    print(f"Avg Perceived Lag:     {clean_df['drift_ms'].mean():.2f} ms")
    print(f"Max Perceived Lag:     {clean_df['drift_ms'].max():.2f} ms")

    print(f"\n--- SMOOTHNESS & JITTER ---")
    print(f"Avg Throughput:      {avg_tput:.2f} Hz")
    print(f"Throughput Jitter:   {tput_jitter:.2f} Hz ")
    print(f"Inter-Packet Jitter: {ipd_jitter:.2f} ms")
    print(f"Total Video Stalls:  {stall_count} ")
    print(f"Stall Percentage:    {stall_percent:.2f}% of test duration")

if __name__ == "__main__":
    analyze()