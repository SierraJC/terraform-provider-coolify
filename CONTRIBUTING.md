# Contributing

This project is a community-driven initiative to enable Terraform to manage Coolify v4 (beta) resources. We warmly welcome contributions!

Please note that this provider is currently limited by the Coolify API, which is still under development. As the API matures, more resources can be added to this provider. The project follows Semantic Versioning and is currently in its `0.x.x` phase, meaning the API should be considered unstable and potentially subject to breaking changes.

To contribute, please follow these steps:

## Development Workflow

1.  **Fork the repository:** Click the "Fork" button on the top right of the repository page on GitHub. This will create a copy of the repository in your own GitHub account.
2.  **Clone your fork:** Clone the forked repository to your local machine:
    ```bash
    git clone https://github.com/YOUR_USERNAME/terraform-provider-coolify.git
    cd terraform-provider-coolify
    ```
3.  **Create a new branch:** Create a new branch for your changes. Choose a descriptive branch name.
    ```bash
    git checkout -b my-feature-branch
    ```
4.  **Make your changes:** Implement your feature or bug fix.
5.  **Update generated files:** If you've made changes that affect generated code or documentation, run the following command:
    ```bash
    make generate
    ```
    Commit any changes made by this command.
6.  **Run tests:** Ensure all tests pass before submitting your changes.
    ```bash
    make test
    ```
7.  **Commit your changes:** Commit your changes with a clear and concise commit message.
    ```bash
    git add .
    git commit -m "feat: Describe your feature or fix"
    ```
8.  **Push your changes:** Push your changes to your forked repository.
    ```bash
    git push origin my-feature-branch
    ```
9.  **Open a Pull Request:** Go to the original repository on GitHub and click the "New pull request" button. Select your branch and provide a clear description of your changes.

