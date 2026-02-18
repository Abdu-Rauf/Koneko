import pandas as pd
import os
import numpy as np

# --- CONFIGURATION ---
SIMULATED_RTT_MS = 50.0 

# Circle Constants (Must match Go Code to reconstruct the "Bullseye")
CENTER_X = 960.0
CENTER_Y = 540.0
RADIUS = 100.0
ANGLE_STEP = 0.1       
TICK_INTERVAL_MS = 16.0 

def analyze():
    client_file = 'predictor_stats.csv' 
    resource_file = 'predictor_resources.log'
    output_file = 'predictor_analysis_final.csv'

    if not os.path.exists(client_file) or not os.path.exists(resource_file):
        print(f"Error: Could not find '{client_file}' or '{resource_file}'.")
        return

    print(f"--- Processing {client_file} ---")
    
    # 1. Load Data
    go_df = pd.read_csv(client_file)
    
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

    # 2. Merge
    df = pd.merge_asof(go_df.sort_values('timestamp'), 
                       docker_df.sort_values('timestamp'), 
                       on='timestamp', 
                       direction='nearest')

    # 3. Basic Metrics
    df['inter_packet_ms'] = df['timestamp'].diff() * 1000
    window = 60
    df['throughput_hz'] = df['seq'].diff(window) / df['timestamp'].diff(window)
    df['is_stall'] = df['inter_packet_ms'] > 33.33
    df['is_inversion'] = df['seq'] < df['seq'].shift(1)

    # 4. Steady State (Skip Warmup)
    warmup_period = 10 
    start_ts = df['timestamp'].min()
    steady_df = df[df['timestamp'] >= (start_ts + warmup_period)].copy()
    
    if steady_df.empty: steady_df = df.copy()

    # ==============================================================================
    # [ CORE LOGIC ] CALCULATING "PERCEIVED LAG" VIA SPATIAL ERROR
    # ==============================================================================
    
    # A. Reconstruct the Target (The "Bullseye")
    steady_df['target_angle'] = steady_df['seq'] * ANGLE_STEP
    steady_df['target_x'] = CENTER_X + RADIUS * np.cos(steady_df['target_angle'])
    steady_df['target_y'] = CENTER_Y + RADIUS * np.sin(steady_df['target_angle'])

    # B. Calculate Error (Distance from Bullseye)
    steady_df['error_px'] = np.sqrt(
        (steady_df['x'] - steady_df['target_x'])**2 + 
        (steady_df['y'] - steady_df['target_y'])**2
    )

    # C. Convert Error to Time (Perceived Lag)
    # Mouse Speed = Distance / Time = (Radius * Step) / Tick
    # Speed = (100 * 0.1) / 16 = 0.625 px/ms
    mouse_speed_px_ms = (RADIUS * ANGLE_STEP) / TICK_INTERVAL_MS
    
    # This is your "Effective Latency" renamed to match your Base Report
    steady_df['perceived_lag_ms'] = steady_df['error_px'] / mouse_speed_px_ms

    # D. Physical Network Lag (The Control Variable)
    # This stays 50ms to prove you didn't fake the network
    steady_df['network_lag_ms'] = steady_df['drift_ms'] + SIMULATED_RTT_MS

    # ==============================================================================

    # 5. Aggregation
    avg_cpu = steady_df['cpu_percent'].mean()
    avg_mem = steady_df['mem_mib'].mean()
    avg_injection = steady_df['drift_ms'].mean()
    p99_injection = steady_df['drift_ms'].quantile(0.99)
    
    # The Comparison Metrics
    avg_perceived = steady_df['perceived_lag_ms'].mean()
    p99_perceived = steady_df['perceived_lag_ms'].quantile(0.99)
    max_perceived = steady_df['perceived_lag_ms'].max()
    
    avg_tput = steady_df['throughput_hz'].mean()
    ipd_jitter = steady_df['inter_packet_ms'].std()
    stall_percent = (steady_df['is_stall'].sum() / len(steady_df)) * 100

    # --- FINAL REPORT (MATCHING YOUR BASE FORMAT) ---
    print(f"\n{'='*40}")
    print(f"{'KONEKO PREDICTOR REPORT':^40}")
    print(f"{'='*40}")

    print(f"\n[ RESOURCE CONSUMPTION ]")
    print(f"Avg CPU (Steady):    {avg_cpu:.2f}%")
    print(f"Avg Memory:          {avg_mem:.2f} MiB")

    print(f"\n[ SYSTEM LATENCY (INJECTION ONLY) ]")
    print(f"Avg Injection Delay: {avg_injection:.2f} ms")
    print(f"P99 Injection Delay: {p99_injection:.2f} ms")

    print(f"\n[ END-TO-END USER LATENCY (RTT + INJECTION) ]")
    print(f"Simulated Network:   {SIMULATED_RTT_MS} ms")
    # This is the "Honest" Network Lag
    print(f"Physical Network Lag:{steady_df['network_lag_ms'].mean():.2f} ms") 
    
    # This matches your "Avg Perceived Lag" from the Base Report
    print(f"Avg Perceived Lag:   {avg_perceived:.2f} ms  <-- (AI Corrected)")
    print(f"P99 Tail Latency:    {p99_perceived:.2f} ms")
    print(f"Max Spike:           {max_perceived:.2f} ms")

    print(f"\n[ STABILITY & JITTER ]")
    print(f"Avg Throughput:      {avg_tput:.2f} Hz")
    print(f"Inter-Packet Jitter: {ipd_jitter:.2f} ms")
    print(f"Stall Percentage:    {stall_percent:.2f}%")
    
    # Save for graphing
    steady_df.to_csv(output_file, index=False)
    print(f"\nResults saved to '{output_file}'")

if __name__ == "__main__":
    analyze()