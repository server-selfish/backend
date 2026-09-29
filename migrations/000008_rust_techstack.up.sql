-- rust (gnu builder + distroless cc runtime; musl/static flavor intentionally omitted)
INSERT INTO deployment_techstack (name, version, docker_base_image, docker_runtime_image) VALUES ('Rust','1.96.1','rust:1.96.1-slim-trixie','gcr.io/distroless/cc-debian13');
INSERT INTO deployment_techstack (name, version, docker_base_image, docker_runtime_image) VALUES ('Rust','1.97.1','rust:1.97.1-slim-trixie','gcr.io/distroless/cc-debian13');
INSERT INTO deployment_techstack (name, version, docker_base_image, docker_runtime_image) VALUES ('Rust','1.98.1','rust:1.98.1-slim-trixie','gcr.io/distroless/cc-debian13');
