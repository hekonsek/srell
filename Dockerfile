FROM fedora:30

ADD srell /usr/bin/srell

CMD ["/usr/bin/srell"]