import sys
import pandas as pd
import numpy as np
import matplotlib.pyplot as plt
import seaborn as sns
import os
from collections import defaultdict

# Set style for plots
plt.style.use('ggplot')
sns.set_palette("Set2")

# Define colors for consistency
COLORS = {
    'Random': '#ff7f0e',  # Orange
    'LocalSearch': '#1f77b4'  # Blue
}

def clean_solution(solution_str):
    solution = solution_str.strip('[]').split()
    solution = [int(x) for x in solution]
    return solution


def load_optimal_solution(instance_name):
    base_name = os.path.splitext(instance_name)[0]
    sln_path = f"instances/{base_name}.sln"
    
    try:
        with open(sln_path, 'r') as f:
            lines = f.readlines()
            if len(lines) < 2:
                return None
            
            # First line: instance size and optimal fitness separated by whitespace
            size, optimal_fitness = map(int, lines[0].strip().split())
            
            # Initialize an empty list for the optimal solution
            optimal_solution = []
            
            # Start reading the solution from the second line onwards
            for line in lines[1:]:
                # Split the line into integers and append to the optimal_solution list
                optimal_solution.extend(map(int, line.strip().split()))
            
            return optimal_fitness, optimal_solution
    except FileNotFoundError:
        print(f"Optimal solution file not found for {instance_name}")
        return None
    except ValueError as e:
        print(f"Error parsing file {sln_path}: {e}")
        return None

def analyze_qap_results(df):
    instances = df['Instance'].unique()
    
    best_solution_fitness = {}
    best_solutions = {}
    
    for instance in instances:
        best_solution_fitness[instance], best_solutions[instance] = load_optimal_solution(instance)

    print(best_solutions)
    
    # Add gap to best known solution
    df['GapToBest'] = df.apply(lambda row: (row['FinalFitness'] - best_solution_fitness[row['Instance']]) / best_solution_fitness[row['Instance']] * 100, axis=1)
    
    # Add improvement percentage
    df['ImprovementPercent'] = ((df['InitialFitness'] - df['FinalFitness']) / df['InitialFitness'] * 100)
    
    # Add evaluations per second
    df['EvalsPerSecond'] = df['Evaluations'] / (df['TimeMs'] / 1000)
    
    return df, best_solution_fitness

def plot_gap_to_best(df, output_dir):
    plt.figure(figsize=(14, 8))
    
    instances = df['Instance'].unique()
    
    # Calculate mean gap by instance and solver
    gap_data = df.groupby(['Instance', 'Solver'])['GapToBest'].mean().reset_index()
    
    # For each solver, create a set of bars
    solvers = df['Solver'].unique()
    width = 0.8 / len(solvers)
    
    for i, solver in enumerate(solvers):
        solver_data = gap_data[gap_data['Solver'] == solver]
        x_positions = np.arange(len(instances)) + (i - len(solvers)/2 + 0.5) * width
        plt.bar(x_positions, solver_data['GapToBest'], width=width, label=solver, color=COLORS[solver])
        
    plt.xlabel('Instance')
    plt.ylabel('Gap to Best Known Solution (%)')
    plt.title('Average Gap to Best Known Solution by Solver and Instance')
    plt.xticks(np.arange(len(instances)), instances, rotation=45)
    plt.legend()
    plt.tight_layout()
    plt.savefig(os.path.join(output_dir, 'gap_to_best.png'), dpi=300)
    plt.close()

