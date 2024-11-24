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

import typing
import base64

from dotmap import DotMap
from deckhouse import hook
import yaml


config = """
configVersion: v1
kubernetes:
  - name: cluster_config
    apiVersion: v1
    kind: Secret
    queue: "cache-cluster-config"
    executeHookOnEvent: []
    executeHookOnSynchronization: false
    keepFullObjectsInMemory: false
    namespace:
      nameSelector:
        matchNames: ["kube-system"]
    nameSelector:
      matchNames:
        - d8-provider-cluster-configuration
    jqFilter: '.data."cloud-provider-cluster-configuration.yaml" // ""'
kubernetesCustomResourceConversion:
  - name: alpha1_to_alpha2
    crdName: nodegroups.deckhouse.io
    conversions:
    - fromVersion: deckhouse.io/v1alpha1
      toVersion: deckhouse.io/v1alpha2
  - name: alpha2_to_alpha1
    crdName: nodegroups.deckhouse.io
    conversions:
    - fromVersion: deckhouse.io/v1alpha2
      toVersion: deckhouse.io/v1alpha1
  - name: alpha2_to_v1
    includeSnapshotsFrom: ["cluster_config"]
    crdName: nodegroups.deckhouse.io
    conversions:
    - fromVersion: deckhouse.io/v1alpha2
      toVersion: deckhouse.io/v1
  - name: v1_to_alpha2
    crdName: nodegroups.deckhouse.io
    conversions:
    - fromVersion: deckhouse.io/v1
      toVersion: deckhouse.io/v1alpha2
"""

class ConversionDispatcher:
    def __init__(self, ctx: hook.Context):
        self._binding_context = ctx.binding_context
        self._snapshots = ctx.snapshots
        self.__ctx = ctx


    def run(self):
        binding_name = self._binding_context["binding"]

        try:
            action = getattr(self, binding_name)
        except AttributeError:
            self.__ctx.output.conversions.error("Internal error. Handler for binding {} not found".format(binding_name))
            return

        try:
            errors = []
            for obj in self._binding_context["review"]["request"]["objects"]:
                error_msg, res_obj = action(obj)
                if error_msg is not None:
                    errors.append(error_msg)
                    continue
                self.__ctx.output.conversions.collect(res_obj)
            if errors:
                err_msg = ";".join(errors)
                self.__ctx.output.conversions.error(err_msg)
        except Exception as e:
            self.__ctx.output.conversions.error("Internal error: {}".format(str(e)))
            return


class NodeGroupConversionDispatcher(ConversionDispatcher):
    def __init__(self, ctx: hook.Context):
        super().__init__(ctx)


    def alpha1_to_alpha2(self, o: dict) -> typing.Tuple[str | None, dict]:
        obj = DotMap(o)

        if obj.apiVersion != "deckhouse.io/v1alpha1":
            return None, o

        obj.apiVersion = "deckhouse.io/v1alpha2"
        if "docker" in obj.spec:
            if "cri" not in obj.spec:
                obj.spec.cri = DotMap({})
            obj.spec.cri.docker = obj.spec.docker
            del obj.spec.docker

        if "kubernetesVersion" in obj.spec:
            del obj.spec.kubernetesVersion

        if "static" in obj:
            del obj.static

        return None, obj.toDict()


    def alpha2_to_alpha1(self, o: dict) -> typing.Tuple[str | None, dict]:
        obj = DotMap(o)

        if obj.apiVersion != "deckhouse.io/v1alpha2":
            return None, o

        obj.apiVersion = "deckhouse.io/v1alpha1"

        if "cri" in obj.spec:
            if "docker" in obj.spec.cri:
                obj.spec.docker = obj.spec.cri.docker
                del obj.spec.cri.docker

        return None, obj.toDict()


    def alpha2_to_v1(self, o: dict) -> typing.Tuple[str | None, dict]:
        obj = DotMap(o)

        if obj.apiVersion != "deckhouse.io/v1alpha2":
            return None, o

        obj.apiVersion = "deckhouse.io/v1"

        try:
            provider_config_yaml = base64.standard_b64decode(self._snapshots["cluster_config"][0]["filterResult"])
        except Exception as e:
            return f"Cannot decode provider cluster configuration: {e}", {}

        try:
            provider_config = yaml.safe_load(provider_config_yaml)
        except Exception as e:
            return f"Cannot parse provider cluster configuration: {e}", {}

        ng_name = obj.metadata.name
        ng_type = obj.spec.nodeType
        if ng_type == "Cloud":
            ng_type = "CloudEphemeral"
        elif ng_type == "Hybrid":
            found_in_permanent = False
            if ng_name == "master":
                found_in_permanent = True
            else:
                if "nodeGroups" in provider_config:
                    for ng in provider_config["nodeGroups"]:
                        if ng["name"] == ng_name:
                            found_in_permanent = True
                            break
            ng_type = "CloudPermanent" if found_in_permanent else "CloudStatic"

        obj.spec.nodeType = ng_type

        return None, obj.toDict()


    def v1_to_alpha2(self, o: dict) -> typing.Tuple[str | None, dict]:
        obj = DotMap(o)

        if obj.apiVersion != "deckhouse.io/v1":
            return None, o

        obj.apiVersion = "deckhouse.io/v1alpha2"

        ng_type = obj.spec.nodeType
        if ng_type == "CloudEphemeral":
            ng_type = "Cloud"
        elif ng_type == "CloudPermanent" or ng_type == "CloudStatic":
            ng_type = "Hybrid"

        obj.spec.nodeType = ng_type

        return None, obj.toDict()


def main(ctx: hook.Context):
    NodeGroupConversionDispatcher(ctx).run()


if __name__ == "__main__":
    hook.run(main, config=config)

