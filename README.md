# Stash

[![Build](https://github.com/stashapp/stash/actions/workflows/build.yml/badge.svg?branch=develop&event=push)](https://github.com/stashapp/stash/actions/workflows/build.yml)
[![Docker pulls](https://img.shields.io/docker/pulls/stashapp/stash?logo=docker)](https://hub.docker.com/r/stashapp/stash 'DockerHub')
[![GitHub Sponsors](https://img.shields.io/github/sponsors/stashapp?logo=github)](https://github.com/sponsors/stashapp)
[![Open Collective backers](https://img.shields.io/opencollective/backers/stashapp?logo=opencollective)](https://opencollective.com/stashapp)
[![Go Report Card](https://goreportcard.com/badge/github.com/stashapp/stash)](https://goreportcard.com/report/github.com/stashapp/stash)
[![Discord](https://img.shields.io/discord/559159668438728723.svg?logo=discord)](https://discord.gg/2TsNFKt)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/stashapp/stash?logo=github)](https://github.com/stashapp/stash/releases/latest)
[![Codeberg Translate](https://img.shields.io/weblate/progress/stash?server=https%3A%2F%2Ftranslate.codeberg.org&logo=weblate)](https://translate.codeberg.org/engage/stash/)
[![GitHub issues by-label](https://img.shields.io/github/issues-raw/stashapp/stash/bounty?logo=github)](https://github.com/stashapp/stash/labels/bounty)

<h3>Stash is a self-hosted webapp written in Go which organizes and serves your diverse content collection, catering to both your SFW and NSFW needs.</h3>

![Screenshot of Stash web application interface](docs/readme_assets/demo_image.png)

- Stash gathers information about videos in your collection from the internet, and is extensible through the use of community-built plugins for a large number of content producers and sites.
- Stash supports a wide variety of both video and image formats.
- You can tag videos and find them later.
- Stash provides statistics about performers, tags, studios and more.

You can [watch a SFW demo video](https://vimeo.com/545323354) to see it in action.

For further information see [Support & Resources](#support--resources) section.

> [!note]
> **This fork adds OpenID Connect (OIDC) single sign-on.** Stash can delegate login to an external identity provider such as Authelia, Authentik, Keycloak, Google, Okta or Auth0. See [Single Sign-On with OpenID Connect](#single-sign-on-with-openid-connect-oidc) below.

## Installing Stash

> [!tip]
Step-by-step instructions are available at [docs.stashapp.cc/installation](https://docs.stashapp.cc/installation/).

> [!important]
> **Windows Users**
>
> As of version 0.27.0, Stash no longer supports _Windows 7, 8, Server 2008 and Server 2012._  
> At least Windows 10 or Server 2016 is required.
>
> **macOS Users**
>
> As of version 0.29.0, Stash requires _macOS 11 Big Sur_ or later.  
> As of version 0.32.0, Stash requires _macOS 12 Monterey_ or later.  
> Older versions can still be run through Docker.

<img src="docs/readme_assets/windows_logo.svg" width="100%" height="75"> Windows | <img src="docs/readme_assets/mac_logo.svg" width="100%" height="75"> macOS | <img src="docs/readme_assets/linux_logo.svg" width="100%" height="75"> Linux | <img src="docs/readme_assets/docker_logo.svg" width="100%" height="75"> Docker
:---:|:---:|:---:|:---:
[Latest Release](https://github.com/stashapp/stash/releases/latest/download/stash-win.exe) <br /> <sup><sub>[Development Preview](https://github.com/stashapp/stash/releases/download/latest_develop/stash-win.exe)</sub></sup> | [Latest Release](https://github.com/stashapp/stash/releases/latest/download/Stash.app.zip) <br /> <sup><sub>[Development Preview](https://github.com/stashapp/stash/releases/download/latest_develop/Stash.app.zip)</sub></sup> | [Latest Release (amd64)](https://github.com/stashapp/stash/releases/latest/download/stash-linux) <br /> <sup><sub>[Development Preview (amd64)](https://github.com/stashapp/stash/releases/download/latest_develop/stash-linux)</sub></sup> <br /> [More Architectures...](https://github.com/stashapp/stash/releases/latest) | [Instructions](docker/production/README.md) <br /> <sup><sub>[Sample docker-compose.yml](docker/production/docker-compose.yml)</sub></sup>

Download links for other platforms and architectures are available on the [Releases](https://github.com/stashapp/stash/releases) page.

### First Run

#### Windows/macOS Users: Security Prompt

On Windows or macOS, running the app might present a security prompt since the application binary isn't yet signed. 

- On Windows, bypass this by clicking "more info" and then the "run anyway" button.
- On macOS, Control+Click the app, click "Open", and then "Open" again.

#### ffmpeg

Stash requires FFmpeg. If you don't have it installed, Stash will prompt you to download a copy during setup. It is recommended that Linux users install `ffmpeg` from their distro's package manager.

## Usage

### Quickstart Guide

Stash is a web-based application. Once the application is running, the interface is available (by default) from `http://localhost:9999`.

On first run, Stash will prompt you for some configuration options and media directories to index, called "Scanning" in Stash. After scanning, your media will be available for browsing, curating, editing, and tagging.

Stash can pull metadata (performers, tags, descriptions, studios, and more) directly from many sites through the use of [scrapers](https://github.com/stashapp/stash/blob/develop/ui/v2.5/src/docs/en/Manual/Scraping.md), which integrate directly into Stash. Identifying an entire collection will typically require a mix of multiple sources:
- The stashapp team maintains [StashDB](https://stashdb.org/), a crowd-sourced repository of scene, studio, and performer information. Connecting it to Stash will allow you to automatically identify much of a typical media collection. It runs on our stash-box software and is primarily focused on mainstream digital scenes and studios. Instructions, invite codes, and more can be found in this guide to [Accessing StashDB](https://guidelines.stashdb.org/docs/faq_getting-started/stashdb/accessing-stashdb/).
- Several community-managed stash-box databases can also be connected to Stash in a similar manner. Each one serves a slightly different niche and follows their own methodology. A rundown of each stash-box, their differences, and the information you need to sign up can be found in the [Metadata Sources](https://docs.stashapp.cc/metadata-sources/stash-box-instances/) section of the documentation.
- Many community-maintained scrapers can also be downloaded, installed, and updated from within Stash, allowing you to pull data from a wide range of other websites and databases. They can be found by navigating to `Settings → Metadata Providers → Available Scrapers → Community (stable)`. These can be trickier to use than a stash-box because every scraper works a little differently. For more information, please visit the [CommunityScrapers repository](https://github.com/stashapp/CommunityScrapers).
- All of the above methods of scraping data into Stash are also covered in more detail in our [Guide to Scraping](https://docs.stashapp.cc/beginner-guides/guide-to-scraping/).

<sub>[StashDB](http://stashdb.org) is the canonical instance of our open source metadata API, [stash-box](https://github.com/stashapp/stash-box).</sub>

### Single Sign-On with OpenID Connect (OIDC)

Stash can hand authentication off to any OpenID Connect provider. Once OIDC is configured, the login page shows a **Login with SSO** button.

OIDC works on its own or next to the built-in username/password login:

- **OIDC only:** configure OIDC and leave username/password empty. The password form is hidden and only SSO is offered.
- **OIDC + password:** configure both, and users can sign in either way.

A configured OIDC provider counts as configured authentication, so the usual protections (such as blocking access from public IP addresses while no authentication is set up) apply.

#### 1. Register Stash with your provider

Create an OIDC client (authorization code flow) in your identity provider and set its redirect URI to:

```
https://<your-stash-host>/oidc/callback
```

Note the **issuer URL**, the **client ID** and, for confidential clients, the **client secret**. Public clients without a secret are supported, because Stash always uses PKCE.

#### 2. Configure Stash

> [!important]
> OIDC cannot be configured in the Stash web interface. Use either `config.yml` or environment variables, and restart Stash after every change.

**Where is `config.yml`?**

| How Stash runs | Location |
| --- | --- |
| Docker (`stashapp/stash` or `ghcr.io/azathoth88/stash:oidc`) | `/root/.stash/config.yml` inside the container. On the host it is in the folder mounted to `/root/.stash`, e.g. `./config/config.yml` with the [sample docker-compose.yml](docker/production/docker-compose.yml) |
| Native install | the first one found: `config.yml` in the working directory, then `~/.stash/config.yml` |
| Custom location | whatever `--config <path>` or `STASH_CONFIG_FILE` points to |

The file is created by the setup wizard on first start, so run through the setup once before adding the `oidc` section. If you are not sure which host folder a container uses, list its mounts and look for `/root/.stash`:

```
docker inspect <container> --format '{{range .Mounts}}{{.Source}} -> {{.Destination}}{{"\n"}}{{end}}'
```

(With Podman, replace `docker` with `podman`. `docker ps` / `podman ps` shows the container name.)

Add an `oidc` section to `config.yml`. Only `issuer` and `client_id` are required:

```yaml
oidc:
  issuer: https://auth.example.com
  client_id: stash
  client_secret: super-secret                          # omit for public clients
  redirect_url: https://stash.example.com/oidc/callback  # optional, derived from the request otherwise
  scopes: [profile, email, groups]                     # "openid" is always added; default: profile, email
  username_claim: preferred_username                   # default: sub

  # optional: only allow members of these groups
  groups_claim: groups                                 # default: groups
  allowed_groups:
    - stash-users
```

Alternatively, the scalar settings can be set through environment variables, which is handy for Docker. List settings (`scopes`, `allowed_groups`) can only be set in `config.yml`. Both can be combined; an environment variable takes precedence over the same key in `config.yml`.

| Variable | Config key |
| --- | --- |
| `STASH_OIDC_ISSUER` | `oidc.issuer` |
| `STASH_OIDC_CLIENT_ID` | `oidc.client_id` |
| `STASH_OIDC_CLIENT_SECRET` | `oidc.client_secret` |
| `STASH_OIDC_REDIRECT_URL` | `oidc.redirect_url` |
| `STASH_OIDC_USERNAME_CLAIM` | `oidc.username_claim` |
| `STASH_OIDC_GROUPS_CLAIM` | `oidc.groups_claim` |

Example with Docker Compose, using the image this fork publishes to GHCR:

```yaml
services:
  stash:
    image: ghcr.io/azathoth88/stash:oidc
    environment:
      STASH_OIDC_ISSUER: https://auth.example.com
      STASH_OIDC_CLIENT_ID: stash
      STASH_OIDC_CLIENT_SECRET: super-secret
      STASH_OIDC_REDIRECT_URL: https://stash.example.com/oidc/callback
      STASH_OIDC_USERNAME_CLAIM: preferred_username
```

#### 3. Restart and sign in

Restart Stash (for Docker: `docker compose restart stash`, or recreate the container with `docker compose up -d` after changing environment variables) and open the login page. Click **Login with SSO**, authenticate at your provider, and you are sent back to the page you originally requested.

#### Good to know

- **Security:** the login flow uses a random `state` (CSRF protection), a `nonce` and PKCE. The ID token signature is verified, and the return URL only accepts local paths to prevent open redirects.
- **Group restriction:** if `allowed_groups` is set, the user must be in at least one of them. Make sure the provider actually puts the groups claim into the ID token (often by requesting the `groups` scope).
- **Single-user model:** if password credentials are configured as well, OIDC sessions are bound to the configured username so that features tied to it (for example signed media URLs) keep working.
- **Reverse proxy:** the redirect URI must exactly match what Stash sends. Behind a proxy, either set `redirect_url` explicitly or make sure `X-Forwarded-Proto`, `X-Forwarded-Prefix` and `Host` are forwarded (or set `external_host`).

Full details are in [docs/OIDC.md](docs/OIDC.md).

## Support & Resources

Need help or want to get involved? Start with the documentation, then reach out to the community if you need further assistance.

### Documentation

- [Official documentation](https://docs.stashapp.cc) - official guides guides and troubleshooting.
- [In-app manual](https://docs.stashapp.cc/in-app-manual) press <kbd>Shift</kbd> + <kbd>?</kbd> in the app or view the manual online.
- [FAQ](https://discourse.stashapp.cc/c/support/faq/28) - common questions and answers.
- [Community wiki](https://discourse.stashapp.cc/tags/c/community-wiki/22/stash) - guides, how-to’s and tips.
  
### Community & Discussion

- [Community forum](https://discourse.stashapp.cc) - community support, feature requests and discussions.
- [Discord](https://discord.gg/2TsNFKt) - real-time chat and community support.
- [GitHub discussions](https://github.com/stashapp/stash/discussions) - community support and feature discussions.
- [Lemmy community](https://discuss.online/c/stashapp) - board-style community space.

### Community Scrapers & Plugins

- [Metadata sources](https://docs.stashapp.cc/metadata-sources/)
- [Plugins](https://docs.stashapp.cc/plugins/)
- [Themes](https://docs.stashapp.cc/themes/)
- [Other projects](https://docs.stashapp.cc/other-projects/)

## Architecture

You can find an overview of Stash's architecture in the [ARCHITECTURE.md](docs/ARCHITECTURE.md) document.

## Contributing

We welcome contributions and help from all humans who want to improve the project.

Before contributing, please read the [Contributing](docs/CONTRIBUTING.md) document to understand our guidelines and processes for contributing to the project.

You can learn about setting up a local development environment in the [Development](docs/DEVELOPMENT.md) document. 

## Translation

The widget below shows the current translation status of Stash across all supported languages. If you want to help us translate Stash, you can make an account at [Codeberg Translate](https://translate.codeberg.org/projects/stash/stash/) to contribute to new or existing languages. Thanks!

[![Translation status](https://translate.codeberg.org/widget/stash/stash/multi-auto.svg)](https://translate.codeberg.org/engage/stash/)