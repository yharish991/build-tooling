Extending and building integrations for Tanzu is not easy right now. Someone who wants to add an integration doesn’t
have any standard template/pattern that they can follow or use to get started building that integration.
We want to make it as simple as possible for platform developers to add integrations and extend Tanzu.
This tooling could abstract away much of the build complexity, managing a clear, quick feedback loop for developers,
and enabling release engineers to ship things quickly.

Before we go into how to use the build tooling, lets quickly talk about the contents of this repo


The purpose of this document is to provide a good starting place for someone who wants to use build tooling for Tanzu
integrations.

1. Copy the contents of the Makefile

    This is the first step to consume the build tooling, the `templates` directory in this project's root directory
    contains a Makefile with a bunch of make targets. These make targets are for initializing the build tooling,
    building and publishing the images and packages etc. in a containerized environment so that the builds are
    deterministic and reproducible in any environment.

    To learn more about the make targets in the Makefile, check the documentation [here](./makefile-targets.md).

2. Set COMPONENTS variable

    For the build tooling to understand where your components are located, it needs the `COMPONENTS` variable to be set.
    You can set the `COMPONENTS` variable either in the makefile or as an environment variable. We need to provide the
    component's location, default image name and the package name of the component delimited by a `?` something like
    below:

    ```
    COMPONENTS ?= featuregates?featuregates-controller-manager?featuregates
    ```

    Here `featuregates` is the path to the featuregates component from project's root directory, `featuregates-controller-manager`
    is the default image name, `featuregates` is the name of the package, this should be same as the directory name that
    holds the package definition in `packages` directory that's in the project's root directory.

    If your project has multiple go modules and you want to build images and packages for each of the go module, you can
    do that by setting multiple components to the COMPONENTS variable. For example:

    ```
    COMPONENTS ?= featuregates?featuregates-controller-manager?featuregates capabilities?capabilities-controller-manager?capabilities
    ```

3. Run make init

    To initialize the build tooling we need to run the `init` make target, this fetches the Dockerfile that is used for
    building the image, testing the go module etc. and other templates needed by build tooling.
    It also pulls the packaging image that is needed to build and publish package and repo bundles.

    ```
    make init
    ```

4. Create packages directory and package definition in it
    Check this documentation on how to create a package