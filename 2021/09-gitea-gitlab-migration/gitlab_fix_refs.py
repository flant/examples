#!/usr/bin/env python
# -*- coding: utf-8 -*-
import gitlab
import random
import string
import re
import os
import subprocess
import shutil

# gitlab connection
gl = gitlab.Gitlab('https://gitlab.example.com', private_token='superSecret')

shutil.rmtree('code',ignore_errors=True)

# Adding empty commit everywhere
all_orgs = gl.groups.list()
skip_orgs = ['someorg']
for org in all_orgs:
    print(org.name)
    if org.name in skip_orgs:
        print("Skip group", org.name)
        continue
    projects = org.projects.list(all=True)
    for project in projects:
        print(project.name)
        id=project.id
        mrs=gl.projects.get(id=id).mergerequests.list(state='opened', sort='desc',page=1, per_page=10000)
        os.mkdir('code')
        print(subprocess.run(["git", "clone", project.ssh_url_to_repo, "code"], capture_output=True))
        for mr in mrs:
            print(project.name, id, mr.title, mr.source_branch, '=>', mr.target_branch)
            print(subprocess.run(["git", "checkout", mr.source_branch], cwd='code', capture_output=True))
            print(subprocess.run(["git", "pull"], cwd='code', capture_output=True))
            print(subprocess.run(["git", "commit", "--allow-empty", "-m", "Nothing here"], cwd='code', capture_output=True))
            print(subprocess.run(["git", "push"], cwd='code', capture_output=True))
        shutil.rmtree('code',ignore_errors=True)
