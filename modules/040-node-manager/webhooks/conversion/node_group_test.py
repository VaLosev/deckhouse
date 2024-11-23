#!/usr/bin/env python3

# Copyright 2024 Flant JSC
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

import unittest
import typing

from dotmap import DotMap
from deckhouse import hook
from node_group import main, NodeGroupConversionDispatcher


def _assert_conversion(t: unittest.TestCase, o: hook.Output, res: dict | typing.List[dict], failed_msg: str | None):
    d = o.conversions.data

    t.assertEqual(len(d), 1)

    if not failed_msg is None:
        t.assertEqual(len(d[0]), 1)
        t.assertEqual(d[0]["failedMessage"], failed_msg)
        return

    expected = res
    if isinstance(res, dict):
        expected = [res]

    t.assertEqual(d[0]["convertedObjects"], expected)


def test_dispatcher_for_unit_tests(snapshots: dict | None) -> NodeGroupConversionDispatcher:
    output = hook.Output(
        hook.MetricsCollector(),
        hook.KubeOperationCollector(),
        hook.ValuesPatchesCollector({}),
        hook.ConversionsCollector(),
        hook.ValidationsCollector(),
    )

    bctx = {}
    if snapshots is not None:
        bctx = {"snapshots": snapshots}

    return NodeGroupConversionDispatcher(hook.Context(bctx, {}, {}, output))



class TestUnitAlpha1ToAlpha2Method(unittest.TestCase):
    def test_not_alpha1_should_return_same_object(self):
        obj = {
            "apiVersion": "deckhouse.io/v1alpha2",
            "kind": "NodeGroup",
            "metadata": {
                "name": "worker-static",
            },
            "spec": {
                "disruptions": {
                    "approvalMode": "Automatic"
                },
            }
        }

        err, res_obj = test_dispatcher_for_unit_tests(None).alpha1_to_alpha2(DotMap(obj))

        self.assertIsNone(err)
        self.assertEqual(obj, res_obj)


    def test_should_remove_static_and_spec_kubernetes_version(self):
        obj = {
            "apiVersion": "deckhouse.io/v1alpha1",
            "kind": "NodeGroup",
            "metadata": {
                "name": "worker-static",
            },
            "spec": {
                "disruptions": {
                    "approvalMode": "Automatic"
                },
                "kubernetesVersion": "1.11"
            },
            "static": {
                "internalNetworkCIDR": "127.0.0.1/8"
            }
        }

        err, res_obj = test_dispatcher_for_unit_tests(None).alpha1_to_alpha2(obj)

        self.assertIsNone(err)
        self.assertEqual(res_obj, {
            "apiVersion": "deckhouse.io/v1alpha2",
            "kind": "NodeGroup",
            "metadata": {
                "name": "worker-static",
            },
            "spec": {
                "disruptions": {
                    "approvalMode": "Automatic"
                },
            },
        })


    def test_should_move_docker_to_cri_docker_cri_not_set(self):
        obj = {
            "apiVersion": "deckhouse.io/v1alpha1",
            "kind": "NodeGroup",
            "metadata": {
                "name": "worker-static",
            },
            "spec": {
                "disruptions": {
                    "approvalMode": "Automatic"
                },
                "docker": {
                    "manage": False,
                    "maxConcurrentDownloads": 4
                }
            },

        }

        err, res_obj = test_dispatcher_for_unit_tests(None).alpha1_to_alpha2(obj)

        self.assertIsNone(err)
        self.assertEqual(res_obj, {
            "apiVersion": "deckhouse.io/v1alpha2",
            "kind": "NodeGroup",
            "metadata": {
                "name": "worker-static",
            },
            "spec": {
                "disruptions": {
                    "approvalMode": "Automatic"
                },
                "cri": {
                    "docker": {
                        "manage": False,
                        "maxConcurrentDownloads": 4
                    }
                }
            },
        })


class TestUnitAlpha2ToAlpha1Method(unittest.TestCase):
    def test_not_alpha2_should_return_same_object(self):
        obj = {
            "apiVersion": "deckhouse.io/v1alpha1",
            "kind": "NodeGroup",
            "metadata": {
                "name": "worker-static",
            },
            "spec": {
                "disruptions": {
                    "approvalMode": "Automatic"
                },
            }
        }

        err, res_obj = test_dispatcher_for_unit_tests(None).alpha2_to_alpha1(obj)

        self.assertIsNone(err)
        self.assertEqual(obj, res_obj)


    def test_should_move_spec_cri_docker_to_spec_docker(self):
        obj = {
            "apiVersion": "deckhouse.io/v1alpha2",
            "kind": "NodeGroup",
            "metadata": {
                "name": "worker-static",
            },
            "spec": {
                "disruptions": {
                    "approvalMode": "Automatic"
                },
                "cri": {
                    "docker": {
                        "manage": False,
                        "maxConcurrentDownloads": 4
                    }
                }
            },
        }

        err, res_obj = test_dispatcher_for_unit_tests(None).alpha2_to_alpha1(obj)

        self.assertIsNone(err)
        self.assertEqual(res_obj, {
            "apiVersion": "deckhouse.io/v1alpha1",
            "kind": "NodeGroup",
            "metadata": {
                "name": "worker-static",
            },
            "spec": {
                "disruptions": {
                    "approvalMode": "Automatic"
                },
                "cri": {},
                "docker": {
                    "manage": False,
                    "maxConcurrentDownloads": 4
                }
            },

        })


