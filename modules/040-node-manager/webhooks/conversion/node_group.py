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


from dotmap import DotMap
from deckhouse import hook

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
        self._binding_context = DotMap(ctx.binding_context)
        self._snapshots = DotMap(ctx.snapshots)
        self.__ctx = ctx

    def run(self):
        binding_name = self._binding_context.binding
        try:
            action = getattr(self, binding_name)
        except AttributeError:
            self.__ctx.output.conversions.error("Internal error. Handler for binding {} not found".format(binding_name))
            return

        try:
            error_msg, res_obj = action()
            if error_msg is None:
                if not isinstance(res_obj, dict):
                    raise TypeError("Result of conversion should be dict")
                self.__ctx.output.conversions.collect(res_obj)
            else:
                self.__ctx.output.conversions.error(error_msg)
        except Exception as e:
            self.__ctx.output.conversions.error("Internal error. Catch exception: {}".format(str(e)))
            return

    def _get_objects(self) -> typing.List[DotMap]:
        return self._binding_context.review.request.objects

class NodeGroupConversionDispatcher(ConversionDispatcher):
    def __init__(self, ctx: hook.Context):
        super().__init__(ctx)

    def v1_to_alpha2(self) -> typing.Tuple[str | None, dict]:
        return None, {}

    def alpha2_to_alpha1(self) -> typing.Tuple[str | None, dict]:
        return None, {}

    def alpha2_to_v1(self) -> typing.Tuple[str | None, dict]:
        return None, {}

    def alpha1_to_alpha2(self) -> typing.Tuple[str | None, dict]:
        return None, {}

def main(ctx: hook.Context):
    NodeGroupConversionDispatcher(ctx).run()

if __name__ == "__main__":
    hook.run(main, config=config)

