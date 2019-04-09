---
about: 'Provide all the information about how to proceed with the bump of a version'
name: 'Bump version'
title: 'Bump [major|minor|patch] version'
labels: bump
assignees: mdelapenya
---

## Bump checklist
- [ ] VERSION.txt updated with bumped version
- [ ] Static files bundled into the Go binary (as shown in `./scripts/bind-assets.sh`)
- [ ] Git commit including the bump
- [ ] Git tag using version as tag name
- [ ] Git tag pushed to upstream
- [ ] Github release created for the tag
- [ ] Binaries generated for the new version
- [ ] Binaries uploaded to the proper Github release