class TestUnitAlpha2ToV1Method(unittest.TestCase):
    __snapshots = {
        "cluster_config": [
            {
                "filterResult": "YXBpVmVyc2lvbjogZGVja2hvdXNlLmlvL3YxYWxwaGExCmtpbmQ6IE9wZW5TdGFja0NsdXN0ZXJDb25maWd1cmF0aW9uCmxheW91dDogU3RhbmRhcmQKbWFzdGVyTm9kZUdyb3VwOgogIGluc3RhbmNlQ2xhc3M6CiAgICBhZGRpdGlvbmFsU2VjdXJpdHlHcm91cHM6CiAgICAtIHNhbmRib3gKICAgIC0gc2FuZGJveC1mcm9udGVuZAogICAgLSB0c3Qtc2VjLWdyb3VwCiAgICBldGNkRGlza1NpemVHYjogMTAKICAgIGZsYXZvck5hbWU6IG0xLmxhcmdlLTUwZwogICAgaW1hZ2VOYW1lOiB1YnVudHUtMTgtMDQtY2xvdWQtYW1kNjQKICByZXBsaWNhczogMQogIHZvbHVtZVR5cGVNYXA6CiAgICBub3ZhOiBjZXBoLXNzZApub2RlR3JvdXBzOgotIGluc3RhbmNlQ2xhc3M6CiAgICBhZGRpdGlvbmFsU2VjdXJpdHlHcm91cHM6CiAgICAtIHNhbmRib3gKICAgIC0gc2FuZGJveC1mcm9udGVuZAogICAgY29uZmlnRHJpdmU6IGZhbHNlCiAgICBmbGF2b3JOYW1lOiBtMS54c21hbGwKICAgIGltYWdlTmFtZTogdWJ1bnR1LTE4LTA0LWNsb3VkLWFtZDY0CiAgICBtYWluTmV0d29yazogc2FuZGJveAogICAgcm9vdERpc2tTaXplOiAxNQogIG5hbWU6IGZyb250LW5tCiAgbm9kZVRlbXBsYXRlOgogICAgbGFiZWxzOgogICAgICBhYWE6IGFhYWEKICAgICAgY2NjOiBjY2NjCiAgcmVwbGljYXM6IDIKICB2b2x1bWVUeXBlTWFwOgogICAgbm92YTogY2VwaC1zc2QKcHJvdmlkZXI6CiAgYXV0aFVSTDogaHR0cHM6Ly9jbG91ZC5leGFtcGxlLmNvbS92My8KICBkb21haW5OYW1lOiBEZWZhdWx0CiAgcGFzc3dvcmQ6IHBhc3N3b3JkCiAgcmVnaW9uOiByZWcKICB0ZW5hbnROYW1lOiB1c2VyCiAgdXNlcm5hbWU6IHVzZXIKc3NoUHVibGljS2V5OiBzc2gtcnNhIEFBQQpzdGFuZGFyZDoKICBleHRlcm5hbE5ldHdvcmtOYW1lOiBwdWJsaWMKICBpbnRlcm5hbE5ldHdvcmtDSURSOiAxOTIuMTY4LjE5OC4wLzI0CiAgaW50ZXJuYWxOZXR3b3JrRE5TU2VydmVyczoKICAtIDguOC44LjgKICAtIDEuMS4xLjEKICAtIDguOC40LjQKICBpbnRlcm5hbE5ldHdvcmtTZWN1cml0eTogdHJ1ZQp0YWdzOgogIGE6IGIK"
            }
        ]
    }

    def test_not_alpha2_should_return_same_object(self):
        obj = {
            "apiVersion": "deckhouse.io/v1alpha1",
            "kind": "NodeGroup",
            "metadata": {
                "name": "worker-static",
            },
            "spec": {
                "disruptions": {
                    "approvalMode": "Automatic"
                },
            }
        }

        err, res_obj = test_dispatcher_for_unit_tests(TestUnitAlpha2ToV1Method.__snapshots).alpha2_to_v1(obj)

        self.assertIsNone(err)
        self.assertEqual(obj, res_obj)


    def test_change_node_type_from_cloud_to_cloud_ephemeral(self):
        obj = {
            "apiVersion": "deckhouse.io/v1alpha2",
            "kind": "NodeGroup",
            "metadata": {
                "name": "worker-static",
            },
            "spec": {
                "disruptions": {
                    "approvalMode": "Automatic"
                },
                "cri": {
                    "docker": {
                        "manage": False,
                        "maxConcurrentDownloads": 4
                    }
                },
                "nodeType": "Cloud"
            },
        }

        err, res_obj = test_dispatcher_for_unit_tests(TestUnitAlpha2ToV1Method.__snapshots).alpha2_to_v1(obj)

        self.assertIsNone(err)
        self.assertEqual(res_obj, {
            "apiVersion": "deckhouse.io/v1",
            "kind": "NodeGroup",
            "metadata": {
                "name": "worker-static",
            },
            "spec": {
                "disruptions": {
                    "approvalMode": "Automatic"
                },
                "cri": {
                    "docker": {
                        "manage": False,
                        "maxConcurrentDownloads": 4
                    }
                },
                "nodeType": "CloudEphemeral"
            },

        })


    def test_change_node_type_from_hybrid_to_cloud_permanent_for_master_ng(self):
        obj = {
            "apiVersion": "deckhouse.io/v1alpha2",
            "kind": "NodeGroup",
            "metadata": {
                "name": "master",
            },
            "spec": {
                "disruptions": {
                    "approvalMode": "Automatic"
                },
                "cri": {
                    "docker": {
                        "manage": False,
                        "maxConcurrentDownloads": 4
                    }
                },
                "nodeType": "Hybrid"
            },
        }

        err, res_obj = test_dispatcher_for_unit_tests(TestUnitAlpha2ToV1Method.__snapshots).alpha2_to_v1(obj)

        self.assertIsNone(err)
        self.assertEqual(res_obj, {
            "apiVersion": "deckhouse.io/v1",
            "kind": "NodeGroup",
            "metadata": {
                "name": "master",
            },
            "spec": {
                "disruptions": {
                    "approvalMode": "Automatic"
                },
                "cri": {
                    "docker": {
                        "manage": False,
                        "maxConcurrentDownloads": 4
                    }
                },
                "nodeType": "CloudPermanent"
            },

        })


    def test_change_node_type_from_hybrid_to_cloud_permanent_for_ng_in_provider_cluster_config(self):
        obj = {
            "apiVersion": "deckhouse.io/v1alpha2",
            "kind": "NodeGroup",
            "metadata": {
                "name": "front-nm",
            },
            "spec": {
                "disruptions": {
                    "approvalMode": "Automatic"
                },
                "cri": {
                    "docker": {
                        "manage": False,
                        "maxConcurrentDownloads": 4
                    }
                },
                "nodeType": "Hybrid"
            },
        }

        err, res_obj = test_dispatcher_for_unit_tests(TestUnitAlpha2ToV1Method.__snapshots).alpha2_to_v1(obj)

        self.assertIsNone(err)
        self.assertEqual(res_obj, {
            "apiVersion": "deckhouse.io/v1",
            "kind": "NodeGroup",
            "metadata": {
                "name": "front-nm",
            },
            "spec": {
                "disruptions": {
                    "approvalMode": "Automatic"
                },
                "cri": {
                    "docker": {
                        "manage": False,
                        "maxConcurrentDownloads": 4
                    }
                },
                "nodeType": "CloudPermanent"
            },

        })


    def test_change_node_type_from_hybrid_to_cloud_static_for_ng_not_in_provider_cluster_config(self):
        obj = {
            "apiVersion": "deckhouse.io/v1alpha2",
            "kind": "NodeGroup",
            "metadata": {
                "name": "another",
            },
            "spec": {
                "disruptions": {
                    "approvalMode": "Automatic"
                },
                "cri": {
                    "docker": {
                        "manage": False,
                        "maxConcurrentDownloads": 4
                    }
                },
                "nodeType": "Hybrid"
            },
        }

        err, res_obj = test_dispatcher_for_unit_tests(TestUnitAlpha2ToV1Method.__snapshots).alpha2_to_v1(obj)

        self.assertIsNone(err)
        self.assertEqual(res_obj, {
            "apiVersion": "deckhouse.io/v1",
            "kind": "NodeGroup",
            "metadata": {
                "name": "another",
            },
            "spec": {
                "disruptions": {
                    "approvalMode": "Automatic"
                },
                "cri": {
                    "docker": {
                        "manage": False,
                        "maxConcurrentDownloads": 4
                    }
                },
                "nodeType": "CloudStatic"
            },

        })


