# Add a package to package repository

This document provides guidance on how to add a package to a package repository and test it.

## Steps to add a package to a package repository

Let's take the example of adding a package to `management` package repository.
Below are steps to illustrate how that can be done.

1. Scaffold a new package and update the package files

   ```shell
      PACKAGE_NAME=my-package make create-package
   ```

   This would scaffold the package under the `packages` directory, and the tree structure of the generated package
   would look something like below:

   ```plain
    packages/my-package
    ├── Makefile
    ├── README.md
    ├── bundle
    │   ├── config
    │   │   ├── overlay
    │   │   ├── upstream # This is the directory to add the package contents using ytt templates.
    │   │   └── values.yaml # Package contents can be configured by providing data values in this file.
    ├── vendir.yml # To fetch config files from a different data source.
    ├── metadata.yaml # To provide high level information description about your package.
    └── package.yaml # Update the Package CR spec to add/modify fields such as releaseNotes etc.
   ```

   After scaffolding the package, the files in your package directory should be updated with the package config.
   Significance of each file is provided in the above tree structure.

   The generated Makefile contains `configure-package` and `reset-package` target to configure the package dynamically,
   which is completely optional.

2. Fetch config files from datasource [optional]

   If you have updated the vendir.yaml to fetch the config from a different source, run

   ```shell
      make package-vendir-sync
   ```

3. Update kbld-config.yaml [optional]

   If the container image in your config needs to be replaced by an image at build time, add an entry like below in the
   kbld-config.yaml file in `packages` directory.

   ```yaml
       - image: my-package-manager:latest
         newImage: ""
   ```
4. Add package-values.yaml file to the packages directory

   Package tooling uses package-values.yaml file to build package and repo bundles. If its a new project this file needs
   to be created under `packages` directory, else skip to step 5.

   Below is a sample package-values.yaml file:

   ```
   #@data/values
   ---
   repositories:
     management:
       name: management
       domain: tanzu.vmware.com
       packageSpec:
         syncPeriod: 5m
         deploy:
           kappWaitTimeout: 5m
           kubeAPIQPS: 20
           kubeAPIBurst: 30
       packages:
    ```

5. Update package-values.yaml to add your package details

   `package-values.yaml` contains Ytt data values for all packages and package repositories.

   Add an entry like below to package-values.yaml under `repositories.<packageRepository>.packages`.

   ```yaml
         #! package name
       - name: my-package
         #! Relative path to package bundle
         path: packages/my-package
         domain: tanzu.vmware.com
         version: latest
         #! this should be name:version(my-package:latest), will be replaced at build time
         sha256: "my-package:latest"
   ```

6. Test the package bundle generation
   Run `all` make target by passing all the vars needed to build the images and package bundle generation.

   ```shell
   make all
   ```

