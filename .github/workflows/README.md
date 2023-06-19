# GH action setup

- copy files release.yaml and main.yaml to repo and drop at /.github/workflows
  
## create secret

Lets create 'RELEASE_SECRET'

    1. go to GH Repo in browser and navigate to 'settings'
    2. select 'secrets and variables' from the left navigation
    3. select 'actions'
    4. Find and select tab 'secrets'
    5. Click button 'New repoository secret'

[Refer](./gh-secret.png) diagram

## export secret

on your release.yaml, add below step (see the env)

      - name: Release with Notes
        uses: softprops/action-gh-release@v1
        with:
          files: |
            bin/kpcli_darwin_arm64.bin
            bin/kpcli_darwin_amd64.bin
            bin/kpcli_windows_amd64.exe
            bin/kpcli_linux_amd64.bin
            README.md
        env:
          GITHUB_TOKEN: ${{ secrets.RELEASE_TOKEN }}
