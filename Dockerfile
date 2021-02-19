FROM alpine as base

FROM base as be

RUN mkdir /test
RUN touch test/test
