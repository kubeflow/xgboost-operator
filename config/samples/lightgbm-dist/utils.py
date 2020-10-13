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

import re
import logging
import tempfile
from typing import List, Union

logger = logging.getLogger(__name__)


def generate_machine_list_file(
    master_addr: str, master_port: str, worker_addrs: str, worker_port: str
) -> str:
    logger.info("starting to extract system env")

    filename = tempfile.NamedTemporaryFile(delete=False).name
    with open(filename, "w") as file:
        print(f"{master_addr} {master_port}", file=file)
        for addr in worker_addrs.split(","):
            print(f"{addr} {worker_port}", file=file)

    return filename


def generate_train_conf_file(
    machine_list_file: str,
    world_size: int,
    output_model: str,
    local_port: Union[int, str],
    extra_args: List[str],
) -> str:

    filename = tempfile.NamedTemporaryFile(delete=False).name

    with open(filename, "w") as file:
        print("task = train", file=file)
        print(f"output_model = {output_model}", file=file)
        print(f"num_machines = {world_size}", file=file)
        print(f"local_listen_port = {local_port}", file=file)
        print(f"machine_list_file = {machine_list_file}", file=file)
        for arg in extra_args:
            m = re.match(r"--(.+)=([^\s]+)", arg)
            if m is not None:
                k, v = m.groups()
                print(f"{k} = {v}", file=file)

    return filename
