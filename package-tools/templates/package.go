package templates

const ImagesLockTemplate = `#@ load("@ytt:data", "data")
#@ load("package-helpers.lib.yaml", "get_package_repository")

#@ package_repository = get_package_repository(data.values.packageRepository, "")

---
apiVersion: imgpkg.carvel.dev/v1alpha1
images:
#@ for package in package_repository.packages:
  #@ if not hasattr(package, "packageSubVersion") or package.packageSubVersion == "":
  #@ imagePackageVersion = "v" + package.version
  #@ else:
  #@ imagePackageVersion = "v" + package.version + "_" + package.packageSubVersion
  #@ end
- annotations:
    kbld.carvel.dev/id: #@ "{}/{}:{}".format(data.values.registry, package.name, imagePackageVersion)
  image: #@ "{}/{}@sha256:{}".format(data.values.registry, package.name, package.sha256)
#@ end
kind: ImagesLock
`

const PackageCrOverlay = `#@ load("@ytt:data", "data")
#@ load("@ytt:overlay", "overlay")
#@ load("package-helpers.lib.yaml", "get_package_repository", "get_package", "get_package_spec")

#@ package_repository = get_package_repository(data.values.packageRepository, data.values.packageName)
#@ package = get_package(package_repository, data.values.packageName)
#@ packageSpec = get_package_spec(package_repository, package)

#@ if not hasattr(package, "packageSubVersion") and data.values.subVersion == "":
#@  if package.version == "latest":
#@    packageVersion = data.values.version
#@  else:
#@    packageVersion = package.version
#@  end
#@ else:
#@  subVersion = data.values.subVersion
#@  if subVersion == "":
#@    subVersion = package.packageSubVersion
#@  end
#@  if package.version == "latest":
#@    packageVersion = data.values.version + "+" + subVersion
#@  else:
#@    packageVersion = package.version + "+" + subVersion
#@  end
#@ end

#@ if not hasattr(package, "packageSubVersion") and data.values.subVersion == "":
#@  if package.version == "latest":
#@    imagePackageVersion = "v" + data.values.version
#@  else:
#@    imagePackageVersion = "v" + package.version
#@  end
#@ else:
#@  subVersion = data.values.subVersion
#@  if subVersion == "":
#@    subVersion = package.packageSubVersion
#@  end
#@  if package.version == "latest":
#@    imagePackageVersion = "v" + data.values.version + "_" + subVersion
#@  else:
#@    imagePackageVersion = "v" + package.version + "_" + subVersion
#@  end
#@ end

#@ packageLicense = "VMware’s End User License Agreement (Underlying OSS license: Apache License 2.0)"

#@overlay/match by=overlay.subset({"kind":"Package"}),expects=1
---
metadata:
  name: #@ "{}.{}.{}".format(package.name, package_repository.domain, packageVersion)
  #@overlay/match expects="0+"
  #@overlay/remove
  namespace: ""
spec:
  refName: #@ "{}.{}".format(package.name, package_repository.domain)
  version: #@ packageVersion
  #@overlay/match when=0
  releasedAt: #@ data.values.timestamp
  #@overlay/match missing_ok=True
  #@overlay/replace
  licenses:
    -  #@ packageLicense
  template:
    spec:
      #@ if/end packageSpec:
      #@overlay/match missing_ok=True
      #@overlay/remove
      syncPeriod:
      fetch:
        #@overlay/match by=overlay.index(0)
        - imgpkgBundle:
            image: #@ "{}/{}:{}".format(data.values.registry, package.name, imagePackageVersion)
      template:
        #@overlay/match by=overlay.index(0)
        - ytt:
            #@overlay/match missing_ok=True
            ignoreUnknownComments: true
      deploy:
        #@overlay/match by=overlay.index(0)
        - kapp:
            #@ if packageSpec:
            #@overlay/match missing_ok=True
            rawOptions:
              #@overlay/match by=lambda indexOrKey, left, right: "wait-timeout" in left, missing_ok=True
              -  #@ "--wait-timeout={}".format(packageSpec.deploy.kappWaitTimeout)
              #@overlay/match by=lambda indexOrKey, left, right: "kube-api-qps" in left, missing_ok=True
              -  #@ "--kube-api-qps={}".format(packageSpec.deploy.kubeAPIQPS)
              #@overlay/match by=lambda indexOrKey, left, right: "kube-api-burst" in left, missing_ok=True
              -  #@ "--kube-api-burst={}".format(packageSpec.deploy.kubeAPIBurst)
            #@ end
  #@overlay/match missing_ok=True
  valuesSchema:
    #@overlay/match missing_ok=True
    openAPIv3:
      #@overlay/match missing_ok=True
      title: #@ "{}.{}.{} values schema".format(package.name, package_repository.domain, packageVersion)
`

const PackageHelpersLib = `#@ load("@ytt:data", "data")
#@ load("@ytt:assert", "assert")

#@ def get_package_repository(repository_name, package_name):
#@   if repository_name == "":
#@    for repository in data.values.repositories:
#@      for package in data.values.repositories[repository].packages:
#@        if package.name == package_name:
#@          return data.values.repositories[repository]
#@        end
#@      end
#@    end
#@   else:
#@     return data.values.repositories[repository_name]
#@   end
#@ end

#@ def get_package(package_repository, package_name):
#@  for package in package_repository.packages:
#@    if package.name == package_name:
#@      return package
#@    end
#@  end
#@  return None
#@ end

#@ def get_package_spec(package_repository, package):
#@ if hasattr(package, 'spec'):
#@   return package.spec
#@ elif hasattr(package_repository, 'packageSpec'):
#@   return package_repository.packageSpec
#@ end
#@ return None
#@ end
`

const PackageMetadataCrOverlay = `#@ load("@ytt:data", "data")
#@ load("@ytt:overlay", "overlay")
#@ load("package-helpers.lib.yaml", "get_package_repository", "get_package")

#@ package_repository = get_package_repository(data.values.packageRepository, data.values.packageName)
#@ package = get_package(package_repository, data.values.packageName)

#@overlay/match by=overlay.subset({"kind":"PackageMetadata"}),expects=1
---
metadata:
  name: #@ "{}.{}".format(package.name, package_repository.domain)
  #@overlay/match expects="0+"
  #@overlay/remove
  namespace: ""
`

const PackageRepoTemplate = `
#@ load("@ytt:data", "data")
#@ load("package-helpers.lib.yaml", "get_package_repository")

#@ package_repository = get_package_repository(data.values.packageRepository, "")

---
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageRepository
metadata:
  name: #@ "{}.{}".format(package_repository.name, package_repository.domain)
spec:
  fetch:
    imgpkgBundle:
      image: #@ "{}/{}@sha256:{}".format(data.values.registry, package_repository.name, data.values.sha256)
`
