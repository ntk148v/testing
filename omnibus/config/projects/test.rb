#
# Copyright 2021 YOUR NAME
#
# All Rights Reserved.
#

name "test"
maintainer "CHANGE ME"
homepage "https://CHANGE-ME.com"

# Defaults to C:/test on Windows
# and /opt/test on all other platforms
install_dir "#{default_root}/#{name}"

build_version Omnibus::BuildVersion.semver
build_iteration 1

# Creates required build directories
dependency "preparation"

# test dependencies/components
# dependency "somedep"

exclude "**/.git"
exclude "**/bundler/git"
