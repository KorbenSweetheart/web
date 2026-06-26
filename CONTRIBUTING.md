# Git Cheat Sheet

## The Common Branch Workflow

> [!CAUTION]
> Rule #1: Never commit directly to main. Always work on your own branch.

### Useful commands

```bash
# Check what branches do you have locally
git branch

# Switch to specific branch
git checkout <branch-name>

# It might be good to follow a simplified branch naming.
# Examples: ivan/dijkstra-algorithm, ivan/heap-overflow, ivan/api-documentation

# Create a new branch and switch to it immediately
# Always pull changes from main first, then create a new branch
git pull
git checkout -b <your-name/name-of-the-branch>

# Always before RP run from your branch these commands to get the latest updated from main.
# There might be conflicts, so you need to resolve them, before creating a PR.
git fetch origin
git rebase origin/main

# or
git pull origin main --rebase

# If you branch diverged. e.g. 3/1 (3 behind, and 1 ahead).
# Do the step above and push with --force-with-lease
# Because rebasing rewrites history, a standard git push will be rejected. You must force-push safely:
git push origin <your-branch-name> --force-with-lease

# Get a branch locally
# e.g. if you need to test somebody else code
git fetch origin <branch_name>
git switch -c <branch_name origin/branch_name>

# Delete the branch when done
git branch -d <name-of-the-branch> # Deletes local branch
git push origin --delete <name-of-the-branch> # Deletes remote branch
```

## Merging Code (The Pull Request)

> [!WARNING]
> Always before RP run from your branch `git pull origin main --rebase`

Once your feature is ready:

1. Go to the Gitea web interface and navigate to the `Pull Requests` tab.
2. Click **New Pull Request**. Select `main` as the target and your branch as the source.
3. **Write a good description:** Explicitly state _What_ you changed and _Why_ you changed it. Do not force the reviewer to guess your logic.
4. **Review:** Your partner will review the code, leave comments, and approve it.
5. **Merge:** Once approved, choose the merge option that fits, e.g. "Squash and Merge".
6. **Clean up:** Delete your local and remote branches to keep the repo clean:

```bash
git checkout main
git pull origin main
git branch -d <name-of-the-branch> # Deletes local branch
git push origin --delete <name-of-the-branch> # Deletes remote branch
```

## Commit Message Convention (Optional)

Check [Conventional Commits](https://www.conventionalcommits.org/).

Format: `<type>(<scope>): <subject>`

- `feat:` New feature
- `fix:` Bug fix
- `docs:` Documentation changes
- `style:` Code style changes (formatting, missing semicolons, etc.)
- `refactor:` Refactoring production code (e.g., renaming a variable)
- `test:` Adding or updating tests
- `chore:` Routine tasks (updating dependencies)

**Examples:**
Good: `"feat: implement priority queue for dijkstra algorithm"`
Bad: `"updates"`, `"Add Some Fixes."`

### If you need to move files from one branch to another if needed (e.g. to push some files faster with PR)

[How do I copy a version of a single file from one Git branch to another?](https://stackoverflow.com/questions/307579/how-do-i-copy-a-version-of-a-single-file-from-one-git-branch-to-another)

Run this FROM the branch where you want the file to end up:

```bash
git checkout otherbranch myfile.txt
```

General formulas:

```bash
git checkout <commit_hash> <relative_path_to_file_or_dir>
git checkout <remote_name>/<branch_name> <file_or_dir>
```

Some notes (from comments):

- Using the commit hash, you can pull files from any commit
- This works for files and directories
- Overwrites the file myfile.txt and mydir
- Wildcards don't work, but relative paths do
- Multiple paths can be specified

## Configure Dual Remotes (Optional but Recommended)

We will use Gitea as our main source of truth, but you can sync to your private GitHub repository for your portfolio.

> [!CAUTION]
> Your GitHub repository must be private so other students cannot copy our code.

To avoid forgetting to push to GitHub and causing branch divergence, we can configure Git to push to _both_ servers simultaneously every time you type `git push`:

```bash
# Add GitHub as a secondary push URL
git remote set-url --add --push origin <GITHUB_URL>
```

_(Now, a standard `git push origin main` will update both automatically.)_

### Another option is to configure `git config` in the repository

Run the following command inside the project repository:

```bash
git config --edit
```

Look for the [remote "origin"] section. It likely looks like this:
So you should have 1 fetch and 2 pushurl.
Note: you should edit caps variables.

```text
[remote "origin"]
    url = git@gitea.kood.tech:ivanandreev/PROJECT-NAME.git
    fetch = +refs/heads/*:refs/remotes/origin/*
    pushurl = git@github.com:YOURACCOUNTNAME/PROJECT-NAME.git
    pushurl = git@gitea.kood.tech:ivanandreev/PROJECT-NAME.git
```
