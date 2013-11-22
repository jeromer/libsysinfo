node "default" {
    Exec {
        path => [ "/bin/", "/sbin/" , "/usr/bin/", "/usr/sbin/" ],
    }

    exec { "update":
        command => "apt-get update",
        user    => "root",
    }

    # just in case
    package{"golang":
        ensure => "purged",
    }

    package { "golang-prereq":
        name    => ["lua5.1", "liblua5.1-0-dev","git-core", "bzr", "curl"],
        ensure  => present,
        require => Exec["update"]
    }

    exec { "download go1.1":
        command => "curl -s -o /tmp/go.tar.gz https://go.googlecode.com/files/go1.1.1.linux-amd64.tar.gz",
        require => Package["golang", "golang-prereq"],
    }

    exec { "extract go1.1":
        command => "tar -C /usr/local -xzf /tmp/go.tar.gz",
        user    => "root",
        require => Exec["download go1.1"],
    }

    exec { "link go1.1":
        command => "ln -s /usr/local/go/bin/go /usr/local/bin/go",
        user    => "root",
        require => Exec["extract go1.1"],
    }
}
