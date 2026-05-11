-- go
INSERT INTO deployment_techstack (name, version, docker_base_image, docker_runtime_image) VALUES ('Go','1.20.14','golang:1.20.14-alpine3.19','gcr.io/distroless/static-debian13');
INSERT INTO deployment_techstack (name, version, docker_base_image, docker_runtime_image) VALUES ('Go','1.21.13','golang:1.21.13-alpine3.20','gcr.io/distroless/static-debian13');
INSERT INTO deployment_techstack (name, version, docker_base_image, docker_runtime_image) VALUES ('Go','1.22.12','golang:1.22.12-alpine3.21','gcr.io/distroless/static-debian13');
INSERT INTO deployment_techstack (name, version, docker_base_image, docker_runtime_image) VALUES ('Go','1.23.12','golang:1.23.12-alpine3.22','gcr.io/distroless/static-debian13');
INSERT INTO deployment_techstack (name, version, docker_base_image, docker_runtime_image) VALUES ('Go','1.25.9','golang:1.25.9-alpine3.22','gcr.io/distroless/static-debian13');
INSERT INTO deployment_techstack (name, version, docker_base_image, docker_runtime_image) VALUES ('Go','1.26.2','golang:1.26.2-alpine3.22','gcr.io/distroless/static-debian13');

-- nodejs
INSERT INTO deployment_techstack (name, version, docker_base_image, docker_runtime_image) VALUES ('Node.js','20.20.2','node:20.20.2-slim','gcr.io/distroless/nodejs20-debian13');
INSERT INTO deployment_techstack (name, version, docker_base_image, docker_runtime_image) VALUES ('Node.js','21.7.3','node:21.7.3-alpine3.20','node:21.7.3-alpine3.20');
INSERT INTO deployment_techstack (name, version, docker_base_image, docker_runtime_image) VALUES ('Node.js','22.21.0','node:22.21.0-alpine3.21','gcr.io/distroless/nodejs22-debian13');
INSERT INTO deployment_techstack (name, version, docker_base_image, docker_runtime_image) VALUES ('Node.js','23.11.1','node:23.11.1-alpine3.22','node:23.11.1-alpine3.22');
INSERT INTO deployment_techstack (name, version, docker_base_image, docker_runtime_image) VALUES ('Node.js','24.15.0','node:24.15.0-alpine3.22','gcr.io/distroless/nodejs24-debian13');
INSERT INTO deployment_techstack (name, version, docker_base_image, docker_runtime_image) VALUES ('Node.js','25.9.0','node:25.9.0-trixie-slim','node:25.9.0-trixie-slim');
INSERT INTO deployment_techstack (name, version, docker_base_image, docker_runtime_image) VALUES ('Node.js','26.1.0','node:26.1.0-trixie-slim','node:26.1.0-trixie-slim');

-- python
INSERT INTO deployment_techstack (name, version, docker_base_image, docker_runtime_image) VALUES ('Python','3.9.25','python:3.9.25-slim','gcr.io/distroless/python3-debian13');
INSERT INTO deployment_techstack (name, version, docker_base_image, docker_runtime_image) VALUES ('Python','3.10.20','python:3.10.20-slim','gcr.io/distroless/python3-debian13');
INSERT INTO deployment_techstack (name, version, docker_base_image, docker_runtime_image) VALUES ('Python','3.11.15','python:3.11.15-slim','gcr.io/distroless/python3-debian13');
INSERT INTO deployment_techstack (name, version, docker_base_image, docker_runtime_image) VALUES ('Python','3.12.13','python:3.12.13-slim','gcr.io/distroless/python3-debian13');
INSERT INTO deployment_techstack (name, version, docker_base_image, docker_runtime_image) VALUES ('Python','3.13.13','python:3.13.13-slim','gcr.io/distroless/python3-debian13');
INSERT INTO deployment_techstack (name, version, docker_base_image, docker_runtime_image) VALUES ('Python','3.14.4','python:3.14.4-slim','gcr.io/distroless/python3-debian13');
