#!/usr/bin/env python3
import os
import sys
import requests

def get_module_base_url(line):
    parts = line.split()
    if len(parts) >= 2 and not line.startswith('require (') and not line.startswith(')'):
        # Extract the base module path without versioning
        module_path = parts[0].split("/v")[0]  # Stripping version suffix
        if module_path.startswith('github.com/'):
            repo_path = module_path.replace("github.com/", "raw.githubusercontent.com/")
            return repo_path
    return None

def download_file(url, save_path):
    try:
        response = requests.get(url)
        if response.status_code == 200:
            with open(save_path, 'w') as file:
                file.write(response.text)
            return True
    except requests.exceptions.RequestException as e:
        print(f"Request error: {e}")
    return False

def download_license_files(go_mod_file, save_dir):
    os.makedirs(save_dir, exist_ok=True)
    
    with open(go_mod_file, 'r') as file:
        for line in file:
            if line.startswith('\t') or line.startswith('github.com'):
                repo_path = get_module_base_url(line)
                if repo_path:
                    filename = repo_path.split('/')[-1] + "_LICENSE"
                    file_path = os.path.join(save_dir, filename)
                    
                    # Try downloading from the main branch, then master if main fails
                    for branch in ['main', 'master']:
                        license_url = f"https://{repo_path}/{branch}/LICENSE"
                        license_md_url = f"https://{repo_path}/{branch}/LICENSE.md"
                        
                        if download_file(license_url, file_path) or download_file(license_md_url, file_path):
                            print(f"Downloaded license for {repo_path}")
                            break
                    else:
                        print(f"Failed to download license for {repo_path}")

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python script.py <path_to_go.mod>")
        sys.exit(1)

    go_mod_file = sys.argv[1]
    save_dir = "."
    download_license_files(go_mod_file, save_dir)
