#!/usr/bin/env python
# -*- coding: utf-8 -*-
from __future__ import print_function
import time
import giteapy
import gitlab
from giteapy.rest import ApiException
from pprint import pprint
import re
import random
import string
import hashlib
import os
import sys

def get_random_string(length):
    # choose from all lowercase letter
    letters = string.ascii_lowercase
    result_str = ''.join(random.choice(letters) for i in range(length))
    return result_str


def get_hash(key):
    return hashlib.md5(key.encode('utf-8')).hexdigest()

# Gitea connection
configuration = giteapy.Configuration()
configuration.api_key['access_token'] = 'superSecret'
configuration.host = 'https://gitea.example.com/api/v1'
admin_api_instance = giteapy.AdminApi(giteapy.ApiClient(configuration))
user_api_instance = giteapy.UserApi(giteapy.ApiClient(configuration))
org_api_instance = giteapy.OrganizationApi(giteapy.ApiClient(configuration))

# GitLab connection
gl = gitlab.Gitlab('https://gitlab.example.com', private_token='superSecret')

# clean blocked users keys
for block_gl_user in gl.users.list(blocked=True, page=1, per_page=10000):
    print("Blocked user", block_gl_user.username)
    for block_gl_user_key in block_gl_user.keys.list():
        print("Found a key", block_gl_user_key.title)
        block_gl_user_key.delete()

# inspect GitLab
dict_gl_users = dict()
# get Gitea users
gt_users = admin_api_instance.admin_get_all_users()
pattern = re.compile("^id_.+$")

lg_gt_map = dict()
all_gl_keys_dict = dict()
all_problem_keys_dict = list()

# copy keys
for gt_user in gt_users:
    gt_keys = list()
    try:
        gt_keys = user_api_instance.user_list_keys(gt_user.login)
    except:
        pass
    res = gl.users.list(username=gt_user.login)
    if len(res) > 0:
        if res[0].attributes['state'] == 'blocked':
            print("Skip the blocked user", gt_user.login)
            continue
        dict_gl_users[res[0].username] = res[0]
        gl_keys_dict = dict()
        gt_keys_dict = dict()

        lg_gt_map[gt_user.login] = dict_gl_users[gt_user.login].id
        gl_keys = res[0].keys.list()
        keys_to_delete = list()
        for raw_gl_key in gl_keys:
            if pattern.match(raw_gl_key.title):
                print(gt_user.login, "delete the key", raw_gl_key.title)
                raw_gl_key.delete()
            else:
                gl_key_hash = get_hash(raw_gl_key.key.strip().split(' ')[1])
                gl_keys_dict[gl_key_hash] = raw_gl_key
                keys_to_delete.append(gl_key_hash)
                all_gl_keys_dict[gl_key_hash] = res[0].username
        for raw_gt_key in gt_keys:
            gt_key_hash = get_hash(raw_gt_key.key.strip().split(' ')[1])
            gt_keys_dict[gt_key_hash] = raw_gt_key
        for gt_key in gt_keys_dict:
            if gt_key in gl_keys_dict:
                keys_to_delete.remove(gt_key)
            else:
                print(gt_user.login, "missing a key", gt_keys_dict[gt_key].title)
                try:
                    res[0].keys.create({'title': gt_keys_dict[gt_key].title, 'key': gt_keys_dict[gt_key].key})
                    res[0].save()
                except:
                    all_problem_keys_dict.append(gt_key)
                    print(gt_user.login, "can not add the key", gt_keys_dict[gt_key].title)
        for dkey in keys_to_delete:
            if pattern.match(gl_keys_dict[dkey].title):
                print(gt_user.login, "has an additional key", gl_keys_dict[dkey].title)
                #gl_keys_dict[key].delete()
    else:
        print("New user", gt_user.login)

print("Get problematic key", len(all_problem_keys_dict))
for pkey in all_problem_keys_dict:
    if pkey in all_gl_keys_dict:
        print("This key", pkey, "has user", all_gl_keys_dict[pkey])
    else:
        print("Can not find user for key", pkey)
