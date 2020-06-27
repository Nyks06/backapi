FROM heroku/heroku:18-build as build

COPY . /app
WORKDIR /app

# Setup buildpack
RUN mkdir -p /tmp/buildpack/heroku/go /tmp/build_cache /tmp/env
RUN curl https://codon-buildpacks.s3.amazonaws.com/buildpacks/heroku/go.tgz | tar xz -C /tmp/buildpack/heroku/go

#Execute Buildpack
RUN STACK=heroku-18 GLIDE_SKIP_INSTALL=true /tmp/buildpack/heroku/go/bin/compile /app /tmp/build_cache /tmp/env
RUN mkdir -p /root/go/src/github.com/nyks06/backapi
COPY . /root/go/src/github.com/nyks06/backapi/.
RUN /tmp/build_cache/go1.12.12/go/bin/go build -o server cmd/live/main.go

# Prepare final, minimal image
FROM heroku/heroku:18

COPY --from=build /app /app
ENV HOME /app
WORKDIR /app
RUN useradd -m heroku

CMD ["/bin/sh", "start.sh"]