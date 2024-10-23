# syntax=docker/dockerfile:1

ARG NODE_VERSION=20

################################################################################
FROM node:${NODE_VERSION}-alpine AS base

WORKDIR /usr/src/app

ENV PNPM_HOME="/pnpm"
ENV PATH "$PNPM_HOME:$PATH"

RUN corepack enable

################################################################################
FROM base AS deps

RUN --mount=type=bind,source=package.json,target=package.json \
    --mount=type=bind,source=pnpm-lock.yaml,target=pnpm-lock.yaml \
    --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --prod --frozen-lockfile

################################################################################

FROM base AS build

COPY package.json pnpm-lock.yaml ./
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile

COPY . .
ARG NEXT_PUBLIC_BACKEND_HOST
ENV NEXT_PUBLIC_BACKEND_HOST ${BACKEND_HOST}

RUN pnpm run build

################################################################################
FROM base AS final

ENV NODE_ENV production
RUN chown node:node -R /usr/src/app

USER node

COPY package.json .

COPY --from=deps /usr/src/app/node_modules ./node_modules
COPY --from=build /usr/src/app/next.config.mjs ./next.config.mjs
COPY --from=build /usr/src/app/.next ./.next
# COPY ./public ./public

EXPOSE 3000

CMD [ "pnpm", "start" ]
