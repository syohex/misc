use anyhow::{anyhow, Result};
use std::process::Command;

fn get_current_branch() -> Result<String> {
    let output = Command::new("git").arg("--show-current").output()?;
    let branch = String::from_utf8(output.stdout)?;
    Ok(branch.trim().to_string())
}

fn pull_and_fetch_prone() -> Result<()> {
    let branch = get_current_branch()?;
    let status = Command::new("git")
        .arg("pull")
        .arg("--rebase")
        .arg("origin")
        .arg(&branch)
        .status()?;
    if !status.success() {
        return Err(anyhow!("Failed to 'git pull --rebase origin {branch}'"));
    }

    let status = Command::new("git").arg("fetch").arg("-p").status()?;
    if !status.success() {
        return Err(anyhow!("Failed to 'git fetch -p'"));
    }

    Ok(())
}

fn is_protected_branch(current_branch: &str) -> bool {
    ["main", "master", "develop"].contains(&current_branch)
}

fn merged_branches() -> Result<Vec<String>> {
    let output = Command::new("git").arg("branch").arg("merged").output()?;
    let branches = String::from_utf8(output.stdout)?;

    let mut v = Vec::new();
    for branch in branches.lines() {
        if branch.starts_with('*') {
            // skip current branch
            continue;
        }

        v.push(branch.trim().to_string());
    }

    Ok(v)
}

fn delete_merged_branches() -> Result<()> {
    let merged = merged_branches()?;
    let filtered = merged
        .into_iter()
        .filter(|branch| !["main", "master", "develop"].contains(&branch.as_str()));

    for branch in filtered {
        let status = Command::new("git")
            .arg("branch")
            .arg("-d")
            .arg(&branch)
            .status()?;
        if !status.success() {
            return Err(anyhow!("failed 'git branch -d {branch}'"));
        }
    }

    Ok(())
}

fn main() -> Result<()> {
    pull_and_fetch_prone()?;

    if !is_protected_branch(&get_current_branch().unwrap()) {
        println!("This is not current branch. Skip deleting branches");
        return Ok(());
    }

    delete_merged_branches()?;

    Ok(())
}
