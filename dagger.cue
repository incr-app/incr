package main

import (
	"dagger.io/dagger"
	"universe.dagger.io/docker"
	"universe.dagger.io/go"
)

dagger.#Plan & {
    _cueVersion: "v0.4.3"
    _goVersion: "1.17"

    client: filesystem: ".": read: contents: dagger.#FS
    client: filesystem: "./build": write: contents: actions.build.cmd.output

    _image: actions.cueimage.build.output

    actions: {
        cueimage: {
            build: docker.#Build & {
                _alpineVersion: "3.15.0@sha256:21a3deaa0d32a8057914f36584b5288d2e5ecc984380bc0118285c70fa8c9300"
                steps: [
                    docker.#Pull & {source: "index.docker.io/alpine:\(_alpineVersion)"},
                    docker.#Run & {
                        _script: #"""
                        # Install cue
                        wget -O cue.tar.gz https://github.com/cue-lang/cue/releases/download/${CUE_VERSION}/cue_${CUE_VERSION}_linux_amd64.tar.gz
                        tar xzf cue.tar.gz -C /tmp && \
                        mv /tmp/cue /usr/local/bin/cue && \
                        rm -r /tmp/*
                        """#
                        command: {
                            name: "sh"
                            flags: {
                                "-c": _script
                            }
                        }
                        env: CUE_VERSION: _cueVersion
                    }
                ]
            }
            test: docker.#Run & {
                input: _image
                command: name: "cue"
                command: args: ["version"]

                success: true
            }
        }

        test: {
            cue: {
                vet: {
                    docker.#Run & {
                        input: _image
                        _script: #"""
                        cd /workspace
                        find schemas -name '*.cue' | xargs -n1 dirname | xargs -I{} -n1 cue vet ./{}
                        """#
                        command: name: "sh"
                        command: flags: "-c": _script
                        mounts: "root": contents: client.filesystem.".".read.contents
                        mounts: "root": dest: "/workspace"
                    }
                }
            }
        }

        build: {
            cmd: {
                go.#Build & {
                    _image: go.#Image & {
                        version: _goVersion
                    }
                    container: input: _image.output
                    source: client.filesystem.".".read.contents
                    os: "linux"
                    arch: "amd64"
                    package: "./cmd/incr"
                    env: "CGO_ENABLED": "0"
                }
            }
        }
    }
}
