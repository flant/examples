from functools import partial
from operator import methodcaller

import kubernetes

from copyrator.const import ALLOWED_EVENT_TYPES, CREATE_TYPES_MAP, \
    LIST_TYPES_MAP

__all__ = [
    'handle',
]


def handle_event(v1, specs, event):
    """
    The method for processing one Kubernetes event.
    """
    if event['type'] not in ALLOWED_EVENT_TYPES:
        return

    object_ = event['object']
    labels = object_['metadata'].get('labels', {})

    # Look for the matches using selector
    for key, value in specs['selector'].items():
        if labels.get(key) != value:
            return
    # Get active namespaces
    namespaces = map(
        lambda x: x.metadata.name,
        filter(
            lambda x: x.status.phase == 'Active',
            v1.list_namespace().items
        )
    )
    for namespace in namespaces:
        # Clear the metadata, set the namespace
        object_['metadata'] = {
            'labels': object_['metadata']['labels'],
            'namespace': namespace,
            'name': object_['metadata']['name'],
        }
        # Call the method for creating/updating an object
        methodcaller(
            CREATE_TYPES_MAP[specs['ruleType']],
            namespace,
            object_
        )(v1)


def handle(specs):
    """
    The main method to initiate events processing via operator.
    """
    kubernetes.config.load_incluster_config()
    v1 = kubernetes.client.CoreV1Api()

    # Get the method to watch the objects
    method = getattr(v1, LIST_TYPES_MAP[specs['ruleType']])
    func = partial(method, specs['namespace'])

    w = kubernetes.watch.Watch()
    for event in w.stream(func, _request_timeout=60):
        handle_event(v1, specs, event)
