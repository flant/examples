#!/usr/bin/env python
# -*- coding: utf-8 -*-
from __future__ import print_function
import time
import giteapy
from giteapy.rest import ApiException
from pprint import pprint

import random
import string

# Gitea connection
configuration = giteapy.Configuration()
configuration.api_key['access_token'] = 'superSecret'
configuration.host = 'https://gitea.example.com/api/v1'
admin_api_instance = giteapy.AdminApi(giteapy.ApiClient(configuration))
user_api_instance = giteapy.UserApi(giteapy.ApiClient(configuration))
org_api_instance = giteapy.OrganizationApi(giteapy.ApiClient(configuration))
repo_api_instance = giteapy.RepositoryApi(giteapy.ApiClient(configuration))

# adding RO user for migrations in any repo
all_orgs = admin_api_instance.admin_get_all_orgs()
for org in all_orgs:
    print(org.username)
    for repo in org_api_instance.org_list_repos(org.username):
        print(repo)
        body = giteapy.AddCollaboratorOption()
        repo_api_instance.repo_add_collaborator(repo.owner.login, repo.name, 'migrationuser', body=body)

    teams = org_api_instance.org_list_teams(org.username)
    for team in teams:
        members = org_api_instance.org_list_team_members(team.id)
        for user in members:
            for repo in user_api_instance.user_list_repos(user.login):
                repo_api_instance.repo_add_collaborator(repo.owner.login, repo.name, 'migrationuser', body=body)
