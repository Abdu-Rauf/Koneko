import pandas as pd
import os
import numpy as np

# --- CONFIGURATION ---
# Set this to 50ms to simulate a standard internet connection.
# For your AI version later, you will subtract the prediction offset from this.
SIMULATED_RTT_MS = 50.0 

def analyze():
    client_file = 'legacy_stats.csv' # Or 'docker_exec.csv'
    resource_file = 'legacy_resources.log'

    if not os.path.exists(client_file) or not os.path.exists(resource_file):
        print(f"Error: Could not find '{client_file}' or '{resource_file}'.")
        return

    print("--- Loading and Cleaning Data ---")
    
    # 1. Load Go Client Data
    go_df = pd.read_csv(client_file)

    # 2. Load and Clean Docker Stats
    docker_rows = []
    with open(resource_file, 'r') as f:
        for line in f:
            parts = line.strip().split(',')
            if len(parts) >= 3:
                try:
                    ts = float(parts[0])
                    cpu = float(parts[1].replace('%', ''))
                    mem_str = parts[2].split('/')[0].strip()
                    if 'GiB' in mem_str:
                        mem = float(mem_str.replace('GiB', '')) * 1024
                    else:
                        mem = float(mem_str.replace('MiB', ''))
                    docker_rows.append([ts, cpu, mem])
                except ValueError:
                    continue

    docker_df = pd.DataFrame(docker_rows, columns=['timestamp', 'cpu_percent', 'mem_mib'])

    # 3. Fuzzy Merge
    df = pd.merge_asof(go_df.sort_values('timestamp'), 
                       docker_df.sort_values('timestamp'), 
                       on='timestamp', 
                       direction='nearest')

    # --- 4. ADVANCED METRICS LOGIC ---
    df['inter_packet_ms'] = df['timestamp'].diff() * 1000
    
    # Throughput (Hz)
    window = 60
    df['throughput_hz'] = df['seq'].diff(window) / df['timestamp'].diff(window)
    
    # Stalls & Inversions
    df['is_stall'] = df['inter_packet_ms'] > 33.33
    df['is_inversion'] = df['seq'] < df['seq'].shift(1)
    df['arrival_delta_ms'] = df['timestamp'].diff() * 1000
    df['is_flush'] = (df['arrival_delta_ms'] < 2).rolling(window=5).max()

    # --- 5. STEADY-STATE FILTERING ---
    warmup_period = 10 
    start_ts = df['timestamp'].min()
    steady_df = df[df['timestamp'] >= (start_ts + warmup_period)].copy()
    
    if steady_df.empty:
        steady_df = df.copy()

    # ==============================================================================
    # [NEW LOGIC] ADDING SIMULATED NETWORK LATENCY
    # ==============================================================================
    # "drift_ms" is your Raw Injection Latency (Server -> Container)
    # "total_lag_ms" is what the user feels (Network RTT + Injection Latency)
    
    steady_df['total_lag_ms'] = steady_df['drift_ms'] + SIMULATED_RTT_MS
    
    # ==============================================================================

    # --- 6. AGGREGATE RESULTS ---
    
    # Injection Metrics (System Performance)
    avg_injection = steady_df['drift_ms'].mean()
    p99_injection = steady_df['drift_ms'].quantile(0.99)

    # Perceived Latency Metrics (User Experience)
    avg_perceived = steady_df['total_lag_ms'].mean()
    p95_perceived = steady_df['total_lag_ms'].quantile(0.95)
    p99_perceived = steady_df['total_lag_ms'].quantile(0.99)
    max_perceived = steady_df['total_lag_ms'].max()

    # Smoothness Metrics
    avg_tput = steady_df['throughput_hz'].mean()
    tput_jitter = steady_df['throughput_hz'].std()
    ipd_jitter = steady_df['inter_packet_ms'].std()
    
    flush_packet_count = df['is_flush'].sum()
    flush_percentage = (flush_packet_count / len(df)) * 100
    inversion_count = df['is_inversion'].sum()
    stall_count = steady_df['is_stall'].sum()
    stall_percent = (stall_count / len(steady_df)) * 100

    # --- 7. FINAL REPORT ---
    print(f"\n{'='*40}")
    print(f"{'KONEKO PERFORMANCE REPORT':^40}")
    print(f"{'='*40}")

    print(f"\n[ RESOURCE CONSUMPTION ]")
    print(f"Peak CPU Usage:      {df['cpu_percent'].max():.2f}%")
    print(f"Avg CPU (Steady):    {steady_df['cpu_percent'].mean():.2f}%")
    print(f"Avg Memory:          {steady_df['mem_mib'].mean():.2f} MiB")

    print(f"\n[ SYSTEM LATENCY (INJECTION ONLY) ]")
    print(f"Avg Injection Delay: {avg_injection:.2f} ms")
    print(f"P99 Injection Delay: {p99_injection:.2f} ms")

    print(f"\n[ END-TO-END USER LATENCY (RTT + INJECTION) ]")
    print(f"Simulated Network:   {SIMULATED_RTT_MS} ms")
    print(f"Avg Perceived Lag:   {avg_perceived:.2f} ms")
    print(f"P95 Tail Latency:    {p95_perceived:.2f} ms")
    print(f"P99 Tail Latency:    {p99_perceived:.2f} ms (Worst Case)")
    print(f"Max Spike:           {max_perceived:.2f} ms")

    print(f"\n[ STABILITY & JITTER ]")
    print(f"Avg Throughput:      {avg_tput:.2f} Hz")
    print(f"Throughput Jitter:   {tput_jitter:.2f} Hz")
    print(f"Inter-Packet Jitter: {ipd_jitter:.2f} ms")
    print(f"Sequence Inversions: {inversion_count}")
    print(f"Stall Percentage:    {stall_percent:.2f}%")

    steady_df.to_csv('final_steady_state_analysis.csv', index=False)
    print(f"\nAnalysis complete. Results saved to 'legacy_analysis.csv'.")

if __name__ == "__main__":
    analyze()