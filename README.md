#async-tf-run

### Overview

The ideia is to specify the project names containing `.tf` files under the `task` string while using the async_terraform action. (there is an example below).

Then the action will run all projects concurrently based on the number os workers configured.

### Usage

```
The Action must have four parameters
1. workers -> Specifies the number of workers running concurrently in a worker pool to execute the terraform tasks. The default value is 1.
2. verb -> Specifies the action to be performed by Terraform. Plan, apply or destroy. This is set manually by input using workflow_dispatch as the example below
3. tasks -> The list of tasks (terraform projects) to be read by the action. The default is "."
4. version -> Specifies the terraform version. The default is "1.9.5"
```

**Below is the folder structure example to use the action**
Remember to replace each string in "tasks" with each name of your project. In the example below I'm using terraform_1, terraform_2, terraform_3, terraform_4

### Examples

**The credentials step below can be modified to autenticate in another cloud provider**

Example using the root path if the user need to run only one project
If the user need to run only a single project, there's no need to set workers (default is 1), tasks (default is root directory).

```
├── .github
│   └── workflows
│       └── action.yml
├── main.tf
└── provider.tf
```

Action example:

```yaml
on:
  workflow_dispatch:
    inputs:
      verb:
        description: "Plan, apply or destroy"
        required: true
        type: choice
        options:
          - plan
          - apply
          - -----
          - destroy

name: Async_terraform
jobs:
  terraform:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: AWS Credentials
        uses: aws-actions/configure-aws-credentials@v4
        with:
          aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
          aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          aws-region: us-east-1

      - name: Terraform Task
        uses: adalbertjnr/tf-async-run@v1
        with:
          verb: ${{ inputs.verb }}
          version: "1.9.5"
```

Example for few terraform projects.

```
├── .github
│   └── workflows
│       └── action.yml
├── terraform_1
│   ├── providers.tf
│   └── vpc.tf
├── terraform_2
│   ├── providers.tf
│   └── vpc.tf
├── terraform_3
│   ├── providers.tf
│   └── vpc.tf
├── terraform_4
│   ├── providers.tf
│   └── vpc.tf

```

Action example:

```yaml
on:
  workflow_dispatch:
    inputs:
      verb:
        description: "Plan, apply or destroy"
        required: true
        type: choice
        options:
          - plan
          - apply
          - -----
          - destroy

name: Async_terraform
jobs:
  terraform:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: AWS Credentials
        uses: aws-actions/configure-aws-credentials@v4
        with:
          aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
          aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          aws-region: us-east-1

      - name: Terraform Task
        uses: adalbertjnr/tf-async-run@v1
        with:
          workers: 2
          verb: ${{ inputs.verb }}
          version: "1.9.5"
          tasks: |
            terraform_1,
            terraform_2,
            terraform_3,
            terraform_4
```
