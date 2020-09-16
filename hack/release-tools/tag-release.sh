#!/usr/bin/env bash

# Copyright 2020 the Velero contributors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


# This script will do the necessary checks and actions to create a release of Velero. It will:
# - validate that all prerequisites are met
# - verify the version string is what the user expects.
# - create a git tag
# - push the created git tag to GitHub
# - run GoReleaser

# The following variables are needed:

# - $VELERO_VERSION: defines the tag of Velero that any https://github.com/vmware-tanzu/velero/...
#   links in the docs should redirect to.
# - $publish: TRUE/FALSE value where FALSE (or not including it) will indicate a dry-run, and TRUE, or simply adding 'publish',
#   will tag the release with the $VELERO_VERSION and push the tag to a remote named 'upstream'.
# - $GITHUB_TOKEN: Needed to run the goreleaser process to generate a GitHub release. 
#   Use https://github.com/settings/tokens/new?scopes=repo if you don't already have a token.
#   Regenerate an existing token: https://github.com/settings/tokens.
#   You may regenerate the token for every release if you prefer.
#   See https://goreleaser.com/environment/ for more details.

# This script is meant to be a combination of documentation and executable.
# If you have questions at any point, please stop and ask!

# Directory in which the script itself resides, so we can use it for calling programs that are in the same directory.
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

# Parse out the branch we're on so we can switch back to it at the end of a dry-run, where we delete the tag. Requires git v1.8.1+
upstream_branch=$(git symbolic-ref --short HEAD)

function tag_and_push() {
    echo "Tagging $VELERO_VERSION"
    git tag $VELERO_VERSION || true
    
    if [[ $publish == "TRUE" ]]; then
        echo "Pushing $VELERO_VERSION"
        git push upstream $VELERO_VERSION
    fi
}

# Default to a dry-run mode
publish=FALSE
if [[ "$1" = "publish" ]]; then
    publish=TRUE
fi

# For now, have the person doing the release pass in the VELERO_VERSION variable as an environment variable.
# In the future, we might be able to inspect git via `git describe --abbrev=0` to get a hint for it.
if [[ -z "$VELERO_VERSION" ]]; then
    printf "The \$VELERO_VERSION environment variable is not set. Please set it with\n\texport VELERO_VERSION=v<version.to.release>\nthen try again."
    exit 1
fi

# Make sure the user's provided their github token, so we can give it to goreleaser.
if [[ -z "$GITHUB_TOKEN" ]]; then
    printf "The GITHUB_TOKEN environment variable is not set. Please set it with\n\t export GITHUB_TOKEN=<your github token>\n then try again."
    exit 1
fi

# Ensure that we have a clean working tree before we let any changes happen, especially important for cutting release branches.
if [[ -n $(git status --short) ]]; then 
    echo "Your git working directory is dirty! Please clean up untracked files and stash any changes before proceeding."
    exit 3
fi

# Make sure that there's no issue with the environment variable's format before trying to eval the parsed version.
if ! go run $DIR/chk_version.go --verify;  then
    exit 2
fi
# Since we're past the validation of the VELERO_VERSION, parse the version's individual components.
eval $(go run $DIR/chk_version.go)


printf "To clarify, you've provided a version string of $VELERO_VERSION.\n"
printf "Based on this, the following assumptions have been made: \n"

[[ "$VELERO_PATCH" != 0 ]] && printf "*\t This is a patch release.\n"

# $VELERO_PRERELEASE gets populated by the chk_version.go script that parses and verifies the given version format 
# -n is "string is non-empty"
[[ -n $VELERO_PRERELEASE ]] && printf "*\t This is a pre-release.\n"

# -z is "string is empty"
[[ -z $VELERO_PRERELEASE ]] && printf "*\t This is a GA release.\n"

if [[ $publish == "TRUE" ]]; then
    echo "If this is all correct, press enter/return to proceed to TAG THE RELEASE and UPLOAD THE TAG TO GITHUB."
else
    echo "If this is all correct, press enter/return to proceed to TAG THE RELEASE and PROCEED WITH THE DRY-RUN."
fi

echo "Otherwise, press ctrl-c to CANCEL the process without making any changes."

read -p "Ready to continue? "

echo "Alright, let's go."

echo "Pulling down all git tags and branches before doing any work."
git fetch upstream --tags

# $VELERO_PATCH gets populated by the chk_version.go scrip that parses and verifies the given version format 
# If we've got a patch release, we'll need to create a release branch for it.
if [[ "$VELERO_PATCH" > 0 ]]; then
    release_branch_name=release-$VELERO_MAJOR.$VELERO_MINOR
  
    # Check if the branch exists, creating it if not.
    # The fetch command above should have gotten all the upstream branches, so we can safely assume this check is local & upstream branches.
    if [[ -z $(git branch | grep $release_branch_name) ]]; then
        git checkout -b $release_branch_name
        echo "Release branch made."
    else
        echo "Release branch $release_branch_name exists already."
        git checkout $release_branch_name
    fi

    echo "Now you'll need to cherry-pick any relevant git commits into this release branch."
    echo "Either pause this script with ctrl-z, or open a new terminal window and do the cherry-picking."
    if [[ $publish == "TRUE" ]]; then
        read -p "Press enter when you're done cherry-picking. THIS WILL MAKE A TAG PUSH THE BRANCH TO upstream"
    else
        read -p "Press enter when you're done cherry-picking."
    fi

    # TODO can/should we add a way to review the cherry-picked commits before the push?

    if [[ $publish == "TRUE" ]]; then
        echo "Pushing $release_branch_name to upstream remote"
        git push --set-upstream upstream/$release_branch_name $release_branch_name
    fi
    
    tag_and_push
else
    echo "Checking out upstream/main."
    git checkout upstream/main

    tag_and_push
fi


echo "Invoking Goreleaser to create the GitHub release."
RELEASE_NOTES_FILE=changelogs/CHANGELOG-$VELERO_MAJOR.$VELERO_MINOR.md \
    PUBLISH=$publish \
    make release

if [[ $publish == "FALSE" ]]; then
    # Delete the local tag so we don't potentially conflict when it's re-run for real.
    # This also means we won't have to just ignore existing tags in tag_and_push, which could be a problem if there's an existing tag.
    echo "Dry run complete. Deleting git tag $VELERO_VERSION"
    git checkout $upstream_branch
    git tag -d $VELERO_VERSION
fi

