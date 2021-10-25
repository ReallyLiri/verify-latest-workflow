# Verify Latest Workflow

Github action to verify the successful status of the latest run of a given workflow.

Built using [ezactions](https://github.com/WillAbides/ezactions).

## Usage

Following environment variables are used:

* `GITHUB_TOKEN`
* `OWNER`
* `REPOSITORY`
* `WORKFLOW` - workflow file name in `owner/repostiory`
* `BRANCH` - optional, if unset, all branches will be considered
* `EVENT` - optional, if unset, all triggering events will be considered

Workflow example:

```yaml
name: myflow
on:
  ...
jobs:
  myjob:
    runs-on: ubuntu-latest
    steps:
      - name: verify-status
        uses: "ReallyLiri/verify-latest-workflow@v1.0"
        env:
          GITHUB_TOKEN: "${{ secrets.GITHUB_TOKEN }}"
          OWNER: "ReallyLiri"
          REPOSITORY: "kubescout"
          WORKFLOW: "ci.yaml"
          BRANCH: "main"
```

## Code generation

To generate action code: `go generate .`

To run manually, set the env var `MANUAL` to any value (along with the other required env vars).

```bash
MANUAL=true OWNER=... go run .
```
