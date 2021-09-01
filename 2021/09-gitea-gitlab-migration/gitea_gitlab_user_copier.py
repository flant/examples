#!/usr/bin/env python
# -*- coding: utf-8 -*-
from __future__ import print_function
import time
import giteapy
import gitlab
from giteapy.rest import ApiException
from pprint import pprint

import random
import string

def get_random_string(length):
    # choose from all lowercase letter
    letters = string.ascii_lowercase
    result_str = ''.join(random.choice(letters) for i in range(length))
    return result_str

# gitea connection
configuration = giteapy.Configuration()
configuration.api_key['access_token'] = 'superSecret'
configuration.host = 'https://gitea.example.com/api/v1'
admin_api_instance = giteapy.AdminApi(giteapy.ApiClient(configuration))
user_api_instance = giteapy.UserApi(giteapy.ApiClient(configuration))
org_api_instance = giteapy.OrganizationApi(giteapy.ApiClient(configuration))

# gitlab connection
gl = gitlab.Gitlab('https://gitlab.example.com', private_token='superSecret')

# inspect gitlab
dict_gl_users = dict()
# get gitea users
gt_users = admin_api_instance.admin_get_all_users()
#pprint(api_response)
lg_gt_map = dict()
# copy users
for gt_user in gt_users:
    pprint(gt_user)
    key = None
    try:
        key = user_api_instance.user_current_get_key(gt_user.id)
    except:
        pass
    pprint(key)
    res = gl.users.list(username=gt_user.login)
    if len(res) > 0:
        dict_gl_users[res[0].username] = res[0]
        lg_gt_map[gt_user.login] = dict_gl_users[gt_user.login].id
    else:
        password = get_random_string(16)
        gl_user = gl.users.create({'email': gt_user.email,
                                   'password': password,
                                   'username': gt_user.login,
                                   'name': gt_user.full_name if len(gt_user.full_name) > 0 else gt_user.login,
                                   'admin': gt_user.is_admin})
        gl_user.save()
        print(gt_user.email,password)
        lg_gt_map[gt_user.email] = gl_user.id
        dict_gl_users[gt_user.login] = gl_user
    if key:
        exisint_keys = list(map(lambda x: x.title,dict_gl_users[gt_user.login].keys.list()))
        key_title = 'id_{}'.format(key.id)
        if not key_title in exisint_keys:
            dict_gl_users[gt_user.login].keys.create({'title': key_title, 'key': key.key})

# inspect groups, add users to groups

dict_gl_groups = dict()
map_access = {'Owners': gitlab.OWNER_ACCESS, 'Developers': gitlab.DEVELOPER_ACCESS, 'QA': gitlab.DEVELOPER_ACCESS, 'Manager':gitlab.REPORTER_ACCESS, 'Managers': gitlab.REPORTER_ACCESS, 'Dev': gitlab.DEVELOPER_ACCESS, 'Services': gitlab.REPORTER_ACCESS, 'ml-outsource': gitlab.DEVELOPER_ACCESS, 'services': gitlab.REPORTER_ACCESS}
all_orgs = admin_api_instance.admin_get_all_orgs()
for org in all_orgs:
    print(org.username)
    res = None
    try:
        res = gl.groups.get(org.username)
    except:
        pass

    if res:
        dict_gl_groups[org.username] = res
    else:
        group = gl.groups.create({'name': org.username, 'path': org.username})
        if len(org.description) > 0:
            group.description = org.description
        if len(org.full_name) > 0:
            group.full_name = org.full_name
        group.save()
        dict_gl_groups[org.username] = group

    for group in dict_gl_groups:
        print('GL group',dict_gl_groups[group].name, dict_gl_groups[group].path, dict_gl_groups[group].id)

    teams = org_api_instance.org_list_teams(org.username)
    for team in teams:
        members = org_api_instance.org_list_team_members(team.id)
        for user in members:
            print(org.username,team.name,user.login)
            print(dict_gl_users[user.login].id)
            member = None
            try:
                member = dict_gl_groups[org.username].members.get(dict_gl_users[user.login].id)
            except:
                pass
            if not member:
                member = dict_gl_groups[org.username].members.create({'user_id': dict_gl_users[user.login].id, 'access_level': map_access[team.name]})