def plot_initial_vs_final(df, output_dir):
    instances = df['Instance'].unique()
    
    for instance in instances:
        plt.figure(figsize=(10, 8))
        instance_df = df[df['Instance'] == instance]
        
        for solver in df['Solver'].unique():
            solver_df = instance_df[instance_df['Solver'] == solver]
            plt.scatter(solver_df['InitialFitness'], solver_df['FinalFitness'], 
                       label=solver, alpha=0.7, color=COLORS[solver])
            
        plt.xlabel('Initial Fitness')
        plt.ylabel('Final Fitness')
        plt.title(f'Initial vs Final Fitness for {instance}')
        plt.legend()
        # Add a diagonal line for reference (no improvement line)
        max_val = max(instance_df['InitialFitness'].max(), instance_df['FinalFitness'].max())
        min_val = min(instance_df['InitialFitness'].min(), instance_df['FinalFitness'].min())
        plt.plot([min_val, max_val], [min_val, max_val], 'k--', alpha=0.3)
        plt.tight_layout()
        plt.savefig(os.path.join(output_dir, f'initial_vs_final_{instance}.png'), dpi=300)
        plt.close()
    
    # Create one combined plot with normalized values for comparison
    plt.figure(figsize=(12, 10))
    
    for instance in instances:
        instance_df = df[df['Instance'] == instance]
        
        # Normalize values for this instance
        min_initial = instance_df['InitialFitness'].min()
        max_initial = instance_df['InitialFitness'].max()
        min_final = instance_df['FinalFitness'].min()
        max_final = instance_df['FinalFitness'].max()
        
        for solver in df['Solver'].unique():
            solver_df = instance_df[instance_df['Solver'] == solver]
            
            x_norm = (solver_df['InitialFitness'] - min_initial) / (max_initial - min_initial) if max_initial > min_initial else solver_df['InitialFitness'] / max_initial
            y_norm = (solver_df['FinalFitness'] - min_final) / (max_final - min_final) if max_final > min_final else solver_df['FinalFitness'] / max_final
            plt.scatter(
                x_norm, 
                y_norm, 
                alpha=0.7, 
                marker='o' if solver == 'Random' else '^', 
                color=COLORS[solver],
            )
    
    plt.xlabel('Normalized Initial Fitness')
    plt.ylabel('Normalized Final Fitness')
    plt.title('Initial vs Final Fitness (Normalized) Across All Instances')
    plt.legend(bbox_to_anchor=(1.05, 1), loc='upper left')
    plt.tight_layout()
    plt.savefig(os.path.join(output_dir, 'initial_vs_final_normalized.png'), dpi=300)
    plt.close()

def plot_best_solutions(df, best_solutions, output_dir):
    plt.figure(figsize=(14, 8))
    
    instances = list(best_solutions.keys())
    best_by_solver = defaultdict(list)
    
    for instance in instances:
        instance_df = df[df['Instance'] == instance]
        for solver in df['Solver'].unique():
            solver_df = instance_df[instance_df['Solver'] == solver]
            best_by_solver[solver].append(solver_df['FinalFitness'].min())
    
    # Create grouped bar plot
    x = np.arange(len(instances))
    width = 0.35
    
    _, ax = plt.subplots(figsize=(14, 8))
    
    solvers = list(best_by_solver.keys())
    for i, solver in enumerate(solvers):
        offset = (i - len(solvers)/2 + 0.5) * width
        ax.bar(x + offset, best_by_solver[solver], width, label=solver, color=COLORS[solver])
    
    ax.set_ylabel('Best Fitness Found')
    ax.set_title('Best Solution Found for Each Instance by Solver')
    ax.set_xticks(x)
    ax.set_xticklabels(instances, rotation=45)
    ax.legend()
    
    # Add text with the values
    for i, solver in enumerate(solvers):
        offset = (i - len(solvers)/2 + 0.5) * width
        for j, v in enumerate(best_by_solver[solver]):
            ax.text(j + offset, v, f"{v:.0f}", 
                   ha='center', va='bottom', fontsize=8, rotation=90)
    
    plt.tight_layout()
    plt.yscale("log")
    plt.savefig(os.path.join(output_dir, 'best_solutions.png'), dpi=300)
    plt.close()

def plot_time_efficiency(df, output_dir):
    plt.figure(figsize=(14, 8))
    
    # Group by instance and solver to get mean values
    perf_data = df.groupby(['Instance', 'Solver'])[['TimeMs', 'Evaluations', 'EvalsPerSecond']].mean().reset_index()
    
    # Create bar plot for evaluations per second by instance and solver
    instances = df['Instance'].unique()
    solvers = df['Solver'].unique()
    
    width = 0.8 / len(solvers)
    
    for i, solver in enumerate(solvers):
        solver_data = perf_data[perf_data['Solver'] == solver]
        x_positions = np.arange(len(instances)) + (i - len(solvers)/2 + 0.5) * width
        
        plt.bar(x_positions, solver_data['EvalsPerSecond'], width=width, label=solver, color=COLORS[solver])
    
    plt.xlabel('Instance')
    plt.ylabel('Evaluations per Second')
    plt.title('Algorithm Efficiency: Evaluations per Second')
    plt.xticks(np.arange(len(instances)), instances, rotation=45)
    plt.yscale('log')
    plt.legend()
    plt.tight_layout()
    plt.savefig(os.path.join(output_dir, 'evaluations_per_second.png'), dpi=300)
    plt.close()
    
    # Plot evaluations count
    plt.figure(figsize=(14, 8))
    
    for i, solver in enumerate(solvers):
        solver_data = perf_data[perf_data['Solver'] == solver]
        x_positions = np.arange(len(instances)) + (i - len(solvers)/2 + 0.5) * width
        
        plt.bar(x_positions, solver_data['Evaluations'], width=width, label=solver, color=COLORS[solver])
    
    plt.xlabel('Instance')
    plt.ylabel('Number of Evaluations')
    plt.title('Number of Function Evaluations by Instance and Solver')
    plt.xticks(np.arange(len(instances)), instances, rotation=45)
    plt.yscale('log')
    plt.legend()
    plt.tight_layout()
    plt.savefig(os.path.join(output_dir, 'evaluations_count.png'), dpi=300)
    plt.close()
    
    # Plot time spent
    plt.figure(figsize=(14, 8))
    
    for i, solver in enumerate(solvers):
        solver_data = perf_data[perf_data['Solver'] == solver]
        x_positions = np.arange(len(instances)) + (i - len(solvers)/2 + 0.5) * width
        
        plt.bar(x_positions, solver_data['TimeMs'], width=width, label=solver, color=COLORS[solver])
    
    plt.xlabel('Instance')
    plt.ylabel('Time (ms)')
    plt.title('Execution Time by Instance and Solver')
    plt.xticks(np.arange(len(instances)), instances, rotation=45)
    plt.yscale('log')
    plt.legend()
    plt.tight_layout()
    plt.savefig(os.path.join(output_dir, 'execution_time.png'), dpi=300)
    plt.close()