class TestUnitV1ToAlpha2Method(unittest.TestCase):
    def test_not_alpha2_should_return_same_object(self):
        obj = {
            "apiVersion": "deckhouse.io/v1alpha1",
            "kind": "NodeGroup",
            "metadata": {
                "name": "worker-static",
            },
            "spec": {
                "disruptions": {
                    "approvalMode": "Automatic"
                },
            }
        }

        err, res_obj = test_dispatcher_for_unit_tests(None).v1_to_alpha2(obj)

        self.assertIsNone(err)
        self.assertEqual(obj, res_obj)


    def test_change_node_type_from_cloud_ephemeral_to_cloud(self):
        obj = {
            "apiVersion": "deckhouse.io/v1",
            "kind": "NodeGroup",
            "metadata": {
                "name": "worker-static",
            },
            "spec": {
                "disruptions": {
                    "approvalMode": "Automatic"
                },
                "cri": {
                    "docker": {
                        "manage": False,
                        "maxConcurrentDownloads": 4
                    }
                },
                "nodeType": "CloudEphemeral"
            },

        }

        err, res_obj = test_dispatcher_for_unit_tests(None).v1_to_alpha2(obj)

        self.assertIsNone(err)
        self.assertEqual(res_obj, {
            "apiVersion": "deckhouse.io/v1alpha2",
            "kind": "NodeGroup",
            "metadata": {
                "name": "worker-static",
            },
            "spec": {
                "disruptions": {
                    "approvalMode": "Automatic"
                },
                "cri": {
                    "docker": {
                        "manage": False,
                        "maxConcurrentDownloads": 4
                    }
                },
                "nodeType": "Cloud"
            },
        })


    def test_change_node_type_from_cloud_permanent_to_hybrid(self):
        obj = {
            "apiVersion": "deckhouse.io/v1",
            "kind": "NodeGroup",
            "metadata": {
                "name": "master",
            },
            "spec": {
                "disruptions": {
                    "approvalMode": "Automatic"
                },
                "cri": {
                    "docker": {
                        "manage": False,
                        "maxConcurrentDownloads": 4
                    }
                },
                "nodeType": "CloudPermanent"
            },

        }

        err, res_obj = test_dispatcher_for_unit_tests(None).v1_to_alpha2(obj)

        self.assertIsNone(err)
        self.assertEqual(res_obj, {
            "apiVersion": "deckhouse.io/v1alpha2",
            "kind": "NodeGroup",
            "metadata": {
                "name": "master",
            },
            "spec": {
                "disruptions": {
                    "approvalMode": "Automatic"
                },
                "cri": {
                    "docker": {
                        "manage": False,
                        "maxConcurrentDownloads": 4
                    }
                },
                "nodeType": "Hybrid"
            },
        })


    def test_change_node_type_from_cloud_static_to_hybrid(self):
        obj = {
            "apiVersion": "deckhouse.io/v1",
            "kind": "NodeGroup",
            "metadata": {
                "name": "another",
            },
            "spec": {
                "disruptions": {
                    "approvalMode": "Automatic"
                },
                "cri": {
                    "docker": {
                        "manage": False,
                        "maxConcurrentDownloads": 4
                    }
                },
                "nodeType": "CloudStatic"
            },

        }

        err, res_obj = test_dispatcher_for_unit_tests(None).v1_to_alpha2(obj)

        self.assertIsNone(err)
        self.assertEqual(res_obj, {
            "apiVersion": "deckhouse.io/v1alpha2",
            "kind": "NodeGroup",
            "metadata": {
                "name": "another",
            },
            "spec": {
                "disruptions": {
                    "approvalMode": "Automatic"
                },
                "cri": {
                    "docker": {
                        "manage": False,
                        "maxConcurrentDownloads": 4
                    }
                },
                "nodeType": "Hybrid"
            },
        })



