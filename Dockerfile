ARG BUILD_FROM
FROM $BUILD_FROM

ENV LANG C.UTF-8

# Copy data for add-on
COPY inkplate6-sg-next-bus /
COPY *.ttf /
COPY run.sh /
RUN chmod a+x /run.sh

CMD [ "/run.sh" ]