def plot_improvement_distribution(df, output_dir):
    plt.figure(figsize=(12, 8))
    
    for solver in df['Solver'].unique():
        solver_data = df[df['Solver'] == solver]
        sns.kdeplot(solver_data['ImprovementPercent'], label=solver, fill=True, color=COLORS[solver], alpha=0.5)
    
    plt.xlabel('Improvement Percentage')
    plt.ylabel('Density')
    plt.xlim(left=0)
    plt.title('Distribution of Improvement Percentage by Solver')
    plt.legend()
    plt.tight_layout()
    plt.savefig(os.path.join(output_dir, 'improvement_distribution.png'), dpi=300)
    plt.close()
    
    # Now plot improvement by instance
    plt.figure(figsize=(14, 8))
    
    # Group by instance and solver
    imp_data = df.groupby(['Instance', 'Solver'])['ImprovementPercent'].mean().reset_index()
    
    # For each solver, create a set of bars
    instances = df['Instance'].unique()
    solvers = df['Solver'].unique()
    width = 0.8 / len(solvers)
    
    for i, solver in enumerate(solvers):
        solver_data = imp_data[imp_data['Solver'] == solver]
        x_positions = np.arange(len(instances)) + (i - len(solvers)/2 + 0.5) * width
        plt.bar(x_positions, solver_data['ImprovementPercent'], width=width, label=solver, color=COLORS[solver])
    
    plt.xlabel('Instance')
    plt.ylabel('Average Improvement (%)')
    plt.title('Average Solution Improvement by Instance and Solver')
    plt.xticks(np.arange(len(instances)), instances, rotation=45)
    plt.legend()
    plt.tight_layout()
    plt.savefig(os.path.join(output_dir, 'improvement_by_instance.png'), dpi=300)
    plt.close()

def create_summary_table(df, output_dir):
    """
    Create a summary table with key metrics
    """
    # Group by instance and solver
    summary = df.groupby(['Instance', 'Solver']).agg({
        'InitialFitness': 'mean',
        'FinalFitness': ['mean', 'min'],
        'TimeMs': 'mean',
        'Evaluations': 'mean',
        'EvalsPerSecond': 'mean',
        'ImprovementPercent': 'mean',
        'GapToBest': 'mean'
    }).reset_index()
    
    # Save to CSV
    summary_path = os.path.join(output_dir, 'summary_metrics.csv')
    summary.to_csv(summary_path, index=False)
    
    return summary

def main(csv_path):
    df = pd.read_csv(csv_path)
    output_dir = os.path.splitext(csv_path)[0] + '_plots'
    os.makedirs(output_dir, exist_ok=True)
    
    df["Solution"] = df["Solution"].apply(clean_solution)
    df, best_solutions = analyze_qap_results(df)
    
    # Create visualizations
    plot_gap_to_best(df, output_dir)
    plot_initial_vs_final(df, output_dir)
    plot_best_solutions(df, best_solutions, output_dir)
    plot_time_efficiency(df, output_dir)
    plot_improvement_distribution(df, output_dir)
    
    # Create summary table
    summary = create_summary_table(df, output_dir)
    print(summary)
    
    print(f"Analysis complete. Results saved to {output_dir}")


if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python script.py <path_to_csv>")
        sys.exit(1)
    
    csv_path = sys.argv[1]

    main(csv_path)
