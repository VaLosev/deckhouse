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


def test_dispatcher_for_unit_tests() -> NodeGroupConversionDispatcher:
    output = hook.Output(
        hook.MetricsCollector(),
        hook.KubeOperationCollector(),
        hook.ValuesPatchesCollector({}),
        hook.ConversionsCollector(),
        hook.ValidationsCollector(),
    )
    return NodeGroupConversionDispatcher(hook.Context({}, {}, {}, output))



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

        err, res_obj = test_dispatcher_for_unit_tests().alpha1_to_alpha2(DotMap(obj))

        self.assertIsNone(err)
        self.assertEqual(obj, res_obj.toDict())


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

        err, res_obj = test_dispatcher_for_unit_tests().alpha1_to_alpha2(DotMap(obj))

        self.assertIsNone(err)
        self.assertEqual(res_obj.toDict(), {
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

        err, res_obj = test_dispatcher_for_unit_tests().alpha1_to_alpha2(DotMap(obj))

        self.assertIsNone(err)
        self.assertEqual(res_obj.toDict(), {
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


