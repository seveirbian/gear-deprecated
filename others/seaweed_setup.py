import sys
# package need to be installed, pip install docker
import docker
import time
# package need to be installed, pip install pyyaml
import yaml
import os

class Runner:

    def __init__(self, image, ports, command, remove, waitline):
        self.image = image
        self.ports = ports
        self.command = command
        self.remove = remove
        self.waitline = waitline


    def start(self):
        image_pulled = client.images.pull(self.image)

        container = client.containers.create(image=self.image, ports=self.ports,
                                        command=self.command, detach=True, remove=self.remove)

        container.start()

        while True:
            if container.logs().find(self.waitline) >= 0:
                break


if __name__ == "__main__":
    runner = Runner(image="chrislusf/seaweedfs:latest",
        ports={"9333/tcp":"9333"}, command="server", remove=true, waitline="added volume server")