class TestGroupValidationWebhook(unittest.TestCase):
    def test_should_convert_from_v1_to_alpha2(self):
        ctx = {
            "binding": "v1_to_alpha2",
            "fromVersion": "deckhouse.io/v1",
            "review": {
                "request": {
                    "uid": "54896fc3-89d1-466e-a90d-63a829f8233b",
                    "desiredAPIVersion": "deckhouse.io/v1alpha1",
                    "objects": [
                        {
                            "apiVersion": "deckhouse.io/v1",
                            "kind": "NodeGroup",
                            "metadata": {
                                "creationTimestamp": "2023-10-16T10:03:38Z",
                                "generation": 3,
                                "managedFields": [
                                    {
                                        "apiVersion": "deckhouse.io/v1",
                                        "fieldsType": "FieldsV1",
                                        "fieldsV1": {
                                            "f:spec": {
                                                ".": {},
                                                "f:disruptions": {
                                                    ".": {},
                                                    "f:approvalMode": {}
                                                },
                                                "f:nodeType": {}
                                            }
                                        },
                                        "manager": "kubectl-create",
                                        "operation": "Update",
                                        "time": "2023-10-16T10:03:38Z"
                                    },
                                    {
                                        "apiVersion": "deckhouse.io/v1",
                                        "fieldsType": "FieldsV1",
                                        "fieldsV1": {
                                            "f:spec": {
                                                "f:nodeTemplate": {
                                                    ".": {},
                                                    "f:labels": {
                                                        ".": {},
                                                        "f:aaaa": {}
                                                    }
                                                }
                                            }
                                        },
                                        "manager": "kubectl-edit",
                                        "operation": "Update",
                                        "time": "2023-10-16T10:06:01Z"
                                    },
                                    {
                                        "apiVersion": "deckhouse.io/v1",
                                        "fieldsType": "FieldsV1",
                                        "fieldsV1": {
                                            "f:spec": {
                                                "f:kubelet": {
                                                    "f:resourceReservation": {
                                                        "f:mode": {}
                                                    }
                                                }
                                            }
                                        },
                                        "manager": "deckhouse-controller",
                                        "operation": "Update",
                                        "time": "2024-01-30T07:08:03Z"
                                    },
                                    {
                                        "apiVersion": "deckhouse.io/v1",
                                        "fieldsType": "FieldsV1",
                                        "fieldsV1": {
                                            "f:status": {
                                                ".": {},
                                                "f:conditionSummary": {
                                                    ".": {},
                                                    "f:ready": {},
                                                    "f:statusMessage": {}
                                                },
                                                "f:conditions": {},
                                                "f:deckhouse": {
                                                    ".": {},
                                                    "f:observed": {
                                                        ".": {},
                                                        "f:checkSum": {},
                                                        "f:lastTimestamp": {}
                                                    },
                                                    "f:processed": {
                                                        ".": {},
                                                        "f:checkSum": {},
                                                        "f:lastTimestamp": {}
                                                    },
                                                    "f:synced": {}
                                                },
                                                "f:error": {},
                                                "f:kubernetesVersion": {},
                                                "f:nodes": {},
                                                "f:ready": {},
                                                "f:upToDate": {}
                                            }
                                        },
                                        "manager": "deckhouse-controller",
                                        "operation": "Update",
                                        "subresource": "status",
                                        "time": "2024-11-22T13:51:16Z"
                                    }
                                ],
                                "name": "worker-static",
                                "uid": "5a5e3820-5fdd-4fea-a8cf-032586ae0be5"
                            },
                            "spec": {
                                "disruptions": {
                                    "approvalMode": "Automatic"
                                },
                                "kubelet": {
                                    "containerLogMaxFiles": 4,
                                    "containerLogMaxSize": "50Mi",
                                    "resourceReservation": {
                                        "mode": "Off"
                                    }
                                },
                                "nodeTemplate": {
                                    "labels": {
                                        "aaaa": "bbbb"
                                    }
                                },
                                "nodeType": "CloudStatic"
                            },
                            "status": {
                                "conditionSummary": {
                                    "ready": "True",
                                    "statusMessage": ""
                                },
                                "conditions": [
                                    {
                                        "lastTransitionTime": "2023-10-16T10:03:39Z",
                                        "status": "True",
                                        "type": "Ready"
                                    },
                                    {
                                        "lastTransitionTime": "2023-10-26T12:24:02Z",
                                        "status": "False",
                                        "type": "Updating"
                                    },
                                    {
                                        "lastTransitionTime": "2023-10-16T10:03:39Z",
                                        "status": "False",
                                        "type": "WaitingForDisruptiveApproval"
                                    },
                                    {
                                        "lastTransitionTime": "2023-10-16T10:03:39Z",
                                        "status": "False",
                                        "type": "Error"
                                    }
                                ],
                                "deckhouse": {
                                    "observed": {
                                        "checkSum": "8cfb8c1cda6ce98f0b50e52302d5a871",
                                        "lastTimestamp": "2024-11-22T13:50:06Z"
                                    },
                                    "processed": {
                                        "checkSum": "8cfb8c1cda6ce98f0b50e52302d5a871",
                                        "lastTimestamp": "2024-11-22T13:51:15Z"
                                    },
                                    "synced": "True"
                                },
                                "error": "",
                                "kubernetesVersion": "1.27",
                                "nodes": 0,
                                "ready": 0,
                                "upToDate": 0
                            }
                        }
                    ]
                }
            },
            "toVersion": "deckhouse.io/v1alpha2",
            "type": "Conversion"
        }

        out = hook.testrun(main, [ctx])

        _assert_conversion(self, out, {}, None)


    def test_should_convert_from_alpha2_to_alpha1(self):
        ctx = {
            "binding": "alpha2_to_alpha1",
            "fromVersion": "deckhouse.io/v1alpha2",
            "review": {
                "request": {
                    "uid": "e8a2a24c-2232-480d-a188-2589de68c684",
                    "desiredAPIVersion": "deckhouse.io/v1alpha1",
                    "objects": [
                        {
                            "apiVersion": "deckhouse.io/v1alpha2",
                            "kind": "NodeGroup",
                            "metadata": {
                                "creationTimestamp": "2023-10-16T10:03:38Z",
                                "generation": 3,
                                "managedFields": [
                                    {
                                        "apiVersion": "deckhouse.io/v1",
                                        "fieldsType": "FieldsV1",
                                        "fieldsV1": {
                                            "f:spec": {
                                                ".": {},
                                                "f:disruptions": {
                                                    ".": {},
                                                    "f:approvalMode": {}
                                                },
                                                "f:nodeType": {}
                                            }
                                        },
                                        "manager": "kubectl-create",
                                        "operation": "Update",
                                        "time": "2023-10-16T10:03:38Z"
                                    },
                                    {
                                        "apiVersion": "deckhouse.io/v1",
                                        "fieldsType": "FieldsV1",
                                        "fieldsV1": {
                                            "f:spec": {
                                                "f:nodeTemplate": {
                                                    ".": {},
                                                    "f:labels": {
                                                        ".": {},
                                                        "f:aaaa": {}
                                                    }
                                                }
                                            }
                                        },
                                        "manager": "kubectl-edit",
                                        "operation": "Update",
                                        "time": "2023-10-16T10:06:01Z"
                                    },
                                    {
                                        "apiVersion": "deckhouse.io/v1",
                                        "fieldsType": "FieldsV1",
                                        "fieldsV1": {
                                            "f:spec": {
                                                "f:kubelet": {
                                                    "f:resourceReservation": {
                                                        "f:mode": {}
                                                    }
                                                }
                                            }
                                        },
                                        "manager": "deckhouse-controller",
                                        "operation": "Update",
                                        "time": "2024-01-30T07:08:03Z"
                                    },
                                    {
                                        "apiVersion": "deckhouse.io/v1",
                                        "fieldsType": "FieldsV1",
                                        "fieldsV1": {
                                            "f:status": {
                                                ".": {},
                                                "f:conditionSummary": {
                                                    ".": {},
                                                    "f:ready": {},
                                                    "f:statusMessage": {}
                                                },
                                                "f:conditions": {},
                                                "f:deckhouse": {
                                                    ".": {},
                                                    "f:observed": {
                                                        ".": {},
                                                        "f:checkSum": {},
                                                        "f:lastTimestamp": {}
                                                    },
                                                    "f:processed": {
                                                        ".": {},
                                                        "f:checkSum": {},
                                                        "f:lastTimestamp": {}
                                                    },
                                                    "f:synced": {}
                                                },
                                                "f:error": {},
                                                "f:kubernetesVersion": {},
                                                "f:nodes": {},
                                                "f:ready": {},
                                                "f:upToDate": {}
                                            }
                                        },
                                        "manager": "deckhouse-controller",
                                        "operation": "Update",
                                        "subresource": "status",
                                        "time": "2024-11-22T16:40:04Z"
                                    }
                                ],
                                "name": "worker-static",
                                "uid": "5a5e3820-5fdd-4fea-a8cf-032586ae0be5"
                            },
                            "spec": {
                                "disruptions": {
                                    "approvalMode": "Automatic"
                                },
                                "kubelet": {
                                    "containerLogMaxFiles": 4,
                                    "containerLogMaxSize": "50Mi",
                                    "resourceReservation": {
                                        "mode": "Off"
                                    }
                                },
                                "nodeTemplate": {
                                    "labels": {
                                        "aaaa": "bbbb"
                                    }
                                },
                                "nodeType": "Hybrid"
                            },
                            "status": {
                                "conditionSummary": {
                                    "ready": "True",
                                    "statusMessage": ""
                                },
                                "conditions": [
                                    {
                                        "lastTransitionTime": "2023-10-16T10:03:39Z",
                                        "status": "True",
                                        "type": "Ready"
                                    },
                                    {
                                        "lastTransitionTime": "2023-10-26T12:24:02Z",
                                        "status": "False",
                                        "type": "Updating"
                                    },
                                    {
                                        "lastTransitionTime": "2023-10-16T10:03:39Z",
                                        "status": "False",
                                        "type": "WaitingForDisruptiveApproval"
                                    },
                                    {
                                        "lastTransitionTime": "2023-10-16T10:03:39Z",
                                        "status": "False",
                                        "type": "Error"
                                    }
                                ],
                                "deckhouse": {
                                    "observed": {
                                        "checkSum": "8cfb8c1cda6ce98f0b50e52302d5a871",
                                        "lastTimestamp": "2024-11-22T16:40:04Z"
                                    },
                                    "processed": {
                                        "checkSum": "8cfb8c1cda6ce98f0b50e52302d5a871",
                                        "lastTimestamp": "2024-11-22T16:22:54Z"
                                    },
                                    "synced": "True"
                                },
                                "error": "",
                                "kubernetesVersion": "1.27",
                                "nodes": 0,
                                "ready": 0,
                                "upToDate": 0
                            }
                        }
                    ]
                }
            },
            "toVersion": "deckhouse.io/v1alpha1",
            "type": "Conversion"
        }

        out = hook.testrun(main, [ctx])

        _assert_conversion(self, out, {}, None)


    def test_should_convert_from_alpha2_to_v1(self):
        ctx = {
            "binding": "alpha2_to_v1",
            "fromVersion": "deckhouse.io/v1alpha2",
            "review": {
                "request": {
                    "uid": "e8a2a24c-2232-480d-a188-2589de68c684",
                    "desiredAPIVersion": "deckhouse.io/v1alpha1",
                    "objects": [
                        {
                            "apiVersion": "deckhouse.io/v1alpha2",
                            "kind": "NodeGroup",
                            "metadata": {
                                "creationTimestamp": "2023-10-16T10:03:38Z",
                                "generation": 3,
                                "managedFields": [
                                    {
                                        "apiVersion": "deckhouse.io/v1",
                                        "fieldsType": "FieldsV1",
                                        "fieldsV1": {
                                            "f:spec": {
                                                ".": {},
                                                "f:disruptions": {
                                                    ".": {},
                                                    "f:approvalMode": {}
                                                },
                                                "f:nodeType": {}
                                            }
                                        },
                                        "manager": "kubectl-create",
                                        "operation": "Update",
                                        "time": "2023-10-16T10:03:38Z"
                                    },
                                    {
                                        "apiVersion": "deckhouse.io/v1",
                                        "fieldsType": "FieldsV1",
                                        "fieldsV1": {
                                            "f:spec": {
                                                "f:nodeTemplate": {
                                                    ".": {},
                                                    "f:labels": {
                                                        ".": {},
                                                        "f:aaaa": {}
                                                    }
                                                }
                                            }
                                        },
                                        "manager": "kubectl-edit",
                                        "operation": "Update",
                                        "time": "2023-10-16T10:06:01Z"
                                    },
                                    {
                                        "apiVersion": "deckhouse.io/v1",
                                        "fieldsType": "FieldsV1",
                                        "fieldsV1": {
                                            "f:spec": {
                                                "f:kubelet": {
                                                    "f:resourceReservation": {
                                                        "f:mode": {}
                                                    }
                                                }
                                            }
                                        },
                                        "manager": "deckhouse-controller",
                                        "operation": "Update",
                                        "time": "2024-01-30T07:08:03Z"
                                    },
                                    {
                                        "apiVersion": "deckhouse.io/v1",
                                        "fieldsType": "FieldsV1",
                                        "fieldsV1": {
                                            "f:status": {
                                                ".": {},
                                                "f:conditionSummary": {
                                                    ".": {},
                                                    "f:ready": {},
                                                    "f:statusMessage": {}
                                                },
                                                "f:conditions": {},
                                                "f:deckhouse": {
                                                    ".": {},
                                                    "f:observed": {
                                                        ".": {},
                                                        "f:checkSum": {},
                                                        "f:lastTimestamp": {}
                                                    },
                                                    "f:processed": {
                                                        ".": {},
                                                        "f:checkSum": {},
                                                        "f:lastTimestamp": {}
                                                    },
                                                    "f:synced": {}
                                                },
                                                "f:error": {},
                                                "f:kubernetesVersion": {},
                                                "f:nodes": {},
                                                "f:ready": {},
                                                "f:upToDate": {}
                                            }
                                        },
                                        "manager": "deckhouse-controller",
                                        "operation": "Update",
                                        "subresource": "status",
                                        "time": "2024-11-22T16:40:04Z"
                                    }
                                ],
                                "name": "worker-static",
                                "uid": "5a5e3820-5fdd-4fea-a8cf-032586ae0be5"
                            },
                            "spec": {
                                "disruptions": {
                                    "approvalMode": "Automatic"
                                },
                                "kubelet": {
                                    "containerLogMaxFiles": 4,
                                    "containerLogMaxSize": "50Mi",
                                    "resourceReservation": {
                                        "mode": "Off"
                                    }
                                },
                                "nodeTemplate": {
                                    "labels": {
                                        "aaaa": "bbbb"
                                    }
                                },
                                "nodeType": "Hybrid"
                            },
                            "status": {
                                "conditionSummary": {
                                    "ready": "True",
                                    "statusMessage": ""
                                },
                                "conditions": [
                                    {
                                        "lastTransitionTime": "2023-10-16T10:03:39Z",
                                        "status": "True",
                                        "type": "Ready"
                                    },
                                    {
                                        "lastTransitionTime": "2023-10-26T12:24:02Z",
                                        "status": "False",
                                        "type": "Updating"
                                    },
                                    {
                                        "lastTransitionTime": "2023-10-16T10:03:39Z",
                                        "status": "False",
                                        "type": "WaitingForDisruptiveApproval"
                                    },
                                    {
                                        "lastTransitionTime": "2023-10-16T10:03:39Z",
                                        "status": "False",
                                        "type": "Error"
                                    }
                                ],
                                "deckhouse": {
                                    "observed": {
                                        "checkSum": "8cfb8c1cda6ce98f0b50e52302d5a871",
                                        "lastTimestamp": "2024-11-22T16:40:04Z"
                                    },
                                    "processed": {
                                        "checkSum": "8cfb8c1cda6ce98f0b50e52302d5a871",
                                        "lastTimestamp": "2024-11-22T16:22:54Z"
                                    },
                                    "synced": "True"
                                },
                                "error": "",
                                "kubernetesVersion": "1.27",
                                "nodes": 0,
                                "ready": 0,
                                "upToDate": 0
                            }
                        }
                    ]
                }
            },
            "toVersion": "deckhouse.io/v1",
            "type": "Conversion"
        }

        out = hook.testrun(main, [ctx])

        _assert_conversion(self, out, {}, None)


    def test_should_convert_from_alpha1_to_alpha2(self):
        ctx = {
            "binding": "alpha1_to_alpha2",
            "fromVersion": "deckhouse.io/v1alpha1",
            "review": {
                "request": {
                    "uid": "e8a2a24c-2232-480d-a188-2589de68c684",
                    "desiredAPIVersion": "deckhouse.io/v1alpha1",
                    "objects": [
                        {
                            "apiVersion": "deckhouse.io/v1alpha1",
                            "kind": "NodeGroup",
                            "metadata": {
                                "creationTimestamp": "2023-10-16T10:03:38Z",
                                "generation": 3,
                                "managedFields": [
                                    {
                                        "apiVersion": "deckhouse.io/v1",
                                        "fieldsType": "FieldsV1",
                                        "fieldsV1": {
                                            "f:spec": {
                                                ".": {},
                                                "f:disruptions": {
                                                    ".": {},
                                                    "f:approvalMode": {}
                                                },
                                                "f:nodeType": {}
                                            }
                                        },
                                        "manager": "kubectl-create",
                                        "operation": "Update",
                                        "time": "2023-10-16T10:03:38Z"
                                    },
                                    {
                                        "apiVersion": "deckhouse.io/v1",
                                        "fieldsType": "FieldsV1",
                                        "fieldsV1": {
                                            "f:spec": {
                                                "f:nodeTemplate": {
                                                    ".": {},
                                                    "f:labels": {
                                                        ".": {},
                                                        "f:aaaa": {}
                                                    }
                                                }
                                            }
                                        },
                                        "manager": "kubectl-edit",
                                        "operation": "Update",
                                        "time": "2023-10-16T10:06:01Z"
                                    },
                                    {
                                        "apiVersion": "deckhouse.io/v1",
                                        "fieldsType": "FieldsV1",
                                        "fieldsV1": {
                                            "f:spec": {
                                                "f:kubelet": {
                                                    "f:resourceReservation": {
                                                        "f:mode": {}
                                                    }
                                                }
                                            }
                                        },
                                        "manager": "deckhouse-controller",
                                        "operation": "Update",
                                        "time": "2024-01-30T07:08:03Z"
                                    },
                                    {
                                        "apiVersion": "deckhouse.io/v1",
                                        "fieldsType": "FieldsV1",
                                        "fieldsV1": {
                                            "f:status": {
                                                ".": {},
                                                "f:conditionSummary": {
                                                    ".": {},
                                                    "f:ready": {},
                                                    "f:statusMessage": {}
                                                },
                                                "f:conditions": {},
                                                "f:deckhouse": {
                                                    ".": {},
                                                    "f:observed": {
                                                        ".": {},
                                                        "f:checkSum": {},
                                                        "f:lastTimestamp": {}
                                                    },
                                                    "f:processed": {
                                                        ".": {},
                                                        "f:checkSum": {},
                                                        "f:lastTimestamp": {}
                                                    },
                                                    "f:synced": {}
                                                },
                                                "f:error": {},
                                                "f:kubernetesVersion": {},
                                                "f:nodes": {},
                                                "f:ready": {},
                                                "f:upToDate": {}
                                            }
                                        },
                                        "manager": "deckhouse-controller",
                                        "operation": "Update",
                                        "subresource": "status",
                                        "time": "2024-11-22T16:40:04Z"
                                    }
                                ],
                                "name": "worker-static",
                                "uid": "5a5e3820-5fdd-4fea-a8cf-032586ae0be5"
                            },
                            "spec": {
                                "disruptions": {
                                    "approvalMode": "Manual"
                                },
                                "kubelet": {
                                    "containerLogMaxFiles": 4,
                                    "containerLogMaxSize": "50Mi"
                                },
                                "nodeTemplate": {
                                    "labels": {
                                        "node-role.kubernetes.io/control-plane": "",
                                        "node-role.kubernetes.io/master": ""
                                    },
                                    "taints": [
                                        {
                                            "effect": "NoSchedule",
                                            "key": "node-role.kubernetes.io/control-plane"
                                        }
                                    ]
                                },
                                "nodeType": "Hybrid"
                            },
                            "status": {
                                "conditionSummary": {
                                    "ready": "True",
                                    "statusMessage": ""
                                },
                                "error": "",
                                "kubernetesVersion": "1.27",
                                "nodes": 1,
                                "ready": 1,
                                "upToDate": 1
                            }
                        }
                    ]
                }
            },
            "toVersion": "deckhouse.io/v1alpha2",
            "type": "Conversion"
        }

        out = hook.testrun(main, [ctx])

        _assert_conversion(self, out, {}, None)


