# Contributing to Go Serve

First. Welcome !! If you are here is probably you have something to say about go serve. That's amazing ! However, there are some common
rules we would like to maintain during the contributing process. Keep reading.

## Working with this repo

All the internal repository workflows are described in the [Makefile](Makefile) of this repo.

Regarding Git, we work with [trunk based development](https://trunkbaseddevelopment.com/) principle. We only work with the master or main
branch, and feature branches. There's no concept of "release" branches.

## New releases

All things should be automated to just push a new tag in master after the new version merged to the master/main branch. **Please** remember to
update the [changelog](CHANGELOG.md) . We follow [semver](https://semver.org/) for versioning.

## Issues format

Use the issue templates provided on this repo as a guideline.

## New features proposal

Unless you are going to implement the new feature in your fork anyway, we recommend to first create a discussion by creating an issue with
the label `enhancement`. We need to consider the proposed feature is aligned with the overall idea of this repo:
 * keep it as simple as possible.
 * Preserve feature orthogonality. 

## Bugs

Just follow the bug issue template. If possible, we would love to have some test code that checks for the error.

## Code of conduct

Show empathy to others.
