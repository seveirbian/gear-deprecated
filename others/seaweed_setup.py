import sys
# package need to be installed, pip install docker
import docker
import time
# package need to be installed, pip install pyyaml
import yaml
import os

class Runner:

    def __init__(self, image, ports, command, waitline):
        self.image = image
        self.ports = ports
        self.command = command
        self.waitline = waitline


    def start(self):
        print "creaing client..."
        client = docker.from_env()

        print "pulling image..."
        image_pulled = client.images.pull(self.image)

        print "creating container..."
        container = client.containers.create(image=self.image, ports=self.ports,
                                        command=self.command, detach=True)

        print "starting container..."
        container.start()

        while True:
            if container.logs().find(self.waitline) >= 0:
                break
        finish = input("Please enter any key to stop this container...")

        print "removing container..."
        container.remove(force=True)


if __name__ == "__main__":
    runner = Runner(image="chrislusf/seaweedfs:latest",
        ports={"9333/tcp":"9333"}, command="server", waitline="added volume server")

    runner.start()