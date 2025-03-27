import os
import sys
import ast
import pandas as pd
import numpy as np
import matplotlib.pyplot as plt
import seaborn as sns

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


def clean_solution(solution_str):
    """
    Cleans the solution string by removing brackets and splitting by spaces.
    Converts each element to an integer.
    Args:
        solution_str (str): The raw solution string (e.g., '[1 5 19 9 15]').
    Returns:
        List of integers representing the solution.
    """
    # Remove brackets and split by spaces
    solution = solution_str.strip('[]').split()
    
    # Convert to integers
    solution = [int(x) for x in solution]
    
    return solution


def analyze_results(csv_path, output_dir):
    """
    Analyze and visualize QAP results
    
    Args:
        csv_path (str): Path to the input CSV file
        output_dir (str): Directory to save output plots
        matrix_A (numpy.ndarray): Matrix A (distances)
        matrix_B (numpy.ndarray): Matrix B (flows)
    """
    # Create output directory if it doesn't exist
    os.makedirs(output_dir, exist_ok=True)
    
    # Read CSV
    df = pd.read_csv(csv_path)
    df["Solution"] = df["Solution"].apply(clean_solution)
    print(df.Solution)
    
    # Compute unique instances
    instances = df['Instance'].unique()
    
    for instance in instances:
        print(f"Processing instance {instance}")
        # Filter data for this instance
        instance_data = df[df['Instance'] == instance]
        
        # Create a new figure for each instance with a 2x2 grid of subplots
        fig, axs = plt.subplots(2, 2, figsize=(15, 10))
        
        # Subplot 1: Initial vs Final Fitness
        sns.scatterplot(x='InitialFitness', y='FinalFitness', 
                        hue='Solver', data=instance_data, ax=axs[0,0], 
                        alpha=0.7, palette='deep')
        axs[0,0].set_title(f'{instance}: Initial vs Final Fitness')
        axs[0,0].set_xlabel('Initial Fitness')
        axs[0,0].set_ylabel('Final Fitness')
        
        # Subplot 2: Runs vs Best Solution
        solver_groups = instance_data.groupby('Solver')
        for solver, group in solver_groups:
            sorted_group = group.sort_values('Run')
            best_solutions = sorted_group['FinalFitness'].cummin()
            axs[0,1].plot(sorted_group['Run'], best_solutions, label=solver)
        
        axs[0,1].set_title(f'{instance}: Runs vs Best Solution')
        axs[0,1].set_xlabel('Number of Runs')
        axs[0,1].set_ylabel('Best Solution Found')
        axs[0,1].legend()
        
        # Subplot 3: Performance Metrics (Relative Deviation)
        optimal_info = load_optimal_solution(instance)
        if optimal_info:
            optimal_fitness, optimal_solution = optimal_info
            
            # Compute relative deviation from optimal
            instance_data.loc[:, 'RelativeDeviation'] = (instance_data['FinalFitness'] - optimal_fitness) / optimal_fitness * 100
            
            sns.boxplot(x='Solver', y='RelativeDeviation', 
                        data=instance_data, ax=axs[1,0])
            axs[1,0].set_title(f'{instance}: Relative Deviation from Optimum')
            axs[1,0].set_xlabel('Solver')
            axs[1,0].set_ylabel('Relative Deviation (%)')
        
        plt.tight_layout()
        plt.savefig(os.path.join(output_dir, f'{instance}_qap_analysis.png'))
        plt.close()

def main():
    if len(sys.argv) < 2:
        print("Usage: python script.py <path_to_csv>")
        sys.exit(1)
    
    csv_path = sys.argv[1]
    output_dir = os.path.splitext(csv_path)[0] + '_plots'
    
    analyze_results(csv_path, output_dir)
    print(f"Plots saved to {output_dir}")

if __name__ == "__main__":
    main()

