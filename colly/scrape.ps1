Set-Location "I:\coding\hugo\colly"
# Run the Go program and wait for it to finish
Start-Process -FilePath "go" -ArgumentList "run ." -Wait

# Navigate to your repository directory if not already there
Set-Location ..\

# Add changes to git
git add .

# Commit changes with a message
git commit -m "Automated commit scraping data"

# Push changes to the remote repository
git push